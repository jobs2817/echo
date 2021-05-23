package main

import (
	"database/sql"

	Dao "echo/dao"
	models "echo/models"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var msvr *server

type server struct {
	DB    *sql.DB
	Redis redis.Conn
}
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}
type Psd struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func GetResponse(code int, message string, data interface{}) *models.Response {
	rep := &models.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return rep
}

// login ...
func login(c echo.Context) error {
	phone := c.FormValue("phone")
	password := c.FormValue("password")
	corPassword := &Psd{}
	err := msvr.DB.QueryRow(Dao.GetUserPassword, phone).Scan(&corPassword.Name, &corPassword.Password)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	if password == corPassword.Password {
		// Set custom claims
		claims := &jwtCustomClaims{
			phone,
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			},
		}
		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{
			"token":    t,
			"userName": corPassword.Name,
		})
	}
	return echo.ErrUnauthorized
}

// 查询单条数据 ...
func queryInfo(c echo.Context) error {
	user := &models.User{}
	err := msvr.DB.QueryRow(Dao.UserInfoSql, int(1)).Scan(&user.Name, &user.Age, &user.Id, &user.Phone, &user.Hobby, &user.StartTime, &user.EndTime)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResponse(http.StatusOK, "ok", user))
}

// 查询列表数据 ...
func getList(c echo.Context) error {
	list := []models.User{}
	rows, err := msvr.DB.Query(Dao.GetUserList)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		user := &models.User{}
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
	u := new(models.User)
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
	u := new(models.User)
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
	e.POST("/login", login)
	g := e.Group("/api")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	g.Use(middleware.JWTWithConfig(config))
	g.GET("/queryInfo", queryInfo)
	g.GET("/queryList", getList)
	g.GET("/delUser", delUser)
	g.POST("/addUser", addUser)
	g.POST("/updateUser", updateUser)

	e.Start(":8888")
}
