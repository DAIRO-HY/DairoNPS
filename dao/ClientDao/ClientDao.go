package ClientDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

// Add 添加一条客户端数据
func Add(dto *dto.ClientDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql := "insert into client(name,key,remark,updateDate)values(?,?,?,?)"
	id := DBUtil.InsertIgnoreError(sql, dto.Name, dto.Key, dto.Remark, updateDate)
	dto.Id = int(id)
}

/**
* 通过客户端id获取一条数据
* @param id 客户端id
* @return 客户端Dto
 */
func SelectOne(id int) *dto.ClientDto {
	sql := "select * from client where id = ?"
	return DBUtil.SelectOne[dto.ClientDto](sql, id)
}

/**
 * 通过认证秘钥获取一条数据
 * @param key 认证秘钥
 * @return 客户端Dto
 */
func SelectByKey(key string) *dto.ClientDto {
	sql := "select * from client where key = ?"
	return DBUtil.SelectOne[dto.ClientDto](sql, key)
}

/**
 * 更新一条数据
 */
func Update(dto dto.ClientDto) {
	updateDate := time.Now().UnixNano() / int64(time.Millisecond)
	sql :=
		"update client set name = ?,key = ?,enableState=?,remark=?,updateDate=? where id = ? and updateDate=?"
	DBUtil.Exec(sql, dto.Name, dto.Key, dto.EnableState, dto.Remark, updateDate, dto.Id, dto.UpdateDate)
}

/**
 * 同步入出网流量
 */
func SetDataLen(dto *dto.ClientDto) {
	sql := "update client set inDataTotal = ?,outDataTotal=? where id = ?"
	DBUtil.Exec(sql, dto.InDataTotal, dto.OutDataTotal, dto.Id)
}

/**
 * 设置客户端ip地址信息
 */
func SetClientInfo(dto dto.ClientDto) {
	sql :=
		"update client set ip = ?,version=?,lastLoginDate=CURRENT_TIMESTAMP where id = ?"
	DBUtil.Exec(sql, dto.Ip, dto.Version, dto.Id)
}

/**
 * 通过客户端id删除一条数据
 * @param id 客户端id
 */
func delete(id int) {
	sql := "delete from client where id = ?"
	DBUtil.Exec(sql, id)
}

/**
 * 设置备注信息
 */
func setRemark(id int, remark string) {
	sql :=
		"update client set remark = ? where id = ?"
	DBUtil.Exec(sql, remark, id)
}

/**
 * 获取所有客户端列表
 */
func SelectAll() []*dto.ClientDto {
	query := "select * from client order by id desc"
	return DBUtil.SelectList[dto.ClientDto](query)
}
