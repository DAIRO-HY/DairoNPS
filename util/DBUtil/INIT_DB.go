package DBUtil

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/resources"
	"DairoNPS/util/LogUtil"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// VERSION 数据库版本号
const VERSION = 2

func init() {
DebugTimer.Add427()
	_, err := os.Stat(DB_PATH)
	// 如果错误是os.ErrNotExist，表示文件不存在
	if os.IsNotExist(err) { //文件不存在
DebugTimer.Add428()

		// 创建多层目录
		err := os.MkdirAll(filepath.Dir(DB_PATH), 0700)
		if err != nil {
DebugTimer.Add429()
			LogUtil.Error(fmt.Sprintf("创建文件夹[%s]失败 err:%q", DB_PATH, err))
			log.Fatal(err)
			return
		}
	}

	// 打开数据库连接，没有文件时会自动创建
	//db, _ := sql.Open("sqlite3", DB_PATH)
	upgrade()
}

/**
* 更新表结构
 */
func upgrade() {
DebugTimer.Add430()
	version := SelectSingleOneIgnoreError[int]("PRAGMA USER_VERSION")
	if version == 0 {
DebugTimer.Add431()
		create()

		//第一次创建数据库时往系统配置表插入一条数据
		ExecIgnoreError("insert into system_config(inData, outData) values (0, 0);")
	}
	if version > 0 {
DebugTimer.Add432()
		if version < 2 { //添加error字段
DebugTimer.Add433()
			ExecIgnoreError("alter table channel add error TEXT;")
			ExecIgnoreError("alter table forward add error TEXT;")

			ExecIgnoreError("drop table channel_acl;")
			ExecIgnoreError("drop table forward_acl;")
		}
	}

	//设置数据库版本号
	ExecIgnoreError("PRAGMA USER_VERSION = " + strconv.Itoa(VERSION))
}

func create() {
DebugTimer.Add434()
	sqlFiles := []string{"forward.sql", "client.sql", "channel.sql", "system_config.sql", "date_data_size.sql"}
	for _, fn := range sqlFiles {
DebugTimer.Add435()
		createSql, _ := resources.SqlFolder.ReadFile("sql.create/" + fn)
		ExecIgnoreError(string(createSql))
	}
}
