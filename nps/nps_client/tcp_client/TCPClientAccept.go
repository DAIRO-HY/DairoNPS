package tcp_client

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_client/HeaderUtil"
	"DairoNPS/nps/nps_pool/tcp_pool"
	"DairoNPS/util/LogUtil"
	"DairoNPS/util/TcpUtil"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Accept 监听客户端连接
func Accept() {
	listen, err := net.Listen("tcp", ":"+NPSConstant.TcpPort)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("监听客户端监听失败，请参考错误信息。err:%q", err))
		log.Fatal(err)
	}
	defer listen.Close()
	LogUtil.Info(fmt.Sprintf("端口:%s监听成功。\n", NPSConstant.TcpPort))
	for {
		LogUtil.Debug(fmt.Sprintf("监听客户端连接,端口:%s监听成功。", NPSConstant.TcpPort))

		//等待客户端连接
		tcp, err := listen.Accept()
		if err != nil {
			LogUtil.Error(fmt.Sprintf("监听客户端结束,端口:%s", NPSConstant.TcpPort))
			log.Fatal(err)
		}
		LogUtil.Debug(fmt.Sprintf("接收到客户端连接请求,端口:%s监听成功。", NPSConstant.TcpPort))
		go handleAccept(tcp)
	}
}

/**
 * 分配连接
 * @param socketClient 与客户端的连接
 */
func handleAccept(tcp net.Conn) {

	//读取连接的第一个数据,设置超时,避免恶意连接
	tcp.SetReadDeadline(time.Now().Add(3 * time.Second))

	//读取第一个标记字节,通过该自己判断该连接类型
	flagData, err := TcpUtil.ReadNByte(tcp, 1)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("从客户端读取标识失败,可能来自一个非法客户端,IP:%s,err:%q", tcp.RemoteAddr().String(), err))

		//没必要继续执行，直接关闭客户端连接
		tcp.Close()
		return
	}
	//println("-->读取标记结束${System.currentTimeMillis()}:$flag")
	switch flagData[0] {

	//标记该连接为:服务器端往客户端发送指令的连接
	case HeaderUtil.CLIENT_TO_SERVER_MAIN_CONNECTION:
		validateSession(tcp)

	//创建客户端Socket连接池
	case HeaderUtil.REQUEST_TCP_POOL:
		tcp_pool.Add(tcp)
	}
}

// 验证客户端回话
func validateSession(tcp net.Conn) {

	//得到头部数据
	header, err := HeaderUtil.GetHeader(tcp)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("从客户端读取头部数据失败,IP:%s,err:%q", tcp.RemoteAddr().String(), err))
		tcp.Close()
		return
	}
	headers := strings.Split(header, "|")

	//得到客户端key
	key := headers[0]
	client := ClientDao.SelectByKey(key)
	if client == nil {
		LogUtil.Error(fmt.Sprintf("客户端：%s不存在,IP:%s", key, tcp.RemoteAddr().String()))
		tcp.Close()
		return
	}
	if client.EnableState == 0 {
		LogUtil.Error(fmt.Sprintf("客户端：%s已被停止服务,IP:%s", key, tcp.RemoteAddr().String()))
		tcp.Close()
		return
	}

	//设置客户端登录信息-------------------------------------------------------------------------------START
	remoteAddr := tcp.RemoteAddr().String()

	//客户端ip
	ip := strings.Split(remoteAddr, ":")[0]

	//从头部信息中得到客户端版本号
	version := headers[1]
	loginClientDto := dto.ClientDto{
		Id:      client.Id,
		Ip:      ip,
		Version: version,
	}
	ClientDao.SetClientInfo(loginClientDto)
	//设置客户端登录信息-------------------------------------------------------------------------------END

	holdOnClient(client, tcp)
}
