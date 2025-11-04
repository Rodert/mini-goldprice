package models

import "time"

// Shop 店铺模型
type Shop struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`        // 店铺名称
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"` // 店铺编码
	Address     string    `gorm:"size:255" json:"address"`              // 地址
	Phone       string    `gorm:"size:20" json:"phone"`                 // 固定电话
	Mobile      string    `gorm:"size:20" json:"mobile"`                // 手机号
	Hours       string    `gorm:"size:100" json:"hours"`                // 营业时间
	Latitude    float64   `json:"latitude"`                             // 纬度
	Longitude   float64   `json:"longitude"`                            // 经度
	Description string    `gorm:"type:text" json:"description"`         // 店铺介绍
	Status      int8      `gorm:"default:1" json:"status"`              // 状态 1:营业 0:停业
	Sort        int       `gorm:"default:0" json:"sort"`                // 排序
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Shop) TableName() string {
	return "shops"
}

