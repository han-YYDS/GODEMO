package database

// 用户表模型
type User struct {
	DBModel
	Name string `json:"name" gorm:"not null"` // 姓名（非空）
	Age  uint   `json:"age"`                  // 年龄
}

// 用户数据操作对象（DAO）
type userDAO struct{}

var UserDAO = userDAO{} // 实例化 DAO

// 创建用户
func (u *userDAO) Create(user User) error {
	return DBInstance.Create(&user).Error
}

// 查询首个用户
func (u *userDAO) First() (User, error) {
	var user User
	err := DBInstance.First(&user).Error
	return user, err
}

// 更新用户（通过 ID 匹配）
func (u *userDAO) Update(user User) error {
	return DBInstance.Model(&User{}).Where("id = ?", user.ID).Updates(user).Error
}

// 删除用户（软删除，GORM 自动处理）
func (u *userDAO) Delete(id uint) error {
	return DBInstance.Delete(&User{}, id).Error
}
