package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// tagRepository 标签仓储实现
type tagRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db database.Database, logger logger.Logger) TagRepository {
	return &tagRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建标签
func (r *tagRepository) Create(ctx context.Context, tag *model.Tag) error {
	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		r.logger.Error("Failed to create tag", "error", err)
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取标签
func (r *tagRepository) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found with id %d", id)
		}
		r.logger.Error("Failed to get tag by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

// Update 更新标签
func (r *tagRepository) Update(ctx context.Context, tag *model.Tag) error {
	if err := r.db.WithContext(ctx).Save(tag).Error; err != nil {
		r.logger.Error("Failed to update tag", "id", tag.ID, "error", err)
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

// Delete 删除标签
func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Tag{}, id).Error; err != nil {
		r.logger.Error("Failed to delete tag", "id", id, "error", err)
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

// List 获取标签列表
func (r *tagRepository) List(ctx context.Context, opts ListOptions) ([]*model.Tag, int64, error) {
	var tags []*model.Tag
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Tag{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count tags", "error", err)
		return nil, 0, fmt.Errorf("failed to count tags: %w", err)
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

	if err := query.Find(&tags).Error; err != nil {
		r.logger.Error("Failed to get tags", "error", err)
		return nil, 0, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, total, nil
}

// GetBySlug 根据 slug 获取标签
func (r *tagRepository) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found with slug %s", slug)
		}
		r.logger.Error("Failed to get tag by slug", "slug", slug, "error", err)
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

// GetByName 根据名称获取标签
func (r *tagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found with name %s", name)
		}
		r.logger.Error("Failed to get tag by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

// GetByNames 根据名称列表获取标签
func (r *tagRepository) GetByNames(ctx context.Context, names []string) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&tags).Error; err != nil {
		r.logger.Error("Failed to get tags by names", "names", names, "error", err)
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

// applyFilters 应用过滤器
func (r *tagRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "name":
			query = query.Where("name = ?", value)
		case "slug":
			query = query.Where("slug = ?", value)
		case "color":
			query = query.Where("color = ?", value)
		}
	}

	return query
}
