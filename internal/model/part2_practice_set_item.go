package model

import (
    "gorm.io/gorm"
)

// Part2PracticeSetItem is the detail record referencing a Part2 question
// Only holds the question ID information and link to the master set
type Part2PracticeSetItem struct {
    BaseModel

    SetID      uint `json:"set_id" gorm:"column:set_id;not null;index"`
    QuestionID uint `json:"question_id" gorm:"column:question_id;not null;index"`
    OrderIndex int32 `json:"order_index" gorm:"column:order_index;type:int;default:0"`
}

func (Part2PracticeSetItem) TableName() string { return "part2_practice_set_items" }
func (m *Part2PracticeSetItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part2PracticeSetItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

