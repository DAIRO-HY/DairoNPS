package client

import (
	"DairoNPS/bridge"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/dao/dto"
	"DairoNPS/pool"
	"DairoNPS/proxy"
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
var clientSessionMapLock sync.Mutex

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

/**
 * 保持客户端连接
 */
func holdOnClient(client *dto.ClientDto, tcp net.Conn) {
	//先移除之前的连接
	clientSessionMapLock.Lock()
	oldSession := clientSessionMap[client.Id]
	clientSessionMapLock.Unlock()
	if oldSession != nil { //如果存在
		oldSession.Close()
	}

	//新的回话
	session := &ClientSession{
		Client: client,
		tcp:    tcp,
	}
	clientSessionMapLock.Lock()
	clientSessionMap[client.Id] = session
	clientSessionMapLock.Unlock()

	//初始化客户端连接池
	pool.InitEmptyPoolByClient(client.Id)

	//开启该客户端下所有隧道监听
	proxy.AcceptClient(client)
	session.Start()
}

/**
 * 向客户端申请TCP连接池请求
 * @param clientId 客户端ID
 * @param count 申请数量
 */
func (csm *ClientSessionManager) SendTCPPoolRequest(clientId int, count int) {
	send(clientId, HeaderUtil.SERVER_TCP_POOL_REQUEST, strconv.Itoa(count))
}

/**
 * 向客户端申请UDP连接池请求
 * @param clientID 客户端ID
 * @param count 申请数量
 */
//suspend fun sendUDPPoolRequest(clientID: Int, count: Int) {
//    this.send(clientID, HeaderUtil.SERVER_UDP_POOL_REQUEST, count.toString())
//}

/**
 * 往客户端发送数据
 * @param clientID 客户端ID
 * @param flag 头部标记
 * @param message 头部消息
 */
func send(clientId int, flag byte, message string) {
	session := clientSessionMap[clientId]
	err := session.SendHead(flag, message)
	if err != nil {
		session.Close()
	}
}

// 关闭客户端
// - closeSession 当前关闭的对象
func removeSession(closeSession *ClientSession) {
	closeProxyAndPoolAndBridge(closeSession.Client.Id)
	clientId := closeSession.Client.Id
	clientSessionMapLock.Lock()
	session := clientSessionMap[clientId]
	if session != nil { //客户端ID回话如果存在
		if session == closeSession { //当前没有加入新的回话
			delete(clientSessionMap, clientId)
		} else { //由于关闭延迟,有新的回话加入,但是在之前已经关掉了所有的代理监听,所以这里需要再次开启代理监听
			proxy.AcceptClient(session.Client)
		}
	}
	clientSessionMapLock.Unlock()
}

/**
 * 关闭与内网穿透客户端的会话连接
 */
func closeProxyAndPoolAndBridge(clientId int) {

	//关闭代理监听
	proxy.CloseByClient(clientId)

	//关闭所有TCP连接池
	pool.CloseByClient(clientId)

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

	//关闭客户端所有正在通信的连接
	bridge.CloseByClient(clientId)
}
