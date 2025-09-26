package service

import (
    "context"
    "database/sql"
    "fmt"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/repository"
    "vibe-coding-starter/pkg/logger"
)

// Part4PracticeItemService defines operations for user's Part4 practice records
type Part4PracticeItemService interface {
    Create(ctx context.Context, userID uint, req *CreatePart4PracticeItemRequest) (*model.Part4PracticeItem, error)
    GetByID(ctx context.Context, userID uint, id uint) (*model.Part4PracticeItem, error)
    Update(ctx context.Context, userID uint, id uint, req *UpdatePart4PracticeItemRequest) (*model.Part4PracticeItem, error)
    Delete(ctx context.Context, userID uint, id uint) error
    List(ctx context.Context, userID uint, opts *ListPart4PracticeItemOptions) ([]*model.Part4PracticeItem, int64, error)
    Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error)
}

type part4PracticeItemService struct {
    repo   repository.Part4PracticeItemRepository
    logger logger.Logger
}

func NewPart4PracticeItemService(repo repository.Part4PracticeItemRepository, logger logger.Logger) Part4PracticeItemService {
    return &part4PracticeItemService{repo: repo, logger: logger}
}

// DTOs

type CreatePart4PracticeItemRequest struct {
    TalkID            uint  `json:"talk_id" validate:"required,min=1"`
    QuestionIndex     int32 `json:"question_index" validate:"required,min=1,max=3"`
    ScenarioId        int32 `json:"scenario_id"`
    DifficultyLevelId int32 `json:"difficulty_level_id"`
    SelectedOptionId  *int32 `json:"selected_option_id"`
    IsCorrect         bool   `json:"is_correct"`
}

type UpdatePart4PracticeItemRequest struct {
    SelectedOptionId  *int32 `json:"selected_option_id"`
    IsCorrect         *bool  `json:"is_correct"`
    ScenarioId        *int32 `json:"scenario_id"`
    DifficultyLevelId *int32 `json:"difficulty_level_id"`
}

type ListPart4PracticeItemOptions struct {
    Page     int
    PageSize int
    Sort     string
    Order    string
    Filters  map[string]interface{}
}

// Impl

func (s *part4PracticeItemService) Create(ctx context.Context, userID uint, req *CreatePart4PracticeItemRequest) (*model.Part4PracticeItem, error) {
    e := &model.Part4PracticeItem{
        UserID:         userID,
        TalkID:         req.TalkID,
        QuestionIndex:  req.QuestionIndex,
        ScenarioId:     sql.NullInt32{Int32: req.ScenarioId, Valid: req.ScenarioId != 0},
        DifficultyLevelId: sql.NullInt32{Int32: req.DifficultyLevelId, Valid: req.DifficultyLevelId != 0},
        SelectedOptionId: sql.NullInt32{},
        IsCorrect:      req.IsCorrect,
    }
    if req.SelectedOptionId != nil { e.SelectedOptionId = sql.NullInt32{Int32: *req.SelectedOptionId, Valid: *req.SelectedOptionId != 0} }
    if err := s.repo.Create(ctx, e); err != nil { return nil, fmt.Errorf("create part4 practice item failed: %w", err) }
    return e, nil
}

func (s *part4PracticeItemService) GetByID(ctx context.Context, userID uint, id uint) (*model.Part4PracticeItem, error) {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }
    if e.UserID != userID { return nil, fmt.Errorf("not found") }
    return e, nil
}

func (s *part4PracticeItemService) Update(ctx context.Context, userID uint, id uint, req *UpdatePart4PracticeItemRequest) (*model.Part4PracticeItem, error) {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }
    if e.UserID != userID { return nil, fmt.Errorf("not found") }

    if req.SelectedOptionId != nil { e.SelectedOptionId = sql.NullInt32{Int32: *req.SelectedOptionId, Valid: *req.SelectedOptionId != 0} }
    if req.IsCorrect != nil { e.IsCorrect = *req.IsCorrect }
    if req.ScenarioId != nil { e.ScenarioId = sql.NullInt32{Int32: *req.ScenarioId, Valid: *req.ScenarioId != 0} }
    if req.DifficultyLevelId != nil { e.DifficultyLevelId = sql.NullInt32{Int32: *req.DifficultyLevelId, Valid: *req.DifficultyLevelId != 0} }

    if err := s.repo.Update(ctx, e); err != nil { return nil, fmt.Errorf("update failed: %w", err) }
    return e, nil
}

func (s *part4PracticeItemService) Delete(ctx context.Context, userID uint, id uint) error {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return err }
    if e.UserID != userID { return fmt.Errorf("not found") }
    return s.repo.Delete(ctx, id)
}

func (s *part4PracticeItemService) List(ctx context.Context, userID uint, opts *ListPart4PracticeItemOptions) ([]*model.Part4PracticeItem, int64, error) {
    if opts == nil { opts = &ListPart4PracticeItemOptions{} }
    if opts.Page <= 0 { opts.Page = 1 }
    if opts.PageSize <= 0 { opts.PageSize = 20 }
    if opts.PageSize > 100 { opts.PageSize = 100 }
    if opts.Filters == nil { opts.Filters = map[string]interface{}{} }
    opts.Filters["user_id"] = userID
    repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
    return s.repo.List(ctx, repoOpts)
}

func (s *part4PracticeItemService) Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error) {
    if filters == nil { filters = map[string]interface{}{} }
    filters["user_id"] = userID

    totalFilters := make(map[string]interface{})
    for k,v := range filters { totalFilters[k] = v }
    repoOptsAll := repository.ListOptions{Filters: totalFilters}
    _, total, err := s.repo.List(ctx, repoOptsAll)
    if err != nil { return nil, err }

    correctFilters := make(map[string]interface{})
    for k,v := range filters { correctFilters[k] = v }
    correctFilters["is_correct"] = true
    repoOptsCorrect := repository.ListOptions{Filters: correctFilters}
    _, correct, err := s.repo.List(ctx, repoOptsCorrect)
    if err != nil { return nil, err }

    acc := 0.0
    if total > 0 { acc = float64(correct) / float64(total) }
    return &PracticeStats{Total: total, Correct: correct, Accuracy: acc}, nil
}

