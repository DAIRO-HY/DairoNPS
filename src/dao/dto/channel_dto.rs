pub struct ChannelDto {

    ///隧道ID
    pub id: i64,

    /// 客户端id
    pub client_id: i64,

    /// 隧道名
    pub name: String,

    /// 隧道模式, 1:TCP  2:UDP
    pub mode: i8,

    /// 服务端端口
    pub server_port: i16,

    /// 目标端口(ip:端口)
    pub target_port: String,

    /// 入网流量
    pub in_data_total: i64,

    /// 出网流量
    pub out_data_total: i64,

    /// 启用状态 1:开启  0:停止
    pub enable_state: i8,

    /// 是否加密传输
    pub security_state: i8,

    /// 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
    pub acl_state: i8,

    /// 创建时间
    pub create_date: i64,

    /// 最后一次更新时间
    pub update_date: i64,

    /// 一些备注信息,错误信息等
    pub remark: String,
}