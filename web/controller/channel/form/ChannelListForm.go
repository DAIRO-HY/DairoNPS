package form

// 隧道信息
type ChannelListForm struct {

	// 隧道id
	Id int

	// 客户端ID
	ClientId int

	// 客户端名
	ClientName string

	// 隧道名
	Name string

	// 隧道模式, 1:TCP  2:UDP
	Mode string

	// 服务端口
	ServerPort int

	// 目标端口(ip:端口)
	TargetPort string

	// 入网流量
	InData string

	// 出网流量
	OutData string

	// 启用状态 1:开启  0:停止
	EnableStateText string

	// 启用状态 1:开启  0:停止
	EnableState int

	// 是否加密传输
	SecurityState int
}
