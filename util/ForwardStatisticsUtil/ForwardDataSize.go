package ForwardStatisticsUtil

// 端口转发出入流量
type ForwardDataSize struct {

	//入方向
	InData int64

	//上次统计
	PreInData int64

	//出方向
	OutData int64

	//上次统计
	PreOutData int64
}
