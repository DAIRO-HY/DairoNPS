package cn.dairo.cls.dao.dto

import java.util.Date

/**
 * 代理服务器Dto
 */
open class ForwardDto {

    /**
     * 代理服务ID
     */
    var id: Int? = null

    /**
     * 代理服务名
     */
    var name: String? = null

    /**
     * 服务端端口
     */
    var port: Int? = null

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
     * 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
     */
    var aclState: Int? = null

    /**
     * 创建时间
     */
    var createDate: Date? = null

    /**
     * 最后一次更新时间
     */
    var updateDate: Long? = null

    /**
     * 一些备注信息,错误信息等
     */
    var remark: String? = null
}

