package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

type part4PracticeSetRepository struct { db *gorm.DB; logger logger.Logger }
func NewPart4PracticeSetRepository(db database.Database, logger logger.Logger) Part4PracticeSetRepository { return &part4PracticeSetRepository{db: db.GetDB(), logger: logger} }
func (r *part4PracticeSetRepository) Create(ctx context.Context, e *model.Part4PracticeSet) error { return r.db.WithContext(ctx).Create(e).Error }
func (r *part4PracticeSetRepository) GetByID(ctx context.Context, id uint) (*model.Part4PracticeSet, error) { var e model.Part4PracticeSet; if err:=r.db.WithContext(ctx).First(&e, id).Error; err!=nil { return nil, err }; return &e, nil }
func (r *part4PracticeSetRepository) Update(ctx context.Context, e *model.Part4PracticeSet) error { return r.db.WithContext(ctx).Save(e).Error }
func (r *part4PracticeSetRepository) Delete(ctx context.Context, id uint) error { return r.db.WithContext(ctx).Delete(&model.Part4PracticeSet{}, id).Error }
func (r *part4PracticeSetRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part4PracticeSet, int64, error) {
    var list []*model.Part4PracticeSet; var total int64
    q := r.db.WithContext(ctx).Model(&model.Part4PracticeSet{})
    q = applyCommonSetFilters(q, opts.Filters)
    if err := q.Count(&total).Error; err != nil { return nil,0, err }
    if opts.Sort != "" { ord := "ASC"; if opts.Order=="desc" { ord = "DESC" }; q = q.Order(fmt.Sprintf("%s %s", opts.Sort, ord)) } else { q = q.Order("created_at DESC") }
    if opts.Page>0 && opts.PageSize>0 { off := (opts.Page-1)*opts.PageSize; q = q.Offset(off).Limit(opts.PageSize) }
    if err := q.Find(&list).Error; err != nil { return nil,0, err }
    return list, total, nil
}

