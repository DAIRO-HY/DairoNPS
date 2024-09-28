-- 代理服务黑白名单IP配置
CREATE TABLE forward_acl
(
    forwardId INTEGER     NOT NULL, -- 代理服务id
    ip         VARCHAR(15) NOT NULL, -- ip地址
    remark     VARCHAR(32)           -- 备注
);
CREATE INDEX index_forward_id ON forward_acl (forwardId);