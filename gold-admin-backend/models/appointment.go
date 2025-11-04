package models

import (
	"time"
)

// Appointment 预约模型
type Appointment struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ShopID          *uint      `gorm:"index" json:"shop_id"`                                 // 店铺ID
	UserID          *uint      `gorm:"index" json:"user_id"`                                 // 小程序用户ID
	Openid          string     `gorm:"size:100" json:"openid"`                               // 微信openid
	MetalType       string     `gorm:"size:50" json:"metal_type"`                            // 品种
	ServiceType     string     `gorm:"size:20" json:"service_type"`                          // 服务类型（store:到店 home:上门）
	AppointmentTime *time.Time `gorm:"index" json:"appointment_time"`                        // 预约时间
	Name            string     `gorm:"size:50" json:"name"`                                  // 姓名
	Phone           string     `gorm:"size:20" json:"phone"`                                 // 电话
	Address         string     `gorm:"size:255" json:"address"`                              // 地址（上门回收）
	Note            string     `gorm:"type:text" json:"note"`                                // 客户备注
	AdminRemark     string     `gorm:"type:text" json:"admin_remark"`                        // 管理员备注
	Status          string     `gorm:"size:20;default:pending;index" json:"status"`          // 状态（pending/confirmed/completed/cancelled）
	ConfirmedAt     *time.Time `json:"confirmed_at"`                                         // 确认时间
	CompletedAt     *time.Time `json:"completed_at"`                                         // 完成时间
	CancelledAt     *time.Time `json:"cancelled_at"`                                         // 取消时间
	HandlerID       *uint      `json:"handler_id"`                                           // 处理人ID
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Appointment) TableName() string {
	return "appointments"
}

// AppointmentStatus 预约状态
const (
	AppointmentStatusPending   = "pending"   // 待确认
	AppointmentStatusConfirmed = "confirmed" // 已确认
	AppointmentStatusCompleted = "completed" // 已完成
	AppointmentStatusCancelled = "cancelled" // 已取消
)

