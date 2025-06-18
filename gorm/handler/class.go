package handler

import (
	"fmt"
	service "gormdemo/service/class"
	"net/http"

	"github.com/gin-gonic/gin"
)

type classHandler struct{}

var ClassHandler = classHandler{}

func (h *classHandler) Create(c *gin.Context) {
	fmt.Println("111")
	// 请求参数 -> 结构体
	var form service.ClassForm
	if err := c.ShouldBindJSON(&form); err != nil {
		Error(c, err, http.StatusBadRequest)
		return
	}

	// 结构体 -> 数据库
	if err := service.CreateClass(form); err != nil {
		Error(c, err, http.StatusInternalServerError)
		return
	}
	Success(c, "User created successfully")
}
