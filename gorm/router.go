package main

import (
	"gormdemo/handler" // 替换为实际路径
	"log"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	log.Println("Registering routers...")
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		userGroup := api.Group("/users")
		{
			userGroup.POST("", handler.UserHandler.Create)       // 创建用户
			userGroup.GET("", handler.UserHandler.Get)           // 查询首个用户
			userGroup.PUT("/:id", handler.UserHandler.Update)    // 更新用户（通过 ID）
			userGroup.DELETE("/:id", handler.UserHandler.Delete) // 删除用户（通过 ID）
		}

		classGroup := api.Group("/class")
		{
			classGroup.POST("", handler.ClassHandler.Create)
		}
	}
	return r
}
