package cn.dairo.cls.dao

import cn.dairo.cls.dao.dto.DataLogDto
import java.text.SimpleDateFormat
import java.util.*

/**
 * 数据统计操作
 */
object DataLogDao {

    /**
     * 添加一条统计
     * @param targetId 要统计的目标ID
     * @param type 统计类型, 1:隧道  2:数据转发
     * @param inData 入网流量
     * @param outData 出网流量
     */
    fun add(targetId: Int, type: Int, inData: Long, outData: Long) {

        //时间年月日时分
        val ymdhm = SimpleDateFormat("yyyyMMddHHmm").format(Date())
        val year = ymdhm.substring(0, 4).toInt()//年
        val ym = ymdhm.substring(0, 6).toInt()//年月
        val ymd = ymdhm.substring(0, 8).toInt()//年月日
        val ymdh = ymdhm.substring(0, 10).toInt()//年月日时
        val sql =
            "insert into data_log(targetId,year,ym,ymd,ymdh,ymdhm,inData,outData,type)values(?,?,?,?,?,?,?,?,?)"
        NPSDB.db.exec(
            sql, targetId, year, ym, ymd, ymdh, ymdhm, inData, outData, type
        )
    }

    /**
     * 获取数据流量日志列表
     * @param targetId 要统计的目标ID
     * @param type 统计类型, 1:隧道  2:数据转发
     * @param ymdhmStart 统计时间（分）开始
     * @param ymdhmEnd 统计时间（分）结束
     * @param groupBy 分组
     */
    fun getList(
        targetId: Int,
        type: Int,
        ymdhmStart: Long,
        ymdhmEnd: Long,
        groupBy: String
    ): Map<Long, DataLogDto> {
        val sql =
            "select $groupBy,sum(inData),sum(outData) from data_log where targetId = ? and type = ? and ymdhm between ? and ? group by $groupBy"
        val ymdhmToDataLog = HashMap<Long, DataLogDto>()
        NPSDB.db.selectResult(
            sql, {
                val date = it.getLong(1)
                val inData = it.getLong(2)
                val outData = it.getLong(3)
                val dto = DataLogDto()
                dto.inData = inData
                dto.outData = outData
                ymdhmToDataLog[date] = dto
            },
            targetId, type, ymdhmStart, ymdhmEnd
        )
        return ymdhmToDataLog
    }
}