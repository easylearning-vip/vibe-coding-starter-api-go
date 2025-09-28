package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// toeicAiPromptRepository ToeicAiPrompt仓储实现
type toeicAiPromptRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewToeicAiPromptRepository 创建ToeicAiPrompt仓储
func NewToeicAiPromptRepository(db database.Database, logger logger.Logger) ToeicAiPromptRepository {
	return &toeicAiPromptRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建ToeicAiPrompt
func (r *toeicAiPromptRepository) Create(ctx context.Context, entity *model.ToeicAiPrompt) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create toeic ai prompt", "error", err)
		return fmt.Errorf("failed to create toeic ai prompt: %w", err)
	}
	return nil
}

// GetByID 根据ID获取ToeicAiPrompt
func (r *toeicAiPromptRepository) GetByID(ctx context.Context, id uint) (*model.ToeicAiPrompt, error) {
	var entity model.ToeicAiPrompt
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("toeic ai prompt not found")
		}
		r.logger.Error("Failed to get toeic ai prompt by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get toeic ai prompt: %w", err)
	}
	return &entity, nil
}

// Update 更新ToeicAiPrompt
func (r *toeicAiPromptRepository) Update(ctx context.Context, entity *model.ToeicAiPrompt) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update toeic ai prompt", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update toeic ai prompt: %w", err)
	}
	return nil
}

// Delete 删除ToeicAiPrompt
func (r *toeicAiPromptRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.ToeicAiPrompt{}, id).Error; err != nil {
		r.logger.Error("Failed to delete toeic ai prompt", "id", id, "error", err)
		return fmt.Errorf("failed to delete toeic ai prompt: %w", err)
	}
	return nil
}

// List 获取ToeicAiPrompt列表
func (r *toeicAiPromptRepository) List(ctx context.Context, opts ListOptions) ([]*model.ToeicAiPrompt, int64, error) {
	var entities []*model.ToeicAiPrompt
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ToeicAiPrompt{})

	// 搜索条件
	if opts.Search != "" {
		searchPattern := "%" + opts.Search + "%"
		query = query.Where("part2_prompt LIKE ? OR part3_prompt LIKE ? OR part4_prompt LIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}

	// 过滤条件
	for key, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count toeic ai prompts", "error", err)
		return nil, 0, fmt.Errorf("failed to count toeic ai prompts: %w", err)
	}

	// 排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 查询数据
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to list toeic ai prompts", "error", err)
		return nil, 0, fmt.Errorf("failed to list toeic ai prompts: %w", err)
	}

	return entities, total, nil
}
