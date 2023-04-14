package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // 只导入包（执行里面的init函数)
)

// 简单的查询操作，并且获取查询结果
func SimpleQuery(db *sql.DB) {
	fmt.Println("SimpleQuery...")
	var rows *sql.Rows
	var err error
	rows, err = db.Query("select * from student")
	defer func() { _ = rows.Close() }() // 记得Close掉查询结果
	if err != nil {
		fmt.Printf("simple query err: %v\n", err)
		return
	}
	// 获取查询结果
	// *sql.Rows.Next() 如果有下一行，则Next()方法返回true，并且可以用Scan()方法读取
	for rows.Next() {
		// 需要直到表结构，并且定义好接受的变量
		var sid string
		var sname string
		var age time.Time
		var gender string
		err = rows.Scan(&sid, &sname, &age, &gender)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s, %s, %v, %s\n", sid, sname, age, gender)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// 可以用动态的查询参数
func QueryWithParam(db *sql.DB) {
	fmt.Println("QueryWithParam...")
	var target string = "张三"
	// 用问号来作为占位符
	rows, err := db.Query("select * from student where Sname = ?", target)
	defer func() { _ = rows.Close() }()
	if err != nil {
		fmt.Printf("query err: %v\n", err)
		return
	}
	for rows.Next() {
		var sid, sname, gender string
		var age time.Time
		err := rows.Scan(&sid, &sname, &age, &gender)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s, %s, %v, %s\n", sid, sname, age, gender)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	db, err := sql.Open("mysql", "ryan:123456@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local")
	defer func() { _ = db.Close() }()
	if err != nil {
		fmt.Printf("open db err: %v\n", err)
		return
	}
	SimpleQuery(db)
	QueryWithParam(db)
	fmt.Println("========================================")
}
