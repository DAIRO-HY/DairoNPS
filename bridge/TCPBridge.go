package bridge

import (
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/dao/dto"
	"DairoNPS/util/StatisticsUtil"
	"DairoNPS/util/TcpUtil"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type RemoveBridgeFunc func(bridge *TCPBridge)

/**
 * TCP桥接管理
 * @param client 客户端DTO
 * @param channel 隧道信息
 * @param proxySocket TCP代理服务端Socket
 * @param clientSocket 内网穿透客户端Socket
 */
type TCPBridge struct {
	ClientId     int
	Channel      *dto.ChannelDto
	ProxySocket  net.Conn
	ClientSocket net.Conn

	///**
	// * 连接关闭同步锁
	// */
	//private val closeLock = Mutex()
	//
	/**
	 * 最后一次读取到数据的时间,用来判断Socket是否存活
	 */
	lastSessionTime int64

	/**
	 * 代理连接入方向是否被关闭
	 */
	proxyInIsClosed bool

	/**
	 * 客户端连接入方向是否被关闭
	 */
	clientInIsClosed bool

	//val tag = System.currentTimeMillis().toString() + (Math.random() * 1000).toInt()

	//隧道流量统计
	channelDataSize *StatisticsUtil.ChannelDataSizeLog
}

/**
* 开始桥接传输数据
 */
func (mine *TCPBridge) start() {
	channelId := mine.Channel.Id
	mine.channelDataSize = StatisticsUtil.ChannelDataSizeMap[channelId]
	go func() {

		//发送目标端口信息
		mine.sendHeaderToClient()
		mine.receiveByProxySendToClient()
	}()
	go mine.receiveByClientSendToProxy()
}

/**
* 发送目标端口信息
 */
func (mine *TCPBridge) sendHeaderToClient() {

	//将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	header := strconv.Itoa(mine.Channel.SecurityState) + "|" + mine.Channel.TargetPort
	headerData := []byte(header)

	//写入数据长度标识
	data := []byte{byte(len(headerData))}

	//写入数据
	data = append(data, headerData...)
	err := TcpUtil.WriteAll(mine.ClientSocket, data)
	if err != nil {
		mine.ClientSocket.Close()
	}
}

/**
* 从代理服务接收数据发送到客户端
 */
func (mine *TCPBridge) receiveByProxySendToClient() {
	data := make([]byte, CLSConfig.READ_CACHE_SIZE)
	for {
		length, err := mine.ProxySocket.Read(data)
		if err != nil {
			break
		}

		//TODO:每次都计时可能影响性能
		mine.lastSessionTime = time.Now().UnixNano() / int64(time.Millisecond) //标记最后一次读取到数据的时间

		//入网统计
		//mine.inDataTotal = mine.inDataTotal + int64(length)
		//bridge.Channel.InDataTotal = bridge.Channel.InDataTotal + length
		//bridge.Client.InDataTotal = bridge.Client.InDataTotal + length
		//CLSConfig.systemConfig.InDataTotal = CLSConfig.systemConfig.InDataTotal + length

		//原子递增
		atomic.AddInt64(&mine.channelDataSize.InData, int64(length))

		if mine.Channel.SecurityState == 1 { //加密数据
			//SecurityUtil.mapping(data, len)
		}

		//将读取到的数据立即发送客户端
		err = TcpUtil.WriteAll(mine.ClientSocket, data[:length])
		if err != nil {
			break
		}
	}

	//关闭客户端的输出流
	mine.ClientSocket.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	mine.ProxySocket.(*net.TCPConn).CloseRead()

	//标记代理连接读操作已经关闭
	mine.proxyInIsClosed = true
	mine.recycle()
}

/**
* 从客户端接收发送到代理服务器
 */
func (mine *TCPBridge) receiveByClientSendToProxy() {
	data := make([]byte, CLSConfig.READ_CACHE_SIZE)
	for {
		length, err := mine.ClientSocket.Read(data)
		if err != nil {
			break
		}

		//TODO:每次都计时可能影响性能
		mine.lastSessionTime = time.Now().UnixNano() / int64(time.Millisecond) //标记最后一次读取到数据的时间

		//入网统计
		//mine.outDataTotal = mine.outDataTotal + int64(length)
		//bridge.Channel.OutDataTotal = bridge.Channel.OutDataTotal + length
		//bridge.Client.OutDataTotal = bridge.Client.OutDataTotal + length
		//CLSConfig.systemConfig.OutDataTotal = CLSConfig.systemConfig.OutDataTotal + length

		//出网统计 原子递增
		atomic.AddInt64(&mine.channelDataSize.OutData, int64(length))
		if mine.Channel.SecurityState == 1 { //加密数据
			//SecurityUtil.mapping(data, len)
		}

		//将读取到的数据立即发送客户端
		err = TcpUtil.WriteAll(mine.ProxySocket, data[:length])
		if err != nil {
			break
		}
	}

	//关闭客户端的输出流
	mine.ProxySocket.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	mine.ClientSocket.(*net.TCPConn).CloseRead()

	//标记客户端读操作已经关闭
	mine.clientInIsClosed = true
	mine.recycle()
}

/**
* 资源回收
 */
func (mine *TCPBridge) recycle() {
	if mine.proxyInIsClosed && mine.clientInIsClosed {
		mine.ClientSocket.Close()
		mine.ProxySocket.Close()
		remove(mine)
	}
}

/**
* 关闭连接
 */
func (mine *TCPBridge) shutdown() {
	mine.ClientSocket.Close()
	mine.ProxySocket.Close()
}
