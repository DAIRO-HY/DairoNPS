package dto

// 流量统计表DTO
type DataLogDto struct {

	//目标id
	targetId int

	//年
	year int64

	//年月
	ym int64

	//年月日
	ymd int64

	//年月日时
	ymdh int64

	//年月日时分
	ymdhm int64

	//入网流量
	inData int64

	//出网流量
	outData int64

	//统计类型, 1:隧道  2:数据转发
	mode int
}
