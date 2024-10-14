package form

// 客户端信息
type ChannelEditForm struct {

	// 客户端id
	ClientId int

	// 客户端名
	ClientName string

	// 隧道id
	Id int

	// 隧道名
	Name string

	// 隧道模式, 1:TCP  2:UDP
	Mode int

	// 服务端端口
	ServerPort int

	// 客户端端口(ip:端口)
	TargetPort string

	// 入网流量
	InData string

	// 出网流量
	OutData string

	// 启用状态 1:开启  0:停止
	EnableState string

	// 是否加密传输
	SecurityState int

	// 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
	AclState int

	// 创建时间
	Date string

	// 一些备注信息,错误信息等
	Remark string

	// 黑白名单IP列表
	//aclIp = arrayOf<String>()
}
