package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleListRequest 角色列表请求参数
type RoleListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   *int8  `form:"status"`
}

// RoleCreateRequest 创建角色请求参数
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int8   `json:"status"`
	MenuIDs     []uint `json:"menu_ids"` // 菜单ID列表
}

// RoleUpdateRequest 更新角色请求参数
type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Sort        *int   `json:"sort"`
	Status      *int8  `json:"status"`
	MenuIDs     []uint `json:"menu_ids"`
}

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	var req RoleListRequest
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
	query := models.DB.Model(&models.Role{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页查询
	var roles []models.Role
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Order("sort ASC, id DESC").Find(&roles)

	utils.SuccessWithPage(c, roles, total, req.Page, req.PageSize)
}

// GetAllRoles 获取所有角色（不分页，用于下拉选择）
func GetAllRoles(c *gin.Context) {
	var roles []models.Role
	models.DB.Where("status = ?", 1).Order("sort ASC").Find(&roles)
	utils.Success(c, roles)
}

// GetRole 获取角色详情
func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := models.DB.First(&role, id).Error; err != nil {
		utils.ErrorNotFound(c, "角色不存在")
		return
	}

	// 查询角色菜单
	var menuIDs []uint
	models.DB.Model(&models.RoleMenu{}).
		Where("role_id = ?", id).
		Pluck("menu_id", &menuIDs)

	result := map[string]interface{}{
		"role":     role,
		"menu_ids": menuIDs,
	}

	utils.Success(c, result)
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查角色代码是否已存在
	var count int64
	models.DB.Model(&models.Role{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		utils.ErrorBadRequest(c, "角色代码已存在")
		return
	}

	// 创建角色
	role := models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := models.DB.Create(&role).Error; err != nil {
		utils.ErrorInternalServer(c, "创建角色失败")
		return
	}

	// 分配菜单
	if len(req.MenuIDs) > 0 {
		for _, menuID := range req.MenuIDs {
			roleMenu := models.RoleMenu{
				RoleID: role.ID,
				MenuID: menuID,
			}
			models.DB.Create(&roleMenu)
		}
	}

	utils.SuccessWithMessage(c, role, "创建成功")
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查角色是否存在
	var role models.Role
	if err := models.DB.First(&role, id).Error; err != nil {
		utils.ErrorNotFound(c, "角色不存在")
		return
	}

	// 不允许修改超级管理员
	if role.Code == "super_admin" {
		utils.ErrorForbidden(c, "不能修改超级管理员角色")
		return
	}

	// 更新角色信息（只更新非空字段）
	updates := make(map[string]interface{})
	
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// 只在有更新字段时才执行更新
	if len(updates) > 0 {
		if err := models.DB.Model(&role).Updates(updates).Error; err != nil {
			utils.ErrorInternalServer(c, "更新失败")
			return
		}
	}

	// 更新菜单
	if req.MenuIDs != nil {
		// 删除旧菜单
		models.DB.Where("role_id = ?", id).Delete(&models.RoleMenu{})

		// 添加新菜单
		for _, menuID := range req.MenuIDs {
			roleMenu := models.RoleMenu{
				RoleID: uint(id),
				MenuID: menuID,
			}
			models.DB.Create(&roleMenu)
		}
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := models.DB.First(&role, id).Error; err != nil {
		utils.ErrorNotFound(c, "角色不存在")
		return
	}

	// 不允许删除超级管理员角色
	if role.Code == "super_admin" {
		utils.ErrorForbidden(c, "不能删除超级管理员角色")
		return
	}

	// 检查是否有用户使用该角色
	var userCount int64
	models.DB.Model(&models.UserRole{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		utils.ErrorBadRequest(c, "该角色已被用户使用，不能删除")
		return
	}

	// 删除角色
	if err := models.DB.Delete(&role).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	// 删除角色菜单关联
	models.DB.Where("role_id = ?", id).Delete(&models.RoleMenu{})

	utils.SuccessWithMessage(c, nil, "删除成功")
}


