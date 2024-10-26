package nps_channel_proxy

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_bridge"
	"DairoNPS/util/ChannelStatisticsUtil"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// 隧道id对应的服务端口监听
var proxyAcceptMap = make(map[int]*ProxyAccept)

// proxyAcceptMap操作互斥锁
var proxyAcceptLock sync.Mutex

func init() {
	ChannelDao.ClearError()
}

// 隧道代理端口数量
func GetProxyCount() int {
	count := 0
	proxyAcceptLock.Lock()
	count = len(proxyAcceptMap)
	proxyAcceptLock.Unlock()
	return count
}

// 开始客户端的所有监听
func AcceptClient(clientDto *dto.ClientDto) {

	//加载统计数据
	ChannelStatisticsUtil.Init()

	//开启NPS客户端ID下所有的隧道
	activeList := ChannelDao.SelectActiveByClientId(clientDto.Id)
	for _, it := range activeList {
		acceptChannel(clientDto, it)
	}
}

// 开始监听某个隧道
func acceptChannel(client *dto.ClientDto, channel *dto.ChannelDto) {
	proxyAcceptLock.Lock()
	oldProxyTCPAccept := proxyAcceptMap[channel.Id]
	if oldProxyTCPAccept != nil { //若该隧道已经在监听,则先停止
		shutdown(oldProxyTCPAccept)
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(channel.ServerPort))
	if err != nil {
		errMsg := fmt.Sprintf("端口:%d 监听失败。err:%q\n", channel.ServerPort, err)
		ChannelDao.SetError(channel.Id, &errMsg)
		LogUtil.Error(errMsg)
		proxyAcceptLock.Unlock()
		return
	}
	ChannelDao.SetError(channel.Id, nil)
	LogUtil.Info(fmt.Sprintf("端口:%d 监听开始\n", channel.ServerPort))
	proxyAccept := &ProxyAccept{
		Client:  client,
		Channel: channel,
		listen:  listener,
	}
	proxyAcceptMap[channel.Id] = proxyAccept
	proxyAcceptLock.Unlock()

	//开启监听
	go proxyAccept.accept()
}

// 关闭监听
// - channelId 隧道id
func ShutdownByChannel(channelId int) {
	proxyAcceptLock.Lock()
	proxyTCPAccept := proxyAcceptMap[channelId]
	if proxyTCPAccept != nil {
		shutdown(proxyTCPAccept)
	}
	proxyAcceptLock.Unlock()

	//关闭隧道所有正在通信的连接
	nps_bridge.ShutdownByChannel(channelId)
}

// 关闭某个客户端下所有的隧道
func ShutdownByClient(clientId int) {

	//关闭客户端所有隧道
	channelIdList := ChannelDao.SelectIdByClientId(clientId)
	for _, it := range channelIdList {
		ShutdownByChannel(it)
	}

	//关闭客户端所有正在通信的连接
	nps_bridge.ShutdownByClient(clientId)
}

// 停止监听端口
func shutdown(proxyTCPAccept *ProxyAccept) {
	proxyTCPAccept.listen.Close()
	channelId := proxyTCPAccept.Channel.Id
	if proxyAcceptMap[channelId] != nil {
		delete(proxyAcceptMap, channelId)
	}
}
