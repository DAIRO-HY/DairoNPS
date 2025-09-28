package HeaderUtil

import (
	"io"
	"net"
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
const REQUEST_TCP_POOL = 3

/**
 * 向客户端申请UDP连接池请求
 */
const REQUEST_UDP_POOL = 4

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
func GetHeader(clientSocket net.Conn) (string, error) {

	//读取一个字节,该字节代表key长度
	headerLenData := make([]byte, 1)
	if _, err := io.ReadFull(clientSocket, headerLenData); err != nil {
		return "", err
	}

	//得到头部部分数据长度
	headerLen := headerLenData[0]
	headerData := make([]byte, headerLen)
	if _, err := io.ReadFull(clientSocket, headerData); err != nil {
		return "", err
	}
	return string(headerData), nil
}
