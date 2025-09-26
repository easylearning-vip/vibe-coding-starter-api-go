package model

import (
    "database/sql"
    "gorm.io/gorm"
)

// Part3PracticeSet master for user's Part3 practice
type Part3PracticeSet struct {
    BaseModel
    UserID            uint          `json:"user_id" gorm:"column:user_id;not null;index"`
    ScenarioId        sql.NullInt32 `json:"scenario_id" gorm:"column:scenario_id;type:int;index"`
    DifficultyLevelId sql.NullInt32 `json:"difficulty_level_id" gorm:"column:difficulty_level_id;type:int;index"`

    TotalQuestions   int32   `json:"total_questions" gorm:"column:total_questions;not null;default:0"`
    CompletedCount   int32   `json:"completed_count" gorm:"column:completed_count;not null;default:0"`
    CorrectCount     int32   `json:"correct_count" gorm:"column:correct_count;not null;default:0"`
    Accuracy         float64 `json:"accuracy" gorm:"column:accuracy;type:double;not null;default:0"`
}

func (Part3PracticeSet) TableName() string { return "part3_practice_sets" }
func (m *Part3PracticeSet) BeforeCreate(tx *gorm.DB) error { return m.BaseModel.BeforeCreate(tx) }
func (m *Part3PracticeSet) BeforeUpdate(tx *gorm.DB) error { return m.BaseModel.BeforeUpdate(tx) }

