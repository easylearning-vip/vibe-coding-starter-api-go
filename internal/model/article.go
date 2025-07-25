package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	BaseModel
	Title       string     `gorm:"size:200;not null" json:"title" validate:"required,max=200"`
	Slug        string     `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Content     string     `gorm:"type:text" json:"content" validate:"required"`
	Summary     string     `gorm:"column:excerpt;size:500" json:"summary" validate:"max=500"`
	CoverImage  string     `gorm:"column:featured_image;size:255" json:"cover_image" validate:"url"`
	Status      string     `gorm:"size:20;default:draft" json:"status" validate:"oneof=draft published archived"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	LikeCount   int        `gorm:"-" json:"like_count"`
	AuthorID    uint       `gorm:"not null" json:"author_id" validate:"required"`
	Author      User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CategoryID  *uint      `json:"category_id"`
	Category    *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag      `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Comments    []Comment  `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
	PublishedAt *time.Time `json:"published_at"`
}

// ArticleStatus 文章状态常量
const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
	ArticleStatusArchived  = "archived"
)

// Category 分类模型
type Category struct {
	BaseModel
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name" validate:"required,max=50"`
	Slug        string    `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Description string    `gorm:"size:200" json:"description" validate:"max=200"`
	Color       string    `gorm:"size:7" json:"color" validate:"hexcolor"`
	Icon        string    `gorm:"size:50" json:"icon"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	Articles    []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// Tag 标签模型
type Tag struct {
	BaseModel
	Name        string    `gorm:"uniqueIndex;size:30;not null" json:"name" validate:"required,max=30"`
	Slug        string    `gorm:"uniqueIndex;size:30;not null" json:"slug"`
	Description string    `gorm:"size:100" json:"description" validate:"max=100"`
	Color       string    `gorm:"size:7" json:"color" validate:"hexcolor"`
	Articles    []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required"`
	Status    string    `gorm:"size:20;default:pending" json:"status" validate:"oneof=pending approved rejected"`
	ArticleID uint      `gorm:"not null" json:"article_id" validate:"required"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	AuthorID  uint      `gorm:"not null" json:"author_id" validate:"required"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	ParentID  *uint     `json:"parent_id"`
	Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// CommentStatus 评论状态常量
const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
)

// TableName 获取表名
func (Article) TableName() string {
	return "articles"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate GORM 钩子：创建前
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := a.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 生成 slug
	if a.Slug == "" {
		a.Slug = generateSlug(a.Title)
	}

	// 设置默认状态
	if a.Status == "" {
		a.Status = ArticleStatusDraft
	}

	// 如果发布，设置发布时间
	if a.Status == ArticleStatusPublished && a.PublishedAt == nil {
		now := time.Now()
		a.PublishedAt = &now
	}

	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if err := c.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	if c.Slug == "" {
		c.Slug = generateSlug(c.Name)
	}

	return nil
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if err := t.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	if t.Slug == "" {
		t.Slug = generateSlug(t.Name)
	}

	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if err := c.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	if c.Status == "" {
		c.Status = CommentStatusPending
	}

	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	if err := a.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 如果状态改为发布且没有发布时间，设置发布时间
	if tx.Statement.Changed("Status") && a.Status == ArticleStatusPublished && a.PublishedAt == nil {
		now := time.Now()
		a.PublishedAt = &now
	}

	return nil
}

// IsPublished 检查是否已发布
func (a *Article) IsPublished() bool {
	return a.Status == ArticleStatusPublished
}

// IsDraft 检查是否为草稿
func (a *Article) IsDraft() bool {
	return a.Status == ArticleStatusDraft
}

// IsArchived 检查是否已归档
func (a *Article) IsArchived() bool {
	return a.Status == ArticleStatusArchived
}

// IncrementViewCount 增加浏览次数
func (a *Article) IncrementViewCount() {
	a.ViewCount++
}

// IncrementLikeCount 增加点赞次数
func (a *Article) IncrementLikeCount() {
	a.LikeCount++
}

// DecrementLikeCount 减少点赞次数
func (a *Article) DecrementLikeCount() {
	if a.LikeCount > 0 {
		a.LikeCount--
	}
}

// IsApproved 检查评论是否已批准
func (c *Comment) IsApproved() bool {
	return c.Status == CommentStatusApproved
}

// IsPending 检查评论是否待审核
func (c *Comment) IsPending() bool {
	return c.Status == CommentStatusPending
}

// IsRejected 检查评论是否被拒绝
func (c *Comment) IsRejected() bool {
	return c.Status == CommentStatusRejected
}

// generateSlug 生成 slug
func generateSlug(title string) string {
	// 简单的 slug 生成逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// 这里可以添加更复杂的 slug 生成逻辑
	return slug
}
