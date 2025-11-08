package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码 200:成功 其他:失败
	Data    interface{} `json:"data"`    // 返回数据
	Message string      `json:"message"` // 提示信息
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: message,
	})
}

// Error 失败响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

// ErrorBadRequest 400 错误请求
func ErrorBadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// ErrorUnauthorized 401 未授权
func ErrorUnauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// ErrorForbidden 403 禁止访问
func ErrorForbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// ErrorNotFound 404 未找到
func ErrorNotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// ErrorInternalServer 500 服务器错误
func ErrorInternalServer(c *gin.Context, message string) {
	Error(c, 500, message)
}

// PageResult 分页结果
type PageResult struct {
	List  interface{} `json:"list"`  // 数据列表
	Total int64       `json:"total"` // 总数
	Page  int         `json:"page"`  // 当前页
	Size  int         `json:"size"`  // 每页大小
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, size int) {
	Success(c, PageResult{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}








