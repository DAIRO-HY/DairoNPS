package ChannelAclDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

//隧道访问控制IP数据操作

/**
 * 添加访问控制的ip列表
 */
func add(channelId int, dtoList []dto.ChannelAclDto) {
	DBUtil.Exec("delete from channel_acl where channelId=?", channelId)
	for _, item := range dtoList {

		sql :=
			"insert into channel_acl(channelId,ip,remark)values(?,?,?);"
		DBUtil.Exec(
			sql,
			channelId,
			item.Ip,
			item.Remark,
		)
	}
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func selectByChannelId(channelId int) []*dto.ChannelAclDto {
	sql := "select ip,remark from channel_acl where channelId = ?"
	list := DBUtil.SelectList[dto.ChannelAclDto](sql, channelId)
	return list
}
