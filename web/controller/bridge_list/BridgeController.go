package bridge_list

import (
	"DairoNPS/dao/ClientDao"
	"DairoNPS/extension/Number"
	"DairoNPS/forward"
	"DairoNPS/nps/nps_bridge/tcp_bridge"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/web"
	"DairoNPS/web/controller/bridge_list/form"
	"net/http"
	"strings"
	"time"
)

// 初始化
func init() {
	http.HandleFunc("/bridge_list/load_data", web.ApiHandler(loadData))
}

// 获取数据
func loadData(search form.BridgeInForm) []form.BridgeOutForm {

	//生成客户端ID对应的客户端名，用来匹配------------------------------------------------------------START
	clientList := ClientDao.SelectAll()
	clientId2Name := make(map[int]string)
	for _, item := range clientList {
		clientId2Name[item.Id] = item.Name
	}
	//生成客户端ID对应的客户端名，用来匹配------------------------------------------------------------END

	outFormList := []form.BridgeOutForm{}

	//当前时间戳
	nowTime := time.Now().UnixMilli()

	//TCP隧道桥接列表统计------------------------------------------------------------START
	tcpBridgeList := tcp_bridge.GetBridgeList()
	for _, it := range tcpBridgeList {
		if search.ClientId != 0 && search.ClientId != it.ClientId {
			continue
		}
		if search.ChannelId != 0 && search.ChannelId != it.Channel.Id {
			continue
		}
		remoteAddr := it.ProxyTCP.RemoteAddr().String()
		outFormList = append(outFormList, form.BridgeOutForm{

			// 客户端名
			ClientName: clientId2Name[it.ClientId],

			// 隧道名
			ChannelName: it.Channel.Name,

			// 隧道模式
			Mode: "TCP",

			// 创建时间
			CreateTime: Number.ToTimeFormat((nowTime - it.CreateTime) / 1000),

			// 在线时间
			LastRWTime: Number.ToTimeFormat((nowTime - it.LastRWTime) / 1000),

			// 用户端ip
			Ip: remoteAddr,
		})
	}
	//TCP隧道桥接列表统计------------------------------------------------------------END

	//UDP隧道桥接列表统计------------------------------------------------------------START
	udpBridgeList := udp_bridge.GetBridgeList()
	for _, it := range udpBridgeList {
		if search.ClientId != 0 && search.ClientId != it.ClientId {
			continue
		}
		if search.ChannelId != 0 && search.ChannelId != it.Channel.Id {
			continue
		}
		remoteAddr := it.ProxyUDPInfo.Key()
		outFormList = append(outFormList, form.BridgeOutForm{

			// 客户端名
			ClientName: clientId2Name[it.ClientId],

			// 隧道名
			ChannelName: it.Channel.Name,

			// 隧道模式
			Mode: "UDP",

			// 创建时间
			CreateTime: Number.ToTimeFormat((nowTime - it.CreateTime) / 1000),

			// 在线时间
			LastRWTime: Number.ToTimeFormat((nowTime - it.LastRWTime) / 1000),

			// 用户端ip
			Ip: remoteAddr,
		})
	}
	//UDP隧道桥接列表统计------------------------------------------------------------END

	//端口转发桥接列表统计------------------------------------------------------------START
	forwardBridgeList := forward.GetBridgeList()
	for _, it := range forwardBridgeList {
		remoteAddr := it.ProxyTCP.RemoteAddr().String()
		ip := strings.Split(remoteAddr, ":")[0]
		outFormList = append(outFormList, form.BridgeOutForm{

			// 客户端名
			ClientName: "端口转发",

			// 隧道名
			ChannelName: it.ForwardDto.Name,

			// 隧道模式
			Mode: "TCP",

			// 创建时间
			CreateTime: Number.ToTimeFormat((nowTime - it.CreateTime) / 1000),

			// 在线时间
			LastRWTime: Number.ToTimeFormat((nowTime - it.LastRWTime) / 1000),

			// 用户端ip
			Ip: ip,
		})
	}
	//隧道桥接列表统计------------------------------------------------------------END
	return outFormList
}
