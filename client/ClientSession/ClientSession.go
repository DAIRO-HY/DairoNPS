package ClientSession

import (
	"DairoNPS/bridge/TCPBridgeManager"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/dao/dto"
	"DairoNPS/pool/TCPPoolManager"
	"DairoNPS/proxy/ProxyAcceptManager"
	"bufio"
	"net"
	"sync"
	"time"
)

/**
 * 服务端与客户端通信类
 * @param client 客户端DTO
 * @param clientSocket 与客户端的连接
 */
type ClientSession struct {
	Client       *dto.ClientDto
	ClientSocket net.Conn

	/**
	 * 最后一次收到客户端心跳时间
	 */
	lastHeartBeatTime int64
}

/**
 * 发送数据互斥锁
 */
var sendLock sync.Mutex

/**
 * 开始
 */
func Start(session ClientSession) {
	receive(session)
}

/**
 * 接收从客户端发来的数据
 */
func receive(session ClientSession) {
	for {
		flagData := make([]byte, 1)
		len, err := session.ClientSocket.Read(flagData)
		if len == 0 || err != nil {
			session.ClientSocket.Close()
			break
		}
		handle(session, flagData[0])
	}
}

/**
 * 处理从客户端收到的消息
 */
func handle(session ClientSession, flag byte) {
	switch flag {

	//客户端心跳
	case HeaderUtil.MAIN_HEART_BEAT:
		//println("-->接收到客户端的心跳数据${Date()}")

		//记录与客户端最后一次心跳时间戳
		session.lastHeartBeatTime = time.Now().UnixNano() / int64(time.Millisecond)

		writer := bufio.NewWriter(session.ClientSocket)

		//回复客户端心跳,TODO: 这里要考虑异常
		writer.WriteByte(HeaderUtil.MAIN_HEART_BEAT)
		writer.Flush()
	}
}

/**
 * 往客户端发送数据
 * @param flag 头部标记
 * @param message 头部消息
 */
func SendHead(session ClientSession, flag byte, message string) {
	if len(message) == 0 {
		return
	}
	data := []byte(message)
	//if (data.size > Byte.MAX_VALUE) {
	//   throw RuntimeException("一次发送数据长度不能超过${Byte.MAX_VALUE}字节")
	//}
	writer := bufio.NewWriter(session.ClientSocket)
	writer.WriteByte(flag)            //第一个字节代表类型
	writer.WriteByte(byte(len(data))) //第二个字节代表数据长度
	writer.Write(data)
	writer.Flush()
}

/**
 * 往客户端发送数据
 * @param data 要发送的数据
 * @param len 数据长度
 */
func Send(session ClientSession, data []byte, len int) {
	writer := bufio.NewWriter(session.ClientSocket)
	writer.Write(data[:len-1])
	writer.Flush()
}

/**
 * 关闭与内网穿透客户端的会话连接
 */
func Close(session ClientSession) {

	//关闭所有TCP连接池
	TCPPoolManager.CloseByClient(session.Client.Id)

	//关闭所有UDP连接池
	//try {
	//    UDPPoolManager.closeByClient(this.client.id!!)
	//} catch (e: Exception) {
	//    e.printStackTrace()
	//}
	//
	//try {
	//    //关闭正在通信的UDP连接
	//    UDPBridgeManager.closeByClient(this.client.id!!)
	//} catch (e: Exception) {
	//    e.printStackTrace()
	//}

	//关闭代理监听
	ProxyAcceptManager.CloseByClient(session.Client.Id)

	//关闭客户端所有正在通信的连接
	TCPBridgeManager.CloseByClient(session.Client.Id)

	//关闭连接
	session.ClientSocket.Close()
}
