package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuCreateRequest 创建菜单请求参数
type MenuCreateRequest struct {
	ParentID   int    `json:"parent_id"`
	Type       int8   `json:"type" binding:"required,oneof=1 2 3"`
	Name       string `json:"name" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Sort       int    `json:"sort"`
	Visible    int8   `json:"visible"`
	Status     int8   `json:"status"`
}

// MenuUpdateRequest 更新菜单请求参数
type MenuUpdateRequest struct {
	ParentID   *int   `json:"parent_id"`
	Type       *int8  `json:"type"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Sort       *int   `json:"sort"`
	Visible    *int8  `json:"visible"`
	Status     *int8  `json:"status"`
}

// GetMenuList 获取菜单列表（树形结构）
func GetMenuList(c *gin.Context) {
	var menus []models.Menu
	models.DB.Order("sort ASC, id ASC").Find(&menus)

	// 构建树形结构
	tree := BuildMenuTree(menus, 0)

	utils.Success(c, tree)
}

// GetMenuTree 获取菜单树（用于选择父菜单）
func GetMenuTree(c *gin.Context) {
	var menus []models.Menu
	// 只获取目录和菜单类型
	models.DB.Where("type IN ?", []int8{1, 2}).
		Where("status = ?", 1).
		Order("sort ASC, id ASC").
		Find(&menus)

	tree := BuildMenuTree(menus, 0)
	utils.Success(c, tree)
}

// GetMenu 获取菜单详情
func GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		utils.ErrorNotFound(c, "菜单不存在")
		return
	}

	utils.Success(c, menu)
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var req MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查菜单名称是否已存在
	var count int64
	models.DB.Model(&models.Menu{}).
		Where("name = ? AND parent_id = ?", req.Name, req.ParentID).
		Count(&count)
	if count > 0 {
		utils.ErrorBadRequest(c, "同级菜单名称已存在")
		return
	}

	// 创建菜单
	menu := models.Menu{
		ParentID:   req.ParentID,
		Type:       req.Type,
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
	}

	if err := models.DB.Create(&menu).Error; err != nil {
		utils.ErrorInternalServer(c, "创建菜单失败")
		return
	}

	utils.SuccessWithMessage(c, menu, "创建成功")
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查菜单是否存在
	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		utils.ErrorNotFound(c, "菜单不存在")
		return
	}

	// 不能将父菜单设置为自己或子菜单
	if req.ParentID != nil && *req.ParentID == int(id) {
		utils.ErrorBadRequest(c, "不能将父菜单设置为自己")
		return
	}

	// 更新菜单信息
	updates := map[string]interface{}{
		"name":       req.Name,
		"title":      req.Title,
		"icon":       req.Icon,
		"path":       req.Path,
		"component":  req.Component,
		"permission": req.Permission,
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Visible != nil {
		updates["visible"] = *req.Visible
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&menu).Updates(updates).Error; err != nil {
		utils.ErrorInternalServer(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查是否有子菜单
	var childCount int64
	models.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		utils.ErrorBadRequest(c, "该菜单有子菜单，不能删除")
		return
	}

	// 删除菜单
	if err := models.DB.Delete(&models.Menu{}, id).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	// 删除角色菜单关联
	models.DB.Where("menu_id = ?", id).Delete(&models.RoleMenu{})

	utils.SuccessWithMessage(c, nil, "删除成功")
}

// BuildMenuTree 构建菜单树（导出函数，供其他包使用）
func BuildMenuTree(menus []models.Menu, parentID int) []*models.MenuTree {
	var tree []*models.MenuTree

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := &models.MenuTree{
				ID:         menu.ID,
				ParentID:   menu.ParentID,
				Type:       menu.Type,
				Name:       menu.Name,
				Title:      menu.Title,
				Icon:       menu.Icon,
				Path:       menu.Path,
				Component:  menu.Component,
				Permission: menu.Permission,
				Sort:       menu.Sort,
			}
			// 递归查找子菜单
			node.Children = BuildMenuTree(menus, int(menu.ID))
			tree = append(tree, node)
		}
	}

	return tree
}


