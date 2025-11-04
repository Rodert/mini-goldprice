package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShopListRequest 店铺列表请求参数
type ShopListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   *int8  `form:"status"`
}

// ShopCreateRequest 创建店铺请求参数
type ShopCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Address     string  `json:"address"`
	Phone       string  `json:"phone"`
	Mobile      string  `json:"mobile"`
	Hours       string  `json:"hours"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	Status      int8    `json:"status"`
	Sort        int     `json:"sort"`
}

// ShopUpdateRequest 更新店铺请求参数
type ShopUpdateRequest struct {
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	Phone       string   `json:"phone"`
	Mobile      string   `json:"mobile"`
	Hours       string   `json:"hours"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Description string   `json:"description"`
	Status      *int8    `json:"status"`
	Sort        *int     `json:"sort"`
}

// GetShopList 获取店铺列表
func GetShopList(c *gin.Context) {
	var req ShopListRequest
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
	query := models.DB.Model(&models.Shop{})

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
	var shops []models.Shop
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Order("sort ASC, id DESC").Find(&shops)

	utils.SuccessWithPage(c, shops, total, req.Page, req.PageSize)
}

// GetAllShops 获取所有店铺（不分页，用于下拉选择）
func GetAllShops(c *gin.Context) {
	var shops []models.Shop
	models.DB.Where("status = ?", 1).Order("sort ASC").Find(&shops)
	utils.Success(c, shops)
}

// GetShop 获取店铺详情
func GetShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var shop models.Shop
	if err := models.DB.First(&shop, id).Error; err != nil {
		utils.ErrorNotFound(c, "店铺不存在")
		return
	}

	utils.Success(c, shop)
}

// CreateShop 创建店铺
func CreateShop(c *gin.Context) {
	var req ShopCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查店铺代码是否已存在
	var count int64
	models.DB.Model(&models.Shop{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		utils.ErrorBadRequest(c, "店铺代码已存在")
		return
	}

	// 创建店铺
	shop := models.Shop{
		Name:        req.Name,
		Code:        req.Code,
		Address:     req.Address,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		Hours:       req.Hours,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
	}

	if err := models.DB.Create(&shop).Error; err != nil {
		utils.ErrorInternalServer(c, "创建店铺失败")
		return
	}

	utils.SuccessWithMessage(c, shop, "创建成功")
}

// UpdateShop 更新店铺
func UpdateShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req ShopUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查店铺是否存在
	var shop models.Shop
	if err := models.DB.First(&shop, id).Error; err != nil {
		utils.ErrorNotFound(c, "店铺不存在")
		return
	}

	// 更新店铺信息
	updates := map[string]interface{}{
		"name":        req.Name,
		"address":     req.Address,
		"phone":       req.Phone,
		"mobile":      req.Mobile,
		"hours":       req.Hours,
		"description": req.Description,
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}

	if err := models.DB.Model(&shop).Updates(updates).Error; err != nil {
		utils.ErrorInternalServer(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeleteShop 删除店铺
func DeleteShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var shop models.Shop
	if err := models.DB.First(&shop, id).Error; err != nil {
		utils.ErrorNotFound(c, "店铺不存在")
		return
	}

	// 检查是否有价格关联
	var priceCount int64
	models.DB.Model(&models.Price{}).Where("shop_id = ?", id).Count(&priceCount)
	if priceCount > 0 {
		utils.ErrorBadRequest(c, "该店铺有关联的价格数据，不能删除")
		return
	}

	// 删除店铺
	if err := models.DB.Delete(&shop).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "删除成功")
}

