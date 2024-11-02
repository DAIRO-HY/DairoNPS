package udp_proxy

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/dto"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/util/ChannelStatisticsUtil"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// 隧道id对应的服务端口监听
var proxyAcceptMap = make(map[int]*UDPProxyAccept)

// proxyAcceptMap操作互斥锁
var proxyAcceptLock sync.Mutex

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
		if it.Mode == 2 { //只监听UDP隧道
			acceptChannel(clientDto, it)
		}
	}
}

// 开始监听某个隧道
func acceptChannel(client *dto.ClientDto, channel *dto.ChannelDto) {
	proxyAcceptLock.Lock()
	oldProxyUDPAccept := proxyAcceptMap[channel.Id]
	if oldProxyUDPAccept != nil { //若该隧道已经在监听,则先停止
		shutdown(oldProxyUDPAccept)
	}

	// 创建一个 UDP 地址
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(channel.ServerPort))
	if err != nil {
		errMsg := fmt.Sprintf("端口:%d 监听失败。err:%q\n", channel.ServerPort, err)
		ChannelDao.SetError(channel.Id, &errMsg)
		LogUtil.Error(errMsg)
		proxyAcceptLock.Unlock()
		return
	}
	ChannelDao.SetError(channel.Id, nil)
	LogUtil.Info(fmt.Sprintf("端口:%d 监听开始\n", channel.ServerPort))
	proxyAccept := &UDPProxyAccept{
		Client:  client,
		Channel: channel,
		udpAddr: addr,
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
	proxyUDPAccept := proxyAcceptMap[channelId]
	if proxyUDPAccept != nil {
		shutdown(proxyUDPAccept)
	}
	proxyAcceptLock.Unlock()

	//关闭隧道所有正在通信的连接
	udp_bridge.ShutdownByChannel(channelId)
}

// 关闭某个客户端下所有的隧道
func ShutdownByClient(clientId int) {

	//关闭客户端所有隧道
	channelIdList := ChannelDao.SelectIdByClientId(clientId)
	for _, it := range channelIdList {
		ShutdownByChannel(it)
	}

	//关闭客户端所有正在通信的连接
	udp_bridge.ShutdownByClient(clientId)
}

// 停止监听端口
func shutdown(proxyUDPAccept *UDPProxyAccept) {
	//proxyUDPAccept.listen.Close()
	channelId := proxyUDPAccept.Channel.Id
	if proxyAcceptMap[channelId] != nil {
		delete(proxyAcceptMap, channelId)
	}
}
