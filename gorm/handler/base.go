package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// 错误响应
func Error(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"code": code,
		"msg":  err.Error(),
		"data": nil,
	})
}
