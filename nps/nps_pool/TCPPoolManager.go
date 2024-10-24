package nps_pool

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps/nps_client/ClientSessionManagerInterface"
	"DairoNPS/nps/nps_client/HeaderUtil"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

/**
 * 客户端连接池
 */
var poolMap = make(map[int]*[]*TCPPool)
var poolLock sync.Mutex

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
	poolMap[clientID] = &([]*TCPPool{})
	poolLock.Unlock()
}

/**
 * 添加TCP连接池
 * @param clientSocket tcp连接
 */
func Add(clientSocket net.Conn) {

	//从头部信息中得到客户端id
	clientIdStr, err := HeaderUtil.GetHeader(clientSocket)
	if err != nil { //无效的连接
		clientSocket.Close()
		return
	}
	clientIdInt64, err := strconv.ParseInt(clientIdStr, 10, 32)
	if err != nil {
		clientSocket.Close()
		return
	}
	clientId := int(clientIdInt64)
	poolLock.Lock()
	poolList := poolMap[clientId]
	poolLock.Unlock()
	if len(*poolList) >= NPSConstant.MAX_POOL_COUNT { //已经达到最大连接数,拒绝新连接
		clientSocket.Close()
		return
	}

	//设置长连接
	clientSocket.SetReadDeadline(time.Time{})
	poolLock.Lock()
	pool := &TCPPool{
		CreateTime: time.Now().UnixMilli(),
		ClientID:   clientId,
		Socket:     clientSocket,
	}
	*poolList = append(*poolList, pool)
	log.Printf("客户端ID：%d 当前连接数为：%d", clientId, len(*poolList))
	poolLock.Unlock()
}

/**
 * 通过客户端ID获取一个连接
 * @param clientID 客户端ID
 */
func get(clientID int) net.Conn {
	poolLock.Lock()

	//客户端连接池
	poolList := poolMap[clientID]
	if poolList == nil || len(*poolList) == 0 {
		poolLock.Unlock()
		return nil
	}
	var resultTcp net.Conn = nil
	for len(*poolList) > 0 {

		//取最后一次添加到连接池的连接
		pool := (*poolList)[len(*poolList)-1]

		//移除最后一个元素
		*poolList = (*poolList)[:len(*poolList)-1]

		//试探性发送一个数据，检测连接是否已经失效
		//TODO:客户端还未支持
		//_, err := pool.Socket.Write([]uint8{0})
		//if err != nil {
		//	pool.Shutdown()
		//	continue
		//}
		resultTcp = pool.Socket
		break
	}
	poolLock.Unlock()
	return resultTcp
}

/**
 * 从连接池获取一个连接,并请求添加连接池
 * @param clientID 客户端ID
 */
func GetAndAddPool(clientID int) net.Conn {
	var tcp net.Conn
	for i := 0; i < 5; i++ {
		tcp = get(clientID)
		if tcp != nil {
			break
		}

		//申请创建连接池
		poolRequest(clientID, NPSConstant.ADD_POOL_COUNT)
		time.Sleep(1 * time.Second)
	}

	//每消耗一个连接,申请创建两个连接,直到达到最大连接池数量
	go poolRequest(clientID, 2)
	return tcp
}

/**
 * 发起连接池申请请求
 * 每取走一个连接,则请求创建2个新的连接,直到达到最大连接数
 * @param clientID 客户端ID
 */
func poolRequest(clientId int, count int) {
	poolLock.Lock()
	poolList := poolMap[clientId]
	poolLock.Unlock()
	if poolList == nil {
		return
	}
	if len(*poolList) < NPSConstant.MAX_POOL_COUNT {
		Csmi.SendTCPPoolRequest(clientId, count)
	}
}

var Csmi ClientSessionManagerInterface.ClientSessionManagerInterface

/**
 * 移除某个客户端所有的连接池
 * @param clientID 客户端ID
 */
func ShutdownByClient(clientID int) {
	poolLock.Lock()
	poolList := poolMap[clientID]
	poolLock.Unlock()
	for _, it := range *poolList {
		it.Shutdown()
	}
	poolLock.Lock()
	*poolList = []*TCPPool{}
	poolLock.Unlock()
}

/**
 * 从连接池列表中移除
 */
func timeOutRemove(clientID int, pool TCPPool) {
	//val poolList = this.poolMap[clientID] as ClientPoolList
	//poolList.synchronized {
	//    poolList.remove(pool)
	//}
	//if(poolList.isEmpty()){//连接池被掏空,则申请一个连接池备用
	//    this.poolRequest(clientID, 1)
	//}
}
