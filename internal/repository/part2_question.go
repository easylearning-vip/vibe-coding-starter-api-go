package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// part2QuestionRepository Part2Question仓储实现
type part2QuestionRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPart2QuestionRepository 创建Part2Question仓储
func NewPart2QuestionRepository(db database.Database, logger logger.Logger) Part2QuestionRepository {
	return &part2QuestionRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Part2Question
func (r *part2QuestionRepository) Create(ctx context.Context, entity *model.Part2Question) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Part2Question
func (r *part2QuestionRepository) GetByID(ctx context.Context, id uint) (*model.Part2Question, error) {
	var entity model.Part2Question
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part2Question not found")
		}
		r.logger.Error("Failed to get Part2Question by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Part2Question: %w", err)
	}
	return &entity, nil
}

// GetByQuestionText 根据问题文本获取Part2Question
func (r *part2QuestionRepository) GetByQuestionText(ctx context.Context, questionText string) (*model.Part2Question, error) {
	var entity model.Part2Question
	if err := r.db.WithContext(ctx).Where("question_text = ?", questionText).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part2Question not found")
		}
		r.logger.Error("Failed to get Part2Question by question text", "question_text", questionText, "error", err)
		return nil, fmt.Errorf("failed to get Part2Question: %w", err)
	}
	return &entity, nil
}

// Update 更新Part2Question
func (r *part2QuestionRepository) Update(ctx context.Context, entity *model.Part2Question) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Part2Question", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Part2Question: %w", err)
	}
	return nil
}

// Delete 删除Part2Question
func (r *part2QuestionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Part2Question{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Part2Question", "id", id, "error", err)
		return fmt.Errorf("failed to delete Part2Question: %w", err)
	}
	return nil
}

// List 获取Part2Question列表
func (r *part2QuestionRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part2Question, int64, error) {
	var entities []*model.Part2Question
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Part2Question{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("question_text LIKE ?", "%"+opts.Search+"%")
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
func (r *part2QuestionRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "question_text":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("question_text = ?", v)
			}
		case "test_id":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("test_id = ?", v)
			}
		case "scenario_id":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("scenario_id = ?", v)
			}
		case "difficulty_level_id":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("difficulty_level_id = ?", v)
			}
			// 在这里添加更多过滤器
		}
	}
	return query
}
