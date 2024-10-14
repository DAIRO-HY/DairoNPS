package main

import (
	"DairoNPS/client"
	"DairoNPS/pool"
	"DairoNPS/web"
	//初始化Controller
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
)

func init() {

	// 初始化共享接口
	pool.Csmi = &client.ClientSessionManager{}
}

func main() {

	//启动web管理
	go web.Start()

	//启动客户端监听
	client.Accept()
}
