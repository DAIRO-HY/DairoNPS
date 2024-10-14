package ForwardDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

//数据转发操作

/**
 * 添加一条隧道
 */
func Add(dto dto.ForwardDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql :=
		"insert into forward(name,port,targetPort,aclState,enableState,updateDate)values(?,?,?,?,?,?)"
	id := DBUtil.InsertIgnoreError(
		sql,
		dto.Name,
		dto.Port,
		dto.TargetPort,
		dto.AclState,
		dto.EnableState,
		updateDate,
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
func selectAll() []*dto.ForwardDto {
	sql := "select * from forward"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectActive() []*dto.ForwardDto {
	sql := "select * from forward where enableState := 1"
	return DBUtil.SelectList[dto.ForwardDto](sql)
}

/**
 * 更新一条数据
 */
func Update(dto dto.ForwardDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql :=
		"update forward set name = ?,port=?,targetPort=?,enableState=?,aclState=?,remark=?,updateDate=? where id = ? and updateDate=?"
	DBUtil.Exec(
		sql,
		dto.Name,
		dto.Port,
		dto.TargetPort,
		dto.EnableState,
		dto.AclState,
		dto.Remark,
		updateDate,
		dto.Id,
		dto.UpdateDate,
	)
}

/**
 * 同步入出网流量
 */
func SetDataLen(dto dto.ForwardDto) {
	sql := "update forward set inData = ?,outDataTotal=? where id = ?"
	DBUtil.Exec(sql, dto.InDataTotal, dto.OutDataTotal, dto.Id)
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
