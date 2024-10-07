package ProxyTCPAccept

import (
	"DairoNPS/bridge/TCPBridgeManager"
	"DairoNPS/dao/dto"
	"DairoNPS/pool/TCPPoolManager"
	"fmt"
	"net"
	"time"
)

/**
 * TCP隧道代理
 */
type ProxyTCPAccept struct {
	Client  *dto.ClientDto
	Channel *dto.ChannelDto

	/**
	 * 最后一次统计到入网流量
	 * 流量统计时用到
	 */
	lastInDataTotal int64

	/**
	 * 最后一次统计到出网流量
	 * 流量统计时用到
	 */
	lastOutDataTotal int64

	/**
	 * 标记监听已经结束
	 */
	isFinished bool

	/**
	 * 代理SockerServer
	 */
	proxySocketServer net.Listener
}

/**
 * 访问控制的IP地址
 */
//private val aclIpSet = ChannelAclDao.selectByChannelId(channel.id!!).map {
//    it.ip!!
//}.toSet()

/**
 * 开始监听端口
 */
func Start(tcpProxy ProxyTCPAccept) {
	accept(tcpProxy)
}

/**
 * 等待客户端连接
 */
func accept(tcpProxy ProxyTCPAccept) {
	channel := tcpProxy.Channel
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", channel.ServerPort))
	if err != nil {
		fmt.Sprintf("端口:%d 监听失败", channel.ServerPort)
		return
	}
	tcpProxy.proxySocketServer = listener
	for {

		//代理服务端Socket
		proxySocket, err := listener.Accept()
		if err != nil {
			break
		}

		if !hasAccess(tcpProxy, proxySocket) { //判断是否有访问权限
			continue
		}

		//NPS客户端Socket
		clientSocket := TCPPoolManager.GetAndAddPool(channel.ClientId)
		if clientSocket == nil { //没有可用的Socket
			proxySocket.Close()
			continue
		}
		TCPBridgeManager.Start(tcpProxy.Client, channel, proxySocket, clientSocket)
	}
	fmt.Printf("-->端口:d%监听结束", channel.ServerPort)
	tcpProxy.isFinished = true
}

/**
 * 判断是否有访问权限
 */
func hasAccess(tcpProxy ProxyTCPAccept, proxySocket net.Conn) bool {
	//if tcpProxy.channel.AclState == 0 { //访问权限处于关闭状态
	//    return true
	//}
	//
	////获取IP地址
	////不能使用packet.address.hostName,会出现延迟
	//ip := proxySocket.inetAddress.hostAddress
	//if(this.channel.aclState == 1){//白名单模式
	//    if(this.aclIpSet.contains(ip)){
	//        return true
	//    }
	//}
	//if(this.channel.aclState == 2){//黑名单模式
	//    if(!this.aclIpSet.contains(ip)){
	//        return true
	//    }
	//}
	//println("ip:${ip}被拒绝访问")
	//proxySocket.getOutputStream().write("您当前的ip地址处于黑名单状态,被禁止访问".toByteArray())
	//proxySocket.close()
	//return false
	return true
}

/**
 * 停止监听端口
 */
func Close(tcpProxy *ProxyTCPAccept) {
	for {
		time.Sleep(100 * time.Millisecond)
		if tcpProxy.proxySocketServer == nil {
			continue
		}
		tcpProxy.proxySocketServer.Close()
		if tcpProxy.isFinished {
			break
		}
	}
	//ProxyAcceptManager.removeByChannelId(this.channel.id)
}
