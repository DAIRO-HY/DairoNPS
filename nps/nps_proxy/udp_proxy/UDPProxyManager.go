package udp_proxy

				import (
					"DairoNPS/DebugTimer"
				)

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
DebugTimer.Add363()
	count := 0
	proxyAcceptLock.Lock()
	count = len(proxyAcceptMap)
	proxyAcceptLock.Unlock()
	return count
}

// 开始客户端的所有监听
func AcceptClient(clientDto *dto.ClientDto) {
DebugTimer.Add364()

	//加载统计数据
	ChannelStatisticsUtil.Init()

	//开启NPS客户端ID下所有的隧道
	activeList := ChannelDao.SelectActiveByClientId(clientDto.Id)
	for _, it := range activeList {
DebugTimer.Add365()
		if it.Mode == 2 { //只监听UDP隧道
DebugTimer.Add366()
			acceptChannel(clientDto.Id, it)
		}
	}
}

// 开始监听某个隧道
func acceptChannel(ClientId int, channel *dto.ChannelDto) {
DebugTimer.Add367()
	proxyAcceptLock.Lock()
	oldProxyUDPAccept := proxyAcceptMap[channel.Id]
	if oldProxyUDPAccept != nil { //若该隧道已经在监听,则先停止
DebugTimer.Add368()
		shutdown(oldProxyUDPAccept)
	}

	// 创建一个 UDP 地址
	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(channel.ServerPort))

	//代理服务端Socket
	proxySocket, err := net.ListenUDP("udp", addr)
	if err != nil {
DebugTimer.Add369()
		errMsg := fmt.Sprintf("UDP端口:%d 监听失败。err:%q\n", channel.ServerPort, err)
		ChannelDao.SetError(channel.Id, &errMsg)
		LogUtil.Error(errMsg)
		proxyAcceptLock.Unlock()
		return
	}
	ChannelDao.SetError(channel.Id, nil)
	LogUtil.Info(fmt.Sprintf("UDP端口:%d 监听开始\n", channel.ServerPort))
	proxyAccept := &UDPProxyAccept{
		ClientId: ClientId,
		Channel:  channel,
		ProxyUDP: proxySocket,
	}
	proxyAcceptMap[channel.Id] = proxyAccept
	proxyAcceptLock.Unlock()

	//开启监听
	go proxyAccept.accept()
}

// 关闭监听
// - channelId 隧道id
func ShutdownByChannel(channelId int) {
DebugTimer.Add370()
	proxyAcceptLock.Lock()
	proxyUDPAccept := proxyAcceptMap[channelId]
	if proxyUDPAccept != nil {
DebugTimer.Add371()
		shutdown(proxyUDPAccept)
	}
	proxyAcceptLock.Unlock()

	//关闭隧道所有正在通信的连接
	udp_bridge.ShutdownByChannel(channelId)
}

// 关闭某个客户端下所有的隧道
func ShutdownByClient(clientId int) {
DebugTimer.Add372()

	//关闭客户端所有隧道
	channelIdList := ChannelDao.SelectIdByClientId(clientId)
	for _, it := range channelIdList {
DebugTimer.Add373()
		ShutdownByChannel(it)
	}

	//关闭客户端所有正在通信的连接
	udp_bridge.ShutdownByClient(clientId)
}

// 停止监听端口
func shutdown(proxyUDPAccept *UDPProxyAccept) {
DebugTimer.Add374()
	proxyUDPAccept.ProxyUDP.Close()
	channelId := proxyUDPAccept.Channel.Id
	if proxyAcceptMap[channelId] != nil {
DebugTimer.Add375()
		delete(proxyAcceptMap, channelId)
	}
}
