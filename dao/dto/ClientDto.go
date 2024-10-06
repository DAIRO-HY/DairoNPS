package dto

// ClientDto 客户端信息
type ClientDto struct {

	//id
	Id int

	//名称
	Name string

	//版本号
	Version string

	//连接认证秘钥
	Key string

	//ip地址
	Ip string

	//入网流量
	InDataTotal int64

	//出网流量
	OutDataTotal int64

	//在线状态,0:离线 1:在线
	OnlineState int

	//启用状态
	EnableState int

	//最后一次连接时间
	LastLoginDate int64

	//创建时间
	CreateDate int64

	//最后一次更新时间戳
	UpdateDate int64

	//一些备注信息,错误信息等
	Remark string
}
