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
	sql := "select inDataTotal,outDataTotal from system_config"
	return DBUtil.SelectOne[dto.SystemConfigDto](sql)
}

/**
 * 同步入出网流量
 */
func SetDataLen(dto dto.SystemConfigDto) {
	sql := "update system_config set inDataTotal = ?,outDataTotal=?"
	DBUtil.Exec(sql, dto.InDataTotal, dto.OutDataTotal)
}
