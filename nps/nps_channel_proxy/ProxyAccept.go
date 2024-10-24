package nps_channel_proxy

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_bridge"
	"DairoNPS/nps/nps_pool"
	"fmt"
	"net"
)

/**
 * TCP隧道代理
 */
type ProxyAccept struct {
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
func (mine *ProxyAccept) accept() {
	for {

		//代理服务端Socket
		proxySocket, err := mine.listen.Accept()
		fmt.Printf("端口:%d 监听到一个连接\n", mine.Channel.ServerPort)
		if err != nil {
			fmt.Printf("-->端口:%d 监听结束\n", mine.Channel.ServerPort)
			break
		}

		if !mine.hasAccess(proxySocket) { //判断是否有访问权限
			proxySocket.Close()
			continue
		}

		//NPS客户端Socket
		clientSocket := nps_pool.GetAndAddPool(mine.Channel.ClientId)
		if clientSocket == nil { //没有可用的Socket
			proxySocket.Close()
			continue
		}
		nps_bridge.MakeBridge(mine.Client, mine.Channel, proxySocket, clientSocket)
	}
	mine.isFinished = true
}

/**
 * 判断是否有访问权限
 */
func (mine *ProxyAccept) hasAccess(proxySocket net.Conn) bool {
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
