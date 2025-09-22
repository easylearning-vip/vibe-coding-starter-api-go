package model

import (

	"database/sql"


	"gorm.io/gorm"
)

// Scenario Scenario模型
type Scenario struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"` // Name 字符串

	Description string `json:"description" gorm:"column:description;type:text;not null"` // Description 字符串

	Examples sql.NullString `json:"examples" gorm:"column:examples;type:longtext"` // Examples 可空字符串

}

// TableName 获取表名
func (Scenario) TableName() string {
	return "scenarios"
}

// BeforeCreate GORM 钩子：创建前
func (scenario *Scenario) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := scenario.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (scenario *Scenario) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := scenario.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
