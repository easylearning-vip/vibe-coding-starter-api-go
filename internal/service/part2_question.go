package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// Part2QuestionService Part2Question服务接口
type Part2QuestionService interface {
	Create(ctx context.Context, req *CreatePart2QuestionRequest) (*model.Part2Question, error)
	GetByID(ctx context.Context, id uint) (*model.Part2Question, error)
	Update(ctx context.Context, id uint, req *UpdatePart2QuestionRequest) (*model.Part2Question, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListPart2QuestionOptions) ([]*model.Part2Question, int64, error)
}

// part2QuestionService Part2Question服务实现
type part2QuestionService struct {
	part2QuestionRepo repository.Part2QuestionRepository

	logger logger.Logger
}

// NewPart2QuestionService 创建Part2Question服务
func NewPart2QuestionService(
	part2QuestionRepo repository.Part2QuestionRepository,

	logger logger.Logger,
) Part2QuestionService {
	return &part2QuestionService{
		part2QuestionRepo: part2QuestionRepo,

		logger: logger,
	}
}

// CreatePart2QuestionRequest 创建Part2Question请求
type CreatePart2QuestionRequest struct {
	TestId            int32          `json:"test_id" validate:"required,min=0"`
	QuestionNumber    int32          `json:"question_number" validate:"required,min=0"`
	QuestionText      string         `json:"question_text" validate:"required,min=1,max=255"`
	OptionA           string         `json:"option_a" validate:"required,min=1,max=255"`
	OptionB           string         `json:"option_b" validate:"required,min=1,max=255"`
	OptionC           string         `json:"option_c" validate:"required,min=1,max=255"`
	CorrectAnswer     sql.NullString `json:"correct_answer" validate:"required"`
	ScenarioId        sql.NullInt32  `json:"scenario_id" validate:"required"`
	DifficultyLevelId sql.NullInt32  `json:"difficulty_level_id" validate:"required"`
}

// UpdatePart2QuestionRequest 更新Part2Question请求
type UpdatePart2QuestionRequest struct {
	TestId            *int32          `json:"test_id,omitempty" validate:"omitempty,min=0"`
	QuestionNumber    *int32          `json:"question_number,omitempty" validate:"omitempty,min=0"`
	QuestionText      *string         `json:"question_text,omitempty" validate:"omitempty,min=1,max=255"`
	OptionA           *string         `json:"option_a,omitempty" validate:"omitempty,min=1,max=255"`
	OptionB           *string         `json:"option_b,omitempty" validate:"omitempty,min=1,max=255"`
	OptionC           *string         `json:"option_c,omitempty" validate:"omitempty,min=1,max=255"`
	CorrectAnswer     *sql.NullString `json:"correct_answer,omitempty" validate:"omitempty"`
	ScenarioId        *sql.NullInt32  `json:"scenario_id,omitempty" validate:"omitempty"`
	DifficultyLevelId *sql.NullInt32  `json:"difficulty_level_id,omitempty" validate:"omitempty"`
}

// ListPart2QuestionOptions 列表查询选项
type ListPart2QuestionOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Part2Question
func (s *part2QuestionService) Create(ctx context.Context, req *CreatePart2QuestionRequest) (*model.Part2Question, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Part2Question{
		TestId:            req.TestId,
		QuestionNumber:    req.QuestionNumber,
		QuestionText:      req.QuestionText,
		OptionA:           req.OptionA,
		OptionB:           req.OptionB,
		OptionC:           req.OptionC,
		CorrectAnswer:     req.CorrectAnswer,
		ScenarioId:        req.ScenarioId,
		DifficultyLevelId: req.DifficultyLevelId,
	}

	// 保存到数据库
	if err := s.part2QuestionRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create part2question", "error", err)
		return nil, fmt.Errorf("failed to create part2question: %w", err)
	}

	s.logger.Info("Part2Question created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Part2Question
func (s *part2QuestionService) GetByID(ctx context.Context, id uint) (*model.Part2Question, error) {

	// 从数据库获取
	entity, err := s.part2QuestionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part2question: %w", err)
	}

	return entity, nil
}

// Update 更新Part2Question
func (s *part2QuestionService) Update(ctx context.Context, id uint, req *UpdatePart2QuestionRequest) (*model.Part2Question, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.part2QuestionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part2question: %w", err)
	}

	// 更新字段
	if req.TestId != nil {
		entity.TestId = *req.TestId
	}
	if req.QuestionNumber != nil {
		entity.QuestionNumber = *req.QuestionNumber
	}
	if req.QuestionText != nil {
		entity.QuestionText = *req.QuestionText
	}
	if req.OptionA != nil {
		entity.OptionA = *req.OptionA
	}
	if req.OptionB != nil {
		entity.OptionB = *req.OptionB
	}
	if req.OptionC != nil {
		entity.OptionC = *req.OptionC
	}
	if req.CorrectAnswer != nil {
		entity.CorrectAnswer = *req.CorrectAnswer
	}
	if req.ScenarioId != nil {
		entity.ScenarioId = *req.ScenarioId
	}
	if req.DifficultyLevelId != nil {
		entity.DifficultyLevelId = *req.DifficultyLevelId
	}

	// 保存更新
	if err := s.part2QuestionRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update part2question", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update part2question: %w", err)
	}

	s.logger.Info("Part2Question updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Part2Question
func (s *part2QuestionService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.part2QuestionRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get part2question: %w", err)
	}

	// 删除实体
	if err := s.part2QuestionRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete part2question", "id", id, "error", err)
		return fmt.Errorf("failed to delete part2question: %w", err)
	}

	s.logger.Info("Part2Question deleted successfully", "id", id)
	return nil
}

// List 获取Part2Question列表
func (s *part2QuestionService) List(ctx context.Context, opts *ListPart2QuestionOptions) ([]*model.Part2Question, int64, error) {
	// 设置默认值
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}

	// 转换为仓储选项
	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	}

	// 获取列表
	entities, total, err := s.part2QuestionRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list part2questions", "error", err)
		return nil, 0, fmt.Errorf("failed to list part2questions: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *part2QuestionService) validateCreateRequest(req *CreatePart2QuestionRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *part2QuestionService) validateUpdateRequest(req *UpdatePart2QuestionRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
