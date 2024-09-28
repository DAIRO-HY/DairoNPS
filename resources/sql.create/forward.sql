
-- 代理服务器
CREATE TABLE forward
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name            VARCHAR(32)                       NOT NULL, -- 转发名
    port    INTEGER     NOT NULL UNIQUE,                    -- 服务器端端口
    targetPort    VARCHAR(32) NOT NULL,                           -- 目标端口(ip:端口)
    aclState      INTEGER     NOT NULL DEFAULT 0,                 -- 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
    inDataTotal   BIGINT                            NOT NULL DEFAULT 0,                 -- 入网流量
    outDataTotal  BIGINT                            NOT NULL DEFAULT 0,                 -- 出网流量
    enableState    INTEGER                               NOT NULL DEFAULT 1,                 -- 启用状态 1:开启  0:停止
    createDate     DATETIME                          NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updateDate     BIGINT                          NOT NULL, -- 最后一次更新时间
    remark          TEXT                                                                  -- 一些备注信息,错误信息等
);
CREATE INDEX forward_index_update_date ON forward (updateDate);