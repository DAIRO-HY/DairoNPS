package DataLogDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"strconv"
)

/**
 * 数据统计操作
 */

/**
 * 添加一条统计
 * @param targetId 要统计的目标ID
 * @param type 统计类型, 1:隧道  2:数据转发
 * @param inData 入网流量
 * @param outData 出网流量
 */
func add(targetId int, mode int, inData int64, outData int64) {

	////时间年月日时分
	//val ymdhm = SimpleDateFormat("yyyyMMddHHmm").format(Date())
	//val year = ymdhm.substring(0, 4).toInt()//年
	//val ym = ymdhm.substring(0, 6).toInt()//年月
	//val ymd = ymdhm.substring(0, 8).toInt()//年月日
	//val ymdh = ymdhm.substring(0, 10).toInt()//年月日时
	//val sql =
	//    "insert into data_log(targetId,year,ym,ymd,ymdh,ymdhm,inData,outData,type)values(?,?,?,?,?,?,?,?,?)"
	//NPSDB.db.exec(
	//    sql, targetId, year, ym, ymd, ymdh, ymdhm, inData, outData, type
	//)
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
) map[int64]dto.DataLogDto {
	sql :=
		"select $groupBy,sum(inData),sum(outData) from data_log where targetId = ? and type = ? and ymdhm between ? and ? group by ?"

	list := DBUtil.SelectToListMap(sql, targetId, mode, ymdhmStart, ymdhmEnd, groupBy)

	dateToData := make(map[int64]dto.DataLogDto)
	for _, item := range list {
		inData, _ := strconv.ParseInt(item["inData"], 10, 64)
		outData, _ := strconv.ParseInt(item["outData"], 10, 64)
		dto := dto.DataLogDto{
			InData:  inData,
			OutData: outData,
		}
		date, _ := strconv.ParseInt(item["date"], 10, 64)
		dateToData[date] = dto
	}
	return dateToData
}
