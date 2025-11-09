package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	// 允许的源地址（根据环境配置）
	// 开发环境：本地地址
	// 生产环境：需要添加实际的前端域名和小程序域名
	allowOrigins := []string{
		"http://localhost:9527",        // 开发环境前端
		"http://127.0.0.1:9527",        // 开发环境前端
		"http://localhost:8080",        // 开发环境API
		"https://gold.javapub.net.cn",  // 生产环境API
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}















