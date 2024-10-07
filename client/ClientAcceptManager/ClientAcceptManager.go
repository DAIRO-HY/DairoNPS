package ClientAcceptManager

import (
	"DairoNPS/client/ClientSessionManager"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/pool/TCPPoolManager"
	"fmt"
	"net"
	"time"
)

// 接受并HOLD住客户端的socket

/**
 * Socket服务端
 */
var socketServer net.Listener

/**
 * 开启服务
 */
func Start() {
	go accept()
}

/**
 * 监听客户端连接
 */
func accept() {
	lisen, err := net.Listen("tcp", fmt.Sprintf(":%d", CLSConfig.CLIENT_TO_SERVER_TCP_PORT))
	socketServer = lisen
	if err != nil {
		//TODO:启动监听失败
		return
	}
	//(CLSConfig.CLIENT_TO_SERVER_TCP_PORT)
	for {

		//等待客户端连接
		//println("-->监听客户端连接开始")
		socketClient, err := socketServer.Accept()
		if err != nil {
			break
		}
		//println("-->接收到客户端连接请求")
		handle(socketClient)
	}
	socketServer.Close()
	socketServer = nil
	//println("-->监听客户端连接结束")
}

/**
 * 分配连接
 * @param socketClient 与客户端的连接
 */
func handle(socketClient net.Conn) {

	//保持长连接
	//socketClient.keepAlive = true

	//得到输入流
	//val clientIStream = socketClient.inputStream

	//读取连接的第一个数据,设置超时,避免恶意连接
	socketClient.SetReadDeadline(time.Now().Add(5 * time.Second))

	//println("-->读取标记开始${System.currentTimeMillis()}")
	//读取第一个标记字节,通过该自己判断该连接类型

	flagData := make([]byte, 1)
	len, err := socketClient.Read(flagData)
	if len == 0 || err != nil {
		socketClient.Close()
		return
	}
	//println("-->读取标记结束${System.currentTimeMillis()}:$flag")
	switch flagData[0] {

	//标记该连接为:服务器端往客户端发送指令的连接
	case HeaderUtil.CLIENT_TO_SERVER_MAIN_CONNECTION:
		{
			ClientSessionManager.Validate(socketClient)
		}

	//创建客户端Socket连接池
	case HeaderUtil.SERVER_TCP_POOL_REQUEST:
		{
			TCPPoolManager.Add(socketClient)
		}
	}
}
