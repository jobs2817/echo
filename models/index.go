package models

type User struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Id        int    `json:"id"`
	Phone     string `json:"phone"`
	Hobby     string `json:"hobby"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ...
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// auth ...
type Psd struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// token ...
type JwtToken struct {
	Token string `json:"token"`
}
