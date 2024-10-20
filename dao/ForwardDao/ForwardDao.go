package ForwardDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

//数据转发操作

/**
 * 添加一条隧道
 */
func Add(dto *dto.ForwardDto) {
	sql :=
		"insert into forward(name,port,targetPort,aclState)values(?,?,?,?)"
	id := DBUtil.InsertIgnoreError(
		sql,
		dto.Name,
		dto.Port,
		dto.TargetPort,
		dto.AclState,
	)
	dto.Id = int(id)
}

/**
 * 通过id获取一条数据
 * @param id 隧道id
 * @return 隧道Dto
 */
func SelectOne(id int) *dto.ForwardDto {
	sql := "select * from forward where id = ?"
	return DBUtil.SelectOne[dto.ForwardDto](sql, id)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectAll() []*dto.ForwardDto {
	sql := "select * from forward"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectActive() []*dto.ForwardDto {
	sql := "select * from forward where enableState = 1"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

// 通过端口查询一条数据
func SelectByPort(port int) *dto.ForwardDto {
	sql := "select * from forward where port = ?"
	return DBUtil.SelectOne[dto.ForwardDto](sql, port)
}

/**
 * 更新一条数据
 */
func Update(dto *dto.ForwardDto) {
	sql :=
		"update forward set name = ?,port=?,targetPort=?,aclState=?,remark=? where id = ?"
	DBUtil.ExecIgnoreError(
		sql,
		dto.Name,
		dto.Port,
		dto.TargetPort,
		dto.AclState,
		dto.Remark,
		dto.Id,
	)
}

// 设置可用状态
func SetEnableState(id int, state int) {
	sql := "update forward set enableState = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, state, id)
}

/**
 * 同步入出网流量
 */
func SetDataSize(id int, inData int64, outData int64) {
	sql := "update forward set inData = ?,outData=? where id = ?"
	DBUtil.ExecIgnoreError(sql, inData, outData, id)
}

/**
 * 通过id删除一条数据
 * @param id 隧道id
 */
func Delete(id int) {
	sql := "delete from forward where id = ?"
	DBUtil.Exec(sql, id)
}

/**
 * 设置备注信息
 */
func SetRemark(id int, remark string) {
	sql := "update forward set remark = ? where id = ?"
	DBUtil.Exec(sql, remark, id)
}
