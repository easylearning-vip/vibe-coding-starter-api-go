package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// scenarioRepository Scenario仓储实现
type scenarioRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewScenarioRepository 创建Scenario仓储
func NewScenarioRepository(db database.Database, logger logger.Logger) ScenarioRepository {
	return &scenarioRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Scenario
func (r *scenarioRepository) Create(ctx context.Context, entity *model.Scenario) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Scenario
func (r *scenarioRepository) GetByID(ctx context.Context, id uint) (*model.Scenario, error) {
	var entity model.Scenario
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Scenario not found")
		}
		r.logger.Error("Failed to get Scenario by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Scenario: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Scenario
func (r *scenarioRepository) GetByName(ctx context.Context, name string) (*model.Scenario, error) {
	var entity model.Scenario
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Scenario not found")
		}
		r.logger.Error("Failed to get Scenario by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Scenario: %w", err)
	}
	return &entity, nil
}

// Update 更新Scenario
func (r *scenarioRepository) Update(ctx context.Context, entity *model.Scenario) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Scenario", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Scenario: %w", err)
	}
	return nil
}

// Delete 删除Scenario
func (r *scenarioRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Scenario{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Scenario", "id", id, "error", err)
		return fmt.Errorf("failed to delete Scenario: %w", err)
	}
	return nil
}

// List 获取Scenario列表
func (r *scenarioRepository) List(ctx context.Context, opts ListOptions) ([]*model.Scenario, int64, error) {
	var entities []*model.Scenario
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Scenario{})

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
func (r *scenarioRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
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
