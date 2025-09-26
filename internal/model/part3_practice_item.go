package model

import (
    "database/sql"

    "gorm.io/gorm"
)

// Part3PracticeItem records a user's practice for a single question inside a Part3 conversation
// QuestionIndex is 1..3 to indicate which question within the conversation
// SelectedOptionId references part3_answer_options.id if available
// ScenarioId/DifficultyLevelId capture the selected filters for stats.
type Part3PracticeItem struct {
    BaseModel

    UserID            uint          `json:"user_id" gorm:"column:user_id;not null;index"`
    ConversationID    uint          `json:"conversation_id" gorm:"column:conversation_id;not null;index"`
    QuestionIndex     int32         `json:"question_index" gorm:"column:question_index;type:int;not null;index"`

    ScenarioId        sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int;index"`
    DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int;index"`

    SelectedOptionId  sql.NullInt32 `json:"selected_option_id" gorm:"column:selected_option_id;type:int"`
    IsCorrect         bool          `json:"is_correct" gorm:"column:is_correct;not null;default:false;index"`
}

func (Part3PracticeItem) TableName() string { return "part3_practice_items" }

func (m *Part3PracticeItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part3PracticeItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

