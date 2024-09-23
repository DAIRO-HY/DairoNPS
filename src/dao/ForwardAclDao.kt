package cn.dairo.cls.dao

import cn.dairo.cls.dao.dto.ForwardAclDto

/**
 * 代理服务访问控制IP数据操作
 */
object ForwardAclDao {

    /**
     * 添加访问控制的ip列表
     */
    fun add(forwardId: Int, dtoList: List<ForwardAclDto>) {
        NPSDB.db.exec("delete from forward_acl where forwardId=?", forwardId)
        dtoList.forEach {
            val sql =
                "insert into channel_acl(forwardId,ip,remark)values(?,?,?);"
            NPSDB.db.exec(
                sql,
                forwardId,
                it.ip,
                it.remark
            )
        }
    }

    /**
     * 获取所有数据
     * @return 代理服务Dto
     */
    fun selectByForwardId(forwardId: Int): List<ForwardAclDto> {
        val sql = "select ip,remark from forward_acl where forwardId = ?"
        val list = NPSDB.db.selectList(ForwardAclDto::class.java, sql, forwardId)
        return list
    }

}