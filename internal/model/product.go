package model

import (
	"gorm.io/gorm"
)

// Product Product模型
type Product struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);uniqueIndex"` // Name 字符串

	Description string `json:"description" gorm:"column:description;type:varchar(255)"` // Description 字符串

	CategoryId uint `json:"category_id" gorm:"column:category_id;type:int unsigned"` // CategoryId 32位无符号整数

	Sku string `json:"sku" gorm:"column:sku;type:varchar(255)"` // Sku 字符串

	Price float64 `json:"price" gorm:"column:price;type:decimal(10,2)"` // Price 64位浮点数

	CostPrice float64 `json:"cost_price" gorm:"column:cost_price;type:decimal(10,2)"` // CostPrice 64位浮点数

	StockQuantity int `json:"stock_quantity" gorm:"column:stock_quantity;type:int"` // StockQuantity 32位整数

	MinStock int `json:"min_stock" gorm:"column:min_stock;type:int"` // MinStock 32位整数

	IsActive bool `json:"is_active" gorm:"column:is_active;type:boolean;default:false"` // IsActive 布尔值

	Weight float64 `json:"weight" gorm:"column:weight;type:decimal(10,2)"` // Weight 64位浮点数

	Dimensions string `json:"dimensions" gorm:"column:dimensions;type:varchar(255)"` // Dimensions 字符串

}

// TableName 获取表名
func (Product) TableName() string {
	return "products"
}

// BeforeCreate GORM 钩子：创建前
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := product.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (product *Product) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := product.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
