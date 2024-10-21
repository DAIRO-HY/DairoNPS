package ClientDao

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/DBUtil"
	"time"
)

// Add 添加一条客户端数据
func Add(dto *dto.ClientDto) {
	sql := "insert into client(name,key,remark)values(?,?,?)"
	id := DBUtil.InsertIgnoreError(sql, dto.Name, dto.Key, dto.Remark)
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
func Update(dto *dto.ClientDto) {
	sql :=
		"update client set name = ?,key = ?,remark=? where id = ?"
	DBUtil.ExecIgnoreError(sql, dto.Name, dto.Key, dto.Remark, dto.Id)
}

// 统计入出网流量
func SetDataSize(id int, inAdd int64, outAdd int64) {
	sql := "update client set inData = inData + ?,outData = outData + ? where id = ?"
	DBUtil.ExecIgnoreError(sql, inAdd, outAdd, id)
}

/**
 * 设置客户端ip地址信息
 */
func SetClientInfo(dto dto.ClientDto) {
	lastLoginDate := time.Now().UnixMilli()
	sql :=
		"update client set ip = ?, version=?, lastLoginDate = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, dto.Ip, dto.Version, lastLoginDate, dto.Id)
}

/**
 * 通过客户端id删除一条数据
 * @param id 客户端id
 */
func Delete(id int) {
	sql := "delete from client where id = ?"
	DBUtil.ExecIgnoreError(sql, id)
}

/**
 * 设置备注信息
 */
func setRemark(id int, remark string) {
	sql :=
		"update client set remark = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, remark, id)
}

// 设置可用状态
func SetEnableState(id int, state int) {
	sql := "update client set enableState = ? where id = ?"
	DBUtil.ExecIgnoreError(sql, state, id)
}

/**
 * 获取所有客户端列表
 */
func SelectAll() []*dto.ClientDto {
	query := "select * from client order by id desc"
	return DBUtil.SelectList[dto.ClientDto](query)
}
