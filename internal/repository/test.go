package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// testRepository Test仓储实现
type testRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewTestRepository 创建Test仓储
func NewTestRepository(db database.Database, logger logger.Logger) TestRepository {
	return &testRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Test
func (r *testRepository) Create(ctx context.Context, entity *model.Test) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Test
func (r *testRepository) GetByID(ctx context.Context, id uint) (*model.Test, error) {
	var entity model.Test
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Test not found")
		}
		r.logger.Error("Failed to get Test by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Test: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Test
func (r *testRepository) GetByName(ctx context.Context, name string) (*model.Test, error) {
	var entity model.Test
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Test not found")
		}
		r.logger.Error("Failed to get Test by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Test: %w", err)
	}
	return &entity, nil
}

// Update 更新Test
func (r *testRepository) Update(ctx context.Context, entity *model.Test) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Test", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Test: %w", err)
	}
	return nil
}

// Delete 删除Test
func (r *testRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Test{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Test", "id", id, "error", err)
		return fmt.Errorf("failed to delete Test: %w", err)
	}
	return nil
}

// List 获取Test列表
func (r *testRepository) List(ctx context.Context, opts ListOptions) ([]*model.Test, int64, error) {
	var entities []*model.Test
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Test{})

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
func (r *testRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
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
