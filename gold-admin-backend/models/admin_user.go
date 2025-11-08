package models

import (
	"time"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password      string    `gorm:"size:255;not null" json:"-"` // 不返回给前端
	RealName      string    `gorm:"size:50" json:"real_name"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Email         string    `gorm:"size:100" json:"email"`
	Avatar        string    `gorm:"size:255" json:"avatar"`
	Status        int8      `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIP   string    `gorm:"size:50" json:"last_login_ip"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}







