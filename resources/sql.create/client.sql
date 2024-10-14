-- 客户端表
CREATE TABLE client
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name          VARCHAR(32) NOT NULL,
    version       VARCHAR(10),                                                 -- 客户端版本号
    key VARCHAR (32) NOT NULL UNIQUE,                                          -- 客户端验证秘钥
    ip            VARCHAR(128),                                                -- 客户端ip地址
    inData        BIGINT      NOT NULL DEFAULT 0,                              -- 入网流量
    outData       BIGINT      NOT NULL DEFAULT 0,                              -- 出网流量
    onlineState   INTEGER     NOT NULL DEFAULT 0,                              -- 在线状态
    enableState   INTEGER     NOT NULL DEFAULT 1,                              -- 启用状态 1:开启  0:停止
    lastLoginDate BIGINT      NOT NULL DEFAULT (strftime('%s', 'now') * 1000), -- 最后一次连接时间
    Date          BIGINT      NOT NULL DEFAULT (strftime('%s', 'now') * 1000), -- 创建时间
    remark        TEXT                                                         -- 一些备注信息,错误信息等
);
