package cn.dairo.cls.dao.dto

/**
 * 流量统计表DTO
 */
class DataLogDto {

    /**
     * 目标id
     */
    var targetId: Int? = null

    /**
     * 年
     */
    var year: Long? = null

    /**
     * 年月
     */
    var ym: Long? = null

    /**
     * 年月日
     */
    var ymd: Long? = null

    /**
     * 年月日时
     */
    var ymdh: Long? = null

    /**
     * 年月日时分
     */
    var ymdhm: Long? = null

    /**
     * 入网流量
     */
    var inData: Long? = null

    /**
     * 出网流量
     */
    var outData: Long? = null

    /**
     * 统计类型, 1:隧道  2:数据转发
     */
    var type: Int? = null
}

