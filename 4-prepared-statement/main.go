package main

// 演示prepared statement怎样使用

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StmtUsage1(db *sql.DB) {
	// 创建一个Stmt对象
	// 这个Stmt对象是并发安全的
	stmt, err := db.Prepare("insert into Student(SId, Sname, Sage, Ssex) values (?, ?, ?, ?)")
	// 记得要关闭stmt对象
	defer func() { _ = stmt.Close() }()
	if err != nil {
		fmt.Printf("db prepare err: %v\n", err)
		return
	}
	// 然后这个stmt对象就可以反复使用了
	stus := []struct {
		id     string
		name   string
		age    time.Time
		gender string
	}{
		{"14", "john", time.Date(1990, 7, 11, 0, 0, 0, 0, time.Local), "男"},
		{"15", "mike", time.Date(1998, 9, 4, 0, 0, 0, 0, time.Local), "男"},
	}

	for i, stu := range stus {
		// 然后再执行操作的时候填充参数
		fmt.Printf("i=%d\n", i)
		sqlResult, err := stmt.Exec(stu.id, stu.name, stu.age, stu.gender)
		// sqlResult类型是 sql.Result
		if err != nil {
			fmt.Printf("stmt Exec err: %v, at %d\n", err, i)
			continue
		}
		rowsAffected, err := sqlResult.RowsAffected()
		if err != nil {
			fmt.Printf("err : %v\n", err)
			continue
		}
		fmt.Printf("rowsAffected=%d\n", rowsAffected)
	}
}

func StmtUsage2(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, "select SId, Sname from Student where SId = ?")
	if err == nil {
		defer stmt.Close()
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	ids := []string{"14", "15"}
	for _, id := range ids {
		if rows, err := stmt.Query(id); err == nil {
			for rows.Next() {
				var sid, sname string
				if rows.Scan(&sid, &sname) == nil {
					fmt.Printf("SId=%s, Sname=%s at %s\n", sid, sname, id)
				}
			}
		}
	}
}

func StmtUsage3(db *sql.DB) {
	stmt, err := db.Prepare("delete from Student where SId = ?")
	if err == nil {
		defer stmt.Close()
	}
	ids := []string{"14", "15"}
	for _, id := range ids {
		stmt.Exec(id)
	}
}

func main() {
	db, err := sql.Open("mysql", "ryan:123456@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local")
	defer func() { _ = db.Close() }()
	if err != nil {
		log.Fatal(err)
		return
	}
	StmtUsage1(db)
	StmtUsage2(db)
	StmtUsage3(db)
}
