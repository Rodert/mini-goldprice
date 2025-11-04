package main

import (
	"fmt"
	"gold-admin-backend/config"
	"gold-admin-backend/models"
	"gold-admin-backend/router"
	"log"
	"os"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("./config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建数据目录
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	// 创建日志目录
	if err := os.MkdirAll(config.AppConfig.Log.Path, 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化基础数据
	if err := models.InitData(); err != nil {
		log.Fatalf("初始化基础数据失败: %v", err)
	}

	// 设置 Gin 模式
	//gin.SetMode(config.AppConfig.Server.Mode)

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("服务器启动成功，监听地址: %s", addr)
	log.Printf("数据库文件: %s", config.AppConfig.Database.Path)
	log.Println("访问 http://localhost:8080/ping 测试服务")

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
