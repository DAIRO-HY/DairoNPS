package ChannelDataStatisticsDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"strconv"
	"time"
)

/**
 * 添加一条统计
 * @param 隧道id
 * @param inData 入网流量
 * @param outData 出网流量
 */
func Add(channelId int, inData int64, outData int64) {

	//年月日时分
	ymdhm, _ := strconv.ParseInt(time.Now().Format("200601021504"), 10, 64)
	sql :=
		"insert into channel_data_statistics(channelId,ymdhm,inData,outData)values(?,?,?,?)"
	DBUtil.ExecIgnoreError(sql, channelId, ymdhm, inData, outData)
}

/**
 * 获取数据流量日志列表
 * @param targetId 要统计的目标ID
 * @param type 统计类型, 1:隧道  2:数据转发
 * @param ymdhmStart 统计时间（分）开始
 * @param ymdhmEnd 统计时间（分）结束
 * @param groupBy 分组
 */
func GetList(
	targetId int,
	mode int,
	ymdhmStart int64,
	ymdhmEnd int64,
	groupBy string,
) map[int64]dto.ChannelDataStatisticsDto {
	sql :=
		"select $groupBy,sum(inData),sum(outData) from data_log where targetId = ? and type = ? and ymdhm between ? and ? group by ?"

	list := DBUtil.SelectToListMap(sql, targetId, mode, ymdhmStart, ymdhmEnd, groupBy)

	dateToData := make(map[int64]dto.ChannelDataStatisticsDto)
	for _, item := range list {
		inData, _ := strconv.ParseInt(item["inData"], 10, 64)
		outData, _ := strconv.ParseInt(item["outData"], 10, 64)
		dto := dto.ChannelDataStatisticsDto{
			InData:  inData,
			OutData: outData,
		}
		date, _ := strconv.ParseInt(item["date"], 10, 64)
		dateToData[date] = dto
	}
	return dateToData
}
