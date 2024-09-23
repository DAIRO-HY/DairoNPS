
/**
 * 隧道数据操作
 */
object ChannelDao {

    /**
     * 查询所有key
     */
    const val SELECT_KEY = """
     channel.id as id
    ,channel.client_id as clientId
    ,channel.name as name
    ,channel.type as type
    ,channel.server_port as serverPort
    ,channel.target_port as targetPort
    ,channel.in_data_total as inDataTotal
    ,channel.out_data_total as outDataTotal
    ,channel.enable_state as enableState
    ,channel.security_state as securityState
    ,channel.acl_state as aclState
    ,channel.create_date as createDate
    ,channel.update_date as updateDate
    ,channel.remark as remark
    """



    /**
     * 添加一条隧道
     */
    fun add(dto: ChannelDto) {
        val updateDate = System.currentTimeMillis()
        val sql =
            "insert into channel(client_id,name,type,server_port,target_port,security_state,acl_state,enable_state,update_date)values(?,?,?,?,?,?,?,?,?)"
        NPSDB.db.exec(
            sql,
            dto.clientId,
            dto.name,
            dto.type,
            dto.serverPort,
            dto.targetPort,
            dto.securityState,
            dto.aclState,
            dto.enableState,
            updateDate
        )

        val insertIdSql = "select id from channel where update_date = ?"
        val id = NPSDB.db.selectSingleOne<Int>(insertIdSql, updateDate)
        dto.id = id
    }

    /**
     * 通过id获取一条数据
     * @param id 隧道id
     * @return 隧道Dto
     */
    fun selectOne(id: Int): ChannelDto? {
        val sql = "select " + SELECT_KEY + " from channel where id = ?"
        val dto = NPSDB.db.selectOne(ChannelDto::class.java, sql, id)
        return dto
    }

    /**
     * 获取所有数据
     * @return 隧道Dto
     */
    fun selectAll(): List<ChannelDto> {
        val sql = "select " + SELECT_KEY + " from channel"
        val list = NPSDB.db.selectList(ChannelDto::class.java, sql)
        return list
    }

    /**
     * 更新一条数据
     */
    fun update(dto: ChannelDto) {
        val sql =
            "update channel set name = ?,type = ?,server_port=?,target_port=?,enable_state=?,security_state=?,acl_state=?,remark=?,update_date=${System.currentTimeMillis()} where id = ? and update_date=?"
        NPSDB.db.exec(
            sql,
            dto.name,
            dto.type,
            dto.serverPort,
            dto.targetPort,
            dto.enableState,
            dto.securityState,
            dto.aclState,
            dto.remark,
            dto.id,
            dto.updateDate
        )
    }

    /**
     * 同步入出网流量
     */
    fun setDataLen(dto: ChannelDto) {
        val sql = "update channel set in_data_total = ?,out_data_total=? where id = ?"
        NPSDB.db.exec(sql, dto.inDataTotal, dto.outDataTotal, dto.id)
    }

    /**
     * @TODO: 删除数据流量统计信息
     * 通过id删除一条数据
     * @param id 隧道id
     */
    fun delete(id: Int) {
        val sql = "delete from channel where id = ?"
        NPSDB.db.exec(sql, id)
    }

    /**
     * 删除某个客户端下所有的隧道
     * @param clientId 客户端ID
     */
    fun deleteByClient(clientId: Int) {
        val sql = "delete from channel where client_id = ?"
        NPSDB.db.exec(sql, clientId)
    }

    /**
     * 设置备注信息
     */
    fun setRemark(id: Int, remark: String) {
        val sql =
            "update channel set remark = ? where id = ?"
        NPSDB.db.exec(sql, remark, id)
    }

    /**
     * 获取所有隧道列表
     */
    fun search(dto: ChannelListSearchDto): List<ChannelSearchDto> {
        val sql = StringBuilder(
            "select " + SELECT_KEY +
                    ",client.name as clientName" +
                    " from channel left join client on channel.client_id = client.id where 1=1 "
        )

        if (dto.clientId != null) {
            sql.append(" and channel.client_id = ").append(dto.clientId)
        }

        if (dto.type != null) {
            sql.append(" and channel.type = ").append(dto.type)
        }
        sql.append(" order by id desc")
        return NPSDB.db.selectList(ChannelSearchDto::class.java, sql.toString())
    }

    /**
     * 获取所有激活的隧道列表
     */
    fun selectActiveByClientId(clientId: Int): List<ChannelDto> {
        val sql = "select" + SELECT_KEY +
                " from channel left join client on channel.client_id = client.id where channel.client_id = ? and client.enable_state = 1 and channel.enable_state = 1"
        return NPSDB.db.selectList(ChannelDto::class.java, sql, clientId)
    }

    /**
     * 获取客户端下所有的隧道id列表
     */
    fun selectIdByClientId(clientId: Int): List<Int> {
        val sql = "select id from channel where client_id = ?"
        return NPSDB.db.selectList(ChannelDto::class.java, sql, clientId).map {
            it.id!!
        }
    }
}