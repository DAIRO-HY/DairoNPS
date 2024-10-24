package ForwardAclDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

//代理服务访问控制IP数据操作

/**
 * 添加访问控制的ip列表
 */
func Add(forwardId int, dtoList []dto.ForwardAclDto) {
	DBUtil.Exec("delete from forward_acl where forwardId=?", forwardId)
	for _, it := range dtoList {
		sql :=
			"insert into channel_acl(forwardId,ip,remark)values(?,?,?);"
		DBUtil.Exec(
			sql,
			forwardId,
			it.Ip,
			it.Remark,
		)
	}
}

/**
 * 获取所有数据
 * @return 代理服务Dto
 */
func SelectByForwardId(forwardId int) []*dto.ForwardAclDto {
	sql := "select ip,remark from forward_acl where forwardId = ?"
	return DBUtil.SelectList[dto.ForwardAclDto](sql, forwardId)
}
