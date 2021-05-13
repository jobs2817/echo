package main

import (
	"database/sql"
	Dao "echo/dao"
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

var msvr *server

type server struct {
	DB    *sql.DB
	Redis redis.Conn
}

type User struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Id        int    `json:"id"`
	Phone     string `json:"phone"`
	Hobby     string `json:"hobby"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
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

// 查询单条数据 ...
func queryInfo(c echo.Context) error {
	user := &User{}
	err := msvr.DB.QueryRow(Dao.UserInfoSql, int(1)).Scan(&user.Name, &user.Age, &user.Id, &user.Phone, &user.Hobby, &user.StartTime, &user.EndTime)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", user))
}

// 查询列表数据 ...
func getList(c echo.Context) error {
	list := []User{}
	rows, err := msvr.DB.Query(Dao.GetUserList)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		user := &User{}
		if err = rows.Scan(&user.Name, &user.Age, &user.Id, &user.Phone, &user.Hobby, &user.StartTime, &user.EndTime); err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		list = append(list, *user)
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", list))
}

// 新增数据
func addUser(c echo.Context) error {
	u := new(User)
	c.Bind(u)
	rst, err := msvr.DB.Exec(Dao.AddUser, u.Name, u.Age, u.Phone, u.Hobby)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	lastId, err := rst.LastInsertId()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", lastId))
}

// 删除数据
func delUser(c echo.Context) error {
	id := c.QueryParam("id")
	rst, err := msvr.DB.Exec(Dao.DelUser, id)
	if err != nil {
		return err
	}
	delId, err := rst.RowsAffected()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", delId))
}

// 更新数据
func updateUser(c echo.Context) error {
	u := new(User)
	c.Bind(u)
	fmt.Printf("user: %v\n", u)
	rst, err := msvr.DB.Exec(Dao.Updateuser, u.Name, u.Age, u.Phone, u.Hobby, u.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	efId, err := rst.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", efId))
}

//导入echo包
func main() {
	e := echo.New()

	redis := Dao.New()
	msvr = &server{
		DB:    Dao.Init(),
		Redis: redis,
	}
	_, err := redis.Do("Set", "abc", 100)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	e.GET("/queryInfo", queryInfo)
	e.GET("/queryList", getList)
	e.GET("/delUser", delUser)
	e.POST("/addUser", addUser)
	e.POST("/updateUser", updateUser)

	e.Start(":8888")
}
