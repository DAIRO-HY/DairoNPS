package TCPPoolManager

import (
	"DairoNPS/constant/CLSConfig"
	"DairoNPS/pool/TCPPool"
	"bufio"
	"net"
	"sync"
	"time"
)

//客户端TCP连接池管理

/**
 * 客户端连接池
 */
//private val clientScoketPoolMap = HashMap<Int, ClientPoolList<TCPPool>>()
var clientScoketPoolMap = make(map[int][]TCPPool.TCPPool)
var (
	clientScoketPoolMapLock sync.Mutex
)

/**
 * 客户端连接池最后一次请求时间
 */
//private val clientLastRequestTimeMap = ConcurrentHashMap<Int, Long>()
var clientLastRequestTimeMap = make(map[int]int64)

//init {
//    GlobalScope.launch {
//        //recyle()
//    }
//}

/**
 * 当前连接池数量
 */
func getCount() int {
	var total = 0
	//this.clientScoketPoolMapLock.withLock {
	//    this.clientScoketPoolMap.forEach { (_, v) ->
	//        total += v.count()
	//    }
	//}
	return total
}

/**
 * 为客户端创建一个空的连接池
 * @param clientID 客户端ID
 */
func initEmptyPoolByClient(clientID int) {

	// 使用值和 ok 变量判断 key 是否存在
	if _, ok := clientScoketPoolMap[clientID]; ok {
		return
	}
	clientScoketPoolMapLock.Lock()
	if _, ok := clientScoketPoolMap[clientID]; ok {

	} else { //若连接池不存在则创建
		clientScoketPoolMap[clientID] = []TCPPool.TCPPool{}
	}
	clientScoketPoolMapLock.Unlock()
}

/**
 * 添加TCP连接池
 * @param clientSocket tcp连接
 */
func Add(clientSocket net.Conn) {

	////从头部信息中得到客户端id
	//clientId := HeaderUtil.getHeader(clientSocket)?.toInt()
	//if clientId == nil {
	//
	//   //无效的连接
	//   clientSocket.Close()
	//   return
	//}
	//initEmptyPoolByClient(clientId)
	//poolList := clientScoketPoolMap[clientId]
	//
	////TODO: 这里应该每个客户端创建一把锁
	//clientScoketPoolMapLock.Lock()
	//   if (len(poolList) >= CLSConfig.MAX_POOL_COUNT) {//已经达到最大连接数,拒绝新连接
	//       clientSocket.Close()
	//       clientScoketPoolMapLock.Unlock()
	//       return
	//   }
	//   pool := TCPPool.TCPPool{
	//       ClientID : clientId,
	//       Socket:clientSocket,
	//   }
	//poolList = append(poolList, pool)
	//clientScoketPoolMapLock.Unlock()
}

/**
 * 通过客户端ID获取一个连接
 * @param clientID 客户端ID
 */
func get(clientID int) net.Conn {
	initEmptyPoolByClient(clientID)

	//客户端连接池
	poolList := clientScoketPoolMap[clientID]

	var resultPool TCPPool.TCPPool

	//TODO: 这里应该每个客户端创建一把锁
	clientScoketPoolMapLock.Lock()
	if len(poolList) == 0 {
		clientScoketPoolMapLock.Unlock()
		return nil
	}
	for len(poolList) > 0 {

		//取最后一次添加到连接池的连接
		pool := poolList[len(poolList)-1]

		//移除最后一个元素
		poolList = poolList[:len(poolList)-1]
		//println("-->从连接池获取到一个连接,连接池剩余:${poolList.size}")
		//发送紧急数据来判断客户端是否在线,客户端的oobInline设置为false会忽律紧急数据
		writer := bufio.NewWriter(pool.Socket)
		if writer.WriteByte(0) != nil {
			TCPPool.Close(pool)
			continue
		}
		if writer.Flush() != nil {
			TCPPool.Close(pool)
			continue
		}
		resultPool = pool
		break
	}
	clientScoketPoolMapLock.Unlock()
	return resultPool.Socket
}

/**
 * 从连接池获取一个连接,并请求添加连接池
 * @param clientID 客户端ID
 */
func GetAndAddPool(clientID int) net.Conn {

	//记录该客户端最后一次请求连接池时间
	clientLastRequestTimeMap[clientID] = time.Now().UnixNano() / int64(time.Millisecond)
	var socket net.Conn
	for i := 0; i < 5; i++ {
		socket = get(clientID)
		if socket != nil {
			break
		}

		//申请创建连接池
		poolRequest(clientID, CLSConfig.ADD_POOL_COUNT)
		time.Sleep(1 * time.Second)
	}

	//每消耗一个连接,申请创建两个连接,直到达到最大连接池数量
	go poolRequest(clientID, 2)
	return socket
}

/**
 * 发起连接池申请请求
 * 每取走一个连接,则请求创建2个新的连接,直到达到最大连接数
 * @param clientID 客户端ID
 */
func poolRequest(clientID int, count int) {
	//if len(clientScoketPoolMap[clientID]) < CLSConfig.MAX_POOL_COUNT {
	//    ClientSessionManager.sendTCPPoolRequest(clientID, count)
	//}
}

/**
 * 移除某个客户端所有的连接池
 * @param clientID 客户端ID
 */
func CloseByClient(clientID int) {
	poolList := clientScoketPoolMap[clientID]
	//TODO: 这里应该每个客户端创建一把锁
	clientScoketPoolMapLock.Lock()
	for _, it := range poolList {
		TCPPool.Close(it)
	}
	poolList = []TCPPool.TCPPool{}
	clientScoketPoolMapLock.Unlock()
}

/**
 * 从连接池列表中移除
 */
func timeOutRemove(clientID int, pool TCPPool.TCPPool) {
	//val poolList = this.clientScoketPoolMap[clientID] as ClientPoolList
	//poolList.synchronized {
	//    poolList.remove(pool)
	//}
	//if(poolList.isEmpty()){//连接池被掏空,则申请一个连接池备用
	//    this.poolRequest(clientID, 1)
	//}
}
