package SystemConfigDao

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

//系统配置数据操作

/**
 * 获取系统配置
 */
func SelectOne() *dto.SystemConfigDto {
DebugTimer.Add106()
	sql := "select inData,outData from system_config"
	return DBUtil.SelectOne[dto.SystemConfigDto](sql)
}

/**
 * 同步入出网流量
 */
func AddDataSize(inData int64, outData int64) {
DebugTimer.Add107()
	sql := "update system_config set inData = inData + ?,outData = outData + ?"
	DBUtil.ExecIgnoreError(sql, inData, outData)
}
