package udp_bridge

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps"
	"net"
	"sync"
)

// UDP桥接管理

/**
 * DUP代理(IP:PORT)对应的会话信息
 */
var proxyUDPInfoToBridge = make(map[string]*UDPBridge)
var proxyUDPInfoToBridgeLock sync.Mutex

/**
 * 内网穿透客户端(IP:PORT)对应DUP代理(IP:PORT)
 */
var clientUDPInfoToProxyUDPInfo = make(map[string]string)

func init() {

	//定时回收资源
	recyle()
}

/**
 * 当前桥接数量
 */
//val bridgeCount: Int
//    get() = this.proxyUDPInfoToBridge.count()

///**
// * 获取当前桥接列表
// */
//func getBridgeList(): List<UDPBridge> {
//    var result: List<UDPBridge>? = null
//    this.proxyUDPInfoToBridgeLock.withLock {
//        result = this.proxyUDPInfoToBridge.map {
//            it.value
//        }
//    }
//    return result!!
//}

/**
 * 通过代理服务端信息获取会话
 */
func ByProxy(addr *net.UDPAddr) *UDPBridge {
	key := addr.String()
	proxyUDPInfoToBridgeLock.Lock()
	bridge := proxyUDPInfoToBridge[key]
	proxyUDPInfoToBridgeLock.Unlock()
	return bridge
}

/**
 * 通过内网穿透客户端信息获取会话
 * @param ip 内网穿透客户端ip
 * @param port 内网穿透客户端端口
 * @return 代理服务端与内网穿透客户端会话
 */
func ByClient(addr *net.UDPAddr) *UDPBridge {
	key := addr.String()
	proxyUDPInfoToBridgeLock.Lock()
	sourceUDPInfoKey, isExists := clientUDPInfoToProxyUDPInfo[key]
	if !isExists { //不存在
		proxyUDPInfoToBridgeLock.Unlock()
		return nil
	}
	bridge := proxyUDPInfoToBridge[sourceUDPInfoKey]
	proxyUDPInfoToBridgeLock.Unlock()
	return bridge
}

/**
 * 创建桥接会话
 * @param client 客户端DTO
 * @param channel 隧道信息
 * @param proxyUDPInfo 代理服务端UDP信息
 * @param clientUDPInfo 内网穿透客户端UDP信息
 * @return 代理服务端与内网穿透客户端会话
 */
func CreateBridge(
	client *dto.ClientDto,
	channel *dto.ChannelDto,
	proxyUDPInfo *nps.UDPInfo,
	clientUDPInfo *nps.UDPInfo,
) *UDPBridge {
	bridge := &UDPBridge{
		client:        client,
		channel:       channel,
		proxyUDPInfo:  proxyUDPInfo,
		clientUDPInfo: clientUDPInfo,
	}
	proxyUDPInfoToBridgeLock.Lock()
	proxyUDPInfoToBridge["IP:PORT"] = bridge
	clientUDPInfoToProxyUDPInfo["client IP:PORT"] = "proxy IP:PORT"
	//println("-->当前UDP数量:" + this.proxyUDPInfoToBridge.size)

	proxyUDPInfoToBridgeLock.Unlock()
	return bridge
}

/**
 * 关闭会话
 */
func closeByBridge(bridge *UDPBridge) {
	proxyUDPInfoToBridgeLock.Lock()
	delete(clientUDPInfoToProxyUDPInfo, bridge.clientUDPInfo.Key())
	delete(proxyUDPInfoToBridge, bridge.proxyUDPInfo.Key())
	proxyUDPInfoToBridgeLock.Unlock()
}

/**
 * 关闭某个隧道所有的连接
 */
func ShutdownByChannel(channelId int) {

	//this.filter{ //筛选出指定隧道id的会话
	//	it.value.channel.id == channelId,
	//}.forEach
	//{
	//	this.closeByBridge(it)
	//}
}

/**
 * 关闭客户端所有的连接
 */
func ShutdownByClient(clientId int) {
	//this.filter{ //筛选出指定隧道id的会话
	//	it.value.client.id == clientId,
	//}.forEach
	//{
	//	this.closeByBridge(it)
	//}
}

/**
 * 长时间不用的连接回收
 */
func recyle() {
	//while(true)
	//{
	//	delay(CLSConfig.RECYLE_UDP_TIME)
	//	val
	//	now = System.currentTimeMillis()
	//	try{
	//
	//		//要关闭的连接
	//		val, closeList = ArrayList<UDPBridge>()
	//
	//		//当前存活的连接
	//		val activeList = ArrayList<UDPBridge>()
	//		this.proxyUDPInfoToBridgeLock.withLock{
	//		this.proxyUDPInfoToBridge.forEach{
	//		if ((now - it.value.lastSessionTime) > CLSConfig.RECYLE_UDP_TIME){ //筛选出指定隧道id的会话
	//		closeList.add(it.value)
	//	} else{
	//		activeList.add(it.value)
	//	}
	//	}
	//	}
	//		closeList.forEach{ //关闭
	//		this.closeByBridge(it)
	//	}
	//
	//		//按照客户端ID分组
	//		val activeClientToBridgeList = activeList.groupBy{
	//		it.client.id
	//	}
	//
	//		//遍历当前连接中的客户端
	//		ClientSessionManager.getClientList().forEach{
	//		val clientID = it.id!!
	//		val bridgeList = activeClientToBridgeList[clientID]
	//
	//		//整理当前在线的端口
	//		var activePorts = bridgeList?.joinToString(","){bridge->
	//		bridge.clientInfo.port.toString()
	//	}
	//		if (activePorts.isNullOrEmpty()){ //连接全部关闭
	//		activePorts = "0"
	//	}
	//		ClientSessionManager.send(
	//		clientID,
	//		HeaderUtil.SYNC_ACTIVE_BRIDGE_UDP_PORT,
	//		activePorts
	//	)
	//	}
	//	}
	//	catch(e: Exception) {
	//	e.printStackTrace()
	//	println("-->连接回收报错")
	//}
	//}
}

/**
 * 帅选出指定数据
 * 由于HashMap不允边遍历边修改数据,所以遍历时需要加锁
 */
//func filter(predicate: (Map.Entry<String, UDPBridge>) -> Boolean): List<UDPBridge> {
//val result = ArrayList<UDPBridge>()
//this.proxyUDPInfoToBridgeLock.withLock {
//this.proxyUDPInfoToBridge.forEach {
//if (predicate(it)) {
//result.add(it.value)
//}
//}
//}
//return result
//}
//}
