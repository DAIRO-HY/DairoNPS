package udp_pool

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps"
	"DairoNPS/nps/nps_client/ClientSessionManagerInterface"
	"sync"
	"time"
)

/**
 * 客户端TCP连接池管理
 */

/**
 * 客户端ID对应的UDP连接池
 */
var clientScoketPoolMap = make(map[int]*[]*nps.UDPInfo)
var lock sync.Mutex

/**
 * 客户端连接池最后一次请求时间
 */
//private val clientLastRequestTimeMap = ConcurrentHashMap<Int, Long>()

func init() {
	recyle()
}

/**
 * 当前连接池数量
 */
//func getCount() int {
//    var total = 0
//    this.clientScoketPoolMapLock.withLock {
//        this.clientScoketPoolMap.forEach { (_, v) ->
//            total += v.count()
//        }
//    }
//    return total
//}

/**
 * 为客户端创建一个空的连接池
 * @param clientID 客户端ID
 */
func InitEmptyPoolByClient(clientID int) {
	lock.Lock()
	if clientScoketPoolMap[clientID] != nil {
		lock.Unlock()
		return
	}

	//创建连接池列表
	clientScoketPoolMap[clientID] = &([]*nps.UDPInfo{})
	lock.Unlock()
}

/**
 * 加入到连接池
 */
func Add(udpInfo *nps.UDPInfo, clientId int) {
	lock.Lock()
	poolList := clientScoketPoolMap[clientId]

	if len(*poolList) >= NPSConstant.MAX_POOL_COUNT { //已经达到最大连接数,拒绝新连接
		lock.Unlock()

		//发送通知到客户端,该连接关闭
		closeNotify(udpInfo)
		return
	}

	*poolList = append(*poolList, udpInfo)
	lock.Unlock()
}

/**
 * 通知客户端关闭该连接池
 */
func closeNotify(clientUdp *nps.UDPInfo) {
	closeData := []byte(NPSConstant.CLOSE_UDP_POOL_FLAG)
	clientUdp.Send(closeData, len(closeData))
}

/**
 * 通过客户端ID获取一个连接
 */
func get(clientID int) *nps.UDPInfo {
	lock.Lock()

	//客户端连接池
	poolList := clientScoketPoolMap[clientID]
	var resultUDPInfo *nps.UDPInfo

	if len(*poolList) > 0 {

		//取最后一次添加到连接池的连接
		pool := (*poolList)[len(*poolList)-1]

		//移除最后一个元素
		*poolList = (*poolList)[:len(*poolList)-1]

		resultUDPInfo = pool
		//LogUtil.Debug(fmt.Sprintf("客户端ID:%d 当前连接池已经创建时间：%d秒\n", clientID, (time.Now().UnixMilli()-pool.CreateTime)/1000))
	}
	lock.Unlock()
	return resultUDPInfo
}

/**
 * 从连接池获取一个连接,并请求添加连接池
 */
func GetAndAddPool(clientID int) *nps.UDPInfo {

	var udpInfo *nps.UDPInfo
	for i := 0; i < 5; i++ {
		udpInfo = get(clientID)
		if udpInfo != nil {
			break
		}

		//申请创建连接池
		poolRequest(clientID, NPSConstant.ADD_POOL_COUNT)
		time.Sleep(1 * time.Second)
	}

	//每消耗一个连接,申请创建两个连接,直到达到最大连接池数量
	go poolRequest(clientID, 2)
	return udpInfo
}

var Csmi ClientSessionManagerInterface.ClientSessionManagerInterface

/**
 * 每取走一个连接,则请求创建2个新的连接,直到达到最大连接数
 */
func poolRequest(clientId int, count int) {
	lock.Lock()
	size := len(*clientScoketPoolMap[clientId])
	lock.Unlock()
	if size < NPSConstant.MAX_POOL_COUNT {
		Csmi.SendUDPPoolRequest(clientId, count)
	}
}

/**
 * 移除某个客户端所有的连接池
 */
func ShutdownByClient(clientID int) {
	lock.Lock()
	clientPool := clientScoketPoolMap[clientID]
	if clientPool == nil {
		lock.Unlock()
		return
	}
	for _, udpInfo := range *clientPool {
		udpInfo.Socket.Close()
	}
	clientPool = &[]*nps.UDPInfo{}

	////0代表移除所有连接池
	//val activePorts = "0"
	//ClientSessionManager.send(
	//   clientID,
	//   HeaderUtil.SYNC_ACTIVE_POOL_UDP_PORT,
	//   activePorts
	//)
	lock.Unlock()
}

/**
 * 回收长时间不用的连接
 */
func recyle() {
	//while (true) {
	//    delay(CLSConfig.RECYLE_POOL_TIME)
	//    try {
	//        val currentTime = System.currentTimeMillis()
	//        this.clientLastRequestTimeMap.forEach { (clientID, lastRequestTime) ->
	//            if ((currentTime - lastRequestTime) > CLSConfig.RECYLE_POOL_TIME) {//超过指定时间没有通信的连接
	//                val poolList = this.clientScoketPoolMap[clientID] as ClientPoolList
	//                poolList.synchronized {
	//                    while (poolList.size > CLSConfig.MIN_POOL_COUNT) {//移除多余的连接池
	//                        poolList.removeAt(0)
	//                    }
	//                }
	//
	//                //当前激活的端口
	//                var activePorts = poolList.joinToString(",") {
	//                    it.port.toString()
	//                }
	//                if (activePorts.isEmpty()) {//服务商连接池为0
	//                    activePorts = "0"
	//                }
	//                ClientSessionManager.send(
	//                    clientID,
	//                    HeaderUtil.SYNC_ACTIVE_POOL_UDP_PORT,
	//                    activePorts
	//                )
	//            }
	//        }
	//    } catch (e: Exception) {
	//        e.printStackTrace()
	//    }
	//}
}
