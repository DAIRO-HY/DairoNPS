package dto

// 代理服务器Dto
type ForwardDto struct {

	//代理服务ID
	id int

	//代理服务名
	name string

	//服务端端口
	port int

	//目标端口(ip:端口)
	targetPort string

	//入网流量
	inDataTotal int64

	//出网流量
	outDataTotal int64

	//启用状态 1:开启  0:停止
	enableState int

	//黑白名单开启状态 0:关闭 1:白名单 2:黑名单
	aclState int

	//创建时间
	createDate int64

	//最后一次更新时间
	updateDate int64

	//一些备注信息,错误信息等
	remark string
}
