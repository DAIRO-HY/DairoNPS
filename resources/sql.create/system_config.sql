-- 系统配置
CREATE TABLE system_config
(
    in_data_total  BIGINT NOT NULL DEFAULT 0, -- 入网流量
    out_data_total BIGINT NOT NULL DEFAULT 0  -- 出网流量
);