package Dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

func New() redis.Conn {
	c, err := redis.Dial("tcp", "47.101.203.226:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		panic(err)
	}
	return c
}
func Init() (DB *sql.DB) {
	DB, err := sql.Open("mysql", "root:262817@tcp(www.haofeifu.com:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	return DB
}
