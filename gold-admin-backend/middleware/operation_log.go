package middleware

import (
	"bytes"
	"encoding/json"
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 自定义响应写入器，用于捕获响应状态码
type responseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建自定义响应写入器
		customWriter := &responseWriter{
			ResponseWriter: c.Writer,
			statusCode:     200, // 默认200
		}
		c.Writer = customWriter

		// 处理请求
		c.Next()

		// 计算执行时长
		duration := time.Since(startTime).Milliseconds()

		// 获取用户信息
		var userID *uint
		var username string
		if claims, exists := c.Get("claims"); exists {
			if customClaims, ok := claims.(*utils.Claims); ok {
				userID = &customClaims.UserID
				username = customClaims.Username
			}
		}

		// 判断是否需要记录日志
		if !shouldLog(c.Request.Method, c.Request.URL.Path) {
			return
		}

		// 构建操作日志
		log := models.OperationLog{
			UserID:         userID,
			Username:       username,
			Module:         getModule(c.Request.URL.Path),
			Action:         getAction(c.Request.Method, c.Request.URL.Path),
			Description:    getDescription(c.Request.Method, c.Request.URL.Path, string(requestBody)),
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			IP:             c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
			ResponseStatus: customWriter.statusCode,
			Duration:       duration,
			CreatedAt:      time.Now(),
		}

		// 记录请求参数（排除敏感信息）
		if len(requestBody) > 0 && shouldLogRequestBody(c.Request.URL.Path) {
			log.RequestParams = filterSensitiveData(requestBody)
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			log.ErrorMessage = c.Errors.String()
		}

		// 异步保存日志到数据库
		go func() {
			models.DB.Create(&log)
		}()
	}
}

// shouldLog 判断是否需要记录日志
func shouldLog(method, path string) bool {
	// 不记录的路径
	excludePaths := []string{
		"/health",
		"/api/logs", // 查询日志本身不记录
	}

	for _, excludePath := range excludePaths {
		if strings.HasPrefix(path, excludePath) {
			return false
		}
	}

	// 只记录特定方法
	recordMethods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, m := range recordMethods {
		if method == m {
			return true
		}
	}

	// 登录操作也记录
	if path == "/api/login" && method == "POST" {
		return true
	}

	return false
}

// shouldLogRequestBody 判断是否记录请求体
func shouldLogRequestBody(path string) bool {
	// 不记录请求体的路径（如密码相关）
	excludePaths := []string{
		"/api/login",
		"/password",
	}

	for _, excludePath := range excludePaths {
		if strings.Contains(path, excludePath) {
			return false
		}
	}

	return true
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(data []byte) string {
	var params map[string]interface{}
	if err := json.Unmarshal(data, &params); err != nil {
		return string(data)
	}

	// 过滤敏感字段
	sensitiveFields := []string{"password", "pwd", "token", "secret"}
	for _, field := range sensitiveFields {
		if _, exists := params[field]; exists {
			params[field] = "******"
		}
	}

	filtered, _ := json.Marshal(params)
	return string(filtered)
}

// getModule 根据路径获取模块名称
func getModule(path string) string {
	if strings.Contains(path, "/users") {
		return models.ModuleUser
	}
	if strings.Contains(path, "/roles") {
		return models.ModuleRole
	}
	if strings.Contains(path, "/menus") {
		return models.ModuleMenu
	}
	if strings.Contains(path, "/prices") {
		return models.ModulePrice
	}
	if strings.Contains(path, "/shops") {
		return models.ModuleShop
	}
	if strings.Contains(path, "/appointments") {
		return models.ModuleAppointment
	}
	if strings.Contains(path, "/dashboard") {
		return models.ModuleDashboard
	}
	if strings.Contains(path, "/login") || strings.Contains(path, "/logout") {
		return models.ModuleAuth
	}
	return "其他"
}

// getAction 根据方法和路径获取操作动作
func getAction(method, path string) string {
	if strings.Contains(path, "/login") {
		return models.ActionLogin
	}
	if strings.Contains(path, "/logout") {
		return models.ActionLogout
	}

	switch method {
	case "POST":
		return models.ActionCreate
	case "PUT", "PATCH":
		return models.ActionUpdate
	case "DELETE":
		return models.ActionDelete
	case "GET":
		if strings.Contains(path, "/export") {
			return models.ActionExport
		}
		return models.ActionQuery
	default:
		return "未知操作"
	}
}

// getDescription 生成操作描述
func getDescription(method, path, body string) string {
	module := getModule(path)
	action := getAction(method, path)

	// 尝试从body中提取关键信息
	var params map[string]interface{}
	keyInfo := ""
	if len(body) > 0 && json.Unmarshal([]byte(body), &params) == nil {
		// 提取一些关键字段
		if name, ok := params["name"].(string); ok && name != "" {
			keyInfo = name
		} else if username, ok := params["username"].(string); ok && username != "" {
			keyInfo = username
		} else if title, ok := params["title"].(string); ok && title != "" {
			keyInfo = title
		} else if realName, ok := params["real_name"].(string); ok && realName != "" {
			keyInfo = realName
		}
	}

	if keyInfo != "" {
		return action + module + ": " + keyInfo
	}

	return action + module
}














