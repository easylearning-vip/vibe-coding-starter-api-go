package service

import (
    "context"
    "database/sql"
    "fmt"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/repository"
    "vibe-coding-starter/pkg/logger"
)

// Part3PracticeItemService defines operations for user's Part3 practice records
type Part3PracticeItemService interface {
    Create(ctx context.Context, userID uint, req *CreatePart3PracticeItemRequest) (*model.Part3PracticeItem, error)
    GetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeItem, error)
    Update(ctx context.Context, userID uint, id uint, req *UpdatePart3PracticeItemRequest) (*model.Part3PracticeItem, error)
    Delete(ctx context.Context, userID uint, id uint) error
    List(ctx context.Context, userID uint, opts *ListPart3PracticeItemOptions) ([]*model.Part3PracticeItem, int64, error)
    Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error)
}

type part3PracticeItemService struct {
    repo   repository.Part3PracticeItemRepository
    logger logger.Logger
}

func NewPart3PracticeItemService(repo repository.Part3PracticeItemRepository, logger logger.Logger) Part3PracticeItemService {
    return &part3PracticeItemService{repo: repo, logger: logger}
}

// DTOs

type CreatePart3PracticeItemRequest struct {
    ConversationID    uint  `json:"conversation_id" validate:"required,min=1"`
    QuestionIndex     int32 `json:"question_index" validate:"required,min=1,max=3"`
    ScenarioId        int32 `json:"scenario_id"`
    DifficultyLevelId int32 `json:"difficulty_level_id"`
    SelectedOptionId  *int32 `json:"selected_option_id"`
    IsCorrect         bool   `json:"is_correct"`
}

type UpdatePart3PracticeItemRequest struct {
    SelectedOptionId  *int32 `json:"selected_option_id"`
    IsCorrect         *bool  `json:"is_correct"`
    ScenarioId        *int32 `json:"scenario_id"`
    DifficultyLevelId *int32 `json:"difficulty_level_id"`
}

type ListPart3PracticeItemOptions struct {
    Page     int
    PageSize int
    Sort     string
    Order    string
    Filters  map[string]interface{}
}

// Impl

func (s *part3PracticeItemService) Create(ctx context.Context, userID uint, req *CreatePart3PracticeItemRequest) (*model.Part3PracticeItem, error) {
    e := &model.Part3PracticeItem{
        UserID:         userID,
        ConversationID: req.ConversationID,
        QuestionIndex:  req.QuestionIndex,
        ScenarioId:     sql.NullInt32{Int32: req.ScenarioId, Valid: req.ScenarioId != 0},
        DifficultyLevelId: sql.NullInt32{Int32: req.DifficultyLevelId, Valid: req.DifficultyLevelId != 0},
        SelectedOptionId: sql.NullInt32{},
        IsCorrect:      req.IsCorrect,
    }
    if req.SelectedOptionId != nil { e.SelectedOptionId = sql.NullInt32{Int32: *req.SelectedOptionId, Valid: *req.SelectedOptionId != 0} }
    if err := s.repo.Create(ctx, e); err != nil { return nil, fmt.Errorf("create part3 practice item failed: %w", err) }
    return e, nil
}

func (s *part3PracticeItemService) GetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeItem, error) {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }
    if e.UserID != userID { return nil, fmt.Errorf("not found") }
    return e, nil
}

func (s *part3PracticeItemService) Update(ctx context.Context, userID uint, id uint, req *UpdatePart3PracticeItemRequest) (*model.Part3PracticeItem, error) {
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

func (s *part3PracticeItemService) Delete(ctx context.Context, userID uint, id uint) error {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return err }
    if e.UserID != userID { return fmt.Errorf("not found") }
    return s.repo.Delete(ctx, id)
}

func (s *part3PracticeItemService) List(ctx context.Context, userID uint, opts *ListPart3PracticeItemOptions) ([]*model.Part3PracticeItem, int64, error) {
    if opts == nil { opts = &ListPart3PracticeItemOptions{} }
    if opts.Page <= 0 { opts.Page = 1 }
    if opts.PageSize <= 0 { opts.PageSize = 20 }
    if opts.PageSize > 100 { opts.PageSize = 100 }
    if opts.Filters == nil { opts.Filters = map[string]interface{}{} }
    opts.Filters["user_id"] = userID
    repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
    return s.repo.List(ctx, repoOpts)
}

func (s *part3PracticeItemService) Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error) {
    if filters == nil { filters = map[string]interface{}{} }
    filters["user_id"] = userID

    // total
    totalFilters := make(map[string]interface{})
    for k,v := range filters { totalFilters[k] = v }
    repoOptsAll := repository.ListOptions{Filters: totalFilters}
    _, total, err := s.repo.List(ctx, repoOptsAll)
    if err != nil { return nil, err }

    // correct
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

