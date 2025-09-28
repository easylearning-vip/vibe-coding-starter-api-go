package model

import (
	"database/sql"

	"gorm.io/gorm"
)

// Part4PracticeSetItem references a single question within a talk (question_index 1..3)
type Part4PracticeSetItem struct {
	BaseModel

	SetID          uint           `json:"set_id" gorm:"column:set_id;not null;index"`
	TalkID         uint           `json:"talk_id" gorm:"column:talk_id;not null;index"`
	QuestionIndex  int32          `json:"question_index" gorm:"column:question_index;type:int;not null;index"`
	OrderIndex     int32          `json:"order_index" gorm:"column:order_index;type:int;default:0"`
	SelectedAnswer sql.NullString `json:"selected_answer" gorm:"column:selected_answer;type:char(1)"`
	IsCorrect      bool           `json:"is_correct" gorm:"column:is_correct;not null;default:false;index"`
}

func (Part4PracticeSetItem) TableName() string                 { return "part4_practice_set_items" }
func (m *Part4PracticeSetItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part4PracticeSetItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }
