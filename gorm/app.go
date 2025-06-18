package main

import (
	"gormdemo/database" // 替换为实际路径
	"log"
	"net/http"
)

func main() {
	// 初始化数据库
	database.InitDB()
	defer database.CloseDB() // 程序退出时关闭连接

	// 初始化路由
	router := Router()

	// 启动服务
	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Server starting on :8080...")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("Server error: ", err)
	}
}
