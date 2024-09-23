package cn.dairo.cls.dao

import cn.dairo.cls.dao.dto.SystemConfigDto

/**
 * 系统配置数据操作
 */
object SystemConfigDao {

    /**
     * 获取系统配置
     */
    fun selectOne(): SystemConfigDto {
        val sql = "select " +
                " in_data_total as inDataTotal" +
                " ,out_data_total as outDataTotal" +
                " from system_config"
        val dto = NPSDB.db.selectOne(SystemConfigDto::class.java, sql)
        return dto
    }

    /**
     * 同步入出网流量
     */
    fun setDataLen(dto: SystemConfigDto) {
        val sql = "update system_config set in_data_total = ?,out_data_total=?"
        NPSDB.db.exec(sql, dto.inDataTotal, dto.outDataTotal)
    }
}