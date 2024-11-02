package tcp_proxy

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_bridge/tcp_bridge"
	"DairoNPS/nps/nps_pool/tcp_pool"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
)

/**
 * TCP隧道代理
 */
type TCPProxyAccept struct {
	Client  *dto.ClientDto
	Channel *dto.ChannelDto

	/**
	 * 最后一次统计到入网流量
	 * 流量统计时用到
	 */
	lastInData int64

	/**
	 * 最后一次统计到出网流量
	 * 流量统计时用到
	 */
	lastOutData int64

	//标记监听已经结束
	isFinished bool

	//代理端口监听服务
	listen net.Listener
}

/**
 * 访问控制的IP地址
 */
//private val aclIpSet = ChannelAclDao.selectByChannelId(channel.id!!).map {
//    it.ip!!
//}.toSet()

/**
 * 等待客户端连接
 */
func (mine *TCPProxyAccept) accept() {
	for {

		//代理服务端Socket
		proxySocket, err := mine.listen.Accept()
		if err != nil {
			LogUtil.Info(fmt.Sprintf("端口:%d 监听结束\n", mine.Channel.ServerPort))
			break
		}
		LogUtil.Info(fmt.Sprintf("端口:%d 监听到一个连接\n", mine.Channel.ServerPort))

		//NPS客户端Socket
		clientSocket := tcp_pool.GetAndAddPool(mine.Channel.ClientId)
		if clientSocket == nil {
			LogUtil.Error(fmt.Sprintf("客户端: %d没有可用的连接池。", mine.Channel.ClientId))
			proxySocket.Close()
			continue
		}
		tcp_bridge.MakeBridge(mine.Client, mine.Channel, proxySocket, clientSocket)
	}
	mine.isFinished = true
}
