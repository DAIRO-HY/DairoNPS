package main

import (
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/NPSDB"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	NPSDB.Init()
	//clientDto:=dto.ClientDto{
	//	Name:"sd",
	//	Key:"dsfdsf",
	//	Remark:"6t7uyghjbmnlkkj",
	//}
	//ClientDao.Add(clientDto)

	clientDto := ClientDao.SelectOne(1)
	fmt.Println(clientDto.Name)
	return

	//// 打开数据库连接，没有文件时会自动创建
	//db, err := sql.Open("sqlite3", "./example.sqlite")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//
	//// 创建表
	//sqlStmt := `
	//CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);
	//`
	//_, err = db.Exec(sqlStmt)
	//if err != nil {
	//	log.Fatalf("%q: %s\n", err, sqlStmt)
	//}
	//
	//// 插入数据
	//_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 查询数据
	//rows, err := db.Query("SELECT id, name FROM users")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//
	//// 打印查询结果
	//for rows.Next() {
	//	var id int
	//	var name string
	//	err = rows.Scan(&id, &name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("ID: %d, Name: %s\n", id, name)
	//}
	//
	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

