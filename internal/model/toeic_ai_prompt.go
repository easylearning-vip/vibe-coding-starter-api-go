package model

import (
	"gorm.io/gorm"
)

// ToeicAiPrompt TOEIC AI提示词模型
type ToeicAiPrompt struct {
	BaseModel

	Part2Prompt string `json:"part2_prompt" gorm:"column:part2_prompt;type:text;not null"` // Part2Prompt TOEIC Part2 AI提示词
	Part3Prompt string `json:"part3_prompt" gorm:"column:part3_prompt;type:text;not null"` // Part3Prompt TOEIC Part3 AI提示词
	Part4Prompt string `json:"part4_prompt" gorm:"column:part4_prompt;type:text;not null"` // Part4Prompt TOEIC Part4 AI提示词
}

// TableName 获取表名
func (ToeicAiPrompt) TableName() string {
	return "toeic_ai_prompts"
}

// BeforeCreate GORM 钩子：创建前
func (toeicAiPrompt *ToeicAiPrompt) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := toeicAiPrompt.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (toeicAiPrompt *ToeicAiPrompt) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := toeicAiPrompt.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
