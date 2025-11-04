package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardStats 看板统计数据
type DashboardStats struct {
	TodayAppointments   int64 `json:"today_appointments"`    // 今日预约数
	PendingAppointments int64 `json:"pending_appointments"`  // 待处理预约数
	MonthAppointments   int64 `json:"month_appointments"`    // 本月预约数
	TotalShops          int64 `json:"total_shops"`           // 店铺总数
	TotalPrices         int64 `json:"total_prices"`          // 价格品种数
	TotalUsers          int64 `json:"total_users"`           // 管理员数
}

// RecentActivity 最近动态
type RecentActivity struct {
	Type      string    `json:"type"`       // 类型: appointment/price/user
	Title     string    `json:"title"`      // 标题
	Content   string    `json:"content"`    // 内容
	Time      time.Time `json:"time"`       // 时间
}

// GetDashboardStats 获取首页统计数据
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 今日预约数
	today := time.Now().Format("2006-01-02")
	models.DB.Model(&models.Appointment{}).
		Where("DATE(appointment_time) = ?", today).
		Count(&stats.TodayAppointments)

	// 待处理预约数
	models.DB.Model(&models.Appointment{}).
		Where("status = ?", models.AppointmentStatusPending).
		Count(&stats.PendingAppointments)

	// 本月预约数
	firstDay := time.Now().Format("2006-01") + "-01"
	models.DB.Model(&models.Appointment{}).
		Where("appointment_time >= ?", firstDay).
		Count(&stats.MonthAppointments)

	// 店铺总数
	models.DB.Model(&models.Shop{}).
		Where("status = ?", 1).
		Count(&stats.TotalShops)

	// 价格品种数
	models.DB.Model(&models.Price{}).
		Where("status = ?", 1).
		Count(&stats.TotalPrices)

	// 管理员数
	models.DB.Model(&models.AdminUser{}).
		Where("status = ?", 1).
		Count(&stats.TotalUsers)

	utils.Success(c, stats)
}

// GetRecentActivities 获取最近动态
func GetRecentActivities(c *gin.Context) {
	var activities []RecentActivity

	// 获取最近的预约（最多5条）
	var appointments []models.Appointment
	models.DB.Order("created_at DESC").Limit(5).Find(&appointments)

	for _, appt := range appointments {
		activity := RecentActivity{
			Type:    "appointment",
			Title:   "新预约",
			Content: appt.Name + " 预约了" + appt.MetalType + "回收",
			Time:    appt.CreatedAt,
		}
		activities = append(activities, activity)
	}

	// 获取最近更新的价格（最多3条）
	var prices []models.Price
	models.DB.Order("updated_at DESC").Limit(3).Find(&prices)

	for _, price := range prices {
		activity := RecentActivity{
			Type:    "price",
			Title:   "价格更新",
			Content: price.Name + " 价格已更新",
			Time:    price.UpdatedAt,
		}
		activities = append(activities, activity)
	}

	// 按时间排序
	// TODO: 实现更完善的排序逻辑

	utils.Success(c, activities)
}

// GetAppointmentTrend 获取预约趋势（最近7天）
func GetAppointmentTrend(c *gin.Context) {
	type TrendData struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	var trends []TrendData

	// 获取最近7天的数据
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		models.DB.Model(&models.Appointment{}).
			Where("DATE(appointment_time) = ?", date).
			Count(&count)

		trends = append(trends, TrendData{
			Date:  date,
			Count: count,
		})
	}

	utils.Success(c, trends)
}

