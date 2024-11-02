package udp_proxy

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/dto"
	"DairoNPS/nps"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
)

/**
 * UDP隧道代理
 */
type UDPProxyAccept struct {
	Client  *dto.ClientDto
	Channel *dto.ChannelDto

	/**
	 * 最后一次统计到入网流量
	 * 流量统计时用到
	 */
	lastInData int64

	/**
	 * 最后一次统计到出网流量
	 * 流量统计时用到
	 */
	lastOutData int64

	//标记监听已经结束
	isFinished bool

	//代理端口监听服务
	udpAddr *net.UDPAddr
}

/**
 * 接收数据
 */
func (mine *UDPProxyAccept) accept() {
	for {

		//代理服务端Socket
		proxySocket, err := net.ListenUDP("udp", mine.udpAddr)
		if err != nil {
			LogUtil.Info(fmt.Sprintf("端口:%d 监听结束\n", mine.Channel.ServerPort))
			break
		}
		LogUtil.Info(fmt.Sprintf("端口:%d 监听到一个连接\n", mine.Channel.ServerPort))
		data := make([]byte, NPSConstant.READ_UDP_CACHE_SIZE)

		//从代理客户端读取数据
		length, addr, err := proxySocket.ReadFromUDP(data)
		if err != nil {
			LogUtil.Info(fmt.Sprintf("端口:%d 读取数据失败\n", mine.Channel.ServerPort))
			continue
		}
		bridge := udp_bridge.ByProxy(addr)
		if bridge == nil { //会话不存在,创建会话
			proxyUDPInfo := &nps.UDPInfo{
				Socket: proxySocket,
				Addr:   addr,
			}

			//NPS客户端Socket
			clientUDPInfo := udp_pool.GetAndAddPool(mine.Channel.ClientId)
			if clientUDPInfo == nil {
				LogUtil.Error(fmt.Sprintf("客户端: %d没有可用的连接池。", mine.Channel.ClientId))
				proxySocket.Close()
				continue
			}
			bridge = udp_bridge.CreateBridge(mine.Client, mine.Channel, proxyUDPInfo, clientUDPInfo)
		}

		//发送数据到客户端
		bridge.SendToClient(data, length)
		//println("-->端口:${channel.serverPort}收到一条连接请求")

	}
	LogUtil.Info(fmt.Sprintf("端口:%d监听结束", mine.Channel.ServerPort))
	mine.isFinished = true
}
