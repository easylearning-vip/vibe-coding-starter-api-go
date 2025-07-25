package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// commentRepository 评论仓储实现
type commentRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db database.Database, logger logger.Logger) CommentRepository {
	return &commentRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建评论
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		r.logger.Error("Failed to create comment", "error", err)
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取评论
func (r *commentRepository) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Article").
		Preload("Parent").
		First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment not found with id %d", id)
		}
		r.logger.Error("Failed to get comment by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

// Update 更新评论
func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	if err := r.db.WithContext(ctx).Save(comment).Error; err != nil {
		r.logger.Error("Failed to update comment", "id", comment.ID, "error", err)
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

// Delete 删除评论
func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error; err != nil {
		r.logger.Error("Failed to delete comment", "id", id, "error", err)
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// List 获取评论列表
func (r *commentRepository) List(ctx context.Context, opts ListOptions) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Preload("Author").
		Preload("Article").
		Preload("Parent")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("content LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count comments", "error", err)
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&comments).Error; err != nil {
		r.logger.Error("Failed to get comments", "error", err)
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, total, nil
}

// GetByArticle 根据文章获取评论列表
func (r *commentRepository) GetByArticle(ctx context.Context, articleID uint, opts ListOptions) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("article_id = ?", articleID).
		Preload("Author").
		Preload("Parent")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count comments by article", "article_id", articleID, "error", err)
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at ASC")
	}

	// 应用分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&comments).Error; err != nil {
		r.logger.Error("Failed to get comments by article", "article_id", articleID, "error", err)
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, total, nil
}

// GetByAuthor 根据作者获取评论列表
func (r *commentRepository) GetByAuthor(ctx context.Context, authorID uint, opts ListOptions) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("author_id = ?", authorID).
		Preload("Author").
		Preload("Article").
		Preload("Parent")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count comments by author", "author_id", authorID, "error", err)
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&comments).Error; err != nil {
		r.logger.Error("Failed to get comments by author", "author_id", authorID, "error", err)
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, total, nil
}

// GetReplies 获取评论的回复列表
func (r *commentRepository) GetReplies(ctx context.Context, parentID uint, opts ListOptions) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("parent_id = ?", parentID).
		Preload("Author").
		Preload("Article")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count replies", "parent_id", parentID, "error", err)
		return nil, 0, fmt.Errorf("failed to count replies: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at ASC")
	}

	// 应用分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&comments).Error; err != nil {
		r.logger.Error("Failed to get replies", "parent_id", parentID, "error", err)
		return nil, 0, fmt.Errorf("failed to get replies: %w", err)
	}

	return comments, total, nil
}

// applyFilters 应用过滤器
func (r *commentRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "article_id":
			query = query.Where("article_id = ?", value)
		case "author_id":
			query = query.Where("author_id = ?", value)
		case "parent_id":
			if value == nil {
				query = query.Where("parent_id IS NULL")
			} else {
				query = query.Where("parent_id = ?", value)
			}
		}
	}

	return query
}
