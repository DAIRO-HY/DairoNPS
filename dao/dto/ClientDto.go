package dto

// 客户端信息
type ClientDto struct {

	//id
	Id int

	//名称
	Name string

	//版本号
	version string

	//连接认证秘钥
	Key string

	//ip地址
	ip string

	//入网流量
	inDataTotal int64

	//出网流量
	outDataTotal int64

	//在线状态,0:离线 1:在线
	onlineState int

	//启用状态
	enableState int

	//最后一次连接时间
	lastLoginDate int64

	//创建时间
	createDate int64

	//最后一次更新时间戳
	updateDate int64

	//一些备注信息,错误信息等
	Remark string
}
