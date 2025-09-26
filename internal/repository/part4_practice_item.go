package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

// part4PracticeItemRepository implementation
type part4PracticeItemRepository struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewPart4PracticeItemRepository(db database.Database, logger logger.Logger) Part4PracticeItemRepository {
    return &part4PracticeItemRepository{db: db.GetDB(), logger: logger}
}

func (r *part4PracticeItemRepository) Create(ctx context.Context, entity *model.Part4PracticeItem) error {
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        r.logger.Error("Failed to create Part4PracticeItem", "error", err)
        return fmt.Errorf("failed to create Part4PracticeItem: %w", err)
    }
    return nil
}

func (r *part4PracticeItemRepository) GetByID(ctx context.Context, id uint) (*model.Part4PracticeItem, error) {
    var entity model.Part4PracticeItem
    if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("Part4PracticeItem not found")
        }
        r.logger.Error("Failed to get Part4PracticeItem by ID", "id", id, "error", err)
        return nil, fmt.Errorf("failed to get Part4PracticeItem: %w", err)
    }
    return &entity, nil
}

func (r *part4PracticeItemRepository) Update(ctx context.Context, entity *model.Part4PracticeItem) error {
    if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
        r.logger.Error("Failed to update Part4PracticeItem", "id", entity.ID, "error", err)
        return fmt.Errorf("failed to update Part4PracticeItem: %w", err)
    }
    return nil
}

func (r *part4PracticeItemRepository) Delete(ctx context.Context, id uint) error {
    if err := r.db.WithContext(ctx).Delete(&model.Part4PracticeItem{}, id).Error; err != nil {
        r.logger.Error("Failed to delete Part4PracticeItem", "id", id, "error", err)
        return fmt.Errorf("failed to delete Part4PracticeItem: %w", err)
    }
    return nil
}

func (r *part4PracticeItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part4PracticeItem, int64, error) {
    var entities []*model.Part4PracticeItem
    var total int64

    query := r.db.WithContext(ctx).Model(&model.Part4PracticeItem{})
    query = r.applyFilters(query, opts.Filters)

    if err := query.Count(&total).Error; err != nil {
        r.logger.Error("Failed to count Part4PracticeItem", "error", err)
        return nil, 0, fmt.Errorf("failed to count Part4PracticeItem: %w", err)
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
        r.logger.Error("Failed to list Part4PracticeItem", "error", err)
        return nil, 0, fmt.Errorf("failed to list Part4PracticeItem: %w", err)
    }

    return entities, total, nil
}

func (r *part4PracticeItemRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
    for key, value := range filters {
        switch key {
        case "user_id":
            if v, ok := value.(uint); ok && v > 0 { query = query.Where("user_id = ?", v) }
        case "talk_id":
            if v, ok := value.(uint); ok && v > 0 { query = query.Where("talk_id = ?", v) }
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

