
-- 客户端表
CREATE TABLE client
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name            VARCHAR(32)                       NOT NULL,
    version         VARCHAR(10),                                                          -- 客户端版本号
    key             VARCHAR(32)                       NOT NULL UNIQUE,                    -- 客户端验证秘钥
    ip              VARCHAR(128),                                                         -- 客户端ip地址
    in_data_total   BIGINT                            NOT NULL DEFAULT 0,                 -- 入网流量
    out_data_total  BIGINT                            NOT NULL DEFAULT 0,                 -- 出网流量
    online_state    INTEGER                               NOT NULL DEFAULT 0,                 -- 在线状态
    enable_state    INTEGER                               NOT NULL DEFAULT 1,                 -- 启用状态 1:开启  0:停止
    last_login_date DATETIME                          NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 最后一次连接时间
    create_date     DATETIME                          NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    update_date     BIGINT                          NOT NULL, -- 最后一次更新时间
    remark          TEXT                                                                  -- 一些备注信息,错误信息等
);
CREATE INDEX client_index_update_date ON client (update_date);