package form

// 客户端信息
type ClientEditForm struct {

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
	InData string

	// 出网流量
	OutData string

	// 在线状态
	OnlineState string

	// 启用状态
	EnableState string

	// 最后一次连接时间
	LastLoginDate string

	// 创建时间
	Date string

	// 一些备注信息,错误信息等
	Remark string
}
