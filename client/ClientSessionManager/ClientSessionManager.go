package ClientSessionManager

import (
	"DairoNPS/client/ClientSession"
	"DairoNPS/client/HeaderUtil"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/pool/TCPPoolManager"
	"DairoNPS/proxy/ProxyAcceptManager"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type ClientSessionManager struct{}

//往客户端发送指令的专用连接

/**
 * 客户端ID对应的Socket连接
 */
var clientSessionMap = make(map[int]*ClientSession.ClientSession)

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
		log.Printf("-->客户端：%s不存在\n", key)
		clientSocket.Close()
		return
	}
	if client.EnableState == 0 {
		log.Printf("-->客户端：%s的客户端已被停止服务\n", key)
		clientSocket.Close()
		return
	}
	holdOnClient(client, clientSocket)

	////客户端ip
	//ip := clientSocket.RemoteAddr().String()
	//
	////从头部信息中得到客户端版本号
	//version := headers[0]
	//loginClientDto := dto.ClientDto{
	//	Id:      client.Id,
	//	Ip:      ip,
	//	Version: version,
	//}
	//ClientDao.SetClientInfo(loginClientDto)

	//开启该客户端下所有隧道监听
	ProxyAcceptManager.Accept(client)
}

/**
 * 保持客户端连接
 */
func holdOnClient(client *dto.ClientDto, clientSocket net.Conn) {
	//先移除之前的连接
	Close(client.Id)
	clientSessionMapLock.Lock()
	session := &ClientSession.ClientSession{
		Client:       client,
		ClientSocket: clientSocket,
	}
	clientSessionMap[client.Id] = session

	//初始化客户端连接池
	TCPPoolManager.InitEmptyPoolByClient(client.Id)
	go session.Start(Close)
	clientSessionMapLock.Unlock()
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
	session.SendHead(flag, message)
}

/**
 * 关闭客户端
 */
func Close(clientId int) {
	clientSessionMapLock.Lock()
	if _, ok := clientSessionMap[clientId]; ok { //如果存在
		clientSession := clientSessionMap[clientId]

		//去关闭连接
		clientSession.Close()
		ClientDao.SetDataLen(clientSession.Client)
		delete(clientSessionMap, clientId)
	}
	clientSessionMapLock.Unlock()
}
