package model

import (
    "database/sql"

    "gorm.io/gorm"
)

// Part2PracticeItem records a user's practice for a single Part2 question
// It references the original question by ID and stores the user's selection and correctness
// ScenarioId and DifficultyLevelId capture the selected filters when generating the practice
// so we can aggregate history and weaknesses by scenario/difficulty.
type Part2PracticeItem struct {
    BaseModel

    UserID           uint          `json:"user_id" gorm:"column:user_id;not null;index"`
    QuestionID       uint          `json:"question_id" gorm:"column:question_id;not null;index"`
    ScenarioId       sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int;index"`
    DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int;index"`

    SelectedAnswer   sql.NullString `json:"selected_answer" gorm:"column:selected_answer;type:char(1)"`
    IsCorrect        bool           `json:"is_correct" gorm:"column:is_correct;not null;default:false;index"`
}

func (Part2PracticeItem) TableName() string { return "part2_practice_items" }

func (m *Part2PracticeItem) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part2PracticeItem) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

