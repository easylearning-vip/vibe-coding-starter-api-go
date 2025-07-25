package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 获取表名
func (BaseModel) TableName() string {
	return ""
}

// BeforeCreate GORM 钩子：创建前
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// IsDeleted 检查是否已删除
func (m *BaseModel) IsDeleted() bool {
	return m.DeletedAt.Valid
}

// GetID 获取 ID
func (m *BaseModel) GetID() uint {
	return m.ID
}

// GetCreatedAt 获取创建时间
func (m *BaseModel) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (m *BaseModel) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

// GetDeletedAt 获取删除时间
func (m *BaseModel) GetDeletedAt() *time.Time {
	if m.DeletedAt.Valid {
		return &m.DeletedAt.Time
	}
	return nil
}
