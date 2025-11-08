package models

import "time"

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `gorm:"not null;index" json:"role_id"`
	MenuID    uint      `gorm:"not null;index" json:"menu_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}








