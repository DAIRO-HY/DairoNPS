package dto

type ChannelSearchDto struct {
	id int

	//客户端id
	clientId int

	//隧道名
	name string

	//隧道模式, 1:TCP  2:UDP
	mode int

	//服务端端口
	serverPort int

	//目标端口(ip:端口)
	targetPort string

	//入网流量
	inDataTotal int64

	//出网流量
	outDataTotal int64

	//启用状态 1:开启  0:停止
	enableState int

	//是否加密传输
	securityState int

	//创建时间
	createDate int64

	//最后一次更新时间
	updateDate int64

	//客户端名
	clientName string
}
