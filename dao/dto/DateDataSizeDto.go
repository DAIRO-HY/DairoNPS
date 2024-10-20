package dto

// 流量统计表DTO
type DateDataSizeDto struct {

	//隧道id
	ChannelId int

	//年月日时分秒
	Date int64

	//入网流量
	InData int64

	//出网流量
	OutData int64
}
