package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

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

	// New capabilities
	AutoGenerateSet(ctx context.Context, userID uint, req *GeneratePracticeSetRequest) (*model.Part2PracticeSet, []*model.Part2PracticeSetItem, error)
	ListSetQuestions(ctx context.Context, userID uint, setID uint) ([]Part2QuestionDetail, error)
	GetItemDetail(ctx context.Context, userID uint, itemID uint) (*Part2QuestionDetail, error)
	SubmitAnswer(ctx context.Context, userID uint, itemID uint, req *SubmitPart2AnswerRequest) error
}

type part2PracticeSetService struct {
	setRepo  repository.Part2PracticeSetRepository
	itemRepo repository.Part2PracticeSetItemRepository
	qRepo    repository.Part2QuestionRepository
}

func NewPart2PracticeSetService(
	setRepo repository.Part2PracticeSetRepository,
	itemRepo repository.Part2PracticeSetItemRepository,
	qRepo repository.Part2QuestionRepository,
) Part2PracticeSetService {
	return &part2PracticeSetService{setRepo: setRepo, itemRepo: itemRepo, qRepo: qRepo}
}

// Shared DTOs

type CreatePracticeSetRequest struct {
	ScenarioId        int32 `json:"scenario_id"`
	DifficultyLevelId int32 `json:"difficulty_level_id"`
	TotalQuestions    int32 `json:"total_questions"`
}

type GeneratePracticeSetRequest struct {
	ScenarioId        int32  `json:"scenario_id"`
	DifficultyLevelId int32  `json:"difficulty_level_id"`
	Mode              string `json:"mode"`        // sequential | random
	Deduplicate       *bool  `json:"deduplicate"` // default true
	TotalQuestions    int32  `json:"total_questions"`
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

type Part2QuestionDetail struct {
	ItemID     uint                 `json:"item_id"`
	SetID      uint                 `json:"set_id"`
	OrderIndex int32                `json:"order_index"`
	Question   *model.Part2Question `json:"question"`
}

type SubmitPart2AnswerRequest struct {
	SelectedAnswer string `json:"selected_answer"` // "A"|"B"|"C"
	IsCorrect      bool   `json:"is_correct"`
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

// Auto-generate a practice set and populate items from question bank
func (s *part2PracticeSetService) AutoGenerateSet(ctx context.Context, userID uint, req *GeneratePracticeSetRequest) (*model.Part2PracticeSet, []*model.Part2PracticeSetItem, error) {
	if req.TotalQuestions < 10 || req.TotalQuestions > 30 {
		return nil, nil, errors.New("total_questions must be between 10 and 30")
	}
	mode := req.Mode
	if mode == "" {
		mode = "sequential"
	}
	// create the set first
	set, err := s.CreateSet(ctx, userID, &CreatePracticeSetRequest{
		ScenarioId:        req.ScenarioId,
		DifficultyLevelId: req.DifficultyLevelId,
		TotalQuestions:    req.TotalQuestions,
	})
	if err != nil {
		return nil, nil, err
	}

	// gather candidate questions
	filters := map[string]interface{}{}
	if req.ScenarioId != 0 {
		filters["scenario_id"] = req.ScenarioId
	}
	if req.DifficultyLevelId != 0 {
		filters["difficulty_level_id"] = req.DifficultyLevelId
	}
	// large page to fetch enough candidates
	qs, _, err := s.qRepo.List(ctx, repository.ListOptions{Page: 1, PageSize: int(req.TotalQuestions)*5 + 100, Sort: "question_number", Order: "asc", Filters: filters})
	if err != nil {
		return nil, nil, err
	}
	if len(qs) == 0 {
		return set, nil, nil
	}

	// select questions
	candidates := make([]*model.Part2Question, 0, len(qs))
	for _, q := range qs {
		candidates = append(candidates, q)
	}
	if mode == "random" {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	}
	// trim to requested size
	n := int(req.TotalQuestions)
	if len(candidates) < n {
		n = len(candidates)
	}
	selected := candidates[:n]

	// create items with order_index 1..n
	created := make([]*model.Part2PracticeSetItem, 0, n)
	for i, q := range selected {
		it := &model.Part2PracticeSetItem{SetID: set.ID, QuestionID: q.ID, OrderIndex: int32(i + 1)}
		if err := s.itemRepo.Create(ctx, it); err != nil {
			return set, created, err
		}
		created = append(created, it)
	}
	return set, created, nil
}

// List all question details under a set
func (s *part2PracticeSetService) ListSetQuestions(ctx context.Context, userID uint, setID uint) ([]Part2QuestionDetail, error) {
	set, err := s.setRepo.GetByID(ctx, setID)
	if err != nil {
		return nil, err
	}
	if userID != 0 && set.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	items, _, err := s.itemRepo.List(ctx, repository.ListOptions{Page: 1, PageSize: 10000, Sort: "order_index", Order: "asc", Filters: map[string]interface{}{"set_id": setID}})
	if err != nil {
		return nil, err
	}
	res := make([]Part2QuestionDetail, 0, len(items))
	for _, it := range items {
		q, err := s.qRepo.GetByID(ctx, it.QuestionID)
		if err != nil {
			return nil, err
		}
		res = append(res, Part2QuestionDetail{ItemID: it.ID, SetID: it.SetID, OrderIndex: it.OrderIndex, Question: q})
	}
	return res, nil
}

// Get single question detail by set item id
func (s *part2PracticeSetService) GetItemDetail(ctx context.Context, userID uint, itemID uint) (*Part2QuestionDetail, error) {
	it, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	set, err := s.setRepo.GetByID(ctx, it.SetID)
	if err != nil {
		return nil, err
	}
	if userID != 0 && set.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	q, err := s.qRepo.GetByID(ctx, it.QuestionID)
	if err != nil {
		return nil, err
	}
	res := &Part2QuestionDetail{ItemID: it.ID, SetID: it.SetID, OrderIndex: it.OrderIndex, Question: q}
	return res, nil
}

// Submit answer result for a set item and record user practice history
func (s *part2PracticeSetService) SubmitAnswer(ctx context.Context, userID uint, itemID uint, req *SubmitPart2AnswerRequest) error {
	it, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	set, err := s.setRepo.GetByID(ctx, it.SetID)
	if err != nil {
		return err
	}
	if set.UserID != userID {
		return fmt.Errorf("not found")
	}
	// update item with answer
	it.SelectedAnswer = sql.NullString{String: req.SelectedAnswer, Valid: req.SelectedAnswer != ""}
	it.IsCorrect = req.IsCorrect
	if err := s.itemRepo.Update(ctx, it); err != nil {
		return err
	}
	// update set progress (best-effort)
	set.CompletedCount += 1
	if req.IsCorrect {
		set.CorrectCount += 1
	}
	if set.TotalQuestions > 0 {
		set.Accuracy = float64(set.CorrectCount) / float64(set.TotalQuestions)
	}
	return s.setRepo.Update(ctx, set)
}
