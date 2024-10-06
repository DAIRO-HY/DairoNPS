package dto

// 流量统计表DTO
type DataLogDto struct {

	//目标id
	TargetId int

	//年
	Year int64

	//年月
	Ym int64

	//年月日
	Ymd int64

	//年月日时
	Ymdh int64

	//年月日时分
	Ymdhm int64

	//入网流量
	InData int64

	//出网流量
	OutData int64

	//统计类型, 1:隧道  2:数据转发
	Mode int
}
