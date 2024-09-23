package cn.dairo.cls.dao

import cn.dairo.cls.dao.dto.ChannelAclDto

/**
 * 隧道访问控制IP数据操作
 */
object ChannelAclDao {

    /**
     * 添加访问控制的ip列表
     */
    fun add(channelId:Int, dtoList: List<ChannelAclDto>) {
        NPSDB.db.exec("delete from channel_acl where channelId=?",channelId)
        dtoList.forEach {
            val sql =
                "insert into channel_acl(channelId,ip,remark)values(?,?,?);"
            NPSDB.db.exec(
                sql,
                channelId,
                it.ip,
                it.remark
            )
        }
    }

    /**
     * 获取所有数据
     * @return 隧道Dto
     */
    fun selectByChannelId(channelId:Int): List<ChannelAclDto> {
        val sql = "select ip,remark from channel_acl where channelId = ?"
        val list = NPSDB.db.selectList(ChannelAclDto::class.java, sql, channelId)
        return list
    }

}