package model

import (
	"gorm.io/gorm"
)

// DictCategory 数据字典分类模型
type DictCategory struct {
	BaseModel
	Code        string     `gorm:"uniqueIndex;size:50;not null" json:"code" validate:"required,max=50"`
	Name        string     `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Description string     `gorm:"type:text" json:"description"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	Items       []DictItem `gorm:"foreignKey:CategoryCode;references:Code" json:"items,omitempty"`
}

// DictItem 数据字典项模型
type DictItem struct {
	BaseModel
	CategoryCode string        `gorm:"size:50;not null;index:idx_category_key,unique" json:"category_code" validate:"required,max=50"`
	ItemKey      string        `gorm:"size:50;not null;index:idx_category_key,unique" json:"item_key" validate:"required,max=50"`
	ItemValue    string        `gorm:"size:200;not null" json:"item_value" validate:"required,max=200"`
	Description  string        `gorm:"type:text" json:"description"`
	SortOrder    int           `gorm:"default:0" json:"sort_order"`
	IsActive     *bool         `gorm:"default:true" json:"is_active"`
	Category     *DictCategory `gorm:"foreignKey:CategoryCode;references:Code" json:"category,omitempty"`
}

// TableName 获取DictCategory表名
func (DictCategory) TableName() string {
	return "dict_categories"
}

// TableName 获取DictItem表名
func (DictItem) TableName() string {
	return "dict_items"
}

// BeforeCreate GORM 钩子：DictCategory创建前
func (dc *DictCategory) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := dc.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate GORM 钩子：DictCategory更新前
func (dc *DictCategory) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := dc.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}
	return nil
}

// BeforeCreate GORM 钩子：DictItem创建前
func (di *DictItem) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := di.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// BeforeUpdate GORM 钩子：DictItem更新前
func (di *DictItem) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := di.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}
	return nil
}

// IsEnabled 检查字典项是否启用
func (di *DictItem) IsEnabled() bool {
	return di.IsActive != nil && *di.IsActive
}

// GetFullKey 获取完整的键值（分类代码.项键值）
func (di *DictItem) GetFullKey() string {
	return di.CategoryCode + "." + di.ItemKey
}

// 数据字典分类代码常量
const (
	DictCategoryArticleStatus = "article_status"
	DictCategoryCommentStatus = "comment_status"
	DictCategoryUserRole      = "user_role"
	DictCategoryUserStatus    = "user_status"
	DictCategoryStorageType   = "storage_type"
)

// 数据字典项键值常量
const (
	// 文章状态
	DictItemArticleStatusDraft     = "draft"
	DictItemArticleStatusPublished = "published"
	DictItemArticleStatusArchived  = "archived"

	// 评论状态
	DictItemCommentStatusPending  = "pending"
	DictItemCommentStatusApproved = "approved"
	DictItemCommentStatusRejected = "rejected"

	// 用户角色
	DictItemUserRoleAdmin = "admin"
	DictItemUserRoleUser  = "user"

	// 用户状态
	DictItemUserStatusActive   = "active"
	DictItemUserStatusInactive = "inactive"
	DictItemUserStatusBanned   = "banned"

	// 存储类型
	DictItemStorageTypeLocal = "local"
	DictItemStorageTypeS3    = "s3"
	DictItemStorageTypeOSS   = "oss"
)
