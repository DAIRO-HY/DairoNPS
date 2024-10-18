package ChannelDataStatisticsDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

/**
 * 添加一条统计
 * @param 隧道id
 * @param inData 入网流量
 * @param outData 出网流量
 */
func Add(channelId int, inData int64, outData int64) {

	//当前时间戳（秒）
	date := time.Now().Unix()
	sql :=
		"insert into channel_data_statistics(channelId,date,inData,outData)values(?,?,?,?)"
	DBUtil.ExecIgnoreError(sql, channelId, date, inData, outData)
}

/**
 * 获取数据流量日志列表
 * @param clientId 客户端id
 * @param channelId 隧道ID
 * @param startTime 开始时间
 * @param endTime 结束时间
 * @return 数据流量统计列表
 */
func SelectList(
	clientId int,
	channelId int,
	startTime int64,
	endTime int64,
) []*dto.ChannelDataStatisticsDto {
	var sql string
	if clientId == 0 && channelId == 0 { //所有的统计
		sql = "select date,inData,outData from channel_data_statistics where date between ? and ?"
		return DBUtil.SelectList[dto.ChannelDataStatisticsDto](sql, startTime, endTime)
	} else if clientId != 0 { //统计某个客户端
		sql = "select date,inData,outData from channel_data_statistics where channelId in (select id from channel where clientId = ?) and date between ? and ?"
		return DBUtil.SelectList[dto.ChannelDataStatisticsDto](sql, clientId, startTime, endTime)
	} else { //统计某个隧道
		sql = "select date,inData,outData from channel_data_statistics where channelId = ? and date between ? and ?"
		return DBUtil.SelectList[dto.ChannelDataStatisticsDto](sql, channelId, startTime, endTime)
	}
}
