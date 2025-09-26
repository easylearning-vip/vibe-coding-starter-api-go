package repository

import (
    "context"
    "fmt"

    "gorm.io/gorm"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/pkg/logger"
)

type part4PracticeSetItemRepository struct { db *gorm.DB; logger logger.Logger }
func NewPart4PracticeSetItemRepository(db database.Database, logger logger.Logger) Part4PracticeSetItemRepository { return &part4PracticeSetItemRepository{db: db.GetDB(), logger: logger} }
func (r *part4PracticeSetItemRepository) Create(ctx context.Context, e *model.Part4PracticeSetItem) error { return r.db.WithContext(ctx).Create(e).Error }
func (r *part4PracticeSetItemRepository) GetByID(ctx context.Context, id uint) (*model.Part4PracticeSetItem, error) { var e model.Part4PracticeSetItem; if err:=r.db.WithContext(ctx).First(&e, id).Error; err!=nil { return nil, err }; return &e, nil }
func (r *part4PracticeSetItemRepository) Update(ctx context.Context, e *model.Part4PracticeSetItem) error { return r.db.WithContext(ctx).Save(e).Error }
func (r *part4PracticeSetItemRepository) Delete(ctx context.Context, id uint) error { return r.db.WithContext(ctx).Delete(&model.Part4PracticeSetItem{}, id).Error }
func (r *part4PracticeSetItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.Part4PracticeSetItem, int64, error) {
    var list []*model.Part4PracticeSetItem; var total int64
    q := r.db.WithContext(ctx).Model(&model.Part4PracticeSetItem{})
    for k,v := range opts.Filters { switch k { case "set_id": if vv,ok:=v.(uint); ok && vv>0 { q = q.Where("set_id = ?", vv) }; case "talk_id": if vv,ok:=v.(uint); ok && vv>0 { q = q.Where("talk_id = ?", vv) }; case "question_index": if vv,ok:=v.(int); ok && vv>0 { q = q.Where("question_index = ?", vv) } } }
    if err := q.Count(&total).Error; err != nil { return nil,0, err }
    if opts.Sort != "" { ord := "ASC"; if opts.Order=="desc" { ord = "DESC" }; q = q.Order(fmt.Sprintf("%s %s", opts.Sort, ord)) } else { q = q.Order("order_index ASC, id ASC") }
    if opts.Page>0 && opts.PageSize>0 { off := (opts.Page-1)*opts.PageSize; q = q.Offset(off).Limit(opts.PageSize) }
    if err := q.Find(&list).Error; err != nil { return nil,0, err }
    return list, total, nil
}

