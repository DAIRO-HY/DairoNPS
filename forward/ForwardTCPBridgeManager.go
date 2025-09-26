package forward

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"net"
	"sync"
	"time"
)

/**
 * TCP桥接会话管理
 */

/**
 * 当前正在通信的会话
 */
var bridgeList = make(map[*ForwardBridge]bool)
var bridgeLock sync.Mutex

// 端口转发当前桥接数量
func GetBridgeCount() int {
DebugTimer.Add150()
	count := 0
	bridgeLock.Lock()
	count = len(bridgeList)
	bridgeLock.Unlock()
	return count
}

// 获取当前桥接列表
func GetBridgeList() []ForwardBridge {
DebugTimer.Add151()
	var list []ForwardBridge
	bridgeLock.Lock()
	for item := range bridgeList {
DebugTimer.Add152()
		list = append(list, *item)
	}
	bridgeLock.Unlock()
	return list
}

/**
 * 获取当前桥接列表
 */
//suspend fun getBridgeList(): List<ForwardBridge> {
//    var result: List<ForwardBridge>? = null
//    this.bridgeLock.synchronized {
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
func startBridge(forwardDto *dto.ForwardDto, proxyTCP net.Conn, targetTCP net.Conn) {
DebugTimer.Add153()
	bridge := &ForwardBridge{
		ForwardDto: forwardDto,
		ProxyTCP:   proxyTCP,
		TargetTCP:  targetTCP,
		CreateTime: time.Now().UnixMilli(),
	}
	bridgeLock.Lock()
	bridgeList[bridge] = true
	bridgeLock.Unlock()
	bridge.Start()
}

/**
 * 移除会话
 */
func removeBridge(bridge *ForwardBridge) {
DebugTimer.Add154()
	bridgeLock.Lock()
	delete(bridgeList, bridge)
	bridgeLock.Unlock()
}

/**
 * 关闭隧道所有正在通信的连接
 */
func shutdownBridge(forwardId int) {
DebugTimer.Add155()
	closeList := []*ForwardBridge{}

	bridgeLock.Lock()
	for item := range bridgeList {
DebugTimer.Add156()
		if item.ForwardDto.Id == forwardId {
DebugTimer.Add157()
			closeList = append(closeList, item)
		}
	}
	bridgeLock.Unlock()
	for _, item := range closeList {
DebugTimer.Add158()
		item.shutdown()
	}
}
