package bridge_list

import (
	"DairoNPS/bridge"
	"DairoNPS/web"
	"DairoNPS/web/controller/bridge_list/form"
	"net/http"
)

// 初始化
func init() {
	http.HandleFunc("/bridge_list/load_data", web.ApiHandler(loadData))
}

// 获取数据
func loadData(search form.BridgeInForm) []form.BridgeOutForm {
	outFormList := []form.BridgeOutForm{}

	//val now = System.currentTimeMillis()
	bridgeList := bridge.GetBridgeList()
	for _, it := range bridgeList {
		outFormList = append(outFormList, form.BridgeOutForm{

			// 客户端名
			ClientName: "客户端名",

			// 隧道名
			ChannelName: it.Channel.Name,

			// 隧道模式
			//Mode string

			// 在线时间
			//OnlineTime string

			// 客户端ip
			//Ip:it.Ip
		})
	}
	return outFormList
}
