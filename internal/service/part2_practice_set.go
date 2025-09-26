package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
)

type Part2PracticeSetService interface {
	CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part2PracticeSet, error)
	GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part2PracticeSet, error)
	UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part2PracticeSet, error)
	DeleteSet(ctx context.Context, userID uint, id uint) error
	ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part2PracticeSet, int64, error)

	CreateItem(ctx context.Context, userID uint, req *CreatePart2SetItemRequest) (*model.Part2PracticeSetItem, error)
	ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part2PracticeSetItem, int64, error)
	DeleteItem(ctx context.Context, userID uint, id uint) error
}

type part2PracticeSetService struct {
	setRepo  repository.Part2PracticeSetRepository
	itemRepo repository.Part2PracticeSetItemRepository
}

func NewPart2PracticeSetService(setRepo repository.Part2PracticeSetRepository, itemRepo repository.Part2PracticeSetItemRepository) Part2PracticeSetService {
	return &part2PracticeSetService{setRepo: setRepo, itemRepo: itemRepo}
}

// Shared DTOs

type CreatePracticeSetRequest struct {
	ScenarioId        int32 `json:"scenario_id"`
	DifficultyLevelId int32 `json:"difficulty_level_id"`
	TotalQuestions    int32 `json:"total_questions"`
}

type UpdatePracticeSetRequest struct {
	TotalQuestions *int32 `json:"total_questions"`
	CompletedCount *int32 `json:"completed_count"`
	CorrectCount   *int32 `json:"correct_count"`
}

type ListPracticeSetOptions struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Filters  map[string]interface{}
}

type CreatePart2SetItemRequest struct {
	SetID      uint  `json:"set_id"`
	QuestionID uint  `json:"question_id"`
	OrderIndex int32 `json:"order_index"`
}

type ListPracticeSetItemOptions struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Filters  map[string]interface{}
}

func (s *part2PracticeSetService) CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part2PracticeSet, error) {
	e := &model.Part2PracticeSet{
		UserID:            userID,
		ScenarioId:        sql.NullInt32{Int32: req.ScenarioId, Valid: req.ScenarioId != 0},
		DifficultyLevelId: sql.NullInt32{Int32: req.DifficultyLevelId, Valid: req.DifficultyLevelId != 0},
		TotalQuestions:    req.TotalQuestions,
		CompletedCount:    0,
		CorrectCount:      0,
		Accuracy:          0,
	}
	if err := s.setRepo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *part2PracticeSetService) GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part2PracticeSet, error) {
	e, err := s.setRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	return e, nil
}

func (s *part2PracticeSetService) UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part2PracticeSet, error) {
	e, err := s.setRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	if req.TotalQuestions != nil {
		e.TotalQuestions = *req.TotalQuestions
	}
	if req.CompletedCount != nil {
		e.CompletedCount = *req.CompletedCount
	}
	if req.CorrectCount != nil {
		e.CorrectCount = *req.CorrectCount
	}
	if e.TotalQuestions > 0 {
		e.Accuracy = float64(e.CorrectCount) / float64(e.TotalQuestions)
	} else {
		e.Accuracy = 0
	}
	if err := s.setRepo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *part2PracticeSetService) DeleteSet(ctx context.Context, userID uint, id uint) error {
	e, err := s.setRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if e.UserID != userID {
		return fmt.Errorf("not found")
	}
	return s.setRepo.Delete(ctx, id)
}

func (s *part2PracticeSetService) ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part2PracticeSet, int64, error) {
	if opts == nil {
		opts = &ListPracticeSetOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = map[string]interface{}{}
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	opts.Filters["user_id"] = userID
	repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
	return s.setRepo.List(ctx, repoOpts)
}

func (s *part2PracticeSetService) CreateItem(ctx context.Context, userID uint, req *CreatePart2SetItemRequest) (*model.Part2PracticeSetItem, error) {
	set, err := s.setRepo.GetByID(ctx, req.SetID)
	if err != nil {
		return nil, err
	}
	if set.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	item := &model.Part2PracticeSetItem{SetID: req.SetID, QuestionID: req.QuestionID, OrderIndex: req.OrderIndex}
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *part2PracticeSetService) ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part2PracticeSetItem, int64, error) {
	set, err := s.setRepo.GetByID(ctx, setID)
	if err != nil {
		return nil, 0, err
	}
	if set.UserID != userID {
		return nil, 0, fmt.Errorf("not found")
	}
	if opts == nil {
		opts = &ListPracticeSetItemOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = map[string]interface{}{}
	}
	opts.Filters["set_id"] = setID
	// defaults
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
	return s.itemRepo.List(ctx, repoOpts)
}

func (s *part2PracticeSetService) DeleteItem(ctx context.Context, userID uint, id uint) error {
	// load to check ownership via set
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	set, err := s.setRepo.GetByID(ctx, item.SetID)
	if err != nil {
		return err
	}
	if set.UserID != userID {
		return fmt.Errorf("not found")
	}
	return s.itemRepo.Delete(ctx, id)
}
