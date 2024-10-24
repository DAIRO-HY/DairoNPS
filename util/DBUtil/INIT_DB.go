package DBUtil

import (
	"DairoNPS/resources"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// VERSION 数据库版本号
const VERSION = 1

func init() {
	_, err := os.Stat(DB_PATH)
	// 如果错误是os.ErrNotExist，表示文件不存在
	if os.IsNotExist(err) { //文件不存在

		// 创建多层目录
		err := os.MkdirAll(filepath.Dir(DB_PATH), 0700)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
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
	version := SelectSingleOneIgnoreError[int]("PRAGMA USER_VERSION")
	if version == 0 {
		create()

		//第一次创建数据库时往系统配置表插入一条数据
		ExecIgnoreError("insert into system_config(inData, outData) values (0, 0);")
	}

	//设置数据库版本号
	ExecIgnoreError("PRAGMA USER_VERSION = " + strconv.Itoa(VERSION))
}

func create() {
	sqlFiles := []string{"forward.sql", "forward_acl.sql", "client.sql", "channel.sql", "channel_acl.sql", "system_config.sql", "date_data_size.sql"}
	for _, fn := range sqlFiles {
		createSql, _ := resources.StaticFiles.ReadFile("sql.create/" + fn)
		ExecIgnoreError(string(createSql))
	}
}
