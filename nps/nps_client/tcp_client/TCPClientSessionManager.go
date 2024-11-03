package tcp_client

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_client/HeaderUtil"
	"DairoNPS/nps/nps_pool/tcp_pool"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/nps/nps_proxy/tcp_proxy"
	"DairoNPS/nps/nps_proxy/udp_proxy"
	"net"
	"strconv"
	"sync"
)

type ClientSessionManager struct{}

//往客户端发送指令的专用连接

/**
 * 客户端ID对应的Socket连接
 */
var clientSessionMap = make(map[int]*ClientSession)

/**
 * 添加互斥锁
 */
var clientSessionLock sync.Mutex

/**
 * 获取与客户端的会话
 */
//func getSession(clientID: Int) = this.clientSessionMap[clientID]

/**
 * 当前客户端数量
 */
//val size: Int
//    get() = this.clientSessionMap.size

/**
 * 获取当前客户端会话列表
 */
//fun getSessionList(): List<ClientSession> = this.clientSessionMapLock.withLock {
//    return this.clientSessionMap.map {
//        it.value
//    }
//}

/**
 * 获取当前客户端列表
 */
//suspend fun getClientList(): List<ClientDto> = this.clientSessionMapLock.withLock {
//    return this.clientSessionMap.map {
//        it.value.client
//    }
//}

// 保持客户端连接
func holdOnClient(client *dto.ClientDto, tcp net.Conn) {
	//先移除之前的连接
	clientSessionLock.Lock()
	oldSession := clientSessionMap[client.Id]
	clientSessionLock.Unlock()
	if oldSession != nil { //如果存在
		oldSession.Shutdown()
	}

	//新的回话
	session := &ClientSession{
		Client: client,
		tcp:    tcp,
	}
	clientSessionLock.Lock()
	clientSessionMap[client.Id] = session
	clientSessionLock.Unlock()

	//初始化客户端连接池
	tcp_pool.InitEmptyPoolByClient(client.Id)
	udp_pool.InitEmptyPoolByClient(client.Id)

	//开启该客户端下所有隧道监听
	tcp_proxy.AcceptClient(client)
	udp_proxy.AcceptClient(client)
	session.Start()
}

/**
 * 向客户端申请TCP连接池请求
 * @param clientId 客户端ID
 * @param count 申请数量
 */
func (mine *ClientSessionManager) SendTCPPoolRequest(clientId int, count int) {
	send(clientId, HeaderUtil.REQUEST_TCP_POOL, strconv.Itoa(count))
}

/**
 * 向客户端申请UDP连接池请求
 * @param clientID 客户端ID
 * @param count 申请数量
 */
func (mine *ClientSessionManager) SendUDPPoolRequest(clientId int, count int) {
	send(clientId, HeaderUtil.REQUEST_UDP_POOL, strconv.Itoa(count))
}

/**
 * 向客户端当前激活的UDP端口
 * @param clientID 客户端ID
 * @param count 申请数量
 */
func (mine *ClientSessionManager) SendActiveUDPBridge(clientId int, ports string) {
	send(clientId, HeaderUtil.SYNC_ACTIVE_BRIDGE_UDP_PORT, ports)
}

/**
 * 往客户端发送数据
 * @param clientID 客户端ID
 * @param flag 头部标记
 * @param message 头部消息
 */
func send(clientId int, flag uint8, message string) {
	clientSessionLock.Lock()
	session := clientSessionMap[clientId]
	clientSessionLock.Unlock()
	if session == nil {
		return
	}
	err := session.SendHead(flag, message)
	if err != nil {
		session.Shutdown()
	}
}

// 关闭客户端
// - closeSession 当前关闭的对象
func removeSession(closeSession *ClientSession) {
	shutdownProxyAndPoolAndBridge(closeSession.Client.Id)
	clientId := closeSession.Client.Id
	clientSessionLock.Lock()
	session := clientSessionMap[clientId]
	if session != nil { //客户端ID回话如果存在
		if session == closeSession { //当前没有加入新的回话
			delete(clientSessionMap, clientId)
		} else { //由于关闭延迟,有新的回话加入,但是在之前已经关掉了所有的代理监听,所以这里需要再次开启代理监听,概率很小，但不能排除
			go tcp_proxy.AcceptClient(session.Client)
			go udp_proxy.AcceptClient(session.Client)
		}
	}
	clientSessionLock.Unlock()
}

// 关闭与内网穿透客户端的会话连接
func shutdownProxyAndPoolAndBridge(clientId int) {

	//关闭代理监听
	tcp_proxy.ShutdownByClient(clientId)
	udp_proxy.ShutdownByClient(clientId)

	//关闭所有连接池
	tcp_pool.ShutdownByClient(clientId)
	udp_pool.ShutdownByClient(clientId)

	//关闭所有UDP连接池
	//try {
	//   UDPPoolManager.closeByClient(this.client.id!!)
	//} catch (e: Exception) {
	//   e.printStackTrace()
	//}
	//
	//try {
	//   //关闭正在通信的UDP连接
	//   UDPBridgeManager.closeByClient(this.client.id!!)
	//} catch (e: Exception) {
	//   e.printStackTrace()
	//}
}

// 关闭一个客户端
func Shutdown(clientId int) {

	//先移除之前的连接
	clientSessionLock.Lock()
	oldSession := clientSessionMap[clientId]
	clientSessionLock.Unlock()
	if oldSession != nil { //如果存在
		oldSession.Shutdown()
	}
}

// 客户端是否在线监测
func IsOnline(clientId int) bool {
	clientSessionLock.Lock()
	session := clientSessionMap[clientId]
	clientSessionLock.Unlock()
	if session == nil {
		return false
	}
	return session.IsOnline()
}

// 获取当前在线客户端数量
func OnlineCount() int {
	onlineClientCount := 0
	clientSessionLock.Lock()
	for _, session := range clientSessionMap {
		if session.IsOnline() {
			onlineClientCount++
		}
	}
	clientSessionLock.Unlock()
	return onlineClientCount
}
