package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

type part2PracticeSetItemRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPart2PracticeSetItemRepository(db database.Database, logger logger.Logger) Part2PracticeSetItemRepository {
	return &part2PracticeSetItemRepository{db: db.GetDB(), logger: logger}
}

func (r *part2PracticeSetItemRepository) Create(ctx context.Context, entity *model.Part2PracticeSetItem) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("create Part2PracticeSetItem: %w", err)
	}
	return nil
}
func (r *part2PracticeSetItemRepository) GetByID(ctx context.Context, id uint) (*model.Part2PracticeSetItem, error) {
	var e model.Part2PracticeSetItem
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}
func (r *part2PracticeSetItemRepository) Update(ctx context.Context, entity *model.Part2PracticeSetItem) error {
	// Only update mutable fields to avoid unintended changes
	return r.db.WithContext(ctx).Model(&model.Part2PracticeSetItem{}).
		Where("id = ?", entity.ID).
		Updates(map[string]interface{}{
			"selected_answer": entity.SelectedAnswer,
			"is_correct":      entity.IsCorrect,
		}).Error
}
func (r *part2PracticeSetItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Part2PracticeSetItem{}, id).Error
}
func (r *part2PracticeSetItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part2PracticeSetItem, int64, error) {
	var list []*model.Part2PracticeSetItem
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Part2PracticeSetItem{})
	for k, v := range opts.Filters {
		switch k {
		case "set_id":
			if vv, ok := v.(uint); ok && vv > 0 {
				q = q.Where("set_id = ?", vv)
			}
		case "question_id":
			if vv, ok := v.(uint); ok && vv > 0 {
				q = q.Where("question_id = ?", vv)
			}
		}
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if opts.Sort != "" {
		ord := "ASC"
		if opts.Order == "desc" {
			ord = "DESC"
		}
		q = q.Order(fmt.Sprintf("%s %s", opts.Sort, ord))
	} else {
		q = q.Order("order_index ASC, id ASC")
	}
	if opts.Page > 0 && opts.PageSize > 0 {
		off := (opts.Page - 1) * opts.PageSize
		q = q.Offset(off).Limit(opts.PageSize)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
