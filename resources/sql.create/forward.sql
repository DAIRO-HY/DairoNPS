-- 代理服务器
CREATE TABLE forward
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name        VARCHAR(32) NOT NULL,                                        -- 转发名
    port        INTEGER     NOT NULL UNIQUE,                                 -- 服务器端端口
    targetPort  VARCHAR(32) NOT NULL,                                        -- 目标端口(ip:端口)
    aclState    INTEGER     NOT NULL DEFAULT 0,                              -- 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
    inData      BIGINT      NOT NULL DEFAULT 0,                              -- 入网流量
    outData     BIGINT      NOT NULL DEFAULT 0,                              -- 出网流量
    enableState INTEGER     NOT NULL DEFAULT 1,                              -- 启用状态 1:开启  0:停止
    date        BIGINT      NOT NULL DEFAULT (strftime('%s', 'now') * 1000), -- 创建时间
    remark      TEXT,                                                        -- 一些备注信息,错误信息等
    error       TEXT                                                         -- 错误信息

);
