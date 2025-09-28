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

type Part3PracticeSetService interface {
	CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part3PracticeSet, error)
	GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeSet, error)
	UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part3PracticeSet, error)
	DeleteSet(ctx context.Context, userID uint, id uint) error
	ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part3PracticeSet, int64, error)

	CreateItem(ctx context.Context, userID uint, req *CreatePart3SetItemRequest) (*model.Part3PracticeSetItem, error)
	ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part3PracticeSetItem, int64, error)
	DeleteItem(ctx context.Context, userID uint, id uint) error

	// New capabilities
	AutoGenerateSet(ctx context.Context, userID uint, req *GeneratePracticeSetRequest) (*model.Part3PracticeSet, []*model.Part3PracticeSetItem, error)
	ListSetQuestions(ctx context.Context, userID uint, setID uint) ([]Part3QuestionDetail, error)
	GetItemDetail(ctx context.Context, userID uint, itemID uint) (*Part3QuestionDetail, error)
	SubmitAnswer(ctx context.Context, userID uint, itemID uint, req *SubmitPart3AnswerRequest) error
}

type part3PracticeSetService struct {
	setRepo  repository.Part3PracticeSetRepository
	itemRepo repository.Part3PracticeSetItemRepository
	convRepo repository.Part3ConversationRepository
}

func NewPart3PracticeSetService(setRepo repository.Part3PracticeSetRepository, itemRepo repository.Part3PracticeSetItemRepository, convRepo repository.Part3ConversationRepository) Part3PracticeSetService {
	return &part3PracticeSetService{setRepo: setRepo, itemRepo: itemRepo, convRepo: convRepo}
}

func (s *part3PracticeSetService) CreateSet(ctx context.Context, userID uint, req *CreatePracticeSetRequest) (*model.Part3PracticeSet, error) {
	e := &model.Part3PracticeSet{UserID: userID, ScenarioId: sql.NullInt32{Int32: req.ScenarioId, Valid: req.ScenarioId != 0}, DifficultyLevelId: sql.NullInt32{Int32: req.DifficultyLevelId, Valid: req.DifficultyLevelId != 0}, TotalQuestions: req.TotalQuestions}
	if err := s.setRepo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
func (s *part3PracticeSetService) GetSetByID(ctx context.Context, userID uint, id uint) (*model.Part3PracticeSet, error) {
	e, err := s.setRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	return e, nil
}
func (s *part3PracticeSetService) UpdateSet(ctx context.Context, userID uint, id uint, req *UpdatePracticeSetRequest) (*model.Part3PracticeSet, error) {
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
func (s *part3PracticeSetService) DeleteSet(ctx context.Context, userID uint, id uint) error {
	e, err := s.setRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if e.UserID != userID {
		return fmt.Errorf("not found")
	}
	return s.setRepo.Delete(ctx, id)
}
func (s *part3PracticeSetService) ListSets(ctx context.Context, userID uint, opts *ListPracticeSetOptions) ([]*model.Part3PracticeSet, int64, error) {
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

type CreatePart3SetItemRequest struct {
	SetID          uint  `json:"set_id"`
	ConversationID uint  `json:"conversation_id"`
	QuestionIndex  int32 `json:"question_index"`
	OrderIndex     int32 `json:"order_index"`
}

func (s *part3PracticeSetService) CreateItem(ctx context.Context, userID uint, req *CreatePart3SetItemRequest) (*model.Part3PracticeSetItem, error) {
	set, err := s.setRepo.GetByID(ctx, req.SetID)
	if err != nil {
		return nil, err
	}
	if set.UserID != userID {
		return nil, fmt.Errorf("not found")
	}
	item := &model.Part3PracticeSetItem{SetID: req.SetID, ConversationID: req.ConversationID, QuestionIndex: req.QuestionIndex, OrderIndex: req.OrderIndex}
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
func (s *part3PracticeSetService) ListItems(ctx context.Context, userID uint, setID uint, opts *ListPracticeSetItemOptions) ([]*model.Part3PracticeSetItem, int64, error) {
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
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	opts.Filters["set_id"] = setID
	repoOpts := repository.ListOptions{Page: opts.Page, PageSize: opts.PageSize, Sort: opts.Sort, Order: opts.Order, Filters: opts.Filters}
	return s.itemRepo.List(ctx, repoOpts)
}
func (s *part3PracticeSetService) DeleteItem(ctx context.Context, userID uint, id uint) error {
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

// New DTOs for details and submission

type Part3QuestionDetail struct {
	ItemID        uint                     `json:"item_id"`
	SetID         uint                     `json:"set_id"`
	OrderIndex    int32                    `json:"order_index"`
	QuestionIndex int32                    `json:"question_index"`
	Conversation  *model.Part3Conversation `json:"conversation"`
}

type SubmitPart3AnswerRequest struct {
	SelectedAnswer string `json:"selected_answer"` // "A"|"B"|"C"
	IsCorrect      bool   `json:"is_correct"`
}

// Auto-generate
func (s *part3PracticeSetService) AutoGenerateSet(ctx context.Context, userID uint, req *GeneratePracticeSetRequest) (*model.Part3PracticeSet, []*model.Part3PracticeSetItem, error) {
	if req.TotalQuestions < 10 || req.TotalQuestions > 30 {
		return nil, nil, errors.New("total_questions must be between 10 and 30")
	}
	mode := req.Mode
	if mode == "" {
		mode = "sequential"
	}

	set, err := s.CreateSet(ctx, userID, &CreatePracticeSetRequest{ScenarioId: req.ScenarioId, DifficultyLevelId: req.DifficultyLevelId, TotalQuestions: req.TotalQuestions})
	if err != nil {
		return nil, nil, err
	}

	// list conversations
	filters := map[string]interface{}{}
	if req.ScenarioId != 0 {
		filters["scenario_id"] = req.ScenarioId
	}
	if req.DifficultyLevelId != 0 {
		filters["difficulty_level_id"] = req.DifficultyLevelId
	}
	convs, _, err := s.convRepo.List(ctx, repository.ListOptions{Page: 1, PageSize: int(req.TotalQuestions)*2 + 50, Sort: "conversation_number", Order: "asc", Filters: filters})
	if err != nil {
		return nil, nil, err
	}
	if len(convs) == 0 {
		return set, nil, nil
	}

	// build flat list of (conversation, qIndex)
	type cq struct {
		c   *model.Part3Conversation
		idx int32
	}
	flat := make([]cq, 0, len(convs)*3)
	for _, c := range convs {
		for i := int32(1); i <= 3; i++ {
			flat = append(flat, cq{c: c, idx: i})
		}
	}

	candidates := make([]cq, 0, len(flat))
	for _, x := range flat {
		candidates = append(candidates, x)
	}
	if mode == "random" {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	}
	n := int(req.TotalQuestions)
	if len(candidates) < n {
		n = len(candidates)
	}

	created := make([]*model.Part3PracticeSetItem, 0, n)
	for i := 0; i < n; i++ {
		it := &model.Part3PracticeSetItem{SetID: set.ID, ConversationID: candidates[i].c.ID, QuestionIndex: candidates[i].idx, OrderIndex: int32(i + 1)}
		if err := s.itemRepo.Create(ctx, it); err != nil {
			return set, created, err
		}
		created = append(created, it)
	}
	return set, created, nil
}

func (s *part3PracticeSetService) ListSetQuestions(ctx context.Context, userID uint, setID uint) ([]Part3QuestionDetail, error) {
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
	res := make([]Part3QuestionDetail, 0, len(items))
	for _, it := range items {
		convs, _, err := s.convRepo.ListWithAnswerOptions(ctx, repository.ListOptions{Page: 1, PageSize: 1, Filters: map[string]interface{}{"id": int32(it.ConversationID)}})
		if err != nil {
			return nil, err
		}
		if len(convs) == 0 {
			return nil, fmt.Errorf("not found")
		}
		conv := convs[0]
		res = append(res, Part3QuestionDetail{ItemID: it.ID, SetID: it.SetID, OrderIndex: it.OrderIndex, QuestionIndex: it.QuestionIndex, Conversation: conv})
	}
	return res, nil
}

func (s *part3PracticeSetService) GetItemDetail(ctx context.Context, userID uint, itemID uint) (*Part3QuestionDetail, error) {
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
	convs, _, err := s.convRepo.ListWithAnswerOptions(ctx, repository.ListOptions{Page: 1, PageSize: 1, Filters: map[string]interface{}{"id": int32(it.ConversationID)}})
	if err != nil {
		return nil, err
	}
	if len(convs) == 0 {
		return nil, fmt.Errorf("not found")
	}
	conv := convs[0]
	return &Part3QuestionDetail{ItemID: it.ID, SetID: it.SetID, OrderIndex: it.OrderIndex, QuestionIndex: it.QuestionIndex, Conversation: conv}, nil
}

func (s *part3PracticeSetService) SubmitAnswer(ctx context.Context, userID uint, itemID uint, req *SubmitPart3AnswerRequest) error {
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
	// update set progress
	set.CompletedCount += 1
	if req.IsCorrect {
		set.CorrectCount += 1
	}
	if set.TotalQuestions > 0 {
		set.Accuracy = float64(set.CorrectCount) / float64(set.TotalQuestions)
	}
	return s.setRepo.Update(ctx, set)
}
