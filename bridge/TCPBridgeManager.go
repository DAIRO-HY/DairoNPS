package bridge

import (
	"DairoNPS/dao/dto"
	"net"
	"sync"
)

//TCP桥接会话管理

// 当前正在通信的桥接
var BridgeList1 []*TCPBridge
var BridgeListMap = make(map[*TCPBridge]bool)
var BridgeListLock sync.Mutex

///**
// * 当前桥接数量
// */
//val bridgeCount: Int
//    get() = this.BridgeList.count()

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
func MakeBridge(client *dto.ClientDto, channel *dto.ChannelDto, proxySocket net.Conn, clientSocket net.Conn) {
	bridge := &TCPBridge{
		Channel:      channel,
		Client:       client,
		ProxySocket:  proxySocket,
		ClientSocket: clientSocket,
	}
	//this.BridgeListLock.synchronized {
	//    this.BridgeList.add(bridge)
	//}
	BridgeListLock.Lock()
	BridgeListMap[bridge] = true
	BridgeListLock.Unlock()
	bridge.start()
}

/**
 * 移除会话
 */
func removeBridgeList(bridge *TCPBridge) {
	BridgeListLock.Lock()
	delete(BridgeListMap, bridge)
	BridgeListLock.Unlock()
}

/**
 * 关闭客户端所有正在通信的连接
 */
func CloseByClient(clientId int) {
	BridgeListLock.Lock()

	//帅选出要删除的客户端桥接
	for bridge := range BridgeListMap {
		if bridge.Client.Id == clientId {
			bridge.shutdown()
		}
	}
	BridgeListLock.Unlock()
}

/**
 * 关闭隧道所有正在通信的连接
 */
func CloseByChannel(channelId int) {
	BridgeListLock.Lock()

	//帅选出要删除的客户端桥接
	for bridge := range BridgeListMap {
		if bridge.Channel.Id == channelId {
			bridge.shutdown()
		}
	}
	BridgeListLock.Unlock()
}

/**
 * 回收长时间不用的连接
 */
func Recycle() {
	//while (true) {
	//    delay(CLSConfig.BRIDGE_SESSION_TIMEOUT)
	//    try {
	//
	//        //当前是同时间戳
	//        val now = System.currentTimeMillis()
	//        var result: List<TCPBridge>? = null
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
