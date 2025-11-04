package models

import "time"

// Menu 菜单模型
type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ParentID   int       `gorm:"default:0" json:"parent_id"`          // 父菜单ID，0为顶级
	Type       int8      `gorm:"default:1" json:"type"`               // 类型 1:目录 2:菜单 3:按钮
	Name       string    `gorm:"size:50;not null" json:"name"`        // 菜单名称（英文）
	Title      string    `gorm:"size:50;not null" json:"title"`       // 菜单标题（中文）
	Icon       string    `gorm:"size:50" json:"icon"`                 // 图标
	Path       string    `gorm:"size:100" json:"path"`                // 路由路径
	Component  string    `gorm:"size:100" json:"component"`           // 组件路径
	Permission string    `gorm:"size:100" json:"permission"`          // 权限标识（预留）
	Sort       int       `gorm:"default:0" json:"sort"`               // 排序
	Visible    int8      `gorm:"default:1" json:"visible"`            // 是否显示 1:显示 0:隐藏
	Status     int8      `gorm:"default:1" json:"status"`             // 状态 1:启用 0:禁用
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// MenuTree 菜单树结构（用于前端展示）
type MenuTree struct {
	ID         uint        `json:"id"`
	ParentID   int         `json:"parent_id"`
	Type       int8        `json:"type"`
	Name       string      `json:"name"`
	Title      string      `json:"title"`
	Icon       string      `json:"icon"`
	Path       string      `json:"path"`
	Component  string      `json:"component"`
	Permission string      `json:"permission"`
	Sort       int         `json:"sort"`
	Children   []*MenuTree `json:"children,omitempty"`
}



