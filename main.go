package main

import (
	"DairoNPS/client"
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/NPSDB"
	"DairoNPS/dao/dto"
	"DairoNPS/pool"
	"DairoNPS/util/SecurityUtil"
	"DairoNPS/util/StatisticsUtil"
	"DairoNPS/web"

	//初始化Controller
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
)

func main() {
	initInterface()

	//加密秘钥生成
	SecurityUtil.Init()
	NPSDB.Init()
	//addTestData()
	StatisticsUtil.Init()

	//隧道流量统计
	go StatisticsUtil.Statistics()

	//启动web管理
	go web.Start()

	//启动客户端监听
	client.Accept()
}

// 初始化共享接口
func initInterface() {
	pool.Csmi = &client.ClientSessionManager{}
}

func addTestData() {
	clientDto := dto.ClientDto{
		Name:   "test",
		Key:    "njeHds*fs4tfsd",
		Remark: "6t7uyghjbmnlkkj",
	}
	ClientDao.Add(&clientDto)

	channel := &dto.ChannelDto{
		//客户端id
		ClientId: clientDto.Id,

		//隧道名
		Name: "TCP测试隧道1",

		//隧道模式, 1:TCP  2:UDP
		Mode:        1,
		EnableState: 1,

		//服务端端口
		ServerPort: 19090,

		//目标端口(ip:端口)
		TargetPort: "127.0.0.1:19091",

		//是否加密传输
		SecurityState: 0,

		//一些备注信息,错误信息等
		Remark: "",
	}
	ChannelDao.Add(channel)
}
