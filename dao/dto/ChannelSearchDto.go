package dto

type ChannelSearchDto struct {
	Id int

	//客户端id
	ClientId int

	//隧道名
	Name string

	//隧道模式, 1:TCP  2:UDP
	Mode int

	//服务端端口
	ServerPort int

	//目标端口(ip:端口)
	TargetPort string

	//入网流量
	InData int64

	//出网流量
	OutData int64

	//启用状态 1:开启  0:停止
	EnableState int

	//是否加密传输
	SecurityState int

	//创建时间
	CreateDate int64

	//最后一次更新时间
	UpdateDate int64

	//客户端名
	ClientName string
}
