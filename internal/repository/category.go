package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// categoryRepository 分类仓储实现
type categoryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db database.Database, logger logger.Logger) CategoryRepository {
	return &categoryRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建分类
func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		r.logger.Error("Failed to create category", "error", err)
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取分类
func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found with id %d", id)
		}
		r.logger.Error("Failed to get category by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// Update 更新分类
func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		r.logger.Error("Failed to update category", "id", category.ID, "error", err)
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// Delete 删除分类
func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Category{}, id).Error; err != nil {
		r.logger.Error("Failed to delete category", "id", id, "error", err)
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// List 获取分类列表
func (r *categoryRepository) List(ctx context.Context, opts ListOptions) ([]*model.Category, int64, error) {
	var categories []*model.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Category{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count categories", "error", err)
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("sort_order ASC, created_at DESC")
	}

	// 应用分页
	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&categories).Error; err != nil {
		r.logger.Error("Failed to get categories", "error", err)
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, total, nil
}

// GetBySlug 根据 slug 获取分类
func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found with slug %s", slug)
		}
		r.logger.Error("Failed to get category by slug", "slug", slug, "error", err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// GetByName 根据名称获取分类
func (r *categoryRepository) GetByName(ctx context.Context, name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found with name %s", name)
		}
		r.logger.Error("Failed to get category by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// applyFilters 应用过滤器
func (r *categoryRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "name":
			query = query.Where("name = ?", value)
		case "slug":
			query = query.Where("slug = ?", value)
		}
	}

	return query
}
