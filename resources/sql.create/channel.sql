-- 隧道表
CREATE TABLE channel
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    clientId      INTEGER     NOT NULL,                                        -- 客户端id
    name          VARCHAR(32) NOT NULL,
    mode          INTEGER     NOT NULL,                                        -- 隧道模式, 1:TCP  2:UDP
    serverPort    INTEGER     NOT NULL UNIQUE,                                 -- 服务器端端口
    targetPort    VARCHAR(32) NOT NULL,                                        -- 目标端口(ip:端口)
    inData        BIGINT      NOT NULL DEFAULT 0,                              -- 入网流量
    outData       BIGINT      NOT NULL DEFAULT 0,                              -- 出网流量
    enableState   INTEGER     NOT NULL DEFAULT 1,                              -- 启用状态 1:开启  0:停止
    securityState INTEGER     NOT NULL DEFAULT 0,                              -- 是否加密传输
    aclState      INTEGER     NOT NULL DEFAULT 0,                              -- 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
    date          BIGINT      NOT NULL DEFAULT (strftime('%s', 'now') * 1000), -- 创建时间
    remark        TEXT,                                                        -- 一些备注信息,错误信息等
    error         TEXT                                                         -- 错误信息
);