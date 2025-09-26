package model

import (
    "gorm.io/gorm"
)

// Part3PracticeSetItem references a single question within a conversation (question_index 1..3)
type Part3PracticeSetItem struct {
    BaseModel

    SetID         uint `json:"set_id" gorm:"column:set_id;not null;index"`
    ConversationID uint `json:"conversation_id" gorm:"column:conversation_id;not null;index"`
    QuestionIndex int32 `json:"question_index" gorm:"column:question_index;type:int;not null;index"`
    OrderIndex    int32 `json:"order_index" gorm:"column:order_index;type:int;default:0"`
}

func (Part3PracticeSetItem) TableName() string { return "part3_practice_set_items" }
func (m *Part3PracticeSetItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part3PracticeSetItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

