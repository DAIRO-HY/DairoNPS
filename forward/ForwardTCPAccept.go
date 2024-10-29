package forward

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strings"
)

/**
 * 数据转发Socket请求等待
 */
type ForwardTCPAccept struct {

	//代理端口监听服务
	listen net.Listener

	forwardDto *dto.ForwardDto

	///**
	// * 最后一次统计到入网流量
	// * 流量统计时用到
	// */
	//var lastInDataTotal: Long? = null
	//
	///**
	// * 最后一次统计到出网流量
	// * 流量统计时用到
	// */
	//var lastOutDataTotal: Long? = null

	/**
	 * 标记监听已经结束
	 */
	isFinished bool

	/**
	 * 代理SockerServer
	 */
	forwardSocketServer net.Conn

	/**
	 * 访问控制的IP地址
	 */
	//private val aclIpSet = ForwardAclDao.selectByForwardId(forwardDto.id!!).map {
	//    it.ip!!
	//}.toSet()
}

/**
 * 等待客户端连接
 */
func (mine *ForwardTCPAccept) accept() {
	//val forwardDto = this.forwardDto

	//this.forwardSocketServer = ServerSocket(forwardDto.port!!)
	for {

		//代理服务端Socket
		proxyTCP, err := mine.listen.Accept()
		if err != nil {
			LogUtil.Info(fmt.Sprintf("转发端口:%d 监听结束\n", mine.forwardDto.Port))
			break
		}
		if !mine.hasAccess(proxyTCP) { //判断是否有访问权限
			proxyTCP.Close()
			continue
		}
		targetIpAndPort := mine.forwardDto.TargetPort
		if !strings.Contains(targetIpAndPort, ":") {
			targetIpAndPort = "127.0.0.1:" + targetIpAndPort
		}

		//目标服务器Socket连接
		targetTCP, err := net.Dial("tcp", targetIpAndPort)
		if err != nil {
			proxyTCP.Close()
			LogUtil.Debug(fmt.Sprintf("转发端口:%d 连接失败\n", mine.forwardDto.Port))
			continue
		}

		//开始桥接
		startBridge(mine.forwardDto, proxyTCP, targetTCP)
	}
	mine.listen.Close()
	LogUtil.Debug(fmt.Sprintf("转发端口:%d 监听结束\n", mine.forwardDto.Port))
	mine.isFinished = true
}

/**
 * 判断是否有访问权限
 */
func (mine *ForwardTCPAccept) hasAccess(proxySocket net.Conn) bool {
	//if(this.forwardDto.aclState == 0){//访问权限处于关闭状态
	//    return true
	//}
	//
	////获取IP地址
	////不能使用packet.address.hostName,会出现延迟
	//val ip = proxySocket.inetAddress.hostAddress
	//if(this.forwardDto.aclState == 1){//白名单模式
	//    if(this.aclIpSet.contains(ip)){
	//        return true
	//    }
	//}
	//if(this.forwardDto.aclState == 2){//黑名单模式
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
func (mine *ForwardTCPAccept) shutdown() {
	mine.listen.Close()

	//关闭当前的桥接通信
	shutdownBridge(mine.forwardDto.Id)
	removeAccept(mine.forwardDto.Id)
}
