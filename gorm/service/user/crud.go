package user

import (
	"gormdemo/database" // 替换为实际路径
	"strconv"

	"gorm.io/gorm"
)

// 前端传入的用户表单数据
type UserForm struct {
	Name string `json:"name" binding:"required"` // 绑定验证规则（必填）
	Age  uint   `json:"age"`
}

// 创建用户（业务层调用数据层）
func CreateUser(form UserForm) error {
	user := database.User{
		Name: form.Name,
		Age:  form.Age,
	}
	return database.UserDAO.Create(user)
}

// 查询首个用户
func GetFirstUser() (database.User, error) {
	return database.UserDAO.First()
}

// 更新用户（需传入 ID）
func UpdateUser(form UserForm, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	user := database.User{
		DBModel: database.DBModel{Model: &gorm.Model{ID: uint(id)}},
		Name:    form.Name,
		Age:     form.Age,
	}
	return database.UserDAO.Update(user)
}

// 删除用户
func DeleteUser(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return database.UserDAO.Delete(uint(id))
}
