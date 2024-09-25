// /**
//  * 客户端数据操作
//  */
//
//     /**
//      * 添加一条客户端数据
//      */
//     fun add(dto: ClientDto) {
//         val updateDate = System.currentTimeMillis()
//         val sql = "insert into client(name,key,remark,update_date)values(?,?,?,$updateDate)"
//         NPSDB.db.exec(sql, dto.name, dto.key, dto.remark)
//         val insertIdSql = "select id from client where update_date = ?"
//         val id = NPSDB.db.selectSingleOne<Int>(insertIdSql, updateDate)
//         dto.id = id
//     }

use crate::dao::dto::client_dto::ClientDto;

/**
 * 通过客户端id获取一条数据
 * @param id 客户端id
 * @return 客户端Dto
 */
pub fn selectOne(_id: u64) -> ClientDto {
    let dto = ClientDto {
        // id
        id: 12,

        // 名称
        name: String::from("String"),

        // 版本号
        version: String::from("String"),

        // 连接认证秘钥
        key: String::from("String"),

        // ip地址
        ip: String::from("String"),

        // 入网流量
        in_data_total: 12,

        // 出网流量
        out_data_total: 12,

        // 在线状态,0:离线 1:在线
        online_state: 12,

        // 启用状态
        enable_state: 12,

        // 最后一次连接时间
        last_login_date: 12,

        // 创建时间
        create_date: 12,

        // 最后一次更新时间戳
        update_date: 12,

        // 一些备注信息,错误信息等
        remark: String::from("String"),
    };
    return dto;
}

/**
 * 通过认证秘钥获取一条数据
 * @param key 认证秘钥
 * @return 客户端Dto
 */
pub fn selectByKey(_key: &str) -> Option<ClientDto> {
    let dto = ClientDto {
        // id
        id: 12,

        // 名称
        name: String::from("String"),

        // 版本号
        version: String::from("String"),

        // 连接认证秘钥
        key: String::from("String"),

        // ip地址
        ip: String::from("String"),

        // 入网流量
        in_data_total: 12,

        // 出网流量
        out_data_total: 12,

        // 在线状态,0:离线 1:在线
        online_state: 12,

        // 启用状态
        enable_state: 12,

        // 最后一次连接时间
        last_login_date: 12,

        // 创建时间
        create_date: 12,

        // 最后一次更新时间戳
        update_date: 12,

        // 一些备注信息,错误信息等
        remark: String::from("String"),
    };
    return dto;
}
//
// /**
//  * 更新一条数据
//  */
// fun update(dto: ClientDto) {
//     val sql =
//         "update client set name = ?,key = ?,enable_state=?,remark=?,update_date=${System.currentTimeMillis()} where id = ? and update_date=?"
//     NPSDB.db.exec(sql, dto.name, dto.key, dto.enableState, dto.remark, dto.id, dto.updateDate)
// }
//
// /**
//  * 同步入出网流量
//  */
// fun setDataLen(dto: ClientDto) {
//     val sql = "update client set in_data_total = ?,out_data_total=? where id = ?"
//     NPSDB.db.exec(sql, dto.inDataTotal, dto.outDataTotal, dto.id)
// }
//
/**
 * 设置客户端ip地址信息
 */
pub fn setClientInfo(_id: i64,_ip:String,_version:String) {
    // val sql =
    //     "update client set ip = ?,version=?,last_login_date=CURRENT_TIMESTAMP where id = ?"
    // NPSDB.db.exec(sql, dto.ip, dto.version, dto.id)
}
//
// /**
//  * 通过客户端id删除一条数据
//  * @param id 客户端id
//  */
// fun delete(id: Int) {
//     val sql = "delete from client where id = ?"
//     NPSDB.db.exec(sql, id)
// }
//
// /**
//  * 设置备注信息
//  */
// fun setRemark(id: Int, remark: String) {
//     val sql =
//         "update client set remark = ? where id = ?"
//     NPSDB.db.exec(sql, remark, id)
// }
//
// /**
//  * 获取所有客户端列表
//  */
// fun selectAll(): List<ClientDto> {
//     val sql = "select id as id" +
//             " ,name as name" +
//             " ,version as version" +
//             " ,key as key" +
//             " ,ip as ip" +
//             " ,in_data_total as inDataTotal" +
//             " ,out_data_total as outDataTotal" +
//             " ,online_state as onlineState" +
//             " ,enable_state as enableState" +
//             " ,create_date as createDate" +
//             " ,update_date as updateDate" +
//             " from client order by id desc"
//     return NPSDB.db.selectList(ClientDto::class.java, sql)
// }
