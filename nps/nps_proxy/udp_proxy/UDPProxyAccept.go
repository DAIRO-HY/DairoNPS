package udp_proxy

				import (
					"DairoNPS/DebugTimer"
				)

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

// UDP隧道代理
type UDPProxyAccept struct {
	ClientId int //客户端ID
	Channel  *dto.ChannelDto
	ProxyUDP *net.UDPConn //代理端口监听服务
}

/**
 * 接收数据
 */
func (mine *UDPProxyAccept) accept() {
DebugTimer.Add356()
	for {
DebugTimer.Add357()
		data := make([]byte, NPSConstant.READ_UDP_CACHE_SIZE)

		//从代理客户端读取数据
		length, addr, err := mine.ProxyUDP.ReadFromUDP(data)
		if err != nil {
DebugTimer.Add358()
			LogUtil.Info(fmt.Sprintf("UDP端口:%d 读取数据失败:%q", mine.Channel.ServerPort, err))
			break
		}
		bridge := udp_bridge.ByProxy(addr)
		if bridge == nil { //会话不存在,创建会话
DebugTimer.Add359()
			proxyUDPInfo := &nps.UDPInfo{
				Udp:     mine.ProxyUDP,
				CliAddr: addr,
			}

			//NPS客户端Socket
			pool := udp_pool.GetAndAddPool(mine.Channel.ClientId)
			if pool == nil {
DebugTimer.Add360()
				LogUtil.Error(fmt.Sprintf("UDP端口:%d 没有可用的连接池。", mine.Channel.ServerPort))
				continue
			}
			bridge = udp_bridge.CreateBridge(mine.ClientId, mine.Channel, proxyUDPInfo, pool.UDPInfo)
			if bridge == nil { //创建桥接失败
DebugTimer.Add361()
				continue
			}
		}

		//发送数据到客户端
		sendErr := bridge.SendToClient(data, length)
		if sendErr != nil {
DebugTimer.Add362()
			udp_bridge.RemoveBridge(bridge)
			continue
		}
	}
	LogUtil.Info(fmt.Sprintf("UDP端口:%d监听结束", mine.Channel.ServerPort))
}
