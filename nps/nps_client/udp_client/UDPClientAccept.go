package udp_client

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strconv"
)

/**
 * 接收客户端UDP消息
 */

///**
// * UDP与客户端通信的服务端
// */
//private var datagramSocket: DatagramSocket? = null
//
///**
// * 开启服务
// */
//fun start() = GlobalScope.launch(CLSDispatchers.IO) {
//    receive()
//}

/**
 * 监听客户端UDP连接
 */
func Accept() {

	// 创建一个 UDP 地址
	addr, err := net.ResolveUDPAddr("udp", ":"+NPSConstant.UDPPort)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("UDP端口:%s 监听失败:%q", NPSConstant.UDPPort, err))
		return
	}

	// 监听指定地址
	udp, err := net.ListenUDP("udp", addr)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("UDP端口:%s 监听失败:%q", NPSConstant.UDPPort, err))
		return
	}
	LogUtil.Info(fmt.Sprintf("UDP端口:%s 监听成功.", NPSConstant.UDPPort))
	for {
		data := make([]byte, 10*1024)
		length, addr, err := udp.ReadFromUDP(data)
		if err != nil {
			//LogUtil.Error(fmt.Sprintf("UDP端口:%s 监听失败:%q",NPSConstant.UDPPort,err))
			continue
		}

		//获取IP地址
		////不能使用packet.address.hostName,会出现延迟
		//val ip = packet.address.hostAddress
		//
		////终端设备端口
		//val port = packet.port

		bridge := udp_bridge.ByClient(addr)
		if bridge != nil {
			bridge.SendToProxy(data, length)
		} else {
			clientIdStr := string(data[:length])
			clientId, _ := strconv.ParseInt(clientIdStr, 10, 64)

			udpInfo := &nps.UDPInfo{
				Socket: udp,
				Addr:   addr,
			}

			//加入连接池
			udp_pool.Add(udpInfo, int(clientId))
		}
	}
}
