package tcp_bridge

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"net"
	"sync"
	"time"
)

//TCP桥接会话管理

// 当前正在通信的桥接
var bridgeMap = make(map[*TCPBridge]bool)
var bridgeLock sync.Mutex

// 当前桥接数量
func GetBridgeCount() int {
DebugTimer.Add187()
	count := 0
	bridgeLock.Lock()
	count = len(bridgeMap)
	bridgeLock.Unlock()
	return count
}

// 获取当前桥接列表
func GetBridgeList() []TCPBridge {
DebugTimer.Add188()
	list := []TCPBridge{}
	bridgeLock.Lock()
	for item := range bridgeMap {
DebugTimer.Add189()
		list = append(list, *item)
	}
	bridgeLock.Unlock()
	return list
}

//init {
//    GlobalScope.launch {
//        this@TCPBridgeManager.recycle()
//    }
//
//}

/**
 * 获取当前桥接列表
 */
//func GetBridgeList() []*TCPBridge {
//	return BridgeList
//}

/**
 * 开始会话
 * @param client 客户端DTO
 * @param channel 隧道信息
 * @param proxySocket 代理服务端Socket
 * @param clientSocket 内网穿透客户端Socket
 */
func MakeBridge(ClientId int, channel *dto.ChannelDto, proxySocket net.Conn, clientSocket net.Conn) {
DebugTimer.Add190()
	bridge := &TCPBridge{
		Channel:    channel,
		ClientId:   ClientId,
		ProxyTCP:   proxySocket,
		ClientTCP:  clientSocket,
		CreateTime: time.Now().UnixMilli(),
		LastRWTime: time.Now().UnixMilli(),
	}
	//this.BridgeListLock.synchronized {
	//    this.BridgeList.add(bridge)
	//}
	bridgeLock.Lock()
	bridgeMap[bridge] = true
	bridgeLock.Unlock()
	go bridge.start()
}

// 关闭客户端所有正在通信的连接
func ShutdownByClient(clientId int) {
DebugTimer.Add191()
	bridgeLock.Lock()

	//帅选出要删除的客户端桥接
	for bridge := range bridgeMap {
DebugTimer.Add192()
		if bridge.ClientId == clientId {
DebugTimer.Add193()
			bridge.shutdown()
		}
	}
	bridgeLock.Unlock()
}

// 关闭隧道所有正在通信的连接
func ShutdownByChannel(channelId int) {
DebugTimer.Add194()
	bridgeLock.Lock()

	//帅选出要删除的客户端桥接
	for bridge := range bridgeMap {
DebugTimer.Add195()
		if bridge.Channel.Id == channelId {
DebugTimer.Add196()
			bridge.shutdown()
		}
	}
	bridgeLock.Unlock()
}

// 移除桥接通信
func removeBridge(bridge *TCPBridge) {
DebugTimer.Add197()
	bridgeLock.Lock()
	delete(bridgeMap, bridge)
	bridgeLock.Unlock()
}

/**
 * 回收长时间不用的连接
 */
func Recycle() {
DebugTimer.Add198()
	//while (true) {
	//    delay(CLSConfig.BRIDGE_SESSION_TIMEOUT)
	//    try {
	//
	//        //当前是同时间戳
	//        val now = System.currentTimeMillis()
	//    result: List<TCPBridge>? = null
	//        this.BridgeListLock.synchronized {
	//            result = this.BridgeList.filter {
	//                (now - it.lastSessionTime) > CLSConfig.BRIDGE_SESSION_TIMEOUT
	//            }
	//        }
	//        result?.forEach { //关掉长时间不通信的连接
	//            it.close()
	//        }
	//    } catch (e: Exception) {
	//        //e.printStackTrace()
	//    }
	//}
}
