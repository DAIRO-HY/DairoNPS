package cn.dairo.cls.dao

import cn.dairo.cl.common.util.sqlite.DBBase
import cn.dairo.cl.common.util.sqlite.SqliteTool
import java.io.File
import java.nio.charset.Charset

object NPSDB {

    /**
     * 数据库版本号
     */
    const val VERSION = 11

    /**
     * db连接工具
     */
    val db: DBBase

    /**
     * 文件路径
     */
    private const val dbPath = "./data/cls.sqlite"

    init {

        //db文件是否已经存在
        val isExists = File(dbPath).exists()
        val db = SqliteTool(dbPath)

        if (!isExists) {

            //创建文件夹
            File(dbPath).parentFile.mkdirs()
            create(db)
        } else {
            val oldVersion: Int = db.selectSingleOne("PRAGMA USER_VERSION")
            upgrade(db, oldVersion, VERSION)
        }

        //设置数据库版本号
        db.exec("PRAGMA USER_VERSION = $VERSION")
        this@NPSDB.db = db
    }

    /**
     * 第一次运行时,创建表
     */
    private fun create(db: DBBase) {
        "sql/create/forward.sql".exec(db)
        "sql/create/forward_acl.sql".exec(db)
        "sql/create/client.sql".exec(db)
        "sql/create/channel.sql".exec(db)
        "sql/create/channel_acl.sql".exec(db)
        "sql/create/system_config.sql".exec(db)
        "sql/create/data_log.sql".exec(db)


        //第一次创建数据库时往系统配置表插入一条数据
        db.exec("insert into system_config(in_data_total, out_data_total) values (0, 0);")
    }

    /**
     * 更新表结构
     */
    private fun upgrade(db: DBBase, oldVersion: Int, newVersion: Int) {
        if (oldVersion < 3) {
            this.create(db)
        }
        if (oldVersion < 5) {
            db.exec("ALTER TABLE channel ADD acl_state INTEGER NOT NULL DEFAULT 0")
        }
        if (oldVersion < 8) {
            this.create(db)
        }
        if (oldVersion < 9) {//修改表名
            db.exec("alter table proxyer rename to forward;")
            db.exec("drop table proxyer_acl;")
            this.create(db)
        }
        if (oldVersion < 10) {
            this.create(db)
        }
        if (oldVersion < 11) {
            db.exec("drop table data_log;")
            this.create(db)
        }
    }

    fun String.exec(db: DBBase) = NPSDB.javaClass.classLoader.getResourceAsStream(this).use {
        val sql = String(it.readAllBytes(), Charset.forName("UTF-8"))
        try {
            db.exec(sql)
        } catch (e: Exception) {
            //e.printStackTrace()
        }
    }
}
