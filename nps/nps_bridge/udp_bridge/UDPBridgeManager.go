package udp_bridge

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
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
DebugTimer.Add207()

	//定时回收资源
	go timeoutCheck()
}

// 当前桥接数量
func GetBridgeCount() int {
DebugTimer.Add208()
	count := 0
	bridgeLock.Lock()
	count = len(proxyUDPInfoToBridge)
	bridgeLock.Unlock()
	return count
}

// 获取当前桥接列表
func GetBridgeList() []UDPBridge {
DebugTimer.Add209()
	list := []UDPBridge{}
	bridgeLock.Lock()
	for _, item := range proxyUDPInfoToBridge {
DebugTimer.Add210()
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
DebugTimer.Add211()
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
DebugTimer.Add212()
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
DebugTimer.Add213()
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
DebugTimer.Add214()
	key := addr.String()
	bridgeLock.Lock()
	sourceUDPInfoKey, isExists := clientUDPInfoToProxyUDPInfo[key]
	if !isExists { //不存在
DebugTimer.Add215()
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
DebugTimer.Add216()
	bridgeLock.Lock()

	//要删除的代理key
	proxyUDPInfoKeys := []string{}

	//要删除的客户端key
	clientUDPInfoKeys := []string{}
	for key, bridge := range proxyUDPInfoToBridge {
DebugTimer.Add217()
		if bridge.Channel.Id == channelId {
DebugTimer.Add218()
			proxyUDPInfoKeys = append(proxyUDPInfoKeys, key)
			clientUDPInfoKeys = append(clientUDPInfoKeys, bridge.ClientUDPInfo.Key())
		}
	}
	for _, key := range proxyUDPInfoKeys {
DebugTimer.Add219()
		delete(proxyUDPInfoToBridge, key)
	}
	for _, key := range clientUDPInfoKeys {
DebugTimer.Add220()
		delete(clientUDPInfoToProxyUDPInfo, key)
	}
	bridgeLock.Unlock()
}

/**
 * 关闭客户端所有的连接
 */
func ShutdownByClient(clientId int) {
DebugTimer.Add221()
	bridgeLock.Lock()

	//要删除的代理key
	proxyUDPInfoKeys := []string{}

	//要删除的客户端key
	clientUDPInfoKeys := []string{}
	for key, bridge := range proxyUDPInfoToBridge {
DebugTimer.Add222()
		if bridge.ClientId == clientId {
DebugTimer.Add223()
			proxyUDPInfoKeys = append(proxyUDPInfoKeys, key)
			clientUDPInfoKeys = append(clientUDPInfoKeys, bridge.ClientUDPInfo.Key())
		}
	}
	for _, key := range proxyUDPInfoKeys {
DebugTimer.Add224()
		delete(proxyUDPInfoToBridge, key)
	}
	for _, key := range clientUDPInfoKeys {
DebugTimer.Add225()
		delete(clientUDPInfoToProxyUDPInfo, key)
	}
	bridgeLock.Unlock()
}

// 移除桥接通信
func RemoveBridge(bridge *UDPBridge) {
DebugTimer.Add226()
	bridgeLock.Lock()
	delete(proxyUDPInfoToBridge, bridge.ProxyUDPInfo.Key())
	delete(clientUDPInfoToProxyUDPInfo, bridge.ClientUDPInfo.Key())
	bridgeLock.Unlock()
}

/**
 * 长时间不用的连接回收
 */
func timeoutCheck() {
DebugTimer.Add227()
	for {
DebugTimer.Add228()
		time.Sleep(NPSConstant.UDP_BRIDGE_TIMEOUT*time.Millisecond + 1000)

		//当前时间戳秒
		now := time.Now().UnixMilli()

		//要关闭的连接
		var closeList []*UDPBridge
		bridgeLock.Lock()
		for _, bridge := range proxyUDPInfoToBridge {
DebugTimer.Add229()
			if now-bridge.LastRWTime > NPSConstant.UDP_BRIDGE_TIMEOUT {
DebugTimer.Add230()

				//本次需要关闭的桥接
				closeList = append(closeList, bridge)
			}
		}
		for _, bridge := range closeList { //移除桥接
DebugTimer.Add231()
			delete(proxyUDPInfoToBridge, bridge.ProxyUDPInfo.Key())
			delete(clientUDPInfoToProxyUDPInfo, bridge.ClientUDPInfo.Key())
		}
		bridgeLock.Unlock()
		for _, bridge := range closeList { //发送关闭标识
DebugTimer.Add232()
			closeData := []byte(NPSConstant.UDP_BRIDIGE_CLOSE_FLAG)
			bridge.ClientUDPInfo.Send(closeData, len(closeData))
		}
	}
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
