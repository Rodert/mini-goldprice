package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OperationLogListRequest 操作日志列表请求参数
type OperationLogListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username  string `form:"username"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	IP        string `form:"ip"`
	StartDate string `form:"start_date"` // 开始日期 YYYY-MM-DD
	EndDate   string `form:"end_date"`   // 结束日期 YYYY-MM-DD
}

// GetOperationLogList 获取操作日志列表
func GetOperationLogList(c *gin.Context) {
	var req OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询
	query := models.DB.Model(&models.OperationLog{})

	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.IP != "" {
		query = query.Where("ip = ?", req.IP)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(created_at) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(created_at) <= ?", req.EndDate)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页查询
	var logs []models.OperationLog
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&logs)

	// 返回结果
	utils.SuccessWithPage(c, logs, total, req.Page, req.PageSize)
}

// GetOperationLog 获取操作日志详情
func GetOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var log models.OperationLog
	if err := models.DB.First(&log, id).Error; err != nil {
		utils.ErrorNotFound(c, "日志不存在")
		return
	}

	utils.Success(c, log)
}

// DeleteOperationLog 删除操作日志
func DeleteOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	if err := models.DB.Delete(&models.OperationLog{}, id).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "删除成功")
}

// ClearOperationLogs 清空操作日志（保留最近N天）
func ClearOperationLogs(c *gin.Context) {
	type ClearRequest struct {
		Days int `json:"days" binding:"required,min=1"` // 保留天数
	}

	var req ClearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 删除N天之前的日志
	result := models.DB.Exec("DELETE FROM operation_logs WHERE created_at < datetime('now', '-"+strconv.Itoa(req.Days)+" days')")

	if result.Error != nil {
		utils.ErrorInternalServer(c, "清空失败")
		return
	}

	utils.SuccessWithMessage(c, map[string]interface{}{
		"deleted_count": result.RowsAffected,
	}, "清空成功")
}

// GetOperationLogStats 获取操作日志统计
func GetOperationLogStats(c *gin.Context) {
	type ModuleStat struct {
		Module string `json:"module"`
		Count  int64  `json:"count"`
	}

	type ActionStat struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}

	type UserStat struct {
		Username string `json:"username"`
		Count    int64  `json:"count"`
	}

	type Stats struct {
		TotalCount      int64            `json:"total_count"`
		TodayCount      int64            `json:"today_count"`
		ModuleStats     []ModuleStat     `json:"module_stats"`
		ActionStats     []ActionStat     `json:"action_stats"`
		RecentLogs      []models.OperationLog `json:"recent_logs"`
		TopUsers        []UserStat       `json:"top_users"`
	}

	var stats Stats

	// 总数
	models.DB.Model(&models.OperationLog{}).Count(&stats.TotalCount)

	// 今日数量
	models.DB.Model(&models.OperationLog{}).
		Where("DATE(created_at) = DATE('now')").
		Count(&stats.TodayCount)

	// 按模块统计
	models.DB.Model(&models.OperationLog{}).
		Select("module, COUNT(*) as count").
		Group("module").
		Order("count DESC").
		Limit(10).
		Scan(&stats.ModuleStats)

	// 按操作统计
	models.DB.Model(&models.OperationLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&stats.ActionStats)

	// 最近操作
	models.DB.Model(&models.OperationLog{}).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentLogs)

	// 操作最多的用户
	models.DB.Model(&models.OperationLog{}).
		Select("username, COUNT(*) as count").
		Where("username != ''").
		Group("username").
		Order("count DESC").
		Limit(10).
		Scan(&stats.TopUsers)

	utils.Success(c, stats)
}

// ExportOperationLogs 导出操作日志
func ExportOperationLogs(c *gin.Context) {
	var req OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 构建查询（不分页）
	query := models.DB.Model(&models.OperationLog{})

	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.IP != "" {
		query = query.Where("ip = ?", req.IP)
	}
	if req.StartDate != "" {
		query = query.Where("DATE(created_at) >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("DATE(created_at) <= ?", req.EndDate)
	}

	var logs []models.OperationLog
	query.Order("created_at DESC").Limit(10000).Find(&logs) // 最多导出1万条

	// TODO: 实现CSV或Excel导出
	utils.Success(c, logs)
}

