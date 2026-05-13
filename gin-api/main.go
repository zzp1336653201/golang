package main

import (
	"log"

	"gin-api/config"
	"gin-api/models"
	"gin-api/routes"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移
	models.AutoMigrate()
	log.Println("数据库迁移完成")

	// 配置路由
	r := routes.SetupRouter()

	// 启动服务
	log.Println("服务启动: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
