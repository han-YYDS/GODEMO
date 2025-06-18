package database

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 基础模型（包含 GORM 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt）
type DBModel struct {
	*gorm.Model
}

// 数据库实例
var DBInstance *gorm.DB

// 初始化数据库连接（以 MySQL 为例）
func InitDB() {
	log.Println("Initializing database...")
	dsn := "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DBInstance = db

	// 自动建表
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed to migrate table: %v", err)
	}

	if err := db.AutoMigrate(&Class{}); err != nil {
		log.Fatalf("failed to migrate table: %v", err)
	}

	// 连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间
}

// 关闭数据库连接（可选，用于程序退出时）
func CloseDB() {
	sqlDB, _ := DBInstance.DB()
	sqlDB.Close()
}
