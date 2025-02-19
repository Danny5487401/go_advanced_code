package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

var (
	db *sqlx.DB
)

func initDB() (err error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

type Person struct {
	Id         int       `db:"Id"`
	Name       string    `db:"Name"`
	City       string    `db:"City"`
	AddTime    time.Time `db:"AddTime"`
	UpdateTime time.Time `db:"UpdateTime"`
}

// 查询单条数据示例
func queryRowDemo() {
	p := &Person{}
	err := db.Get(p, "select * from Person where Name=?", "Zhang San")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("person: %+v", p)
}

// 插入数据
func insertRowDemo() {
	insertResult := db.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	lastInsertId, _ := insertResult.LastInsertId()
	log.Println("Insert Id is ", lastInsertId)
}
