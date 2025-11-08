package middleware

import (
	"gold-admin-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorUnauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 提取 token（格式: Bearer xxx）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.ErrorUnauthorized(c, "Authorization 格式错误")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorUnauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}






