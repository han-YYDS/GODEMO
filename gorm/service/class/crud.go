package class

import "gormdemo/database"

type ClassForm struct {
	Id uint `json:"id" binding:"required"` // 绑定验证规则（必填）
}

// 创建用户（业务层调用数据层）
func CreateClass(form ClassForm) error {
	class := database.Class{
		Id: form.Id,
	}
	return database.ClassDAO.Create(class)
}
