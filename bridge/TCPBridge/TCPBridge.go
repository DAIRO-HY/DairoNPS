package TCPBridge

import (
	"DairoNPS/bridge/TCPBridgeManager"
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/dao/dto"
	"bufio"
	"fmt"
	"net"
	"time"
)

/**
 * TCP桥接管理
 * @param client 客户端DTO
 * @param channel 隧道信息
 * @param proxySocket TCP代理服务端Socket
 * @param clientSocket 内网穿透客户端Socket
 */
type TCPBridge struct {
	Client       *dto.ClientDto
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
	//
	///**
	// * 是否加密数据
	// */
	//private val isEncodeData = this.channel.securityState == 1

	/**
	 * 本次连接入网总计
	 */
	inDataTotal int64

	/**
	 * 本次连接出网总计
	 */

	outDataTotal int64

	/**
	 * 代理连接入方向是否被关闭
	 */
	proxyInIsClosed bool

	/**
	 * 客户端连接入方向是否被关闭
	 */
	clientInIsClosed bool

	//val tag = System.currentTimeMillis().toString() + (Math.random() * 1000).toInt()
}

/**
* 开始传输数据
 */
func Start(bridge TCPBridge) {
	go func() {

		//发送目标端口信息
		sendHeaderToClient(bridge)
		receiveByProxySendToClient(bridge)
	}()
	go receiveByClientSendToProxy(bridge)
}

/**
* 发送目标端口信息
 */
func sendHeaderToClient(bridge TCPBridge) {

	//将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	header := fmt.Sprintf("%d|%d", bridge.Channel.SecurityState, bridge.Channel.TargetPort)
	headerData := []byte(header)

	writer := bufio.NewWriter(bridge.ClientSocket)

	//写入数据长度标识
	writer.WriteByte(byte(len(headerData)))

	//写入数据
	writer.Write(headerData)
	writer.Flush()
}

/**
* 从代理服务接收数据发送到客户端
 */
func receiveByProxySendToClient(bridge TCPBridge) {

	////客户端输出流
	//val clientOStream = this.clientSocket.outputStream
	//
	////代理Socket输入流
	//val proxyIStream = this.proxySocket.inputStream
	//val data = ByteArray(CLSConfig.READ_CACHE_SIZE)//缓存大小

	clientWriter := bufio.NewWriter(bridge.ClientSocket)

	data := make([]byte, CLSConfig.READ_CACHE_SIZE)
	for {
		length, err := bridge.ProxySocket.Read(data)
		if err != nil {
			break
		}

		//TODO:每次都计时可能影响性能
		bridge.lastSessionTime = time.Now().UnixNano() / int64(time.Millisecond) //标记最后一次读取到数据的时间

		//入网统计
		//bridge.inDataTotal = bridge.inDataTotal + length
		//bridge.Channel.InDataTotal = bridge.Channel.InDataTotal + length
		//bridge.Client.InDataTotal = bridge.Client.InDataTotal + length
		//CLSConfig.systemConfig.InDataTotal = CLSConfig.systemConfig.InDataTotal + length
		if bridge.Channel.SecurityState == 1 { //加密数据
			//SecurityUtil.mapping(data, len)
		}

		//将读取到的数据立即发送客户端
		_, err1 := clientWriter.Write(data[:length])
		if err1 != nil {
			break
		}

		//每次读取到的内容立刻发送到客户端,否则可能导致死锁
		if clientWriter.Flush() != nil {
			break
		}
	}

	//关闭客户端的输出流
	bridge.ClientSocket.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	bridge.ProxySocket.(*net.TCPConn).CloseRead()

	//标记代理连接入方向是否被关闭
	bridge.proxyInIsClosed = true
	recyle(bridge)
}

/**
* 从客户端接收发送到代理服务器
 */
func receiveByClientSendToProxy(bridge TCPBridge) {

	//客户端输入流
	//val clientIStream = this.clientSocket.inputStream

	//代理Socket输出流
	//val proxyOStream = this.proxySocket.outputStream
	proxyWriter := bufio.NewWriter(bridge.ProxySocket)

	data := make([]byte, CLSConfig.READ_CACHE_SIZE)

	for {
		length, err := bridge.ClientSocket.Read(data)
		if err != nil {
			break
		}

		//TODO:每次都计时可能影响性能
		bridge.lastSessionTime = time.Now().UnixNano() / int64(time.Millisecond) //标记最后一次读取到数据的时间

		//入网统计
		//bridge.outDataTotal = bridge.outDataTotal + length
		//bridge.Channel.OutDataTotal = bridge.Channel.OutDataTotal + length
		//bridge.Client.OutDataTotal = bridge.Client.OutDataTotal + length
		//CLSConfig.systemConfig.OutDataTotal = CLSConfig.systemConfig.OutDataTotal + length
		if bridge.Channel.SecurityState == 1 { //加密数据
			//SecurityUtil.mapping(data, len)
		}

		//将读取到的数据立即发送客户端
		_, err1 := proxyWriter.Write(data[:length])
		if err1 != nil {
			break
		}

		//每次读取到的内容立刻发送到客户端,否则可能导致死锁
		if proxyWriter.Flush() != nil {
			break
		}
	}

	//关闭客户端的输出流
	bridge.ProxySocket.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	bridge.ClientSocket.(*net.TCPConn).CloseRead()

	//标记代理连接入方向是否被关闭
	bridge.clientInIsClosed = true
	recyle(bridge)
}

/**
* 资源回收
 */
func recyle(bridge TCPBridge) {
	if bridge.proxyInIsClosed && bridge.clientInIsClosed {
		bridge.ClientSocket.Close()
		bridge.ProxySocket.Close()
		TCPBridgeManager.RemoveBridgeList(bridge)
	}
}

/**
* 关闭连接
 */
func Close(bridge TCPBridge) {
	bridge.ClientSocket.Close()
	bridge.ProxySocket.Close()
}
