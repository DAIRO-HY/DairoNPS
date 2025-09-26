package DBUtil

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/util/LogUtil"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// DB_PATH 文件路径
const DB_PATH = "./data/dairo-nps.sqlite"

// 执行sql语句,忽略错误
func ExecIgnoreError(query string, args ...any) int64 {
DebugTimer.Add397()
	count, err := Exec(query, args...)
	if err != nil {
DebugTimer.Add398()
		log.Printf("%q: %s\n", err, query)
		return -1
	}
	return count
}

// 执行sql
func Exec(query string, args ...any) (int64, error) {
DebugTimer.Add399()
	db := GetDb()
	defer db.Close()
	rs, err := db.Exec(query, args...)
	if err != nil {
DebugTimer.Add400()
		return -1, err
	}
	count, err := rs.RowsAffected()
	if err != nil {
DebugTimer.Add401()
		return -1, err
	}
	return count, nil
}

// 添加数据,忽略错误
func InsertIgnoreError(query string, args ...any) int64 {
DebugTimer.Add402()
	count, err := Insert(query, args...)
	if err != nil {
DebugTimer.Add403()
		LogUtil.Error(fmt.Sprintf("添加数据失败:%s  err:%q\n", query, err))
		return -1
	}
	return count
}

// 添加数据,并返回最后一次添加的ID
func Insert(insert string, args ...any) (int64, error) {
DebugTimer.Add404()
	db := GetDb()
	defer db.Close()
	rs, err := db.Exec(insert, args...)
	if err != nil {
DebugTimer.Add405()
		return -1, err
	}
	lastInsertId, err := rs.LastInsertId()
	if err != nil {
DebugTimer.Add406()
		return -1, err
	}
	return lastInsertId, nil
}

// SelectSingleOneIgnoreError 查询第一个数据并忽略错误
func SelectSingleOneIgnoreError[T any](query string, args ...any) T {
DebugTimer.Add407()
	value, _ := SelectSingleOne[T](query, args...)
	return value
}

// SelectSingleOne 查询第一个数据
func SelectSingleOne[T any](query string, args ...any) (T, error) {
DebugTimer.Add408()
	db := GetDb()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
DebugTimer.Add409()
		return *new(T), err // 返回默认值和错误
	}
	defer rows.Close()

	if !rows.Next() {
DebugTimer.Add410()
		return *new(T), sql.ErrNoRows // 如果没有结果，返回默认值和 ErrNoRows
	}

	var value T
	err = rows.Scan(&value) // 使用 Scan 将结果赋值给 value
	if err != nil {
DebugTimer.Add411()
		return *new(T), err // 返回默认值和错误
	}

	return value, nil
}

// SelectOne 查询第一个数据
func SelectOne[T any](query string, args ...any) *T {
DebugTimer.Add412()
	dtoList := SelectList[T](query, args...)
	if len(dtoList) == 0 {
DebugTimer.Add413()
		return nil
	}
	return dtoList[0]
}

// SelectList 查询列表
func SelectList[T any](query string, args ...any) []*T {
DebugTimer.Add414()
	list := SelectToListMap(query, args...)

	// 创建一个空切片
	dtoList := make([]*T, 0) // 初始化空切片
	for _, item := range list {
DebugTimer.Add415()
		dtoT := new(T)
		reflectDto := reflect.ValueOf(dtoT).Elem()
		for key := range item {
DebugTimer.Add416()

			//将首字符大写
			field := strings.ToUpper(string(key[0])) + key[1:]
			nameField := reflectDto.FieldByName(field)
			value := item[key]
			kind := nameField.Kind()

			//判断数据类型
			switch kind {
			case reflect.String:
				nameField.SetString(value)
			case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int16, reflect.Int8:
				int64Value, _ := strconv.ParseInt(value, 10, 64)
				nameField.SetInt(int64Value)
			default:
				fmt.Println("未知类型")
			}
		}
		dtoList = append(dtoList, dtoT)
	}
	return dtoList
}

// SelectToListMap 将查询结果以List<Map>的类型返回
func SelectToListMap(query string, args ...any) []map[string]string {
DebugTimer.Add417()
	db := GetDb()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
DebugTimer.Add418()
		LogUtil.Error(fmt.Sprintf("查询数据失败:%s: err:%q", query, err))
		return nil
	}
	defer rows.Close()

	// 获取列的名称
	columns, err := rows.Columns()
	if err != nil {
DebugTimer.Add419()
		log.Printf("%q: %s\n", err, query)
		return nil
	}

	// 创建一个长度与列数相同的slice来存放查询结果
	values := make([]interface{}, len(columns))

	// 创建一个[]interface{}的slice, 每个元素指向values中的对应位置
	valuePtrs := make([]interface{}, len(columns))

	// 创建一个空切片
	list := make([]map[string]string, 0) // 初始化空切片
	for rows.Next() {
DebugTimer.Add420()
		for i := range values {
DebugTimer.Add421()
			valuePtrs[i] = &values[i]
		}

		// 将当前行的数据扫描到valuePtrs中
		if err := rows.Scan(valuePtrs...); err != nil {
DebugTimer.Add422()
			LogUtil.Error(fmt.Sprintf("数据扫描失败:%s: err:%q", query, err))
			return nil
		}

		// 使用map将列名和对应的值关联起来
		rowMap := make(map[string]string)
		for i, col := range columns {
DebugTimer.Add423()
			value := values[i]
			if value == nil {
DebugTimer.Add424()
				continue
			}
			rowMap[col] = fmt.Sprintf("%v", value)
		}
		list = append(list, rowMap)
	}
	return list
}

func GetDb() *sql.DB {
DebugTimer.Add425()
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
DebugTimer.Add426()
		LogUtil.Error(fmt.Sprintf("打开数据库失败 err:%q", err))
		log.Fatal(err)
		return nil
	}
	return db
}
