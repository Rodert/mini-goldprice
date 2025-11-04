package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Token     string           `json:"token"`
	UserInfo  *UserInfo        `json:"user_info"`
	MenuList  []*models.MenuTree `json:"menu_list"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
	Roles    []string `json:"roles"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 查询用户
	var user models.AdminUser
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.ErrorBadRequest(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorBadRequest(c, "用户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		utils.ErrorForbidden(c, "账号已被禁用")
		return
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ErrorInternalServer(c, "生成Token失败")
		return
	}

	// 更新最后登录信息
	models.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_time": time.Now(),
		"last_login_ip":   c.ClientIP(),
	})

	// 获取用户角色
	roles := getUserRoles(user.ID)

	// 获取用户菜单
	menuList := getUserMenus(user.ID)

	// 返回登录信息
	utils.Success(c, LoginResponse{
		Token: token,
		UserInfo: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Roles:    roles,
		},
		MenuList: menuList,
	})
}

// Logout 登出
func Logout(c *gin.Context) {
	// JWT 是无状态的，登出只需要前端删除 Token
	utils.SuccessWithMessage(c, nil, "登出成功")
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.AdminUser
	if err := models.DB.First(&user, userID).Error; err != nil {
		utils.ErrorNotFound(c, "用户不存在")
		return
	}

	// 获取用户角色
	roles := getUserRoles(userID)

	// 获取用户菜单
	menuList := getUserMenus(userID)

	utils.Success(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"avatar":    user.Avatar,
		"roles":     roles,
		"menus":     menuList,
	})
}

// GetUserMenus 获取当前用户菜单
func GetUserMenus(c *gin.Context) {
	userID := c.GetUint("user_id")
	menuList := getUserMenus(userID)
	utils.Success(c, menuList)
}

// getUserRoles 获取用户角色列表
func getUserRoles(userID uint) []string {
	var roleCodes []string

	// 查询用户的角色ID
	var roleIDs []uint
	models.DB.Model(&models.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return []string{}
	}

	// 查询角色代码
	models.DB.Model(&models.Role{}).
		Where("id IN ? AND status = ?", roleIDs, 1).
		Pluck("code", &roleCodes)

	return roleCodes
}

// getUserMenus 获取用户菜单树
func getUserMenus(userID uint) []*models.MenuTree {
	// 查询用户的角色ID
	var roleIDs []uint
	models.DB.Model(&models.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return []*models.MenuTree{}
	}

	// 查询角色关联的菜单ID
	var menuIDs []uint
	models.DB.Model(&models.RoleMenu{}).
		Where("role_id IN ?", roleIDs).
		Distinct("menu_id").
		Pluck("menu_id", &menuIDs)

	if len(menuIDs) == 0 {
		return []*models.MenuTree{}
	}

	// 查询菜单（只返回目录和菜单，不返回按钮）
	var menus []models.Menu
	models.DB.Where("id IN ?", menuIDs).
		Where("type IN ?", []int8{1, 2}).
		Where("status = ? AND visible = ?", 1, 1).
		Order("sort ASC, id ASC").
		Find(&menus)

	// 构建菜单树
	return BuildMenuTree(menus, 0)
}


