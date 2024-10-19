package nps_client

import (
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_client/HeaderUtil"
	"DairoNPS/util/SecurityUtil"
	"DairoNPS/util/TcpUtil"
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
	Client *dto.ClientDto
	tcp    net.Conn

	//最后一次收到客户端心跳时间
	lastHeartBeatTime int64
}

/**
 * 发送数据互斥锁
 */
var sendLock sync.Mutex

/**
 * 开始
 */
func (mine *ClientSession) Start() {
	finallyFunc := func() {
		time.Sleep(1 * time.Second)
		mine.Close()
		removeSession(mine)
	}

	//设置长连接
	mine.tcp.SetReadDeadline(time.Time{})
	if mine.sendClientId() != nil {
		finallyFunc()
		return
	}
	if mine.sendClientSecurityKey() != nil {
		finallyFunc()
		return
	}
	mine.receive()
	finallyFunc()
}

// 将客户端id返回给客户端
func (mine *ClientSession) sendClientId() error {
	clientId := mine.Client.Id
	return mine.SendHead(HeaderUtil.SERVER_TO_CLIENT_ID, strconv.Itoa(clientId))
}

// 将加密秘钥发送到客户端
func (mine *ClientSession) sendClientSecurityKey() error {
	return mine.Send(SecurityUtil.ClientKeyArray)
}

/**
 * 接收从客户端发来的数据
 */
func (mine *ClientSession) receive() {
	for {
		flagData, err := TcpUtil.ReadNByte(mine.tcp, 1)
		if err != nil {
			log.Printf("-->接收客户端数据标识出错。 err:%q\n", err)
			break
		}
		flag := flagData[0]
		if mine.handle(flag) != nil {
			log.Printf("-->处理客户端数据失败。err:%q\n", err)
			break
		}
	}
}

/**
 * 处理从客户端收到的消息
 */
func (mine *ClientSession) handle(flag byte) error {
	switch flag {
	case HeaderUtil.MAIN_HEART_BEAT:
		log.Println("-->接收到客户端的心跳数据")

		//记录与客户端最后一次心跳时间戳
		mine.lastHeartBeatTime = time.Now().UnixNano() / int64(time.Millisecond)
		return mine.Send([]byte{HeaderUtil.MAIN_HEART_BEAT})
	default:
		return errors.New("未知的Flag")
	}
}

/**
 * 往客户端发送数据
 * @param flag 头部标记
 * @param message 头部消息
 */
func (mine *ClientSession) SendHead(flag byte, message string) error {
	data := []byte(message)
	//if (data.size > Byte.MAX_VALUE) {
	//   throw RuntimeException("一次发送数据长度不能超过${Byte.MAX_VALUE}字节")
	//}

	//第1个字节代表类型，第2个字节代表数据长度
	buffer := []byte{flag, byte(len(data))}
	buffer = append(buffer, data...)
	return mine.Send(buffer)
}

/**
 * 往客户端发送数据
 * @param data 要发送的数据
 */
func (mine *ClientSession) Send(data []byte) error {
	sendLock.Lock()
	err := TcpUtil.WriteAll(mine.tcp, data)
	sendLock.Unlock()
	return err
}

/**
 * 关闭与内网穿透客户端的会话连接
 */
func (mine *ClientSession) Close() {
	mine.tcp.Close()
}

// 客户端是否在线监测
func (mine *ClientSession) IsOnline() bool {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	//在指定时间内没有收到客户端心跳,则视为离线
	if now-mine.lastHeartBeatTime > CLSConfig.HEART_TIME*2 {
		return false
	}
	return true
}
