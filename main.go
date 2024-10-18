package main

import (
	"DairoNPS/client"
	"DairoNPS/pool"
	"DairoNPS/web"

	_ "DairoNPS/web/controller/bridge_list"
	//初始化Controller
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
	_ "DairoNPS/web/controller/common"
	_ "DairoNPS/web/controller/data_size_log"
	_ "DairoNPS/web/controller/index"
	_ "DairoNPS/web/controller/speed_chart"
)

func init() {

	// 初始化共享接口
	pool.Csmi = &client.ClientSessionManager{}
}

var list = make([]int, 0)

func main() {

	//启动web管理
	go web.Start()

	//启动客户端监听
	client.Accept()
}
