// 隧道数据操作
package ChannelDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

/**
 * 添加一条隧道
 */
func add(dto dto.ChannelDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql :=
		"insert into channel(client_id,name,type,serverPort,targetPort,securityState,aclState,enableState,updateDate)values(?,?,?,?,?,?,?,?,?)"
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
		"update channel set name = ?,type = ?,server_port=?,target_port=?,enable_state=?,security_state=?,acl_state=?,remark=?,update_date=${System.currentTimeMillis()} where id = ? and update_date=?"
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
func SetDataLen(dto *dto.ChannelDto) {
	sql := "update channel set in_data_total = ?,out_data_total=? where id = ?"
	DBUtil.Exec(sql, dto.InDataTotal, dto.OutDataTotal, dto.Id)
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
	sql := "delete from channel where client_id = ?"
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
//func Search(dto dto.ChannelListSearchDto) []*dto.ChannelSearchDto {
//    sql := "select channel.*,client.name as clientName" +
//                " from channel left join client on channel.client_id = client.id where 1=1 "
//
//    if (dto.ClientId != nil) {
//        sql.append(" and channel.client_id = ").append(dto.clientId)
//    }
//
//    if (dto.type != null) {
//        sql.append(" and channel.type = ").append(dto.type)
//    }
//    sql.append(" order by id desc")
//    return DBUtil.selectList(ChannelSearchDto::class.java, sql.toString())
//}

/**
 * 获取所有激活的隧道列表
 */
func SelectActiveByClientId(clientId int) []int {
	sql := "select channel.* from channel left join client on channel.clientId = client.id where channel.clientId = ? and client.enableState = 1 and channel.enableState = 1"
	return DBUtil.SelectList[dto.ChannelDto](sql, clientId)
}

/**
 * 获取客户端下所有的隧道id列表
 */
func SelectIdByClientId(clientId int) []*int {
	sql := "select id from channel where client_id = ?"
	return DBUtil.SelectList[int](sql, clientId)
}
