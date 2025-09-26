package ForwardDao

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
)

func init() {
DebugTimer.Add93()
	ClearError()
}

/**
 * 添加一条隧道
 */
func Add(dto *dto.ForwardDto) {
DebugTimer.Add94()
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
DebugTimer.Add95()
	sql := "select * from forward where id = ?"
	return DBUtil.SelectOne[dto.ForwardDto](sql, id)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectAll() []*dto.ForwardDto {
DebugTimer.Add96()
	sql := "select * from forward"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectActive() []*dto.ForwardDto {
DebugTimer.Add97()
	sql := "select * from forward where enableState = 1"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

// 通过端口查询一条数据
func SelectByPort(port int) *dto.ForwardDto {
DebugTimer.Add98()
	sql := "select * from forward where port = ?"
	return DBUtil.SelectOne[dto.ForwardDto](sql, port)
}

/**
 * 更新一条数据
 */
func Update(dto *dto.ForwardDto) {
DebugTimer.Add99()
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
DebugTimer.Add100()
	sql := "update forward set enableState = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, state, id)
}

/**
 * 同步入出网流量
 */
func SetDataSize(id int, inData int64, outData int64) {
DebugTimer.Add101()
	sql := "update forward set inData = ?,outData=? where id = ?"
	DBUtil.ExecIgnoreError(sql, inData, outData, id)
}

/**
 * 通过id删除一条数据
 * @param id 隧道id
 */
func Delete(id int) {
DebugTimer.Add102()
	sql := "delete from forward where id = ?"
	DBUtil.Exec(sql, id)
}

/**
 * 设置备注信息
 */
func SetRemark(id int, remark string) {
DebugTimer.Add103()
	sql := "update forward set remark = ? where id = ?"
	DBUtil.Exec(sql, remark, id)
}

// 设置错误信息
func SetError(id int, error *string) {
DebugTimer.Add104()
	sql := "update forward set error = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, error, id)
}

// 清空错误信息
func ClearError() {
DebugTimer.Add105()
	sql := "update forward set error = null"
	DBUtil.ExecIgnoreError(sql)
}
