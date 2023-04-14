package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 插入数据
func InsertIntoDB1(db *sql.DB) {
	res, err := db.Exec("insert into Student values(?, ?, ?, ?)", "14", "菲利克斯", time.Now(), "男")
	if err != nil {
		fmt.Printf("db exec error: %v\n", err)
		return
	}
	// 获取插入的结果
	id, err := res.LastInsertId()	// 如果有主键的话，可以获取新增数据的id
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	rows, err := res.RowsAffected()
	if err !=nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("新插入的行的主键为: %d, 影响的行数为: %d\n", id, rows)
}

// 更新数据
func UpdateFromDB(db *sql.DB) {
	res, err := db.Exec("update Student set Sname = ? where SId = ?", "法外狂徒", "14")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("rows affected = %d\n", rowsAffected)
}

// 删除一条记录
func DeleteFromDB(db *sql.DB) {
	res, err := db.Exec("delete from Student where SId = ?", "14")
	if err != nil {
		fmt.Printf("db exec err: %v\n", err)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("rowAffected = %d\n", rowsAffected)
}

// exec context 这个方法可以传入一个context.Context对象
func ExecWithContext(db *sql.DB) {
	// db.Exec()其实也是对ExecContext的封装，其内部使用的是context.Background()返回一个context
	// 可以给一个语句执行的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	res, err := db.ExecContext(ctx, "insert into Student values(?, ?, ?, ?)", "14", "菲利克斯", time.Now(), "男")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return 
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%d\n", lastId)
	res, err = db.ExecContext(ctx, "delete from Student where SId = ?", "14")
	if err != nil {
		fmt.Printf("%v\n", err)
		return 
	}
	n, _ := res.RowsAffected()
	fmt.Printf("rowAffected = %d\n", n)
}

func main() {
	db, err := sql.Open("mysql", "ryan:123456@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local")
	defer func() { _ = db.Close() }()
	if err != nil {
		log.Fatal(err)
		return
	}

	InsertIntoDB1(db)
	UpdateFromDB(db)
	DeleteFromDB(db)
	ExecWithContext(db)
}
