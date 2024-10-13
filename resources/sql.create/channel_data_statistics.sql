-- 隧道流量统计表
CREATE TABLE channel_data_statistics
(
    channelId INTEGER NOT NULL,           -- 隧道id
    ymdhm    BIGINT NOT NULL,           -- 年月日时分
    inData   BIGINT  NOT NULL DEFAULT 0, -- 入网流量
    outData  BIGINT  NOT NULL DEFAULT 0 -- 出网流量
);
CREATE INDEX index_channelId ON channel_data_statistics (channelId);