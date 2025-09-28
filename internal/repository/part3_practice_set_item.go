package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

type part3PracticeSetItemRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPart3PracticeSetItemRepository(db database.Database, logger logger.Logger) Part3PracticeSetItemRepository {
	return &part3PracticeSetItemRepository{db: db.GetDB(), logger: logger}
}
func (r *part3PracticeSetItemRepository) Create(ctx context.Context, e *model.Part3PracticeSetItem) error {
	return r.db.WithContext(ctx).Create(e).Error
}
func (r *part3PracticeSetItemRepository) GetByID(ctx context.Context, id uint) (*model.Part3PracticeSetItem, error) {
	var e model.Part3PracticeSetItem
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}
func (r *part3PracticeSetItemRepository) Update(ctx context.Context, e *model.Part3PracticeSetItem) error {
	return r.db.WithContext(ctx).Model(&model.Part3PracticeSetItem{}).
		Where("id = ?", e.ID).
		Updates(map[string]interface{}{
			"selected_answer": e.SelectedAnswer,
			"is_correct":      e.IsCorrect,
		}).Error
}
func (r *part3PracticeSetItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Part3PracticeSetItem{}, id).Error
}
func (r *part3PracticeSetItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part3PracticeSetItem, int64, error) {
	var list []*model.Part3PracticeSetItem
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Part3PracticeSetItem{})
	for k, v := range opts.Filters {
		switch k {
		case "set_id":
			if vv, ok := v.(uint); ok && vv > 0 {
				q = q.Where("set_id = ?", vv)
			}
		case "conversation_id":
			if vv, ok := v.(uint); ok && vv > 0 {
				q = q.Where("conversation_id = ?", vv)
			}
		case "question_index":
			if vv, ok := v.(int); ok && vv > 0 {
				q = q.Where("question_index = ?", vv)
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
