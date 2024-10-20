package main

import (
	"DairoNPS/forward"
	"DairoNPS/nps/nps_client"
	"DairoNPS/nps/nps_pool"
	"DairoNPS/web"

	//初始化Controller
	_ "DairoNPS/web/controller/bridge_list"
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
	_ "DairoNPS/web/controller/common"
	_ "DairoNPS/web/controller/data_size_log"
	_ "DairoNPS/web/controller/forward"
	_ "DairoNPS/web/controller/index"
	_ "DairoNPS/web/controller/speed_chart"
)

func init() {

	// 初始化共享接口
	nps_pool.Csmi = &nps_client.ClientSessionManager{}
}

var list = make([]int, 0)

func main() {

	//启动web管理
	go web.Start()

	//启动端口转发
	go forward.StartAcceptAll()

	//启动客户端监听
	nps_client.Accept()

}
