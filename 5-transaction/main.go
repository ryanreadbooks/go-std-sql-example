package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func PrintErr(err error) {
	fmt.Printf("err: %v\n", err)
}

// 演示怎样使用事务
func TxDemo1(db *sql.DB) {
	// 开启一个事务，带有默认的隔离级别
	tx, err := db.Begin()
	if err != nil {
		PrintErr(err)
		return
	}
	// tx的类型为 *sql.Tx
	// 其所包含的方法和*sql.DB中所包含的方法类似
	// 都有 Query Exec Prepare等
	// 但是特有的是，Commit 和 Rollback用于提交事务和事务回滚操作
	// 同样可以在事务上创建一个Stmt对象，不同的是这个对象的生命周期在事务结束后就结束
	stmt, err := tx.Prepare("insert into users(name, email, country, phone_number, updated_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		PrintErr(err)
		return
	}
	defer stmt.Close()
	// 开始通过事务插入内容
	res, err := stmt.Exec("Lily Alister", "lily.alister@google.com", "Brazil", "18612301202", time.Now())
	if err != nil {
		// 如果出错了，就可以回滚了
		PrintErr(err)
		tx.Rollback()
		return
	}
	lastId, _ := res.LastInsertId()
	rowCnt, _ := res.RowsAffected()
	fmt.Printf("lastId=%d, rowCnt=%d\n", lastId, rowCnt)
	// 事务还没有结束，可以继续操作
	res, err = tx.Exec("update users set name = ? where id = ?", "Lily Alister Jr.", lastId)
	if err != nil {
		PrintErr(err)
		// 出错就rollback
		tx.Rollback()
		return
	}
	lastId, _ = res.LastInsertId()
	rowCnt, _ = res.RowsAffected()
	fmt.Printf("lastId=%d, rowCnt=%d\n", lastId, rowCnt)
	// 使用commit提交
	if err := tx.Commit(); err != nil {
		// 提交的时候出错了，回滚
		PrintErr(err)
		tx.Rollback()
		return
	}
}

// 在事务执行的过程中发生panic的处理
// 在recover中rollback
func TxDemo2(db *sql.DB) {
	var tx *sql.Tx
	defer func() {
		if p := recover(); p != nil {
			// 回滚
			fmt.Printf("panic: %v, rollback\n", p)
			tx.Rollback()
		}
	}()
	var err error
	tx, err = db.Begin()
	if err != nil {
		PrintErr(err)
		return
	}
	stmt, err := tx.Prepare("insert into users(name, email, country, phone_number, updated_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		PrintErr(err)
		return
	}
	defer stmt.Close()
	// 开始通过事务插入内容
	res, err := stmt.Exec("Lily Alister", "lily.alister@google.com", "Brazil", "18612301202", time.Now())
	if err != nil {
		// 如果出错了，就可以回滚了
		PrintErr(err)
		tx.Rollback()
		return
	}
	lastId, _ := res.LastInsertId()
	rowCnt, _ := res.RowsAffected()
	fmt.Printf("lastId=%d, rowCnt=%d\n", lastId, rowCnt)
	// 事务还没有结束，可以继续操作
	_, err = tx.Exec("update users set name = ? where id = ?", "Lily Alister Jr.", lastId)
	if err != nil {
		PrintErr(err)
		// 出错就rollback
		tx.Rollback()
		return
	}

	// commit之前强行panic
	panic("Deliberate panic")
}

func main() {
	db, err := sql.Open("mysql", "ryan:123456@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local")
	defer func() { _ = db.Close() }()
	if err != nil {
		fmt.Printf("open db err: %v\n", err)
		return
	}
	// TxDemo1(db)
	TxDemo2(db)
}
