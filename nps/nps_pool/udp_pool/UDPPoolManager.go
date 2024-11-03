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
var poolMap = make(map[int]*[]*UDPPool)
var poolLock sync.Mutex

func init() {
	go timeoutCheck()
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
	poolMap[clientID] = &([]*UDPPool{})
	poolLock.Unlock()
}

/**
 * 加入到连接池
 */
func Add(udpInfo *nps.UDPInfo, clientId int) {
	pool := &UDPPool{
		UDPInfo:    udpInfo,
		CreateTime: time.Now().UnixMilli(),
	}
	poolLock.Lock()
	poolList := poolMap[clientId]
	if len(*poolList) >= NPSConstant.MAX_POOL_COUNT { //已经达到最大连接数,拒绝新连接
		poolLock.Unlock()
		LogUtil.Info("UDP连接池已达到上限")

		//发送通知到客户端,该连接关闭
		pool.CloseNotify()
		return
	}

	*poolList = append(*poolList, pool)
	poolLock.Unlock()
}

/**
 * 通过客户端ID获取一个连接
 */
func get(clientID int) *UDPPool {
	poolLock.Lock()

	//客户端连接池
	poolList := poolMap[clientID]
	var resultUDPInfo *UDPPool

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
func GetAndAddPool(clientID int) *UDPPool {

	var pool *UDPPool
	for i := 0; i < 5; i++ {
		pool = get(clientID)
		if pool != nil {
			break
		}

		//申请创建连接池
		poolRequest(clientID, NPSConstant.ADD_POOL_COUNT)
		time.Sleep(1 * time.Second)
	}

	//每消耗一个连接,申请创建两个连接,直到达到最大连接池数量
	go poolRequest(clientID, 2)
	return pool
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
	for _, pool := range *clientPool {

		//通知客户端关闭
		pool.CloseNotify()
	}
	*clientPool = []*UDPPool{}
	poolLock.Unlock()
}

// 超时连接池整理
func timeoutCheck() {
	for {
		time.Sleep(NPSConstant.RECYLE_POOL_TIME_OUT * time.Millisecond)

		//当前时间戳秒
		now := time.Now().UnixMilli()

		//本次要关闭的连接池
		var closeList []*UDPPool
		poolLock.Lock()
		for clientId, pools := range poolMap { //遍历所有客户端的连接池
			poolList := *pools
			poolSize := len(poolList)
			for i := poolSize - 1; i > -1; i-- {
				pool := (*pools)[i]
				if now-pool.CreateTime > NPSConstant.RECYLE_POOL_TIME_OUT { //连接池超过指定时间
					closeList = append(closeList, pool)
					poolList = poolList[0:i]
				}
			}
			if len(poolList) == 0 { //如果连接池被清空，则请求创建一个新的连接池
				Csmi.SendUDPPoolRequest(clientId, 1)
			}
			*pools = poolList
		}
		poolLock.Unlock()
		for _, pool := range closeList { //发送关闭通知
			pool.CloseNotify()
		}
	}
}
