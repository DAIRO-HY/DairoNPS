package SystemConfigDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

//系统配置数据操作

/**
 * 获取系统配置
 */
func SelectOne() *dto.SystemConfigDto {
	sql := "select " +
		" in_data_total as inDataTotal" +
		" ,out_data_total as outDataTotal" +
		" from system_config"
	return DBUtil.SelectOne[dto.SystemConfigDto](sql)
}

/**
 * 同步入出网流量
 */
func SetDataLen(dto dto.SystemConfigDto) {
	sql := "update system_config set in_data_total = ?,out_data_total=?"
	DBUtil.Exec(sql, dto.InDataTotal, dto.OutDataTotal)
}
