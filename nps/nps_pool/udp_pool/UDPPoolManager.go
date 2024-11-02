package udp_pool

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps"
	"DairoNPS/nps/nps_client/ClientSessionManagerInterface"
	"DairoNPS/util/LogUtil"
	"sync"
	"time"
)

/**
 * 客户端TCP连接池管理
 */

/**
 * 客户端ID对应的UDP连接池
 */
var poolMap = make(map[int]*[]*nps.UDPInfo)
var poolLock sync.Mutex

/**
 * 客户端连接池最后一次请求时间
 */
//private val clientLastRequestTimeMap = ConcurrentHashMap<Int, Long>()

func init() {
	recyle()
}

// 当前连接池数量
func GetPoolCount() int {
	count := 0
	poolLock.Lock()
	for _, pools := range poolMap {
		count += len(*pools)
	}
	poolLock.Unlock()
	return count
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
	poolLock.Lock()
	if poolMap[clientID] != nil {
		poolLock.Unlock()
		return
	}

	//创建连接池列表
	poolMap[clientID] = &([]*nps.UDPInfo{})
	poolLock.Unlock()
}

/**
 * 加入到连接池
 */
func Add(udpInfo *nps.UDPInfo, clientId int) {
	poolLock.Lock()
	poolList := poolMap[clientId]
	if len(*poolList) >= NPSConstant.MAX_POOL_COUNT { //已经达到最大连接数,拒绝新连接
		poolLock.Unlock()
		LogUtil.Info("UDP连接池已达到上限")

		//发送通知到客户端,该连接关闭
		closeNotify(udpInfo)
		return
	}

	*poolList = append(*poolList, udpInfo)
	poolLock.Unlock()
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
	poolLock.Lock()

	//客户端连接池
	poolList := poolMap[clientID]
	var resultUDPInfo *nps.UDPInfo

	if len(*poolList) > 0 {

		//取最后一次添加到连接池的连接
		pool := (*poolList)[len(*poolList)-1]

		//移除最后一个元素
		*poolList = (*poolList)[:len(*poolList)-1]

		resultUDPInfo = pool
		//LogUtil.Debug(fmt.Sprintf("客户端ID:%d 当前连接池已经创建时间：%d秒\n", clientID, (time.Now().UnixMilli()-pool.CreateTime)/1000))
	}
	poolLock.Unlock()
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
	poolLock.Lock()
	size := len(*poolMap[clientId])
	poolLock.Unlock()
	if size < NPSConstant.MAX_POOL_COUNT {
		Csmi.SendUDPPoolRequest(clientId, count)
	}
}

/**
 * 移除某个客户端所有的连接池
 */
func ShutdownByClient(clientID int) {
	poolLock.Lock()
	clientPool := poolMap[clientID]
	if clientPool == nil {
		poolLock.Unlock()
		return
	}
	for _, udpInfo := range *clientPool {

		//通知客户端关闭
		closeNotify(udpInfo)
	}
	*clientPool = []*nps.UDPInfo{}

	////0代表移除所有连接池
	//val activePorts = "0"
	//ClientSessionManager.send(
	//   clientID,
	//   HeaderUtil.SYNC_ACTIVE_POOL_UDP_PORT,
	//   activePorts
	//)
	poolLock.Unlock()
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
