package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// difficultyLevelRepository DifficultyLevel仓储实现
type difficultyLevelRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDifficultyLevelRepository 创建DifficultyLevel仓储
func NewDifficultyLevelRepository(db database.Database, logger logger.Logger) DifficultyLevelRepository {
	return &difficultyLevelRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建DifficultyLevel
func (r *difficultyLevelRepository) Create(ctx context.Context, entity *model.DifficultyLevel) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取DifficultyLevel
func (r *difficultyLevelRepository) GetByID(ctx context.Context, id uint) (*model.DifficultyLevel, error) {
	var entity model.DifficultyLevel
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("DifficultyLevel not found")
		}
		r.logger.Error("Failed to get DifficultyLevel by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get DifficultyLevel: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取DifficultyLevel
func (r *difficultyLevelRepository) GetByName(ctx context.Context, name string) (*model.DifficultyLevel, error) {
	var entity model.DifficultyLevel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("DifficultyLevel not found")
		}
		r.logger.Error("Failed to get DifficultyLevel by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get DifficultyLevel: %w", err)
	}
	return &entity, nil
}

// Update 更新DifficultyLevel
func (r *difficultyLevelRepository) Update(ctx context.Context, entity *model.DifficultyLevel) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update DifficultyLevel", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update DifficultyLevel: %w", err)
	}
	return nil
}

// Delete 删除DifficultyLevel
func (r *difficultyLevelRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DifficultyLevel{}, id).Error; err != nil {
		r.logger.Error("Failed to delete DifficultyLevel", "id", id, "error", err)
		return fmt.Errorf("failed to delete DifficultyLevel: %w", err)
	}
	return nil
}

// List 获取DifficultyLevel列表
func (r *difficultyLevelRepository) List(ctx context.Context, opts ListOptions) ([]*model.DifficultyLevel, int64, error) {
	var entities []*model.DifficultyLevel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DifficultyLevel{})

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
func (r *difficultyLevelRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
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
