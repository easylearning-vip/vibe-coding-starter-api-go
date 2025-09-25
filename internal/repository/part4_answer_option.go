package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// part4AnswerOptionRepository Part4AnswerOption仓储实现
type part4AnswerOptionRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPart4AnswerOptionRepository 创建Part4AnswerOption仓储
func NewPart4AnswerOptionRepository(db database.Database, logger logger.Logger) Part4AnswerOptionRepository {
	return &part4AnswerOptionRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Part4AnswerOption
func (r *part4AnswerOptionRepository) Create(ctx context.Context, entity *model.Part4AnswerOption) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Part4AnswerOption
func (r *part4AnswerOptionRepository) GetByID(ctx context.Context, id uint) (*model.Part4AnswerOption, error) {
	var entity model.Part4AnswerOption
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part4AnswerOption not found")
		}
		r.logger.Error("Failed to get Part4AnswerOption by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Part4AnswerOption: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Part4AnswerOption
func (r *part4AnswerOptionRepository) GetByName(ctx context.Context, name string) (*model.Part4AnswerOption, error) {
	var entity model.Part4AnswerOption
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part4AnswerOption not found")
		}
		r.logger.Error("Failed to get Part4AnswerOption by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Part4AnswerOption: %w", err)
	}
	return &entity, nil
}

// Update 更新Part4AnswerOption
func (r *part4AnswerOptionRepository) Update(ctx context.Context, entity *model.Part4AnswerOption) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Part4AnswerOption", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Part4AnswerOption: %w", err)
	}
	return nil
}

// Delete 删除Part4AnswerOption
func (r *part4AnswerOptionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Part4AnswerOption{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Part4AnswerOption", "id", id, "error", err)
		return fmt.Errorf("failed to delete Part4AnswerOption: %w", err)
	}
	return nil
}

// List 获取Part4AnswerOption列表
func (r *part4AnswerOptionRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part4AnswerOption, int64, error) {
	var entities []*model.Part4AnswerOption
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Part4AnswerOption{})

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
func (r *part4AnswerOptionRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name = ?", v)
			}
		case "talk_id":
			if v, ok := value.(int); ok && v > 0 {
				query = query.Where("talk_id = ?", v)
			}
			// 在这里添加更多过滤器
		}
	}
	return query
}

// DeleteByTalkIDs 根据演讲ID列表删除Part4AnswerOption
func (r *part4AnswerOptionRepository) DeleteByTalkIDs(ctx context.Context, talkIDs []uint) error {
	if len(talkIDs) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Unscoped().Where("talk_id IN ?", talkIDs).Delete(&model.Part4AnswerOption{}).Error; err != nil {
		r.logger.Error("Failed to delete Part4AnswerOption by talk IDs", "talkIDs", talkIDs, "error", err)
		return fmt.Errorf("failed to delete Part4AnswerOption by talk IDs: %w", err)
	}
	return nil
}
