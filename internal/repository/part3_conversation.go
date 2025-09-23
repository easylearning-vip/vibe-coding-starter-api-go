package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// part3ConversationRepository Part3Conversation仓储实现
type part3ConversationRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPart3ConversationRepository 创建Part3Conversation仓储
func NewPart3ConversationRepository(db database.Database, logger logger.Logger) Part3ConversationRepository {
	return &part3ConversationRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Part3Conversation
func (r *part3ConversationRepository) Create(ctx context.Context, entity *model.Part3Conversation) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Part3Conversation
func (r *part3ConversationRepository) GetByID(ctx context.Context, id uint) (*model.Part3Conversation, error) {
	var entity model.Part3Conversation
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part3Conversation not found")
		}
		r.logger.Error("Failed to get Part3Conversation by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Part3Conversation: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Part3Conversation
func (r *part3ConversationRepository) GetByName(ctx context.Context, name string) (*model.Part3Conversation, error) {
	var entity model.Part3Conversation
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Part3Conversation not found")
		}
		r.logger.Error("Failed to get Part3Conversation by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Part3Conversation: %w", err)
	}
	return &entity, nil
}

// Update 更新Part3Conversation
func (r *part3ConversationRepository) Update(ctx context.Context, entity *model.Part3Conversation) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Part3Conversation", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Part3Conversation: %w", err)
	}
	return nil
}

// Delete 删除Part3Conversation
func (r *part3ConversationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Part3Conversation{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Part3Conversation", "id", id, "error", err)
		return fmt.Errorf("failed to delete Part3Conversation: %w", err)
	}
	return nil
}

// List 获取Part3Conversation列表
func (r *part3ConversationRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part3Conversation, int64, error) {
	var entities []*model.Part3Conversation
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Part3Conversation{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR question1 LIKE ? OR question2 LIKE ? OR question3 LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
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
func (r *part3ConversationRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
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
		case "conversation_number":
			if v, ok := value.(int32); ok && v > 0 {
				query = query.Where("conversation_number = ?", v)
			}
		}
	}
	return query
}
