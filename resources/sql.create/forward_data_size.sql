-- 端口转发流量统计表
CREATE TABLE forward_data_size
(
    forwardId INTEGER NOT NULL,           -- 端口转发id
    date    BIGINT NOT NULL,           -- 统计时间（年月日时分秒）
    inData   BIGINT  NOT NULL DEFAULT 0, -- 入网流量
    outData  BIGINT  NOT NULL DEFAULT 0 -- 出网流量
);
CREATE INDEX index_forwardId ON forward_data_size (forwardId);