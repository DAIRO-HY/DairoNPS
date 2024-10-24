package form

/**
 * 客户端信息
 */
type ForwardEditForm struct {

	// 隧道id
	Id int

	// 隧道名
	Name string

	// 服务端端口
	Port int

	// 客户端端口(ip:端口)
	TargetPort string

	// 入网流量
	InDataTotal string

	// 出网流量
	OutDataTotal string

	// 启用状态 1:开启  0:停止
	EnableState string

	// 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
	AclState int

	// 创建时间
	CreateDate string

	// 一些备注信息,错误信息等
	Remark string

	// 黑白名单IP列表
	//aclIp = arrayOf<String>()
}
