package main

import (
	"database/sql"
	Dao "echo/dao"
	"net/http"

	"github.com/labstack/echo"
)

var msvr *server

type server struct {
	DB *sql.DB
}

type User struct {
	Name string
	Age  int
	Id   int
}
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetResponse(code int, message string, data interface{}) *Response {
	rep := &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return rep
}
func getList(c echo.Context) error {
	user := &User{}
	err := msvr.DB.QueryRow(Dao.UserInfoSql, int(1)).Scan(&user.Name, &user.Age, &user.Id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", user))
}

//导入echo包
func main() {
	e := echo.New()

	msvr = &server{
		DB: Dao.Init(),
	}
	e.GET("/getList", getList)
	e.Start(":8888")
}
