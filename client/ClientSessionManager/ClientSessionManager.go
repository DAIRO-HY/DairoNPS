package ClientSessionManager

import (
	"DairoNPS/client/ClientSession"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/proxy/ProxyAcceptManager"
	"net"
	"strconv"
	"strings"
	"sync"
)

//往客户端发送指令的专用连接

/**
 * 客户端ID对应的Socket连接
 */
var clientSessionMap = make(map[int]ClientSession.ClientSession)

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
 * 添加客户端连接
 */
func Validate(clientSocket net.Conn) {

	//得到头部数据
	header := HeaderUtil.GetHeader(clientSocket)
	if len(header) == 0 {
		return
	}

	headers := strings.Split(header, "|")

	//得到客户端key
	key := headers[0]
	client := ClientDao.SelectByKey(key)
	if client == nil {
		//println("key:${key}不存在")
		clientSocket.Close()
		return
	}
	if client.EnableState == 0 {
		//println("key:${key}的客户端已停止服务")
		clientSocket.Close()
		return
	}
	holdOnClient(client, clientSocket)

	//客户端ip
	ip := clientSocket.RemoteAddr().String()

	//从头部信息中得到客户端版本号
	version := headers[0]
	loginClientDto := dto.ClientDto{
		Id:      client.Id,
		Ip:      ip,
		Version: version,
	}
	ClientDao.SetClientInfo(loginClientDto)

	//将客户端ID返回给客户端
	sendClientId(client.Id)

	//将加密秘钥发送到客户端
	sendClientSecurityKey(client.Id)

	//开启该客户端下所有隧道监听
	ProxyAcceptManager.Accept(client)
}

/**
 * 保持客户端连接
 */
func holdOnClient(client *dto.ClientDto, clientSocket net.Conn) {
	clientSessionMapLock.Lock()

	//先移除之前的连接
	Close(client.Id)
	session := ClientSession.ClientSession{
		Client:       client,
		ClientSocket: clientSocket,
	}
	clientSessionMap[client.Id] = session
	go ClientSession.Start(session)
	clientSessionMapLock.Unlock()
}

/**
 * 将客户端ID返回给客户端
 */
func sendClientId(clientID int) {
	send(clientID, HeaderUtil.SERVER_TO_CLIENT_ID, strconv.Itoa(clientID))
}

/**
 * 将加密秘钥发送到客户端
 */
func sendClientSecurityKey(clientID int) {
	//session := clientSessionMap[clientID]
	//ClientSession.Send(session, []byte(SecurityUtil.clientKeyArray))
}

/**
 * 向客户端申请TCP连接池请求
 * @param clientID 客户端ID
 * @param count 申请数量
 */
func sendTCPPoolRequest(clientID int, count int) {
	send(clientID, HeaderUtil.SERVER_TCP_POOL_REQUEST, strconv.Itoa(count))
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
func send(clientID int, flag byte, message string) {
	session := clientSessionMap[clientID]
	ClientSession.SendHead(session, flag, message)
}

/**
 * 关闭客户端
 */
func Close(clientId int) {
	clientSessionMapLock.Lock()
	if _, ok := clientSessionMap[clientId]; ok { //如果存在
		clientSession := clientSessionMap[clientId]
		ClientSession.Close(clientSession)
		ClientDao.SetDataLen(clientSession.Client)
		delete(clientSessionMap, clientId)
	}
	clientSessionMapLock.Unlock()
}
