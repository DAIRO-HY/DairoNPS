package cn.dairo.cls.dao

import cn.dairo.cls.dao.dto.ForwardDto

/**
 * 数据转发操作
 */
object ForwardDao {

    /**
     * 添加一条隧道
     */
    fun add(dto: ForwardDto) {
        val updateDate = System.currentTimeMillis()
        val sql =
            "insert into forward(name,port,targetPort,aclState,enableState,updateDate)values(?,?,?,?,?,?)"
        NPSDB.db.exec(
            sql,
            dto.name,
            dto.port,
            dto.targetPort,
            dto.aclState,
            dto.enableState,
            updateDate
        )

        val insertIdSql = "select id from forward where updateDate = ?"
        val id = NPSDB.db.selectSingleOne<Int>(insertIdSql, updateDate)
        dto.id = id
    }

    /**
     * 通过id获取一条数据
     * @param id 隧道id
     * @return 隧道Dto
     */
    fun selectOne(id: Int): ForwardDto? {
        val sql = "select * from forward where id = ?"
        val dto = NPSDB.db.selectOne(ForwardDto::class.java, sql, id)
        return dto
    }

    /**
     * 获取所有数据
     * @return 隧道Dto
     */
    fun selectAll(): List<ForwardDto> {
        val sql = "select * from forward"
        val list = NPSDB.db.selectList(ForwardDto::class.java, sql)
        return list
    }

    /**
     * 获取所有数据
     * @return 隧道Dto
     */
    fun selectActive(): List<ForwardDto> {
        val sql = "select * from forward where enableState = 1"
        val list = NPSDB.db.selectList(ForwardDto::class.java, sql)
        return list
    }



    /**
     * 更新一条数据
     */
    fun update(dto: ForwardDto) {
        val sql =
            "update forward set name = ?,port=?,targetPort=?,enableState=?,aclState=?,remark=?,updateDate=${System.currentTimeMillis()} where id = ? and updateDate=?"
        NPSDB.db.exec(
            sql,
            dto.name,
            dto.port,
            dto.targetPort,
            dto.enableState,
            dto.aclState,
            dto.remark,
            dto.id,
            dto.updateDate
        )
    }

    /**
     * 同步入出网流量
     */
    fun setDataLen(dto: ForwardDto) {
        val sql = "update forward set inDataTotal = ?,outDataTotal=? where id = ?"
        NPSDB.db.exec(sql, dto.inDataTotal, dto.outDataTotal, dto.id)
    }

    /**
     * 通过id删除一条数据
     * @param id 隧道id
     */
    fun delete(id: Int) {
        val sql = "delete from forward where id = ?"
        NPSDB.db.exec(sql, id)
    }

    /**
     * 设置备注信息
     */
    fun setRemark(id: Int, remark: String) {
        val sql =
            "update forward set remark = ? where id = ?"
        NPSDB.db.exec(sql, remark, id)
    }
}