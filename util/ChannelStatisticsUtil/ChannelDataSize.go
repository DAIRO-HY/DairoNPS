package ChannelStatisticsUtil

// 隧道出入王流量
type ChannelDataSize struct {

	//当前隧道所属客户端ID
	ClientId int

	//入方向
	InData int64

	//上次统计
	PreInData int64

	//出方向
	OutData int64

	//上次统计
	PreOutData int64
}
