package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// productRepository Product仓储实现
type productRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewProductRepository 创建Product仓储
func NewProductRepository(db database.Database, logger logger.Logger) ProductRepository {
	return &productRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Product
func (r *productRepository) Create(ctx context.Context, entity *model.Product) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Product
func (r *productRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Product not found")
		}
		r.logger.Error("Failed to get Product by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Product: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Product
func (r *productRepository) GetByName(ctx context.Context, name string) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Product not found")
		}
		r.logger.Error("Failed to get Product by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Product: %w", err)
	}
	return &entity, nil
}

// Update 更新Product
func (r *productRepository) Update(ctx context.Context, entity *model.Product) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Product", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Product: %w", err)
	}
	return nil
}

// Delete 删除Product
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Product{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Product", "id", id, "error", err)
		return fmt.Errorf("failed to delete Product: %w", err)
	}
	return nil
}

// List 获取Product列表
func (r *productRepository) List(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count <no value>", "error", err)
		return nil, 0, fmt.Errorf("failed to count <no value>: %w", err)
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
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to list <no value>", "error", err)
		return nil, 0, fmt.Errorf("failed to list <no value>: %w", err)
	}

	return entities, total, nil
}

// applyFilters 应用过滤器
func (r *productRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name = ?", v)
			}
		// 在这里添加更多过滤器
		}
	}
	return query
}
