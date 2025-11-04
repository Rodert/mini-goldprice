package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserListRequest 用户列表请求参数
type UserListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username string `form:"username"`
	RealName string `form:"real_name"`
	Status   *int8  `form:"status"`
}

// UserCreateRequest 创建用户请求参数
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"role_ids"` // 角色ID列表
}

// UserUpdateRequest 更新用户请求参数
type UserUpdateRequest struct {
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   *int8  `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

// UserPasswordRequest 修改密码请求参数
type UserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	var req UserListRequest
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
	query := models.DB.Model(&models.AdminUser{})

	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.RealName != "" {
		query = query.Where("real_name LIKE ?", "%"+req.RealName+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页查询
	var users []models.AdminUser
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&users)

	// 返回结果
	utils.SuccessWithPage(c, users, total, req.Page, req.PageSize)
}

// GetUser 获取用户详情
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var user models.AdminUser
	if err := models.DB.First(&user, id).Error; err != nil {
		utils.ErrorNotFound(c, "用户不存在")
		return
	}

	// 查询用户角色
	var roleIDs []uint
	models.DB.Model(&models.UserRole{}).
		Where("user_id = ?", id).
		Pluck("role_id", &roleIDs)

	result := map[string]interface{}{
		"user":     user,
		"role_ids": roleIDs,
	}

	utils.Success(c, result)
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查用户名是否已存在
	var count int64
	models.DB.Model(&models.AdminUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		utils.ErrorBadRequest(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorInternalServer(c, "密码加密失败")
		return
	}

	// 创建用户
	user := models.AdminUser{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		utils.ErrorInternalServer(c, "创建用户失败")
		return
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			userRole := models.UserRole{
				UserID: user.ID,
				RoleID: roleID,
			}
			models.DB.Create(&userRole)
		}
	}

	utils.SuccessWithMessage(c, user, "创建成功")
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查用户是否存在
	var user models.AdminUser
	if err := models.DB.First(&user, id).Error; err != nil {
		utils.ErrorNotFound(c, "用户不存在")
		return
	}

	// 更新用户信息（只更新非空字段）
	updates := make(map[string]interface{})
	
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// 只在有更新字段时才执行更新
	if len(updates) > 0 {
		if err := models.DB.Model(&user).Updates(updates).Error; err != nil {
			utils.ErrorInternalServer(c, "更新失败")
			return
		}
	}

	// 更新角色
	if req.RoleIDs != nil {
		// 删除旧角色
		models.DB.Where("user_id = ?", id).Delete(&models.UserRole{})

		// 添加新角色
		for _, roleID := range req.RoleIDs {
			userRole := models.UserRole{
				UserID: uint(id),
				RoleID: roleID,
			}
			models.DB.Create(&userRole)
		}
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 不允许删除 ID 为 1 的用户（默认管理员）
	if id == 1 {
		utils.ErrorForbidden(c, "不能删除默认管理员")
		return
	}

	// 删除用户
	if err := models.DB.Delete(&models.AdminUser{}, id).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	// 删除用户角色关联
	models.DB.Where("user_id = ?", id).Delete(&models.UserRole{})

	utils.SuccessWithMessage(c, nil, "删除成功")
}

// UpdateUserPassword 修改用户密码
func UpdateUserPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req UserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorInternalServer(c, "密码加密失败")
		return
	}

	// 更新密码
	if err := models.DB.Model(&models.AdminUser{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error; err != nil {
		utils.ErrorInternalServer(c, "修改密码失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "修改成功")
}


