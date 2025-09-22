package model

import (

	"database/sql"


	"gorm.io/gorm"
)

// DifficultyLevel DifficultyLevel模型
type DifficultyLevel struct {
	BaseModel

	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"` // Name 字符串

	Description string `json:"description" gorm:"column:description;type:text;not null"` // Description 字符串

	Characteristics sql.NullString `json:"characteristics" gorm:"column:characteristics;type:longtext"` // Characteristics 可空字符串

}

// TableName 获取表名
func (DifficultyLevel) TableName() string {
	return "difficulty_levels"
}

// BeforeCreate GORM 钩子：创建前
func (difficultyLevel *DifficultyLevel) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := difficultyLevel.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (difficultyLevel *DifficultyLevel) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := difficultyLevel.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
