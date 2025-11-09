package v1

import (
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PriceListRequest 价格列表请求参数
type PriceListRequest struct {
	ShopID *uint  `form:"shop_id"`
	Code   string `form:"code"`
	Name   string `form:"name"`
	Status *int8  `form:"status"`
}

// PriceCreateRequest 创建价格请求参数
type PriceCreateRequest struct {
	ShopID        *uint   `json:"shop_id"`
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Subtitle      string  `json:"subtitle"`
	Icon          string  `json:"icon"`
	IconColor     string  `json:"icon_color"`
	BasePrice     float64 `json:"base_price" binding:"required,gt=0"`
	BuyPriceDiff  float64 `json:"buy_price_diff"`  // 回购差价（可为负）
	SellPriceDiff float64 `json:"sell_price_diff"` // 销售差价（可为正）
	Sort          int     `json:"sort"`
	Status        int8    `json:"status"`
}

// PriceUpdateRequest 更新价格请求参数
type PriceUpdateRequest struct {
	ShopID        *uint    `json:"shop_id"`
	Name          string   `json:"name"`
	Subtitle      string   `json:"subtitle"`
	Icon          string   `json:"icon"`
	IconColor     string   `json:"icon_color"`
	BasePrice     *float64 `json:"base_price"`
	BuyPriceDiff  *float64 `json:"buy_price_diff"`
	SellPriceDiff *float64 `json:"sell_price_diff"`
	Sort          *int     `json:"sort"`
	Status        *int8    `json:"status"`
}

// GetPriceList 获取价格列表
func GetPriceList(c *gin.Context) {
	var req PriceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 构建查询
	query := models.DB.Model(&models.Price{})

	if req.ShopID != nil {
		query = query.Where("shop_id = ?", *req.ShopID)
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 查询价格列表
	var prices []models.Price
	query.Order("sort ASC, id ASC").Find(&prices)

	// 构建带计算字段的结果
	var result []models.PriceWithCalculated
	for _, price := range prices {
		result = append(result, models.PriceWithCalculated{
			Price:     price,
			BuyPrice:  price.GetBuyPrice(),
			SellPrice: price.GetSellPrice(),
		})
	}

	utils.Success(c, result)
}

// GetPrice 获取价格详情
func GetPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var price models.Price
	if err := models.DB.First(&price, id).Error; err != nil {
		utils.ErrorNotFound(c, "价格不存在")
		return
	}

	result := models.PriceWithCalculated{
		Price:     price,
		BuyPrice:  price.GetBuyPrice(),
		SellPrice: price.GetSellPrice(),
	}

	utils.Success(c, result)
}

// CreatePrice 创建价格
func CreatePrice(c *gin.Context) {
	var req PriceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查代码是否已存在（同一店铺下）
	var count int64
	query := models.DB.Model(&models.Price{}).Where("code = ?", req.Code)
	if req.ShopID != nil {
		query = query.Where("shop_id = ?", *req.ShopID)
	} else {
		query = query.Where("shop_id IS NULL")
	}
	query.Count(&count)

	if count > 0 {
		utils.ErrorBadRequest(c, "该标识已存在")
		return
	}

	// 创建价格
	price := models.Price{
		ShopID:        req.ShopID,
		Code:          req.Code,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		Icon:          req.Icon,
		IconColor:     req.IconColor,
		BasePrice:     req.BasePrice,
		BuyPriceDiff:  req.BuyPriceDiff,
		SellPriceDiff: req.SellPriceDiff,
		Sort:          req.Sort,
		Status:        req.Status,
		UpdatedAt:     time.Now(),
	}

	if err := models.DB.Create(&price).Error; err != nil {
		utils.ErrorInternalServer(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, price, "创建成功")
}

// UpdatePrice 更新价格
func UpdatePrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	var req PriceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	// 检查价格是否存在
	var price models.Price
	if err := models.DB.First(&price, id).Error; err != nil {
		utils.ErrorNotFound(c, "价格不存在")
		return
	}

	// 更新价格信息
	updates := map[string]interface{}{
		"name":      req.Name,
		"subtitle":  req.Subtitle,
		"icon":      req.Icon,
		"icon_color": req.IconColor,
		"updated_at": time.Now(),
	}
	if req.ShopID != nil {
		updates["shop_id"] = *req.ShopID
	}
	if req.BasePrice != nil {
		updates["base_price"] = *req.BasePrice
	}
	if req.BuyPriceDiff != nil {
		updates["buy_price_diff"] = *req.BuyPriceDiff
	}
	if req.SellPriceDiff != nil {
		updates["sell_price_diff"] = *req.SellPriceDiff
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&price).Updates(updates).Error; err != nil {
		utils.ErrorInternalServer(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "更新成功")
}

// DeletePrice 删除价格
func DeletePrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorBadRequest(c, "参数错误")
		return
	}

	if err := models.DB.Delete(&models.Price{}, id).Error; err != nil {
		utils.ErrorInternalServer(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, nil, "删除成功")
}

// SyncBasePrice 同步基础价格（从第三方API获取）
func SyncBasePrice(c *gin.Context) {
	// TODO: 从第三方API获取最新金价
	// 这里先返回模拟数据

	utils.SuccessWithMessage(c, nil, "基础价格同步成功")
}















