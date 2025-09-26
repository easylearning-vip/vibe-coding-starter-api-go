package service

import (
    "context"
    "database/sql"
    "fmt"

    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/repository"
    "vibe-coding-starter/pkg/logger"
)

// Part2PracticeItemService defines operations for user's Part2 practice records
type Part2PracticeItemService interface {
    Create(ctx context.Context, userID uint, req *CreatePart2PracticeItemRequest) (*model.Part2PracticeItem, error)
    GetByID(ctx context.Context, userID uint, id uint) (*model.Part2PracticeItem, error)
    Update(ctx context.Context, userID uint, id uint, req *UpdatePart2PracticeItemRequest) (*model.Part2PracticeItem, error)
    Delete(ctx context.Context, userID uint, id uint) error
    List(ctx context.Context, userID uint, opts *ListPart2PracticeItemOptions) ([]*model.Part2PracticeItem, int64, error)
    Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error)
}

type part2PracticeItemService struct {
    repo   repository.Part2PracticeItemRepository
    logger logger.Logger
}

func NewPart2PracticeItemService(repo repository.Part2PracticeItemRepository, logger logger.Logger) Part2PracticeItemService {
    return &part2PracticeItemService{repo: repo, logger: logger}
}

// DTOs

type CreatePart2PracticeItemRequest struct {
    QuestionID        uint  `json:"question_id" validate:"required,min=1"`
    ScenarioId        int32 `json:"scenario_id"`
    DifficultyLevelId int32 `json:"difficulty_level_id"`
    SelectedAnswer    string `json:"selected_answer" validate:"omitempty,max=1"`
    IsCorrect         bool   `json:"is_correct"`
}

type UpdatePart2PracticeItemRequest struct {
    SelectedAnswer    *string `json:"selected_answer" validate:"omitempty,max=1"`
    IsCorrect         *bool   `json:"is_correct"`
    // allow updating tags if needed
    ScenarioId        *int32  `json:"scenario_id"`
    DifficultyLevelId *int32  `json:"difficulty_level_id"`
}

type ListPart2PracticeItemOptions struct {
    Page     int
    PageSize int
    Sort     string
    Order    string
    Filters  map[string]interface{}
}

type PracticeStats struct {
    Total       int64   `json:"total"`
    Correct     int64   `json:"correct"`
    Accuracy    float64 `json:"accuracy"`
}

// Impl

func (s *part2PracticeItemService) Create(ctx context.Context, userID uint, req *CreatePart2PracticeItemRequest) (*model.Part2PracticeItem, error) {
    entity := &model.Part2PracticeItem{
        UserID:     userID,
        QuestionID: req.QuestionID,
        ScenarioId: sql.NullInt32{Int32: req.ScenarioId, Valid: req.ScenarioId != 0},
        DifficultyLevelId: sql.NullInt32{Int32: req.DifficultyLevelId, Valid: req.DifficultyLevelId != 0},
        SelectedAnswer: sql.NullString{String: req.SelectedAnswer, Valid: req.SelectedAnswer != ""},
        IsCorrect:  req.IsCorrect,
    }
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, fmt.Errorf("create part2 practice item failed: %w", err)
    }
    return entity, nil
}

func (s *part2PracticeItemService) GetByID(ctx context.Context, userID uint, id uint) (*model.Part2PracticeItem, error) {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }
    if e.UserID != userID { return nil, fmt.Errorf("not found") }
    return e, nil
}

func (s *part2PracticeItemService) Update(ctx context.Context, userID uint, id uint, req *UpdatePart2PracticeItemRequest) (*model.Part2PracticeItem, error) {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return nil, err }
    if e.UserID != userID { return nil, fmt.Errorf("not found") }

    if req.SelectedAnswer != nil {
        e.SelectedAnswer = sql.NullString{String: *req.SelectedAnswer, Valid: *req.SelectedAnswer != ""}
    }
    if req.IsCorrect != nil { e.IsCorrect = *req.IsCorrect }
    if req.ScenarioId != nil { e.ScenarioId = sql.NullInt32{Int32: *req.ScenarioId, Valid: *req.ScenarioId != 0} }
    if req.DifficultyLevelId != nil { e.DifficultyLevelId = sql.NullInt32{Int32: *req.DifficultyLevelId, Valid: *req.DifficultyLevelId != 0} }

    if err := s.repo.Update(ctx, e); err != nil { return nil, fmt.Errorf("update failed: %w", err) }
    return e, nil
}

func (s *part2PracticeItemService) Delete(ctx context.Context, userID uint, id uint) error {
    e, err := s.repo.GetByID(ctx, id)
    if err != nil { return err }
    if e.UserID != userID { return fmt.Errorf("not found") }
    return s.repo.Delete(ctx, id)
}

func (s *part2PracticeItemService) List(ctx context.Context, userID uint, opts *ListPart2PracticeItemOptions) ([]*model.Part2PracticeItem, int64, error) {
    if opts == nil { opts = &ListPart2PracticeItemOptions{} }
    if opts.Page <= 0 { opts.Page = 1 }
    if opts.PageSize <= 0 { opts.PageSize = 20 }
    if opts.PageSize > 100 { opts.PageSize = 100 }
    if opts.Filters == nil { opts.Filters = map[string]interface{}{} }
    // force user scope
    opts.Filters["user_id"] = userID

    repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
    return s.repo.List(ctx, repoOpts)
}

func (s *part2PracticeItemService) Stats(ctx context.Context, userID uint, filters map[string]interface{}) (*PracticeStats, error) {
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

