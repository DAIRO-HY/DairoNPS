package dto

type ChannelDto struct {

	//隧道ID
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

	//黑白名单开启状态 0:关闭 1:白名单 2:黑名单
	aclState int

	//创建时间
	//createDate int64

	//最后一次更新时间
	updateDate int64

	//一些备注信息,错误信息等
	remark string
}
