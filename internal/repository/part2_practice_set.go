package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

type part2PracticeSetRepository struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewPart2PracticeSetRepository(db database.Database, logger logger.Logger) Part2PracticeSetRepository {
    return &part2PracticeSetRepository{db: db.GetDB(), logger: logger}
}

func (r *part2PracticeSetRepository) Create(ctx context.Context, entity *model.Part2PracticeSet) error {
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        r.logger.Error("Failed to create Part2PracticeSet", "error", err)
        return fmt.Errorf("failed to create Part2PracticeSet: %w", err)
    }
    return nil
}

func (r *part2PracticeSetRepository) GetByID(ctx context.Context, id uint) (*model.Part2PracticeSet, error) {
    var entity model.Part2PracticeSet
    if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("Part2PracticeSet not found")
        }
        r.logger.Error("Failed to get Part2PracticeSet by ID", "id", id, "error", err)
        return nil, fmt.Errorf("failed to get Part2PracticeSet: %w", err)
    }
    return &entity, nil
}

func (r *part2PracticeSetRepository) Update(ctx context.Context, entity *model.Part2PracticeSet) error {
    if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
        r.logger.Error("Failed to update Part2PracticeSet", "id", entity.ID, "error", err)
        return fmt.Errorf("failed to update Part2PracticeSet: %w", err)
    }
    return nil
}

func (r *part2PracticeSetRepository) Delete(ctx context.Context, id uint) error {
    if err := r.db.WithContext(ctx).Delete(&model.Part2PracticeSet{}, id).Error; err != nil {
        r.logger.Error("Failed to delete Part2PracticeSet", "id", id, "error", err)
        return fmt.Errorf("failed to delete Part2PracticeSet: %w", err)
    }
    return nil
}

func (r *part2PracticeSetRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part2PracticeSet, int64, error) {
    var entities []*model.Part2PracticeSet
    var total int64

    query := r.db.WithContext(ctx).Model(&model.Part2PracticeSet{})
    query = applyCommonSetFilters(query, opts.Filters)

    if err := query.Count(&total).Error; err != nil {
        r.logger.Error("Failed to count Part2PracticeSet", "error", err)
        return nil, 0, fmt.Errorf("failed to count Part2PracticeSet: %w", err)
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
        r.logger.Error("Failed to list Part2PracticeSet", "error", err)
        return nil, 0, fmt.Errorf("failed to list Part2PracticeSet: %w", err)
    }

    return entities, total, nil
}

func applyCommonSetFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
    for key, value := range filters {
        switch key {
        case "user_id":
            if v, ok := value.(uint); ok && v > 0 { query = query.Where("user_id = ?", v) }
        case "scenario_id":
            switch vv := value.(type) { case int32: if vv>0 {query= query.Where("scenario_id = ?", vv)}; case int: if vv>0 {query= query.Where("scenario_id = ?", vv)} }
        case "difficulty_level_id":
            switch vv := value.(type) { case int32: if vv>0 {query= query.Where("difficulty_level_id = ?", vv)}; case int: if vv>0 {query= query.Where("difficulty_level_id = ?", vv)} }
        }
    }
    return query
}

