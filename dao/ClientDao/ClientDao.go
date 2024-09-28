package ClientDao

import (
    "DairoNPS/dao/dto"
    "DairoNPS/util/DBUtil"
    "fmt"
    "log"
    "time"
)

/**
 * 客户端数据操作
 */


    //添加一条客户端数据
    func Add(dto dto.ClientDto) {
        updateDate := time.Now().UnixNano() / int64(time.Millisecond)
        sql := "insert into client(name,key,remark,update_date)values(?,?,?,?)"
        id := DBUtil.InsertIgnoreError(sql, dto.Name, dto.Key, dto.Remark,updateDate)
        dto.Id = int(id)
    }

   /**
    * 通过客户端id获取一条数据
    * @param id 客户端id
    * @return 客户端Dto
    */
   func SelectOne(id int) dto.ClientDto {
       //query := "select id" +
       //    " ,name" +
       //    " ,version" +
       //    " ,key" +
       //    " ,ip" +
       //    " ,in_data_total as inDataTotal" +
       //    " ,out_data_total as outDataTotal" +
       //    " ,online_state as onlineState" +
       //    " ,enable_state as enableState" +
       //    " ,last_login_date as lastLoginDate" +
       //    " ,create_date as createDate" +
       //    " ,update_date as updateDate" +
       //    " ,remark" +
       //    " from client where id = ?"
       query := "select id,name from client where id = ?"
       DBUtil.SelectOne(query,id)
       db := DBUtil.GetDb()
       defer db.Close()
       rows,err := db.Query(query, id)
       defer rows.Close()
       if err != nil{
           log.Fatalf("%q: %s\n", err, query)
           //return nil
       }

       if !rows.Next(){
           fmt.Println("no Data")
       }
       //one := DBUtil.SelectSingleOneIgnoreError[any]("select * from client",id)
       //fmt.Println(one)

       clientDto := dto.ClientDto{}
       //rows.Scan(&clientDto.Id,&clientDto.Name)
       id1 := -1
       name := ""
       rows.Scan(&id1,&name)
       fmt.Println(id1)
       fmt.Println(name)

       //tt,_ := rows.
       //fmt.Println(tt[1])
       return clientDto
   }
//
//    /**
//     * 通过认证秘钥获取一条数据
//     * @param key 认证秘钥
//     * @return 客户端Dto
//     */
//    fun selectByKey(key: String): ClientDto? {
//        val sql = "select id as id" +
//                " ,name as name" +
//                " ,version as version" +
//                " ,key as key" +
//                " ,ip as ip" +
//                " ,in_data_total as inDataTotal" +
//                " ,out_data_total as outDataTotal" +
//                " ,online_state as onlineState" +
//                " ,enable_state as enableState" +
//                " ,last_login_date as lastLoginDate" +
//                " ,create_date as createDate" +
//                " ,update_date as updateDate" +
//                " from client where key = ?"
//        val dto = NPSDB.db.selectOne(ClientDto::class.java, sql, key)
//        return dto
//    }
//
//    /**
//     * 更新一条数据
//     */
//    fun update(dto: ClientDto) {
//        val sql =
//            "update client set name = ?,key = ?,enable_state=?,remark=?,update_date=${System.currentTimeMillis()} where id = ? and update_date=?"
//        NPSDB.db.exec(sql, dto.name, dto.key, dto.enableState, dto.remark, dto.id, dto.updateDate)
//    }
//
//    /**
//     * 同步入出网流量
//     */
//    fun setDataLen(dto: ClientDto) {
//        val sql = "update client set in_data_total = ?,out_data_total=? where id = ?"
//        NPSDB.db.exec(sql, dto.inDataTotal, dto.outDataTotal, dto.id)
//    }
//
//    /**
//     * 设置客户端ip地址信息
//     */
//    fun setClientInfo(dto: ClientDto) {
//        val sql =
//            "update client set ip = ?,version=?,last_login_date=CURRENT_TIMESTAMP where id = ?"
//        NPSDB.db.exec(sql, dto.ip, dto.version, dto.id)
//    }
//
//    /**
//     * 通过客户端id删除一条数据
//     * @param id 客户端id
//     */
//    fun delete(id: Int) {
//        val sql = "delete from client where id = ?"
//        NPSDB.db.exec(sql, id)
//    }
//
//    /**
//     * 设置备注信息
//     */
//    fun setRemark(id: Int, remark: String) {
//        val sql =
//            "update client set remark = ? where id = ?"
//        NPSDB.db.exec(sql, remark, id)
//    }
//
//    /**
//     * 获取所有客户端列表
//     */
//    fun selectAll(): List<ClientDto> {
//        val sql = "select id as id" +
//                " ,name as name" +
//                " ,version as version" +
//                " ,key as key" +
//                " ,ip as ip" +
//                " ,in_data_total as inDataTotal" +
//                " ,out_data_total as outDataTotal" +
//                " ,online_state as onlineState" +
//                " ,enable_state as enableState" +
//                " ,create_date as createDate" +
//                " ,update_date as updateDate" +
//                " from client order by id desc"
//        return NPSDB.db.selectList(ClientDto::class.java, sql)
//    }
//}