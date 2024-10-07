package TCPBridgeManager

import (
	"DairoNPS/bridge/TCPBridge"
	"DairoNPS/dao/dto"
	"net"
	"sync"
)

//TCP桥接会话管理

/**
 * 当前正在通信的会话
 */
var bridgeList = []TCPBridge.TCPBridge{}
var (
	bridgeListLock sync.Mutex
)

///**
// * 当前桥接数量
// */
//val bridgeCount: Int
//    get() = this.bridgeList.count()

//init {
//    GlobalScope.launch {
//        this@TCPBridgeManager.recycle()
//    }
//
//}

/**
 * 获取当前桥接列表
 */
func GetBridgeList() []TCPBridge.TCPBridge {
	return bridgeList
}

/**
 * 开始会话
 * @param client 客户端DTO
 * @param channel 隧道信息
 * @param proxySocket 代理服务端Socket
 * @param clientSocket 内网穿透客户端Socket
 */
func Start(client *dto.ClientDto, channel *dto.ChannelDto, proxySocket net.Conn, clientSocket net.Conn) {
	bridge := TCPBridge.TCPBridge{
		Channel:      channel,
		Client:       client,
		ProxySocket:  proxySocket,
		ClientSocket: clientSocket,
	}
	//this.bridgeListLock.synchronized {
	//    this.bridgeList.add(bridge)
	//}
	bridgeListLock.Lock()
	bridgeList = append(bridgeList, bridge)
	bridgeListLock.Unlock()
	TCPBridge.Start(bridge)
}

/**
 * 移除会话
 */
func RemoveBridgeList(bridge TCPBridge.TCPBridge) {
	bridgeListLock.Lock()
	for index, it := range bridgeList {
		if it == bridge {
			bridgeList = append(bridgeList[:index], bridgeList[index+1:]...)
			break
		}
	}
	bridgeListLock.Unlock()
}

/**
 * 关闭客户端所有正在通信的连接
 */
func CloseByClient(clientId int) {
	var closeList = []TCPBridge.TCPBridge{}
	bridgeListLock.Lock()

	//帅选出要删除的客户端桥接
	for _, it := range bridgeList {
		if it.Client.Id == clientId {
			closeList = append(closeList, it)
		}
	}

	//移除即将关闭的桥接
	for _, closeIt := range closeList {
		for index, it := range bridgeList {
			if it == closeIt {
				bridgeList = append(bridgeList[:index], bridgeList[index+1:]...)
				break
			}
		}
	}
	bridgeListLock.Unlock()

	//关闭的桥接
	for _, closeIt := range closeList {
		TCPBridge.Close(closeIt)
	}
}

/**
 * 关闭隧道所有正在通信的连接
 */
func CloseByChannel(channelId int) {

	var closeList = []TCPBridge.TCPBridge{}
	bridgeListLock.Lock()

	//帅选出要删除的客户端桥接
	for _, it := range bridgeList {
		if it.Channel.Id == channelId {
			closeList = append(closeList, it)
		}
	}

	//移除即将关闭的桥接
	for _, closeIt := range closeList {
		for index, it := range bridgeList {
			if it == closeIt {
				bridgeList = append(bridgeList[:index], bridgeList[index+1:]...)
				break
			}
		}
	}
	bridgeListLock.Unlock()

	//关闭的桥接
	for _, closeIt := range closeList {
		TCPBridge.Close(closeIt)
	}
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
	//        this.bridgeListLock.synchronized {
	//            result = this.bridgeList.filter {
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
