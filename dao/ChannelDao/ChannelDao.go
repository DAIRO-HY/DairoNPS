// 隧道数据操作
package ChannelDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"fmt"
	"time"
)

/**
 * 添加一条隧道
 */
func Add(dto dto.ChannelDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql :=
		"insert into channel(clientId,name,mode,serverPort,targetPort,securityState,aclState,enableState,updateDate)values(?,?,?,?,?,?,?,?,?)"
	id := DBUtil.InsertIgnoreError(
		sql,
		dto.ClientId,
		dto.Name,
		dto.Mode,
		dto.ServerPort,
		dto.TargetPort,
		dto.SecurityState,
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
func SelectOne(id int) *dto.ChannelDto {
	sql := "select * from channel where id = ?"
	return DBUtil.SelectOne[dto.ChannelDto](sql, id)
}

/**
 * 获取所有数据
 * @return 隧道Dto
 */
func SelectAll() []*dto.ChannelDto {
	sql := "select * from channel"
	return DBUtil.SelectList[dto.ChannelDto](sql)
}

/**
 * 更新一条数据
 */
func Update(dto dto.ChannelDto) {
	sql :=
		"update channel set name = ?,type = ?,serverPort=?,targetPort=?,enableState=?,securityState=?,aclState=?,remark=?,updateDate=${System.currentTimeMillis()} where id = ? and updateDate=?"
	DBUtil.Exec(
		sql,
		dto.Name,
		dto.Mode,
		dto.ServerPort,
		dto.TargetPort,
		dto.EnableState,
		dto.SecurityState,
		dto.AclState,
		dto.Remark,
		dto.Id,
		dto.UpdateDate,
	)
}

/**
 * 同步入出网流量
 */
func SetDataSize(id int, inData int64, outData int64) {
	sql := "update channel set inDataTotal = ?,outDataTotal=? where id = ?"
	DBUtil.Exec(sql, inData, outData, id)
}

/**
 * @TODO: 删除数据流量统计信息
 * 通过id删除一条数据
 * @param id 隧道id
 */
func Delete(id int) {
	sql := "delete from channel where id = ?"
	DBUtil.Exec(sql, id)
}

/**
 * 删除某个客户端下所有的隧道
 * @param clientId 客户端ID
 */
func DeleteByClient(clientId int) {
	sql := "delete from channel where clientId = ?"
	DBUtil.Exec(sql, clientId)
}

/**
 * 设置备注信息
 */
func SetRemark(id int, remark string) {
	sql :=
		"update channel set remark = ? where id = ?"
	DBUtil.Exec(sql, remark, id)
}

/**
 * 获取所有隧道列表
 */
func Search(searchDto dto.ChannelListSearchDto) []*dto.ChannelSearchDto {
	sql := "select channel.*,client.name as clientName" +
		" from channel left join client on channel.client_id = client.id where 1=1 "

	if searchDto.ClientId != 0 {
		sql += fmt.Sprintf(" and channel.client_id = %d", searchDto.ClientId)
	}

	if searchDto.Mode != 0 {
		sql += fmt.Sprintf(" and channel.type = %d", searchDto.Mode)
	}
	sql += " order by id desc"
	return DBUtil.SelectList[dto.ChannelSearchDto](sql)
}

/**
 * 获取所有激活的隧道列表
 */
func SelectActiveByClientId(clientId int) []*dto.ChannelDto {
	sql := "select channel.* from channel left join client on channel.clientId = client.id where channel.clientId = ? and client.enableState = 1 and channel.enableState = 1"
	return DBUtil.SelectList[dto.ChannelDto](sql, clientId)
}

/**
 * 获取客户端下所有的隧道id列表
 */
func SelectIdByClientId(clientId int) []int {
	sql := "select id from channel where clientId = ?"
	list := DBUtil.SelectList[dto.ChannelDto](sql, clientId)
	ids := make([]int, 0)
	for _, it := range list {
		ids = append(ids, it.Id)
	}
	return ids
}
