package DateDataSizeDao

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

func init() {
DebugTimer.Add86()

	//指定删除时间戳
	deleteBeforceTime := time.Now().AddDate(-2, 0, 0).Unix()
	DBUtil.ExecIgnoreError("delete from date_data_size where date < ?", deleteBeforceTime)
}

/**
 * 添加一条统计
 * @param 隧道id
 * @param inData 入网流量
 * @param outData 出网流量
 */
func Add(channelId int, forwardId int, inData int64, outData int64) {
DebugTimer.Add87()

	//当前时间戳（秒）
	date := time.Now().Unix()
	sql :=
		"insert into date_data_size(channelId,forwardId,date,inData,outData)values(?,?,?,?,?)"
	DBUtil.ExecIgnoreError(sql, channelId, forwardId, date, inData, outData)
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
	forwardId int,
	startTime int64,
	endTime int64,
) []*dto.DateDataSizeDto {
DebugTimer.Add88()
	var sql string
	if clientId == 0 && channelId == 0 && forwardId == 0 { //所有的统计
DebugTimer.Add89()
		sql = "select date,inData,outData from date_data_size where date between ? and ?"
		return DBUtil.SelectList[dto.DateDataSizeDto](sql, startTime, endTime)
	} else if clientId != 0 { //统计某个客户端
		sql = "select date,inData,outData from date_data_size where channelId in (select id from channel where clientId = ?) and date between ? and ?"
		return DBUtil.SelectList[dto.DateDataSizeDto](sql, clientId, startTime, endTime)
	} else if channelId != 0 { //统计某个隧道
		sql = "select date,inData,outData from date_data_size where channelId = ? and date between ? and ?"
		return DBUtil.SelectList[dto.DateDataSizeDto](sql, channelId, startTime, endTime)
	} else if forwardId != 0 { //统计某个端口转发
		sql = "select date,inData,outData from date_data_size where forwardId = ? and date between ? and ?"
		return DBUtil.SelectList[dto.DateDataSizeDto](sql, forwardId, startTime, endTime)
	} else {
		return nil
	}
}

// 通过隧道ID删除
func DeleteByChannelId(channelId int) {
DebugTimer.Add90()
	sql := "delete from date_data_size where channelId = ?"
	DBUtil.ExecIgnoreError(sql, channelId)
}

// 通过转发ID删除
func DeleteByForward(forwardId int) {
DebugTimer.Add91()
	sql := "delete from date_data_size where forwardId = ?"
	DBUtil.ExecIgnoreError(sql, forwardId)
}

// 通过客户端ID删除
func DeleteByClientId(clientId int) {
DebugTimer.Add92()
	sql := "delete from date_data_size where channelId in (select id from channel where clientId = ?)"
	DBUtil.ExecIgnoreError(sql, clientId)
}
