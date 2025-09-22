package model

import (

	"database/sql"


	"gorm.io/gorm"
)

// Part3AnswerOption Part3AnswerOption模型
type Part3AnswerOption struct {
	BaseModel

	ConversationId int32 `json:"conversation_id" gorm:"column:conversation_id;type:int;not null"` // ConversationId 32位整数

	QuestionNumber int32 `json:"question_number" gorm:"column:question_number;type:int;not null"` // QuestionNumber 32位整数

	OptionA string `json:"option_a" gorm:"column:option_a;type:text;not null"` // OptionA 字符串

	OptionB string `json:"option_b" gorm:"column:option_b;type:text;not null"` // OptionB 字符串

	OptionC string `json:"option_c" gorm:"column:option_c;type:text;not null"` // OptionC 字符串

	OptionD string `json:"option_d" gorm:"column:option_d;type:text;not null"` // OptionD 字符串

	CorrectAnswer sql.NullString `json:"correct_answer" gorm:"column:correct_answer;type:char(1)"` // CorrectAnswer 可空字符串

}

// TableName 获取表名
func (Part3AnswerOption) TableName() string {
	return "part3_answer_options"
}

// BeforeCreate GORM 钩子：创建前
func (part3AnswerOption *Part3AnswerOption) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part3AnswerOption.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (part3AnswerOption *Part3AnswerOption) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part3AnswerOption.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
