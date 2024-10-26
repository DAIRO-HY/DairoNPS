package form

/**
 * 端口转发列表返回表单
 */
type ForwardListOutForm struct {

	// 隧道id
	Id int

	// 隧道名
	Name string

	// 服务端口
	Port int

	// 目标端口(ip:端口)
	TargetPort string

	// 入网流量
	InData string

	// 出网流量
	OutData string

	// 启用状态 1:开启  0:停止
	EnableState int

	//错误信息
	Error string
}
