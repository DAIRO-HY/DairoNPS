package tcp_bridge

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/dto"
	"DairoNPS/util/ChannelStatisticsUtil"
	"DairoNPS/util/LogUtil"
	"DairoNPS/util/SecurityUtil"
	"DairoNPS/util/TcpUtil"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type RemoveBridgeFunc func(bridge *TCPBridge)

// TCPBridge TCP桥接信息
type TCPBridge struct {

	// 客户端ID
	ClientId int

	// 当前隧道
	Channel *dto.ChannelDto

	// 隧道代理端TCP链接
	ProxyTCP net.Conn

	// 与客户端的TCP链接
	ClientTCP net.Conn

	//代理连接入方向是否被关闭
	proxyInIsClosed bool

	//客户端连接入方向是否被关闭
	clientInIsClosed bool

	// 创建时间(毫秒)
	CreateTime int64

	// 记录最后通信时间(毫秒)
	LastRWTime int64

	//隧道流量统计
	channelDataSize *ChannelStatisticsUtil.ChannelDataSize
}

/**
* 开始桥接传输数据
 */
func (mine *TCPBridge) start() {
DebugTimer.Add169()
	mine.channelDataSize = ChannelStatisticsUtil.Get(mine.Channel.Id)

	//发送目标端口信息
	mine.sendHeaderToClient()
	go mine.receiveByProxySendToClient()
	mine.receiveByClientSendToProxy()
}

/**
* 发送目标端口信息
 */
func (mine *TCPBridge) sendHeaderToClient() {
DebugTimer.Add170()

	//将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	header := strconv.Itoa(mine.Channel.SecurityState) + "|" + mine.Channel.TargetPort
	headerData := []uint8(header)

	//写入数据长度标识
	data := []uint8{uint8(len(headerData))}

	//写入数据
	data = append(data, headerData...)
	err := TcpUtil.WriteAll(mine.ClientTCP, data)
	if err != nil {
DebugTimer.Add171()
		LogUtil.Error(fmt.Sprintf("往客户端发送头部失败 err:%q", err))
		mine.ClientTCP.Close()
		return
	}
}

/**
* 从代理服务接收数据发送到客户端
 */
func (mine *TCPBridge) receiveByProxySendToClient() {
DebugTimer.Add172()
	data := make([]uint8, NPSConstant.READ_CACHE_SIZE)
	for {
DebugTimer.Add173()
		n, readErr := mine.ProxyTCP.Read(data)
		if n > 0 {
DebugTimer.Add174()

			//记录最后通信时间
			mine.LastRWTime = time.Now().UnixMilli()

			//原子递增
			atomic.AddInt64(&mine.channelDataSize.InData, int64(n))
			if mine.Channel.SecurityState == 1 { //加密数据
DebugTimer.Add175()
				SecurityUtil.Mapping(data, n)
			}

			//将读取到的数据立即发送客户端
			writeErr := TcpUtil.WriteAll(mine.ClientTCP, data[:n])
			if writeErr != nil {
DebugTimer.Add176()
				break
			}
		}
		if readErr != nil {
DebugTimer.Add177()
			break
		}
	}

	//关闭客户端的输出流
	mine.ClientTCP.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	mine.ProxyTCP.(*net.TCPConn).CloseRead()

	//标记代理连接读操作已经关闭
	mine.proxyInIsClosed = true
	mine.recycle()
}

// 从客户端接收发送到代理服务器
func (mine *TCPBridge) receiveByClientSendToProxy() {
DebugTimer.Add178()
	data := make([]uint8, NPSConstant.READ_CACHE_SIZE)
	for {
DebugTimer.Add179()
		n, readErr := mine.ClientTCP.Read(data)
		if n > 0 {
DebugTimer.Add180()

			//记录最后通信时间
			mine.LastRWTime = time.Now().UnixMilli()

			//出网统计 原子递增
			atomic.AddInt64(&mine.channelDataSize.OutData, int64(n))
			if mine.Channel.SecurityState == 1 { //加密数据
DebugTimer.Add181()
				SecurityUtil.Mapping(data, n)
			}

			//将读取到的数据立即发送客户端
			writeErr := TcpUtil.WriteAll(mine.ProxyTCP, data[:n])
			if writeErr != nil {
DebugTimer.Add182()
				break
			}
		}
		if readErr != nil {
DebugTimer.Add183()
			break
		}
	}

	//关闭客户端的输出流
	mine.ProxyTCP.(*net.TCPConn).CloseWrite()

	//关闭代理端的输入流
	mine.ClientTCP.(*net.TCPConn).CloseRead()

	//标记客户端读操作已经关闭
	mine.clientInIsClosed = true
	mine.recycle()
}

/**
* 资源回收
 */
func (mine *TCPBridge) recycle() {
DebugTimer.Add184()
	if mine.proxyInIsClosed && mine.clientInIsClosed {
DebugTimer.Add185()
		mine.ClientTCP.Close()
		mine.ProxyTCP.Close()
		removeBridge(mine)
	}
}

/**
* 关闭连接
 */
func (mine *TCPBridge) shutdown() {
DebugTimer.Add186()
	mine.ClientTCP.Close()
	mine.ProxyTCP.Close()
}
