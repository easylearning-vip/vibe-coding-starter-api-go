package model

import (
    "database/sql"

    "gorm.io/gorm"
)

// Part4PracticeItem records a user's practice for a single question inside a Part4 talk
// QuestionIndex is 1..3 to indicate which question within the talk
// SelectedOptionId references part4_answer_options.id if available
// ScenarioId/DifficultyLevelId capture the selected filters for stats.
type Part4PracticeItem struct {
    BaseModel

    UserID            uint          `json:"user_id" gorm:"column:user_id;not null;index"`
    TalkID            uint          `json:"talk_id" gorm:"column:talk_id;not null;index"`
    QuestionIndex     int32         `json:"question_index" gorm:"column:question_index;type:int;not null;index"`

    ScenarioId        sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int;index"`
    DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int;index"`

    SelectedOptionId  sql.NullInt32 `json:"selected_option_id" gorm:"column:selected_option_id;type:int"`
    IsCorrect         bool          `json:"is_correct" gorm:"column:is_correct;not null;default:false;index"`
}

func (Part4PracticeItem) TableName() string { return "part4_practice_items" }

func (m *Part4PracticeItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part4PracticeItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

