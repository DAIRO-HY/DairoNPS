package udp_client

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/util/LogUtil"
	"fmt"
	"log"
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
		log.Fatalf("UDP端口:%s 监听失败:%q\n", NPSConstant.UDPPort, err)
		return
	}
	LogUtil.Info(fmt.Sprintf("UDP端口:%s 监听成功.", NPSConstant.UDPPort))
	for {
		data := make([]byte, 10*1024)

		//从客户端读取数据
		length, clientAddr, err := udp.ReadFromUDP(data)
		if err != nil {
			LogUtil.Error(fmt.Sprintf("UDP端口:%s 读取数据失败:%q", NPSConstant.UDPPort, err))
			break
		}
		bridge := udp_bridge.ByClient(clientAddr)
		if bridge != nil { //桥接通信已经存在
			bridge.SendToProxy(data, length)
		} else { //这可能是一个连接池
			clientIdStr := string(data[:length])
			clientId, toClientErr := strconv.ParseInt(clientIdStr, 10, 64)
			if toClientErr != nil {
				continue
			}
			udpInfo := &nps.UDPInfo{
				Udp:     udp,
				CliAddr: clientAddr,
			}

			//加入连接池
			udp_pool.Add(udpInfo, int(clientId))
		}
	}
}
