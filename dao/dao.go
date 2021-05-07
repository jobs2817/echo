package Dao

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (DB *sql.DB) {
	DB, err := sql.Open("mysql", "root:262817@tcp(www.haofeifu.com:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	fmt.Printf("lianjiele")
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	return DB
}
