package main

import (
	"DairoNPS/client"
	"DairoNPS/pool"
	"fmt"

	//初始化Controller
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
)

func init() {

	// 初始化共享接口
	pool.Csmi = &client.ClientSessionManager{}
}

var list = make([]int, 0)

func main() {

	for i := 0; i < 10; i++ {
		list = append(list, i)
	}

	list2 := list

	for i := 0; i < 10; i++ {
		list = append(list, i)
	}

	fmt.Println(len(list))
	fmt.Println(len(list2))

	////启动web管理
	//go web.Start()
	//
	////启动客户端监听
	//client.Accept()
}
