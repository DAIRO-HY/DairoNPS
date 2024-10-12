package proxy

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/dto"
	"sync"
	"time"
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
	oldProxyTCPAccept := channelIdToProxyAccept[channel.Id]
	if oldProxyTCPAccept != nil { //若该隧道已经在监听,则先停止
		shutdown(oldProxyTCPAccept)
	}
	//accept := when (channel.type) {
	//    ChannelType.TCP -> ProxyAccept(client, channel)
	//    ChannelType.UDP -> ProxyUDPAccept(client, channel)
	//    else -> return@synchronized
	//}

	proxyAccept := &ProxyAccept{
		Client:  client,
		Channel: channel,
	}
	channelIdToProxyAccept[channel.Id] = proxyAccept
	channelIdToProxyAcceptLock.Unlock()

	//开启监听
	go proxyAccept.Start()
}

/**
 * 关闭监听
 * @param channelId 隧道id
 */
func closeByChannel(channelId int) {
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
		closeByChannel(it)
	}
}

/**
 * 停止监听端口
 */
func shutdown(proxyTCPAccept *ProxyAccept) {
	for {

		//@TODO:这里需要优化,不应该以休眠的方式关闭
		time.Sleep(100 * time.Millisecond)
		if proxyTCPAccept.ProxySocketServer == nil {
			continue
		}
		proxyTCPAccept.ProxySocketServer.Close()
		if proxyTCPAccept.IsFinished {
			break
		}
	}
	removeByChannelId(proxyTCPAccept.Channel.Id)
}

/**
 * 移除隧道监听列表
 */
func removeByChannelId(channelId int) {
	proxyTCPAccept := channelIdToProxyAccept[channelId]
	if proxyTCPAccept != nil {
		proxyTCPAccept := channelIdToProxyAccept[channelId]
		delete(channelIdToProxyAccept, channelId)

		//关闭隧道的时候保存流量
		ChannelDao.SetDataLen(proxyTCPAccept.Channel)
	}
}
