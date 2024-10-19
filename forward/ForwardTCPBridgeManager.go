package forward

import (
	"DairoNPS/dao/dto"
	"net"
	"sync"
)

/**
 * TCP桥接会话管理
 */

/**
 * 当前正在通信的会话
 */
var bridgeList = make(map[*ForwardBridge]bool)
var bridgeListLock sync.Mutex

/**
 * 当前桥接数量
 */
//val bridgeCount: Int
//    get() = this.bridgeList.count()

/**
 * 获取当前桥接列表
 */
//suspend fun getBridgeList(): List<ForwardBridge> {
//    var result: List<ForwardBridge>? = null
//    this.bridgeListLock.synchronized {
//        result = this.bridgeList.map {
//            it
//        }
//    }
//    return result!!
//}

/**
 * 开始会话
 * @param channel 隧道信息
 * @param proxyTCP 代理服务端Socket
 * @param clientSocket 内网穿透客户端Socket
 */
func startBridge(forwardDto dto.ForwardDto, proxyTCP net.Conn, targetTCP net.Conn) {
	bridge := &ForwardBridge{
		ForwardDto: forwardDto,
		ProxyTCP:   proxyTCP,
		TargetTCP:  targetTCP,
	}
	bridgeListLock.Lock()
	bridgeList[bridge] = true
	bridgeListLock.Unlock()
	bridge.Start()
}

/**
 * 移除会话
 */
func removeBridge(bridge *ForwardBridge) {
	bridgeListLock.Lock()
	delete(bridgeList, bridge)
	bridgeListLock.Unlock()
}

/**
 * 关闭隧道所有正在通信的连接
 */
func shutdownBridge(forwardId int) {
	closeList := []*ForwardBridge{}

	bridgeListLock.Lock()
	for item := range bridgeList {
		if item.ForwardDto.Id == forwardId {
			closeList = append(closeList, item)
		}
	}
	bridgeListLock.Unlock()
	for _, item := range closeList {
		item.shutdown()
	}
}
