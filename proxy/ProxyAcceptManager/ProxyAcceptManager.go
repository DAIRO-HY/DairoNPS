package ProxyAcceptManager

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/dto"
	"DairoNPS/proxy/ProxyTCPAccept"
	"sync"
)

//代理服务端口监听管理

/**
 * 隧道id对应的服务端口监听
 */
//val channelIdToProxyAccept = ConcurrentHashMap<Int, ProxyAccept>()
var channelIdToProxyAccept = make(map[int]*ProxyTCPAccept.ProxyTCPAccept)

/**
 * mChannelIdToProxyAccept操作互斥锁
 */
//private val channelIdToProxyAcceptLock = Mutex()
var channelIdToProxyAcceptLock sync.Mutex

/**
 * 开启端口监听
 * @param client 客户端
 */
func Accept(client *dto.ClientDto) {

	//开启NPS客户端ID下所有的隧道
	activeList := ChannelDao.SelectActiveByClientId(client.Id)
	for _, it := range activeList {
		accept(client, it)
	}
}

/**
 * 开启端口监听
 * @param channel 隧道信息
 */
func accept(client *dto.ClientDto, channel *dto.ChannelDto) {
	channelIdToProxyAcceptLock.Lock()
	if _, ok := channelIdToProxyAccept[channel.Id]; ok { //若该隧道已经在监听,则先停止
		ProxyTCPAccept.Close(channelIdToProxyAccept[channel.Id])
	}
	//accept := when (channel.type) {
	//    ChannelType.TCP -> ProxyTCPAccept(client, channel)
	//    ChannelType.UDP -> ProxyUDPAccept(client, channel)
	//    else -> return@synchronized
	//}

	proxyTCPAccept := ProxyTCPAccept.ProxyTCPAccept{
		Client:  client,
		Channel: channel,
	}
	channelIdToProxyAccept[channel.Id] = &proxyTCPAccept

	//开启监听
	go ProxyTCPAccept.Start(proxyTCPAccept)
	channelIdToProxyAcceptLock.Unlock()
}

/**
 * 移除隧道监听列表
 */
func RemoveByChannelId(channelId int) {
	channelIdToProxyAcceptLock.Lock()

	if _, ok := channelIdToProxyAccept[channelId]; ok { //若该隧道已经在监听,则先停止
	} else {
		channelIdToProxyAcceptLock.Unlock()
		return
	}

	proxyTCPAccept := channelIdToProxyAccept[channelId]
	delete(channelIdToProxyAccept, channelId)

	//关闭隧道的时候保存流量
	ChannelDao.SetDataLen(proxyTCPAccept.Channel)
	channelIdToProxyAcceptLock.Unlock()
}

/**
 * 关闭监听
 * @param channelId 隧道id
 */
func closeByChannel(channelId int) {
	proxyTCPAccept := channelIdToProxyAccept[channelId]
	if proxyTCPAccept != nil {
		ProxyTCPAccept.Close(proxyTCPAccept)
	}
}

/**
 * 关闭某个客户端下所有的隧道
 */
func CloseByClient(clientId int) {

	//关闭客户端所有隧道
	channelIdList := ChannelDao.SelectIdByClientId(clientId)
	for _, it := range channelIdList {
		closeByChannel(it)
	}
}
