/// 客户端信息
pub struct ClientDto {
    /// id
    pub id: i64,

    /// 名称
    pub name: String,

    /// 版本号
    pub version: String,

    /// 连接认证秘钥
    pub key: String,

    /// ip地址
    pub ip: String,

    /// 入网流量
    pub in_data_total: i64,

    /// 出网流量
    pub out_data_total: i64,

    /// 在线状态,0:离线 1:在线
    pub online_state: i64,

    /// 启用状态
    pub enable_state: i64,

    /// 最后一次连接时间
    pub last_login_date: i64,

    /// 创建时间
    pub create_date: i64,

    /// 最后一次更新时间戳
    pub update_date: i64,

    /// 一些备注信息,错误信息等
    pub remark: String,
}

