package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

// part3PracticeItemRepository implementation
type part3PracticeItemRepository struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewPart3PracticeItemRepository(db database.Database, logger logger.Logger) Part3PracticeItemRepository {
    return &part3PracticeItemRepository{db: db.GetDB(), logger: logger}
}

func (r *part3PracticeItemRepository) Create(ctx context.Context, entity *model.Part3PracticeItem) error {
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        r.logger.Error("Failed to create Part3PracticeItem", "error", err)
        return fmt.Errorf("failed to create Part3PracticeItem: %w", err)
    }
    return nil
}

func (r *part3PracticeItemRepository) GetByID(ctx context.Context, id uint) (*model.Part3PracticeItem, error) {
    var entity model.Part3PracticeItem
    if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("Part3PracticeItem not found")
        }
        r.logger.Error("Failed to get Part3PracticeItem by ID", "id", id, "error", err)
        return nil, fmt.Errorf("failed to get Part3PracticeItem: %w", err)
    }
    return &entity, nil
}

func (r *part3PracticeItemRepository) Update(ctx context.Context, entity *model.Part3PracticeItem) error {
    if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
        r.logger.Error("Failed to update Part3PracticeItem", "id", entity.ID, "error", err)
        return fmt.Errorf("failed to update Part3PracticeItem: %w", err)
    }
    return nil
}

func (r *part3PracticeItemRepository) Delete(ctx context.Context, id uint) error {
    if err := r.db.WithContext(ctx).Delete(&model.Part3PracticeItem{}, id).Error; err != nil {
        r.logger.Error("Failed to delete Part3PracticeItem", "id", id, "error", err)
        return fmt.Errorf("failed to delete Part3PracticeItem: %w", err)
    }
    return nil
}

func (r *part3PracticeItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part3PracticeItem, int64, error) {
    var entities []*model.Part3PracticeItem
    var total int64

    query := r.db.WithContext(ctx).Model(&model.Part3PracticeItem{})
    query = r.applyFilters(query, opts.Filters)

    if err := query.Count(&total).Error; err != nil {
        r.logger.Error("Failed to count Part3PracticeItem", "error", err)
        return nil, 0, fmt.Errorf("failed to count Part3PracticeItem: %w", err)
    }

    if opts.Sort != "" {
        order := "ASC"
        if opts.Order == "desc" { order = "DESC" }
        query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
    } else {
        query = query.Order("created_at DESC")
    }

    if opts.Page > 0 && opts.PageSize > 0 {
        offset := (opts.Page - 1) * opts.PageSize
        query = query.Offset(offset).Limit(opts.PageSize)
    }

    if err := query.Find(&entities).Error; err != nil {
        r.logger.Error("Failed to list Part3PracticeItem", "error", err)
        return nil, 0, fmt.Errorf("failed to list Part3PracticeItem: %w", err)
    }

    return entities, total, nil
}

func (r *part3PracticeItemRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
    for key, value := range filters {
        switch key {
        case "user_id":
            if v, ok := value.(uint); ok && v > 0 { query = query.Where("user_id = ?", v) }
        case "conversation_id":
            if v, ok := value.(uint); ok && v > 0 { query = query.Where("conversation_id = ?", v) }
        case "question_index":
            switch vv := value.(type) {
            case int32: if vv > 0 { query = query.Where("question_index = ?", vv) }
            case int: if vv > 0 { query = query.Where("question_index = ?", vv) }
            }
        case "scenario_id":
            switch vv := value.(type) {
            case int32: if vv > 0 { query = query.Where("scenario_id = ?", vv) }
            case int: if vv > 0 { query = query.Where("scenario_id = ?", vv) }
            }
        case "difficulty_level_id":
            switch vv := value.(type) {
            case int32: if vv > 0 { query = query.Where("difficulty_level_id = ?", vv) }
            case int: if vv > 0 { query = query.Where("difficulty_level_id = ?", vv) }
            }
        case "is_correct":
            if v, ok := value.(bool); ok { query = query.Where("is_correct = ?", v) }
        }
    }
    return query
}

