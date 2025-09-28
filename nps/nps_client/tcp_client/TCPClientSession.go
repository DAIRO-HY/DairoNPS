package tcp_client

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_client/HeaderUtil"
	"DairoNPS/util/LogUtil"
	"DairoNPS/util/SecurityUtil"
	"DairoNPS/util/WriterUtil"
	"errors"
	"fmt"
	"io"
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

// 开始会话
func (mine *ClientSession) Start() {
	mine.sendServerInfoAndReceive()
	//time.Sleep(1 * time.Second)
	mine.Shutdown()
	removeSession(mine)
}

// 发送服务器端的信息
func (mine *ClientSession) sendServerInfoAndReceive() {

	//设置长连接
	mine.tcp.SetReadDeadline(time.Time{})
	if mine.sendClientId() != nil {
		return
	}
	if mine.sendClientSecurityKey() != nil {
		return
	}
	mine.receive()
}

// 将客户端id返回给客户端
func (mine *ClientSession) sendClientId() error {
	clientId := mine.Client.Id
	return mine.SendHead(HeaderUtil.SERVER_TO_CLIENT_ID, strconv.Itoa(clientId))
}

// 将加密秘钥发送到客户端
func (mine *ClientSession) sendClientSecurityKey() error {
	return mine.Send(SecurityUtil.ClientSecurityKey[:])
}

/**
 * 接收从客户端发来的数据
 */
func (mine *ClientSession) receive() {
	for {
		flagData := make([]byte, 1)
		if _, err := io.ReadFull(mine.tcp, flagData); err != nil {
			LogUtil.Error(fmt.Sprintf("接收客户端数据标识出错。 err:%q", err))
			break
		}
		flag := flagData[0]

		if err := mine.handle(flag); err != nil {
			LogUtil.Error(fmt.Sprintf("处理客户端数据失败。 err:%q", err))
			break
		}
	}
}

/**
 * 处理从客户端收到的消息
 */
func (mine *ClientSession) handle(flag uint8) error {
	switch flag {
	case HeaderUtil.MAIN_HEART_BEAT:
		//log.Println("-->接收到客户端的心跳数据")

		//记录与客户端最后一次心跳时间戳
		mine.lastHeartBeatTime = time.Now().UnixMilli()
		return mine.Send([]uint8{HeaderUtil.MAIN_HEART_BEAT})
	default:
		return errors.New(fmt.Sprintf("未知的Flag:%d", flag))
	}
}

/**
 * 往客户端发送数据
 * @param flag 头部标记
 * @param message 头部消息
 */
func (mine *ClientSession) SendHead(flag uint8, message string) error {
	data := []uint8(message)
	//if (data.size > Byte.MAX_VALUE) {
	//   throw RuntimeException("一次发送数据长度不能超过${Byte.MAX_VALUE}字节")
	//}

	//第1个字节代表类型，第2个字节代表数据长度
	buffer := []uint8{flag, uint8(len(data))}
	buffer = append(buffer, data...)
	return mine.Send(buffer)
}

/**
 * 往客户端发送数据
 * @param data 要发送的数据
 */
func (mine *ClientSession) Send(data []uint8) error {
	sendLock.Lock()
	err := WriterUtil.WriteFull(mine.tcp, data)
	sendLock.Unlock()
	return err
}

/**
 * 关闭与内网穿透客户端的会话连接
 */
func (mine *ClientSession) Shutdown() {
	mine.tcp.Close()
}

// 客户端是否在线监测
func (mine *ClientSession) IsOnline() bool {
	now := time.Now().UnixMilli()

	//在指定时间内没有收到客户端心跳,则视为离线
	if now-mine.lastHeartBeatTime > NPSConstant.HEART_TIME*2 {
		return false
	}
	return true
}
