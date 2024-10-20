package form

type GetDataInForm struct {

	//客户端ID
	ClientId int

	//隧道ID
	ChannelId int

	//端口转发ID
	ForwordId int

	//入网流量
	StartTime int64

	// 出网流量
	EndTime int64
}
