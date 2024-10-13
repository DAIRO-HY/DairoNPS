package form

// 客户端信息
type ClientListForm struct {

	// id
	Id int

	// 名称
	Name string

	// 版本号
	Version string

	// 连接认证秘钥
	Key string

	// ip地址
	Ip string

	// 入网流量
	InDataTotal int64

	// 出网流量
	OutDataTotal int64

	// 在线状态,0:离线 1:在线
	OnlineState bool

	// 启用状态
	EnableState int
}
