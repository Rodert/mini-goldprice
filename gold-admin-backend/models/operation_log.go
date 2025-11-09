package models

import "time"

// OperationLog 操作日志模型
type OperationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         *uint     `gorm:"index" json:"user_id"`                       // 操作用户ID
	Username       string    `gorm:"size:50;index" json:"username"`              // 操作用户名
	Module         string    `gorm:"size:50;index" json:"module"`                // 操作模块（用户管理、角色管理等）
	Action         string    `gorm:"size:50;index" json:"action"`                // 操作动作（创建、更新、删除等）
	Description    string    `gorm:"size:500" json:"description"`                // 操作描述
	Method         string    `gorm:"size:10" json:"method"`                      // HTTP方法（GET、POST、PUT、DELETE）
	Path           string    `gorm:"size:255" json:"path"`                       // 请求路径
	IP             string    `gorm:"size:50" json:"ip"`                          // 操作IP
	UserAgent      string    `gorm:"type:text" json:"user_agent"`                // 用户代理
	RequestParams  string    `gorm:"type:text" json:"request_params,omitempty"`  // 请求参数（JSON）
	ResponseStatus int       `json:"response_status"`                            // 响应状态码
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`   // 错误信息
	Duration       int64     `json:"duration"`                                   // 执行时长（毫秒）
	CreatedAt      time.Time `gorm:"index" json:"created_at"`                    // 创建时间
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"    // 信息
	LogLevelWarning LogLevel = "warning" // 警告
	LogLevelError   LogLevel = "error"   // 错误
)

// 操作模块常量
const (
	ModuleAuth       = "认证管理"
	ModuleUser       = "用户管理"
	ModuleRole       = "角色管理"
	ModuleMenu       = "菜单管理"
	ModulePrice      = "价格管理"
	ModuleShop       = "店铺管理"
	ModuleAppointment = "预约管理"
	ModuleDashboard  = "首页统计"
)

// 操作动作常量
const (
	ActionLogin  = "登录"
	ActionLogout = "登出"
	ActionCreate = "创建"
	ActionUpdate = "更新"
	ActionDelete = "删除"
	ActionQuery  = "查询"
	ActionExport = "导出"
	ActionImport = "导入"
)














