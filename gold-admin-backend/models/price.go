package models

import "time"

// Price 价格模型
type Price struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ShopID        *uint     `gorm:"index" json:"shop_id"`                      // 店铺ID（可为空表示全局价格）
	Code          string    `gorm:"size:50;not null" json:"code"`              // 唯一标识（gold_9999）
	Name          string    `gorm:"size:50;not null" json:"name"`              // 品种名称
	Subtitle      string    `gorm:"size:100" json:"subtitle"`                  // 副标题
	Icon          string    `gorm:"size:10" json:"icon"`                       // 图标（Au, Ag）
	IconColor     string    `gorm:"size:20" json:"icon_color"`                 // 图标颜色
	BasePrice     float64   `gorm:"not null" json:"base_price"`                // 基础价格（元/克）
	BuyPriceDiff  float64   `gorm:"not null;default:0" json:"buy_price_diff"`  // 回购差价（可为负）
	SellPriceDiff float64   `gorm:"not null;default:0" json:"sell_price_diff"` // 销售差价（可为正）
	Sort          int       `gorm:"default:0" json:"sort"`                     // 排序
	Status        int8      `gorm:"default:1" json:"status"`                   // 状态 1:启用 0:禁用
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Price) TableName() string {
	return "prices"
}

// GetBuyPrice 计算回购价
func (p *Price) GetBuyPrice() float64 {
	return p.BasePrice + p.BuyPriceDiff
}

// GetSellPrice 计算销售价
func (p *Price) GetSellPrice() float64 {
	return p.BasePrice + p.SellPriceDiff
}

// PriceWithCalculated 带计算字段的价格（用于返回给前端）
type PriceWithCalculated struct {
	Price
	BuyPrice  float64 `json:"buy_price"`  // 回购价（计算字段）
	SellPrice float64 `json:"sell_price"` // 销售价（计算字段）
}








