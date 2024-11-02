package udp_bridge

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps"
	"DairoNPS/util/ChannelStatisticsUtil"
	"net"
	"sync"
	"time"
)

// UDP桥接管理

/**
 * DUP代理(IP:PORT)对应的会话信息
 */
var proxyUDPInfoToBridge = make(map[string]*UDPBridge)
var bridgeLock sync.Mutex

/**
 * 内网穿透客户端(IP:PORT)对应DUP代理(IP:PORT)
 */
var clientUDPInfoToProxyUDPInfo = make(map[string]string)

func init() {

	//定时回收资源
	recyle()
}

// 当前桥接数量
func GetBridgeCount() int {
	count := 0
	bridgeLock.Lock()
	count = len(proxyUDPInfoToBridge)
	bridgeLock.Unlock()
	return count
}

// 获取当前桥接列表
func GetBridgeList() []UDPBridge {
	list := []UDPBridge{}
	bridgeLock.Lock()
	for _, item := range proxyUDPInfoToBridge {
		list = append(list, *item)
	}
	bridgeLock.Unlock()
	return list
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
	clientId int,
	channel *dto.ChannelDto,
	proxyUDPInfo *nps.UDPInfo,
	clientUDPInfo *nps.UDPInfo,
) *UDPBridge {
	bridge := &UDPBridge{
		ClientId:        clientId,
		Channel:         channel,
		ProxyUDPInfo:    proxyUDPInfo,
		ClientUDPInfo:   clientUDPInfo,
		CreateTime:      time.Now().UnixMilli(),
		LastRWTime:      time.Now().UnixMilli(),
		channelDataSize: ChannelStatisticsUtil.Get(channel.Id),
	}
	bridgeLock.Lock()
	proxyUDPInfoToBridge[proxyUDPInfo.Key()] = bridge
	clientUDPInfoToProxyUDPInfo[clientUDPInfo.Key()] = proxyUDPInfo.Key()
	bridgeLock.Unlock()
	err := bridge.SendHeaderToClient()
	if err != nil {
		bridgeLock.Lock()
		delete(proxyUDPInfoToBridge, proxyUDPInfo.Key())
		delete(clientUDPInfoToProxyUDPInfo, clientUDPInfo.Key())
		bridgeLock.Unlock()
		return nil
	}
	return bridge
}

/**
 * 通过代理服务端信息获取会话
 */
func ByProxy(addr *net.UDPAddr) *UDPBridge {
	key := addr.String()
	bridgeLock.Lock()
	bridge := proxyUDPInfoToBridge[key]
	bridgeLock.Unlock()
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
	bridgeLock.Lock()
	sourceUDPInfoKey, isExists := clientUDPInfoToProxyUDPInfo[key]
	if !isExists { //不存在
		bridgeLock.Unlock()
		return nil
	}
	bridge := proxyUDPInfoToBridge[sourceUDPInfoKey]
	bridgeLock.Unlock()
	return bridge
}

/**
 * 关闭某个隧道所有的连接
 */
func ShutdownByChannel(channelId int) {
	bridgeLock.Lock()

	//要删除的代理key
	proxyUDPInfoKeys := []string{}

	//要删除的客户端key
	clientUDPInfoKeys := []string{}
	for key, bridge := range proxyUDPInfoToBridge {
		if bridge.Channel.Id == channelId {
			proxyUDPInfoKeys = append(proxyUDPInfoKeys, key)
			clientUDPInfoKeys = append(clientUDPInfoKeys, bridge.ClientUDPInfo.Key())
		}
	}
	for _, key := range proxyUDPInfoKeys {
		delete(proxyUDPInfoToBridge, key)
	}
	for _, key := range clientUDPInfoKeys {
		delete(clientUDPInfoToProxyUDPInfo, key)
	}
	bridgeLock.Unlock()
}

/**
 * 关闭客户端所有的连接
 */
func ShutdownByClient(clientId int) {
	bridgeLock.Lock()

	//要删除的代理key
	proxyUDPInfoKeys := []string{}

	//要删除的客户端key
	clientUDPInfoKeys := []string{}
	for key, bridge := range proxyUDPInfoToBridge {
		if bridge.ClientId == clientId {
			proxyUDPInfoKeys = append(proxyUDPInfoKeys, key)
			clientUDPInfoKeys = append(clientUDPInfoKeys, bridge.ClientUDPInfo.Key())
		}
	}
	for _, key := range proxyUDPInfoKeys {
		delete(proxyUDPInfoToBridge, key)
	}
	for _, key := range clientUDPInfoKeys {
		delete(clientUDPInfoToProxyUDPInfo, key)
	}
	bridgeLock.Unlock()
}

// 移除桥接通信
func RemoveBridge(bridge *UDPBridge) {
	bridgeLock.Lock()
	delete(proxyUDPInfoToBridge, bridge.ProxyUDPInfo.Key())
	delete(clientUDPInfoToProxyUDPInfo, bridge.ClientUDPInfo.Key())
	bridgeLock.Unlock()
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
