package models

import "time"

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int8      `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}





