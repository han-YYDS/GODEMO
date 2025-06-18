package handler

import (
	service "gormdemo/service/user" // 替换为实际路径
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct{}

var UserHandler = userHandler{} // 实例化控制器

// 创建用户
func (h *userHandler) Create(c *gin.Context) {
	// 请求参数 -> 结构体
	var form service.UserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		Error(c, err, http.StatusBadRequest)
		return
	}

	// 结构体 -> 数据库
	if err := service.CreateUser(form); err != nil {
		Error(c, err, http.StatusInternalServerError)
		return
	}
	Success(c, "User created successfully")
}

// 查询首个用户
func (h *userHandler) Get(c *gin.Context) {
	user, err := service.GetFirstUser()
	if err != nil {
		Error(c, err, http.StatusNotFound)
		return
	}
	Success(c, user)
}

// 更新用户
func (h *userHandler) Update(c *gin.Context) {
	var form service.UserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		Error(c, err, http.StatusBadRequest)
		return
	}

	// 从 URL 参数中获取 ID（需在路由中定义参数）
	id, exists := c.Params.Get("id")
	if !exists {
		Error(c, nil, http.StatusBadRequest)
		return
	}

	if err := service.UpdateUser(form, id); err != nil {
		Error(c, err, http.StatusInternalServerError)
		return
	}
	Success(c, "User updated successfully")
}

// 删除用户
func (h *userHandler) Delete(c *gin.Context) {
	id, exists := c.Params.Get("id")
	if !exists {
		Error(c, nil, http.StatusBadRequest)
		return
	}

	if err := service.DeleteUser(id); err != nil {
		Error(c, err, http.StatusInternalServerError)
		return
	}
	Success(c, "User deleted successfully")
}
