package forward

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/dto"
	"DairoNPS/util/ForwardStatisticsUtil"
	"DairoNPS/util/TcpUtil"
	"net"
	"sync/atomic"
	"time"
)

/**
 * 数据转发桥接管理
 * @param forwardDto 隧道信息
 * @param ProxyTCP TCP代理服务端Socket
 * @param TargetTCP 内网穿透客户端Socket
 */
type ForwardBridge struct {

	// 转发明细
	ForwardDto *dto.ForwardDto

	// 转发代理端的TCP
	ProxyTCP net.Conn

	// 目标端的TCP
	TargetTCP net.Conn

	// 代理连接入方向是否被关闭
	isProxyReadClosed bool

	// 目标端的du操作关闭标识
	isTargetReadClosed bool

	//隧道流量统计
	dataSize *ForwardStatisticsUtil.ForwardDataSize

	// 创建时间(毫秒)
	CreateTime int64

	// 记录最后通信时间(毫秒)
	LastRWTime int64
}

/**
 * 开始传输数据
 */
func (mine *ForwardBridge) Start() {
DebugTimer.Add136()
	mine.dataSize = ForwardStatisticsUtil.Get(mine.ForwardDto.Id)
	go mine.receiveByForwardSendToTarget()
	go mine.receiveByTargetSendToForward()
}

/**
 * 从代理服务接收数据发送到目标端
 */
func (mine *ForwardBridge) receiveByForwardSendToTarget() {
DebugTimer.Add137()
	data := make([]uint8, NPSConstant.READ_CACHE_SIZE)
	for {
DebugTimer.Add138()
		n, readErr := mine.ProxyTCP.Read(data)
		if n > 0 {
DebugTimer.Add139()

			//记录最后通信时间
			mine.LastRWTime = time.Now().UnixMilli()

			//原子递增
			atomic.AddInt64(&mine.dataSize.InData, int64(n))

			//从代理端读取到的数据立即发送目标端
			writeErr := TcpUtil.WriteAll(mine.TargetTCP, data[:n])
			if writeErr != nil {
DebugTimer.Add140()
				break
			}
		}
		if readErr != nil {
DebugTimer.Add141()
			break
		}
	}

	//关闭代理端的读操作
	mine.ProxyTCP.(*net.TCPConn).CloseRead()

	//关闭目标端的写操作
	mine.TargetTCP.(*net.TCPConn).CloseWrite()

	//标记代理端读操作已经关闭
	mine.isProxyReadClosed = true
	mine.recycle()
}

/**
 * 从目标端接收发送到代理端
 */
func (mine *ForwardBridge) receiveByTargetSendToForward() {
DebugTimer.Add142()
	data := make([]uint8, NPSConstant.READ_CACHE_SIZE)
	for {
DebugTimer.Add143()
		n, readErr := mine.TargetTCP.Read(data)
		if n > 0 {
DebugTimer.Add144()

			//记录最后通信时间
			mine.LastRWTime = time.Now().UnixMilli()

			//原子递增
			atomic.AddInt64(&mine.dataSize.OutData, int64(n))

			//将读取到的数据立即发送客户端
			writeErr := TcpUtil.WriteAll(mine.ProxyTCP, data[:n])
			if writeErr != nil {
DebugTimer.Add145()
				break
			}
		}
		if readErr != nil {
DebugTimer.Add146()
			break
		}
	}

	//关闭目标端的读操作
	mine.TargetTCP.(*net.TCPConn).CloseRead()

	//关闭代理端的写操作
	mine.ProxyTCP.(*net.TCPConn).CloseWrite()

	//标记目标读操作已经关闭
	mine.isTargetReadClosed = true
	mine.recycle()
}

/**
 * 资源回收
 */
func (mine *ForwardBridge) recycle() {
DebugTimer.Add147()
	if mine.isProxyReadClosed && mine.isTargetReadClosed {
DebugTimer.Add148()
		mine.TargetTCP.Close()
		mine.ProxyTCP.Close()
		removeBridge(mine)
	}
}

/**
 * 关闭连接
 */
func (mine *ForwardBridge) shutdown() {
DebugTimer.Add149()
	mine.TargetTCP.Close()
	mine.ProxyTCP.Close()
}
