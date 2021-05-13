package Dao

var (
	// 查询单条数据
	UserInfoSql = "SELECT name, age, id, phone, hobby, start_time, end_time FROM user WHERE id = ?"
	// 获取列表数据
	GetUserList = "SELECT name, age, id, phone, hobby, start_time, end_time FROM user"
	// 新增数据
	AddUser = "INSERT INTO user (name, age, phone, hobby) VALUES (?, ?, ?, ?)"
	// 删除数据
	DelUser = "DELETE FROM user WHERE id = ?"
	// 更新数据
	Updateuser = "UPDATE user SET name = ?, age = ?, phone = ?, hobby = ? WHERE id = ?"
)
