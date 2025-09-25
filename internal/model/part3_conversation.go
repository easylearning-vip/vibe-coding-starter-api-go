package model

import (
	"database/sql"

	"gorm.io/gorm"
)

// Part3Conversation Part3Conversation模型
type Part3Conversation struct {
	BaseModel

	TestId int32 `json:"test_id" gorm:"column:test_id;type:int;not null"` // TestId 32位整数

	ConversationNumber int32 `json:"conversation_number" gorm:"column:conversation_number;type:int;not null"` // ConversationNumber 32位整数

	Title sql.NullString `json:"title" gorm:"column:title;type:varchar(255)"` // Title 可空字符串

	Content string `json:"content" gorm:"column:content;type:text;not null"` // Content 字符串

	Question1 string `json:"question1" gorm:"column:question1;type:text;not null"` // Question1 字符串

	Question2 string `json:"question2" gorm:"column:question2;type:text;not null"` // Question2 字符串

	Question3 string `json:"question3" gorm:"column:question3;type:text;not null"` // Question3 字符串

	ScenarioId sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int"` // ScenarioId 可空自定义类型

	DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int"` // DifficultyLevelId 可空自定义类型

	// 关联的答案选项
	AnswerOptions []Part3AnswerOption `json:"answer_options,omitempty" gorm:"foreignKey:ConversationId;references:ID"`
}

// TableName 获取表名
func (Part3Conversation) TableName() string {
	return "part3_conversations"
}

// BeforeCreate GORM 钩子：创建前
func (part3Conversation *Part3Conversation) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part3Conversation.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (part3Conversation *Part3Conversation) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part3Conversation.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
