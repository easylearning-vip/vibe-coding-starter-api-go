package model

import (
	"time"
	"gorm.io/gorm"
)

// Product Product模型
type Product struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	Name string `json:"name" gorm:"column:name;type:varchar(255)"` // Name 字符串

	Price float64 `json:"price" gorm:"column:price;type:decimal(10,2)"` // Price 64位浮点数


	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

}

// TableName 指定表名
func (product *Product) TableName() string {
	return "products"
}

// BeforeCreate GORM钩子：创建前
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	// 在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (product *Product) BeforeUpdate(tx *gorm.DB) error {
	// 在这里添加更新前的逻辑
	return nil
}

// AfterCreate GORM钩子：创建后
func (product *Product) AfterCreate(tx *gorm.DB) error {
	// 在这里添加创建后的逻辑
	return nil
}

// AfterUpdate GORM钩子：更新后
func (product *Product) AfterUpdate(tx *gorm.DB) error {
	// 在这里添加更新后的逻辑
	return nil
}

// Validate 验证模型数据
func (product *Product) Validate() error {
	// 在这里添加验证逻辑
	return nil
}
