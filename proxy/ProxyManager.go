package proxy

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/dto"
	"DairoNPS/util/StatisticsUtil"
	"fmt"
	"net"
	"sync"
)

//代理服务端口监听管理

/**
 * 隧道id对应的服务端口监听
 */
//val channelIdToProxyAccept = ConcurrentHashMap<Int, ProxyAccept>()
var channelIdToProxyAccept = make(map[int]*ProxyAccept)

/**
 * mChannelIdToProxyAccept操作互斥锁
 */
//private val channelIdToProxyAcceptLock = Mutex()
var channelIdToProxyAcceptLock sync.Mutex

// 开始客户端的所有监听
func AcceptClient(client *dto.ClientDto) {

	//加载统计数据
	StatisticsUtil.LoadChannelDataLog()

	//开启NPS客户端ID下所有的隧道
	activeList := ChannelDao.SelectActiveByClientId(client.Id)
	for _, it := range activeList {
		acceptChannel(client, it)
	}
}

// 开始监听某个隧道
func acceptChannel(client *dto.ClientDto, channel *dto.ChannelDto) {
	channelIdToProxyAcceptLock.Lock()
	oldProxyTCPAccept := channelIdToProxyAccept[channel.Id]
	if oldProxyTCPAccept != nil { //若该隧道已经在监听,则先停止
		shutdown(oldProxyTCPAccept)
	}
	//accept := when (channel.type) {
	//    ChannelType.TCP -> ProxyAccept(client, channel)
	//    ChannelType.UDP -> ProxyUDPAccept(client, channel)
	//    else -> return@synchronized
	//}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", channel.ServerPort))
	if err != nil {
		fmt.Printf("端口:%d 监听失败。err:%p\n", channel.ServerPort, err)
		channelIdToProxyAcceptLock.Unlock()
		return
	}
	fmt.Printf("端口:%d 监听开始\n", channel.ServerPort)
	proxyAccept := &ProxyAccept{
		Client:  client,
		Channel: channel,
		listen:  listener,
	}
	channelIdToProxyAccept[channel.Id] = proxyAccept
	channelIdToProxyAcceptLock.Unlock()

	//开启监听
	go proxyAccept.accept()
}

/**
 * 关闭监听
 * @param channelId 隧道id
 */
func CloseByChannel(channelId int) {
	channelIdToProxyAcceptLock.Lock()
	proxyTCPAccept := channelIdToProxyAccept[channelId]
	if proxyTCPAccept != nil {
		shutdown(proxyTCPAccept)
	}
	channelIdToProxyAcceptLock.Unlock()
}

/**
 * 关闭某个客户端下所有的隧道
 */
func CloseByClient(clientId int) {

	//关闭客户端所有隧道
	channelIdList := ChannelDao.SelectIdByClientId(clientId)
	for _, it := range channelIdList {
		CloseByChannel(it)
	}
}

/**
 * 停止监听端口
 */
func shutdown(proxyTCPAccept *ProxyAccept) {
	proxyTCPAccept.listen.Close()
	removeByChannelId(proxyTCPAccept.Channel.Id)
}

/**
 * 移除隧道监听列表
 */
func removeByChannelId(channelId int) {
	proxyTCPAccept := channelIdToProxyAccept[channelId]
	if proxyTCPAccept != nil {
		delete(channelIdToProxyAccept, channelId)
	}
}
