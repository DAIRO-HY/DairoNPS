package cn.dairo.cls.dao.dto

import java.util.*

class ChannelSearchDto {
    var id: Int? = null

    /**
     * 客户端id
     */
    var clientId: Int? = null

    /**
     * 隧道名
     */
    var name: String? = null

    /**
     * 隧道模式, 1:TCP  2:UDP
     */
    var type: Int? = null

    /**
     * 服务端端口
     */
    var serverPort: Int? = null

    /**
     * 目标端口(ip:端口)
     */
    var targetPort: String? = null

    /**
     * 入网流量
     */
    var inDataTotal: Long? = null

    /**
     * 出网流量
     */
    var outDataTotal: Long? = null

    /**
     * 启用状态 1:开启  0:停止
     */
    var enableState: Int? = null

    /**
     * 是否加密传输
     */
    var securityState: Int? = null

    /**
     * 创建时间
     */
    var createDate: Date? = null

    /**
     * 最后一次更新时间
     */
    var updateDate: Long? = null

    /**
     * 客户端名
     */
    var clientName: String? = null
}

