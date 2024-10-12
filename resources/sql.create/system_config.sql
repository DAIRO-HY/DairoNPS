-- 系统配置
CREATE TABLE system_config
(
    inDataTotal  BIGINT NOT NULL DEFAULT 0, -- 入网流量
    outDataTotal BIGINT NOT NULL DEFAULT 0  -- 出网流量
);