-- 隧道表
CREATE TABLE channel
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    client_id      INTEGER     NOT NULL,                           -- 客户端id
    name           VARCHAR(32) NOT NULL,
    type           INTEGER     NOT NULL,                           -- 隧道模式, 1:TCP  2:UDP
    server_port    INTEGER     NOT NULL UNIQUE,                    -- 服务器端端口
    target_port    VARCHAR(32) NOT NULL,                           -- 目标端口(ip:端口)
    in_data_total  BIGINT      NOT NULL DEFAULT 0,                 -- 入网流量
    out_data_total BIGINT      NOT NULL DEFAULT 0,                 -- 出网流量
    enable_state   INTEGER     NOT NULL DEFAULT 1,                 -- 启用状态 1:开启  0:停止
    security_state INTEGER     NOT NULL DEFAULT 0,                 -- 是否加密传输
    acl_state      INTEGER     NOT NULL DEFAULT 0,                 -- 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
    create_date    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    update_date    BIGINT      NOT NULL,                           -- 最后一次更新时间
    remark         TEXT                                            -- 一些备注信息,错误信息等
);
CREATE INDEX channel_index_update_date ON channel (update_date);