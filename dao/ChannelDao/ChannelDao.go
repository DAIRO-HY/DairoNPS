// 隧道数据操作
package ChannelDao

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"strconv"
)

func init() {
DebugTimer.Add56()
	ClearError()
}

/**
 * 添加一条隧道
 */
func Add(dto *dto.ChannelDto) {
DebugTimer.Add57()
	sql :=
		"insert into channel(clientId,name,mode,serverPort,targetPort,securityState,aclState,remark)values(?,?,?,?,?,?,?,?)"
	id := DBUtil.InsertIgnoreError(
		sql,
		dto.ClientId,
		dto.Name,
		dto.Mode,
		dto.ServerPort,
		dto.TargetPort,
		dto.SecurityState,
		dto.AclState,
		dto.Remark,
	)
	dto.Id = int(id)
}

/**
 * 通过id获取一条数据
 * @param id 隧道id
 * @return 隧道Dto
 */
func SelectOne(id int) *dto.ChannelDto {
DebugTimer.Add58()
	sql := "select * from channel where id = ?"
	return DBUtil.SelectOne[dto.ChannelDto](sql, id)
}

// 通过端口查询一条数据
func SelectByPort(port int) *dto.ChannelDto {
DebugTimer.Add59()
	sql := "select * from channel where serverPort = ?"
	return DBUtil.SelectOne[dto.ChannelDto](sql, port)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectAll() []*dto.ChannelDto {
DebugTimer.Add60()
	sql := "select * from channel"
	return DBUtil.SelectList[dto.ChannelDto](sql)
}

/**
 * 更新一条数据
 */
func Update(dto *dto.ChannelDto) {
DebugTimer.Add61()
	sql :=
		"update channel set name = ?,mode = ?,serverPort=?,targetPort=?,securityState=?,aclState=?,remark=? where id = ?"
	DBUtil.ExecIgnoreError(
		sql,
		dto.Name,
		dto.Mode,
		dto.ServerPort,
		dto.TargetPort,
		dto.SecurityState,
		dto.AclState,
		dto.Remark,
		dto.Id,
	)
}

/**
 * 同步入出网流量
 */
func SetDataSize(id int, inData int64, outData int64) {
DebugTimer.Add62()
	sql := "update channel set inData = ?,outData=? where id = ?"
	DBUtil.ExecIgnoreError(sql, inData, outData, id)
}

/**
 * @TODO: 删除数据流量统计信息
 * 通过id删除一条数据
 * @param id 隧道id
 */
func Delete(id int) {
DebugTimer.Add63()
	sql := "delete from channel where id = ?"
	DBUtil.ExecIgnoreError(sql, id)
}

/**
 * 删除某个客户端下所有的隧道
 * @param clientId 客户端ID
 */
func DeleteByClient(clientId int) {
DebugTimer.Add64()
	sql := "delete from channel where clientId = ?"
	DBUtil.ExecIgnoreError(sql, clientId)
}

/**
 * 设置备注信息
 */
func SetRemark(id int, remark string) {
DebugTimer.Add65()
	sql := "update channel set remark = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, remark, id)
}

/**
 * 获取所有隧道列表
 */
func Search(searchDto dto.ChannelListSearchDto) []*dto.ChannelSearchDto {
DebugTimer.Add66()
	sql := "select channel.*,client.name as clientName" +
		" from channel left join client on channel.client_id = client.id where 1=1 "

	if searchDto.ClientId != 0 {
DebugTimer.Add67()
		sql += " and channel.client_id = " + strconv.Itoa(searchDto.ClientId)
	}

	if searchDto.Mode != 0 {
DebugTimer.Add68()
		sql += " and channel.mode = " + strconv.Itoa(searchDto.Mode)
	}
	sql += " order by id desc"
	return DBUtil.SelectList[dto.ChannelSearchDto](sql)
}

/**
 * 获取所有激活的隧道列表
 */
func SelectActiveByClientId(clientId int) []*dto.ChannelDto {
DebugTimer.Add69()
	sql := "select channel.* from channel left join client on channel.clientId = client.id where channel.clientId = ? and client.enableState = 1 and channel.enableState = 1"
	return DBUtil.SelectList[dto.ChannelDto](sql, clientId)
}

/**
 * 获取客户端下所有的隧道id列表
 */
func SelectByClientId(clientId int) []*dto.ChannelDto {
DebugTimer.Add70()
	sql := "select * from channel where clientId = ?"
	return DBUtil.SelectList[dto.ChannelDto](sql, clientId)
}

/**
 * 获取客户端下所有的隧道id列表
 */
func SelectIdByClientId(clientId int) []int {
DebugTimer.Add71()
	sql := "select id from channel where clientId = ?"
	list := DBUtil.SelectList[dto.ChannelDto](sql, clientId)
	ids := make([]int, 0)
	for _, it := range list {
DebugTimer.Add72()
		ids = append(ids, it.Id)
	}
	return ids
}

// 设置可用状态
func SetEnableState(id int, state int) {
DebugTimer.Add73()
	sql := "update channel set enableState = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, state, id)
}

// 设置错误信息
func SetError(id int, error *string) {
DebugTimer.Add74()
	sql := "update channel set error = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, error, id)
}

// 清空错误信息
func ClearError() {
DebugTimer.Add75()
	sql := "update channel set error = null"
	DBUtil.ExecIgnoreError(sql)
}
