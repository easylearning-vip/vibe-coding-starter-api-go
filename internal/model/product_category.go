package model

import (


	"gorm.io/gorm"
)

// ProductCategory ProductCategory模型
type ProductCategory struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);uniqueIndex"` // Name 字符串

	Description string `json:"description" gorm:"column:description;type:varchar(255)"` // Description 字符串

	ParentId uint `json:"parent_id" gorm:"column:parent_id;type:int unsigned"` // ParentId 32位无符号整数

	SortOrder int `json:"sort_order" gorm:"column:sort_order;type:int"` // SortOrder 32位整数

	IsActive bool `json:"is_active" gorm:"column:is_active;type:boolean;default:false"` // IsActive 布尔值

}

// TableName 获取表名
func (ProductCategory) TableName() string {
	return "product_categories"
}

// BeforeCreate GORM 钩子：创建前
func (productCategory *ProductCategory) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := productCategory.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (productCategory *ProductCategory) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := productCategory.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
