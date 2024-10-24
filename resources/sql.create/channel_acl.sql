-- 隧道表黑白名单IP配置
CREATE TABLE channel_acl
(
    channelId INTEGER     NOT NULL, -- 隧道id
    ip         VARCHAR(15) NOT NULL, -- ip地址
    remark     VARCHAR(32)           -- 备注
);
CREATE INDEX index_channel_id ON channel_acl (channelId);