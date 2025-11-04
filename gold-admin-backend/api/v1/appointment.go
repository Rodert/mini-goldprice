package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AppointmentListRequest 预约列表请求参数
type AppointmentListRequest struct {
	Page      int       `form:"page" binding:"omitempty,min=1"`
	PageSize  int       `form:"page_size" binding:"omitempty,min=1,max=100"`
	ShopID    *uint     `form:"shop_id"`
	Status    string    `form:"status"`
	Name      string    `form:"name"`
	Phone     string    `form:"phone"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02"`
}

// AppointmentCreateRequest 创建预约请求参数
type AppointmentCreateRequest struct {
	ShopID          *uint      `json:"shop_id"`
	Openid          string     `json:"openid"`
	MetalType       string     `json:"metal_type" binding:"required"`
	ServiceType     string     `json:"service_type" binding:"required,oneof=store home"`
	AppointmentTime *time.Time `json:"appointment_time" binding:"required"`
	Name            string     `json:"name" binding:"required"`
	Phone           string     `json:"phone" binding:"required"`
	Address         string     `json:"address"`
	Note            string     `json:"note"`
}

// AppointmentUpdateRequest 更新预约请求参数
type AppointmentUpdateRequest struct {
	ShopID          *uint      `json:"shop_id"`
	AppointmentTime *time.Time `json:"appointment_time"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone"`
	Address         string     `json:"address"`
	Note            string     `json:"note"`
	AdminRemark     string     `json:"admin_remark"`
	Status          string     `json:"status"`
}

// GetAppointmentList 获取预约列表
func GetAppointmentList(c *gin.Context) {
	var req AppointmentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 构建查询
	query := models.DB.Model(&models.Appointment{})

	if req.ShopID != nil {
		query = query.Where("shop_id = ?", *req.ShopID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if !req.StartDate.IsZero() {
		query = query.Where("appointment_time >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		// 结束日期包含当天
		endDate := req.EndDate.Add(24 * time.Hour)
		query = query.Where("appointment_time < ?", endDate)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页查询
	var appointments []models.Appointment
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Order("appointment_time DESC, id DESC").Find(&appointments)

	utils.SuccessWithPage(c, appointments, total, req.Page, req.PageSize)
}

// GetAppointment 获取预约详情
func GetAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var appointment models.Appointment
	if err := models.DB.First(&appointment, id).Error; err != nil {
		utils.ErrorNotFound(c, "预约不存在")
		return
	}

	utils.Success(c, appointment)
}

// CreateAppointment 创建预约（小程序端调用）
func CreateAppointment(c *gin.Context) {
	var req AppointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 上门服务必须填写地址
	if req.ServiceType == "home" && req.Address == "" {
		utils.ErrorBadRequest(c, "上门服务必须填写地址")
		return
	}

	// 创建预约
	appointment := models.Appointment{
		ShopID:          req.ShopID,
		Openid:          req.Openid,
		MetalType:       req.MetalType,
		ServiceType:     req.ServiceType,
		AppointmentTime: req.AppointmentTime,
		Name:            req.Name,
		Phone:           req.Phone,
		Address:         req.Address,
		Note:            req.Note,
		Status:          models.AppointmentStatusPending,
	}

	if err := models.DB.Create(&appointment).Error; err != nil {
		utils.ErrorInternalServer(c, "创建预约失败")
		return
	}

	utils.SuccessWithMessage(c, appointment, "预约成功，工作人员会尽快联系您")
}

// UpdateAppointment 更新预约
func UpdateAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req AppointmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查预约是否存在
	var appointment models.Appointment
	if err := models.DB.First(&appointment, id).Error; err != nil {
		utils.ErrorNotFound(c, "预约不存在")
		return
	}

	// 获取当前用户ID
	userID := c.GetUint("user_id")

	// 构建更新数据
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Note != "" {
		updates["note"] = req.Note
	}
	if req.AdminRemark != "" {
		updates["admin_remark"] = req.AdminRemark
	}
	if req.ShopID != nil {
		updates["shop_id"] = *req.ShopID
	}
	if req.AppointmentTime != nil {
		updates["appointment_time"] = *req.AppointmentTime
	}

	// 状态变更处理
	if req.Status != "" && req.Status != appointment.Status {
		updates["status"] = req.Status
		now := time.Now()

		switch req.Status {
		case models.AppointmentStatusConfirmed:
			updates["confirmed_at"] = now
			updates["handler_id"] = userID
		case models.AppointmentStatusCompleted:
			updates["completed_at"] = now
			updates["handler_id"] = userID
		case models.AppointmentStatusCancelled:
			updates["cancelled_at"] = now
		}
	}

	if err := models.DB.Model(&appointment).Updates(updates).Error; err != nil {
		utils.ErrorInternalServer(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeleteAppointment 删除预约
func DeleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var appointment models.Appointment
	if err := models.DB.First(&appointment, id).Error; err != nil {
		utils.ErrorNotFound(c, "预约不存在")
		return
	}

	// 删除预约
	if err := models.DB.Delete(&appointment).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "删除成功")
}

// GetAppointmentStats 获取预约统计
func GetAppointmentStats(c *gin.Context) {
	// 今日预约数
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	models.DB.Model(&models.Appointment{}).
		Where("DATE(appointment_time) = ?", today).
		Count(&todayCount)

	// 待处理预约数
	var pendingCount int64
	models.DB.Model(&models.Appointment{}).
		Where("status = ?", models.AppointmentStatusPending).
		Count(&pendingCount)

	// 本月预约数
	firstDay := time.Now().Format("2006-01-01")
	var monthCount int64
	models.DB.Model(&models.Appointment{}).
		Where("appointment_time >= ?", firstDay).
		Count(&monthCount)

	stats := map[string]interface{}{
		"today_count":   todayCount,
		"pending_count": pendingCount,
		"month_count":   monthCount,
	}

	utils.Success(c, stats)
}

