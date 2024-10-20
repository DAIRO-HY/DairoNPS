package dto

// 代理服务器Dto
type ForwardDto struct {

	//代理服务ID
	Id int

	//代理服务名
	Name string

	//服务端端口
	Port int

	//目标端口(ip:端口)
	TargetPort string

	//入网流量
	InData int64

	//出网流量
	OutData int64

	//启用状态 1:开启  0:停止
	EnableState int

	//黑白名单开启状态 0:关闭 1:白名单 2:黑名单
	AclState int

	//创建时间
	Date int64

	//一些备注信息,错误信息等
	Remark string
}
