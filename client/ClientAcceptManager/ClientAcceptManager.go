package ClientAcceptManager

import (
	"DairoNPS/client/ClientSessionManager"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/pool/TCPPoolManager"
	"fmt"
	"log"
	"net"
	"time"
)

// 接受并HOLD住客户端的socket

/**
 * Socket服务端
 */
var socketServer net.Listener

// Start 启动
func Start() {
	//go accept()
	//time.Sleep(5 * time.Second)
	accept()
}

/**
 * 监听客户端连接
 */
func accept() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", CLSConfig.CLIENT_TO_SERVER_TCP_PORT))
	socketServer = listen

	//最终关闭监听
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatalf("%q\n", err)
		}
	}(listen)

	//启动失败
	if err != nil {
		println("启动失败，请参考错误信息。")
		log.Fatalf("%q\n", err)
	}

	fmt.Printf("端口:%d监听成功。\n", CLSConfig.CLIENT_TO_SERVER_TCP_PORT)
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	listen.Close()
	//}()
	for {

		//等待客户端连接
		fmt.Println("-->监听客户端连接")
		socketClient, err := socketServer.Accept()
		if err != nil {
			log.Fatalf("%q\n", err)
		}
		fmt.Println("-->接收到客户端连接请求")
		go handle(socketClient)
	}
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
	socketClient.SetReadDeadline(time.Now().Add(300 * time.Second))
	//socketClient.SetReadDeadline(time.Time{})

	//println("-->读取标记开始${System.currentTimeMillis()}")
	//读取第一个标记字节,通过该自己判断该连接类型

	flagData := make([]byte, 1)
	len, err := socketClient.Read(flagData)
	if len == 0 || err != nil {
		log.Println("-->从客户端读取数据超时")

		//没必要继续执行，直接关闭客户端连接
		socketClient.Close()
		return
	}
	//println("-->读取标记结束${System.currentTimeMillis()}:$flag")
	switch flagData[0] {

	//标记该连接为:服务器端往客户端发送指令的连接
	case HeaderUtil.CLIENT_TO_SERVER_MAIN_CONNECTION:
		ClientSessionManager.Validate(socketClient)

	//创建客户端Socket连接池
	case HeaderUtil.SERVER_TCP_POOL_REQUEST:
		TCPPoolManager.Add(socketClient)
	}
}
