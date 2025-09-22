package model

import (

	"database/sql"


	"gorm.io/gorm"
)

// Part2Question Part2Question模型
type Part2Question struct {
	BaseModel

	TestId int32 `json:"test_id" gorm:"column:test_id;type:int;not null"` // TestId 32位整数

	QuestionNumber int32 `json:"question_number" gorm:"column:question_number;type:int;not null"` // QuestionNumber 32位整数

	QuestionText string `json:"question_text" gorm:"column:question_text;type:text;not null"` // QuestionText 字符串

	OptionA string `json:"option_a" gorm:"column:option_a;type:text;not null"` // OptionA 字符串

	OptionB string `json:"option_b" gorm:"column:option_b;type:text;not null"` // OptionB 字符串

	OptionC string `json:"option_c" gorm:"column:option_c;type:text;not null"` // OptionC 字符串

	CorrectAnswer sql.NullString `json:"correct_answer" gorm:"column:correct_answer;type:char(1)"` // CorrectAnswer 可空字符串

	ScenarioId sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int"` // ScenarioId 可空自定义类型

	DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int"` // DifficultyLevelId 可空自定义类型

}

// TableName 获取表名
func (Part2Question) TableName() string {
	return "part2_questions"
}

// BeforeCreate GORM 钩子：创建前
func (part2Question *Part2Question) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part2Question.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 在这里添加特定的创建前逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (part2Question *Part2Question) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := part2Question.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 在这里添加特定的更新前逻辑
	return nil
}
