package model

import (
	"database/sql"

	"gorm.io/gorm"
)

// Test Test模型
type Test struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"` // Name 字符串

	Description sql.NullString `json:"description" gorm:"column:description;type:text"` // Description 可空字符串

}

// TableName 获取表名
func (Test) TableName() string {
	return "tests"
}

// BeforeCreate GORM 钩子：创建前
func (test *Test) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := test.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (test *Test) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := test.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
