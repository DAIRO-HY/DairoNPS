-- 系统配置
CREATE TABLE system_config
(
    inData  BIGINT NOT NULL DEFAULT 0, -- 入网流量
    outData BIGINT NOT NULL DEFAULT 0  -- 出网流量
);