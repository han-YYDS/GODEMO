package database

type Class struct {
	DBModel
	Id uint `json:"id" gorm:"not null"`
}

// 用户数据操作对象（DAO）
type classDAO struct{}

var ClassDAO = classDAO{} // 实例化 DAO

// 创建用户
func (u *classDAO) Create(class Class) error {
	return DBInstance.Create(&class).Error
}
