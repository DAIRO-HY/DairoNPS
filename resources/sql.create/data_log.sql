-- 流量统计表
CREATE TABLE data_log
(
    targetId INTEGER NOT NULL,           -- 目标id
    year     INTEGER NOT NULL,           -- 年
    ym       INTEGER NOT NULL,           -- 年月
    ymd      INTEGER NOT NULL,           -- 年月日
    ymdh     INTEGER NOT NULL,           -- 年月日时
    ymdhm    BIGINT NOT NULL,           -- 年月日时分
    inData   BIGINT  NOT NULL DEFAULT 0, -- 入网流量
    outData  BIGINT  NOT NULL DEFAULT 0, -- 出网流量
    type     INTEGER NOT NULL            -- 统计类型, 1:隧道  2:数据转发
);
CREATE INDEX index_targetId ON data_log (targetId);