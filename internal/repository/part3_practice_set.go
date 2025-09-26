package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

type part3PracticeSetRepository struct { db *gorm.DB; logger logger.Logger }
func NewPart3PracticeSetRepository(db database.Database, logger logger.Logger) Part3PracticeSetRepository { return &part3PracticeSetRepository{db: db.GetDB(), logger: logger} }
func (r *part3PracticeSetRepository) Create(ctx context.Context, e *model.Part3PracticeSet) error { return r.db.WithContext(ctx).Create(e).Error }
func (r *part3PracticeSetRepository) GetByID(ctx context.Context, id uint) (*model.Part3PracticeSet, error) { var e model.Part3PracticeSet; if err:=r.db.WithContext(ctx).First(&e, id).Error; err!=nil { return nil, err }; return &e, nil }
func (r *part3PracticeSetRepository) Update(ctx context.Context, e *model.Part3PracticeSet) error { return r.db.WithContext(ctx).Save(e).Error }
func (r *part3PracticeSetRepository) Delete(ctx context.Context, id uint) error { return r.db.WithContext(ctx).Delete(&model.Part3PracticeSet{}, id).Error }
func (r *part3PracticeSetRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part3PracticeSet, int64, error) {
    var list []*model.Part3PracticeSet; var total int64
    q := r.db.WithContext(ctx).Model(&model.Part3PracticeSet{})
    q = applyCommonSetFilters(q, opts.Filters)
    if err := q.Count(&total).Error; err != nil { return nil,0, err }
    if opts.Sort != "" { ord := "ASC"; if opts.Order=="desc" { ord = "DESC" }; q = q.Order(fmt.Sprintf("%s %s", opts.Sort, ord)) } else { q = q.Order("created_at DESC") }
    if opts.Page>0 && opts.PageSize>0 { off := (opts.Page-1)*opts.PageSize; q = q.Offset(off).Limit(opts.PageSize) }
    if err := q.Find(&list).Error; err != nil { return nil,0, err }
    return list, total, nil
}

