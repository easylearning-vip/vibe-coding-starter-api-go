package service

import (
    "context"
    "database/sql"
    "fmt"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/repository"
)

type Part3PracticeSetService interface {
    CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part3PracticeSet, error)
    GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeSet, error)
    UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part3PracticeSet, error)
    DeleteSet(ctx context.Context, userID uint, id uint) error
    ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part3PracticeSet, int64, error)

    CreateItem(ctx context.Context, userID uint, req *CreatePart3SetItemRequest) (*model.Part3PracticeSetItem, error)
    ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part3PracticeSetItem, int64, error)
    DeleteItem(ctx context.Context, userID uint, id uint) error
}

type part3PracticeSetService struct { setRepo repository.Part3PracticeSetRepository; itemRepo repository.Part3PracticeSetItemRepository }
func NewPart3PracticeSetService(setRepo repository.Part3PracticeSetRepository, itemRepo repository.Part3PracticeSetItemRepository) Part3PracticeSetService { return &part3PracticeSetService{setRepo: setRepo, itemRepo: itemRepo} }

func (s *part3PracticeSetService) CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part3PracticeSet, error) {
    e := &model.Part3PracticeSet{ UserID: userID, ScenarioId: sql.NullInt32{Int32:req.ScenarioId, Valid:req.ScenarioId!=0}, DifficultyLevelId: sql.NullInt32{Int32:req.DifficultyLevelId, Valid:req.DifficultyLevelId!=0}, TotalQuestions: req.TotalQuestions }
    if err := s.setRepo.Create(ctx, e); err != nil { return nil, err }
    return e, nil
}
func (s *part3PracticeSetService) GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeSet, error) { e,err:=s.setRepo.GetByID(ctx,id); if err!=nil { return nil, err }; if e.UserID!=userID { return nil, fmt.Errorf("not found") }; return e, nil }
func (s *part3PracticeSetService) UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part3PracticeSet, error) { e,err:=s.setRepo.GetByID(ctx,id); if err!=nil { return nil, err }; if e.UserID!=userID { return nil, fmt.Errorf("not found") }; if req.TotalQuestions!=nil { e.TotalQuestions=*req.TotalQuestions }; if req.CompletedCount!=nil { e.CompletedCount=*req.CompletedCount }; if req.CorrectCount!=nil { e.CorrectCount=*req.CorrectCount }; if e.TotalQuestions>0 { e.Accuracy=float64(e.CorrectCount)/float64(e.TotalQuestions) } else { e.Accuracy=0 }; if err:=s.setRepo.Update(ctx,e); err!=nil { return nil, err }; return e,nil }
func (s *part3PracticeSetService) DeleteSet(ctx context.Context, userID uint, id uint) error { e,err:=s.setRepo.GetByID(ctx,id); if err!=nil { return err }; if e.UserID!=userID { return fmt.Errorf("not found") }; return s.setRepo.Delete(ctx,id) }
func (s *part3PracticeSetService) ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part3PracticeSet, int64, error) { if opts==nil { opts=&ListPracticeSetOptions{} }; if opts.Filters==nil { opts.Filters=map[string]interface{}{} }; if opts.Page<=0 { opts.Page=1 }; if opts.PageSize<=0 { opts.PageSize=20 }; if opts.PageSize>100 { opts.PageSize=100 }; opts.Filters["user_id"]=userID; repoOpts:=repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}; return s.setRepo.List(ctx, repoOpts) }

type CreatePart3SetItemRequest struct { SetID uint `json:"set_id"`; ConversationID uint `json:"conversation_id"`; QuestionIndex int32 `json:"question_index"`; OrderIndex int32 `json:"order_index"` }
func (s *part3PracticeSetService) CreateItem(ctx context.Context, userID uint, req *CreatePart3SetItemRequest) (*model.Part3PracticeSetItem, error) { set,err:=s.setRepo.GetByID(ctx, req.SetID); if err!=nil { return nil, err }; if set.UserID!=userID { return nil, fmt.Errorf("not found") }; item:=&model.Part3PracticeSetItem{ SetID:req.SetID, ConversationID:req.ConversationID, QuestionIndex:req.QuestionIndex, OrderIndex:req.OrderIndex }; if err:=s.itemRepo.Create(ctx,item); err!=nil { return nil, err }; return item,nil }
func (s *part3PracticeSetService) ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part3PracticeSetItem, int64, error) { set,err:=s.setRepo.GetByID(ctx,setID); if err!=nil { return nil,0, err }; if set.UserID!=userID { return nil,0, fmt.Errorf("not found") }; if opts==nil { opts=&ListPracticeSetItemOptions{} }; if opts.Filters==nil { opts.Filters=map[string]interface{}{} }; if opts.Page<=0 { opts.Page=1 }; if opts.PageSize<=0 { opts.PageSize=20 }; if opts.PageSize>100 { opts.PageSize=100 }; opts.Filters["set_id"]=setID; repoOpts:=repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}; return s.itemRepo.List(ctx, repoOpts) }
func (s *part3PracticeSetService) DeleteItem(ctx context.Context, userID uint, id uint) error { item,err:=s.itemRepo.GetByID(ctx,id); if err!=nil { return err }; set,err:=s.setRepo.GetByID(ctx,item.SetID); if err!=nil { return err }; if set.UserID!=userID { return fmt.Errorf("not found") }; return s.itemRepo.Delete(ctx,id) }

