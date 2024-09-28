package DBUtil

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const DbPath = "./data/dairo-nps.sqlite"

// 执行sql语句,忽略错误
func ExecIgnoreError(query string, args ...any) int64 {
	count, err := Exec(query, args...)
	if err != nil {
		log.Fatalf("%q: %s\n", err, query)
		return -1
	}
	return count
}

// 执行sql
func Exec(query string, args ...any) (int64, error) {
	db := GetDb()
	defer db.Close()
	rs, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	count, err := rs.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil
}

// 添加数据,忽略错误
func InsertIgnoreError(query string, args ...any) int64 {
	count, err := Insert(query, args...)
	if err != nil {
		log.Fatalf("%q: %s\n", err, query)
		return -1
	}
	return count
}

// 添加数据,并返回最后一次添加的ID
func Insert(insert string, args ...any) (int64, error) {
	db := GetDb()
	defer db.Close()
	rs, err := db.Exec(insert, args...)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	lastInsertId, err := rs.LastInsertId()
	if err != nil {
		return -1, err
	}
	return lastInsertId, nil
}

// 查询第一个数据并忽略错误
func SelectSingleOneIgnoreError[T any](query string, args ...any) T {
	value, _ := SelectSingleOne[T](query, args...)
	return value
}

// 查询第一个数据
func SelectSingleOne[T any](query string, args ...any) (T, error) {
	db := GetDb()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return *new(T), err // 返回默认值和错误
	}
	defer rows.Close()

	if !rows.Next() {
		return *new(T), sql.ErrNoRows // 如果没有结果，返回默认值和 ErrNoRows
	}

	var value T
	err = rows.Scan(&value) // 使用 Scan 将结果赋值给 value
	if err != nil {
		return *new(T), err // 返回默认值和错误
	}

	return value, nil
}

// 查询第一个数据
func SelectOne(query string, args ...any) {
	db := GetDb()
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 获取列的名称
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	if !rows.Next() {
	}

	// 创建一个长度与列数相同的slice来存放查询结果
	values := make([]interface{}, len(columns))
	// 创建一个[]interface{}的slice, 每个元素指向values中的对应位置
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 将当前行的数据扫描到valuePtrs中
	if err := rows.Scan(valuePtrs...); err != nil {
		log.Fatal(err)
	}

	// 使用map将列名和对应的值关联起来
	rowMap := make(map[string]interface{})
	for i, col := range columns {
		rowMap[col] = values[i]
	}
	fmt.Println("dfsf")
}

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
