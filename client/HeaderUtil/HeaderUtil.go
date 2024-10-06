package HeaderUtil

import (
	"net"
	"time"
)

//NPS客户端头部标记

/**
 * 客户端与服务器端通信连接标记
 */
const CLIENT_TO_SERVER_MAIN_CONNECTION = 0

/**
 * 与客户端通信心跳标记
 */
const MAIN_HEART_BEAT = 1

/**
 * 向客户端发送clientId
 */
const SERVER_TO_CLIENT_ID = 2

/**
 * 向客户端申请TCP连接池请求
 */
const SERVER_TCP_POOL_REQUEST = 3

/**
 * 向客户端申请UDP连接池请求
 */
const SERVER_UDP_POOL_REQUEST = 4

/**
 * 服务器向客户端同步当前处于激活状态的UDP连接池端口
 */
const SYNC_ACTIVE_POOL_UDP_PORT = 5

/**
 * 向客户端同步当前处于激活状态的UDP连接端口
 */
const SYNC_ACTIVE_BRIDGE_UDP_PORT = 6

/**
 * 向客户端发送clientId
 */
const SECURITY_CLIENT_KEY = 7

/**
 * 获取客户端Socket头部信息
 */
func GetHeader(clientSocket net.Conn) string {

	// 设置读取超时时间为 3 秒
	deadline := time.Now().Add(3 * time.Second)
	clientSocket.SetReadDeadline(deadline)

	data := make([]byte, 1)

	//读取一个字节,该字节代表key长度
	length, err := clientSocket.Read(data)
	if err != nil {
		clientSocket.Close()
		return ""
	}
	if length != 1 {
		clientSocket.Close()
		return ""
	}

	//得到头部部分数据长度
	headerLen := data[0]

	//记录已经读取到的数据大小
	var readedLength byte = 0
	headerData := make([]byte, headerLen)
	for {
		buffer := make([]byte, headerLen-readedLength)
		length, err := clientSocket.Read(buffer)
		if err != nil {
			clientSocket.Close()
			return ""
		}
		copy(headerData[readedLength:readedLength+byte(length)-1], buffer[0:length-1])
		readedLength += byte(length)
		if headerLen == headerLen {
			break
		}
	}

	//设置读取数据超时
	clientSocket.SetReadDeadline(time.Now().Add(99999999999 * time.Second))
	return string(headerData)
}
