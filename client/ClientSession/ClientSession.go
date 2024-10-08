package ClientSession

import (
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/dao/dto"
	"DairoNPS/util/SecurityUtil"
	"bufio"
	"errors"
	"log"
	"net"
	"strconv"
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
func Start(session ClientSession, recycle func(clientId int)) {
	if sendClientId(session) != nil {
		recycle(session.Client.Id)
		return
	}
	if sendClientSecurityKey(session) != nil {
		recycle(session.Client.Id)
		return
	}
	receive(session)
	recycle(session.Client.Id)
}

// 将客户端id返回给客户端
func sendClientId(session ClientSession) error {
	clientId := session.Client.Id
	return SendHead(session, HeaderUtil.SERVER_TO_CLIENT_ID, strconv.Itoa(clientId))
}

// 将加密秘钥发送到客户端
func sendClientSecurityKey(session ClientSession) error {
	return Send(session, SecurityUtil.ClientKeyArray, 128)
}

/**
 * 接收从客户端发来的数据
 */
func receive(session ClientSession) {
	for {
		flagData := make([]byte, 1)
		length, err := session.ClientSocket.Read(flagData)
		if length == 0 || err != nil {
			log.Printf("-->接收客户端数据标识出错。len:%d  err:%q\n", length, err)
			break
		}
		flag := flagData[0]
		if handle(session, flag) != nil {
			log.Printf("-->处理客户端数据失败。err:%q\n", err)
			break
		}
	}
}

/**
 * 处理从客户端收到的消息
 */
func handle(session ClientSession, flag byte) error {
	switch flag {
	case HeaderUtil.MAIN_HEART_BEAT:
		log.Println("-->接收到客户端的心跳数据")

		//记录与客户端最后一次心跳时间戳
		session.lastHeartBeatTime = time.Now().UnixNano() / int64(time.Millisecond)
		return Send(session, []byte{HeaderUtil.MAIN_HEART_BEAT}, 1)
	default:
		return errors.New("未知的Flag")
	}
}

/**
 * 往客户端发送数据
 * @param flag 头部标记
 * @param message 头部消息
 */
func SendHead(session ClientSession, flag byte, message string) error {
	data := []byte(message)
	//if (data.size > Byte.MAX_VALUE) {
	//   throw RuntimeException("一次发送数据长度不能超过${Byte.MAX_VALUE}字节")
	//}

	//第1个字节代表类型，第2个字节代表数据长度
	buffer := []byte{flag, byte(len(data))}
	buffer = append(buffer, data[:]...)
	return Send(session, buffer, len(buffer))
}

/**
 * 往客户端发送数据
 * @param data 要发送的数据
 * @param len 数据长度
 */
func Send(session ClientSession, data []byte, len int) error {
	sendLock.Lock()

	//TODO:这里应该要设置写入超时
	//session.ClientSocket.SetWriteDeadline()
	writer := bufio.NewWriter(session.ClientSocket)
	_, err := writer.Write(data[:len])
	if err != nil {
		sendLock.Unlock()
		return err
	}
	err = writer.Flush()
	sendLock.Unlock()
	return err
}

/**
 * 关闭与内网穿透客户端的会话连接
 */
func Close(session ClientSession) {
	//
	////关闭所有TCP连接池
	//TCPPoolManager.CloseByClient(session.Client.Id)
	//
	////关闭所有UDP连接池
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
	//
	////关闭代理监听
	//ProxyAcceptManager.CloseByClient(session.Client.Id)
	//
	////关闭客户端所有正在通信的连接
	//TCPBridgeManager.CloseByClient(session.Client.Id)

	//关闭连接
	session.ClientSocket.Close()
}

// 回收数据
//func recycle(session ClientSession) {
//	//Close(session)
//}
