package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// part3AnswerOptionRepository Part3AnswerOption仓储实现
type part3AnswerOptionRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPart3AnswerOptionRepository 创建Part3AnswerOption仓储
func NewPart3AnswerOptionRepository(db database.Database, logger logger.Logger) Part3AnswerOptionRepository {
	return &part3AnswerOptionRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Part3AnswerOption
func (r *part3AnswerOptionRepository) Create(ctx context.Context, entity *model.Part3AnswerOption) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Part3AnswerOption
func (r *part3AnswerOptionRepository) GetByID(ctx context.Context, id uint) (*model.Part3AnswerOption, error) {
	var entity model.Part3AnswerOption
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part3AnswerOption not found")
		}
		r.logger.Error("Failed to get Part3AnswerOption by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Part3AnswerOption: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Part3AnswerOption
func (r *part3AnswerOptionRepository) GetByName(ctx context.Context, name string) (*model.Part3AnswerOption, error) {
	var entity model.Part3AnswerOption
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part3AnswerOption not found")
		}
		r.logger.Error("Failed to get Part3AnswerOption by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Part3AnswerOption: %w", err)
	}
	return &entity, nil
}

// Update 更新Part3AnswerOption
func (r *part3AnswerOptionRepository) Update(ctx context.Context, entity *model.Part3AnswerOption) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Part3AnswerOption", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Part3AnswerOption: %w", err)
	}
	return nil
}

// Delete 删除Part3AnswerOption
func (r *part3AnswerOptionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Part3AnswerOption{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Part3AnswerOption", "id", id, "error", err)
		return fmt.Errorf("failed to delete Part3AnswerOption: %w", err)
	}
	return nil
}

// List 获取Part3AnswerOption列表
func (r *part3AnswerOptionRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part3AnswerOption, int64, error) {
	var entities []*model.Part3AnswerOption
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Part3AnswerOption{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("option_a LIKE ? OR option_b LIKE ? OR option_c LIKE ? OR option_d LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
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
func (r *part3AnswerOptionRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "conversation_id":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("conversation_id = ?", v)
			}
		case "question_number":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("question_number = ?", v)
			}
		case "correct_answer":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("correct_answer = ?", v)
			}
		}
	}
	return query
}
