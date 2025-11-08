package main

import (
	"embed"
	"flag"
	"fmt"
	"gold-admin-backend/config"
	"gold-admin-backend/models"
	"gold-admin-backend/router"
	"log"
	"os"
)

//go:embed config/config.yaml
var embeddedConfig embed.FS

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "c", "", "配置文件路径（可选，默认使用嵌入的配置）")
	flag.Parse()

	// 加载配置
	var err error
	if configPath != "" {
		// 使用外部配置文件
		log.Printf("使用外部配置文件: %s", configPath)
		if err = config.LoadConfigFromFile(configPath); err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}
	} else {
		// 使用嵌入的默认配置
		log.Println("使用嵌入的默认配置文件")
		if err = config.LoadConfigFromEmbed(embeddedConfig, "config/config.yaml"); err != nil {
			log.Fatalf("加载嵌入配置失败: %v", err)
		}
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
