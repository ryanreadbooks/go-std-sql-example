package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"	// 必须导入这个驱动包，才可以在Open中指定mysql，不然会报错
)

func main() {
	// 连接格式为
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err := sql.Open("mysql", "ryan:123456@(127.0.0.1:3306)/study?charset=utf8mb4&parsetTime=True&loc=Local")
	if err != nil {
		fmt.Printf("open db error: %v\n", err)
		return
	}
	defer db.Close()
}