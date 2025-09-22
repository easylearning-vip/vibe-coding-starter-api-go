package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// Part3AnswerOptionService Part3AnswerOption服务接口
type Part3AnswerOptionService interface {
	Create(ctx context.Context, req *CreatePart3AnswerOptionRequest) (*model.Part3AnswerOption, error)
	GetByID(ctx context.Context, id uint) (*model.Part3AnswerOption, error)
	Update(ctx context.Context, id uint, req *UpdatePart3AnswerOptionRequest) (*model.Part3AnswerOption, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListPart3AnswerOptionOptions) ([]*model.Part3AnswerOption, int64, error)
}

// part3AnswerOptionService Part3AnswerOption服务实现
type part3AnswerOptionService struct {
	part3AnswerOptionRepo repository.Part3AnswerOptionRepository

	logger logger.Logger
}

// NewPart3AnswerOptionService 创建Part3AnswerOption服务
func NewPart3AnswerOptionService(
	part3AnswerOptionRepo repository.Part3AnswerOptionRepository,

	logger logger.Logger,
) Part3AnswerOptionService {
	return &part3AnswerOptionService{
		part3AnswerOptionRepo: part3AnswerOptionRepo,

		logger: logger,
	}
}

// CreatePart3AnswerOptionRequest 创建Part3AnswerOption请求
type CreatePart3AnswerOptionRequest struct {
	ConversationId int32          `json:"conversation_id" validate:"required,min=0"`
	QuestionNumber int32          `json:"question_number" validate:"required,min=0"`
	OptionA        string         `json:"option_a" validate:"required,min=1,max=255"`
	OptionB        string         `json:"option_b" validate:"required,min=1,max=255"`
	OptionC        string         `json:"option_c" validate:"required,min=1,max=255"`
	OptionD        string         `json:"option_d" validate:"required,min=1,max=255"`
	CorrectAnswer  sql.NullString `json:"correct_answer" validate:"required"`
}

// UpdatePart3AnswerOptionRequest 更新Part3AnswerOption请求
type UpdatePart3AnswerOptionRequest struct {
	ConversationId *int32          `json:"conversation_id,omitempty" validate:"omitempty,min=0"`
	QuestionNumber *int32          `json:"question_number,omitempty" validate:"omitempty,min=0"`
	OptionA        *string         `json:"option_a,omitempty" validate:"omitempty,min=1,max=255"`
	OptionB        *string         `json:"option_b,omitempty" validate:"omitempty,min=1,max=255"`
	OptionC        *string         `json:"option_c,omitempty" validate:"omitempty,min=1,max=255"`
	OptionD        *string         `json:"option_d,omitempty" validate:"omitempty,min=1,max=255"`
	CorrectAnswer  *sql.NullString `json:"correct_answer,omitempty" validate:"omitempty"`
}

// ListPart3AnswerOptionOptions 列表查询选项
type ListPart3AnswerOptionOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Part3AnswerOption
func (s *part3AnswerOptionService) Create(ctx context.Context, req *CreatePart3AnswerOptionRequest) (*model.Part3AnswerOption, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Part3AnswerOption{
		ConversationId: req.ConversationId,
		QuestionNumber: req.QuestionNumber,
		OptionA:        req.OptionA,
		OptionB:        req.OptionB,
		OptionC:        req.OptionC,
		OptionD:        req.OptionD,
		CorrectAnswer:  req.CorrectAnswer,
	}

	// 保存到数据库
	if err := s.part3AnswerOptionRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create part3answeroption", "error", err)
		return nil, fmt.Errorf("failed to create part3answeroption: %w", err)
	}

	s.logger.Info("Part3AnswerOption created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Part3AnswerOption
func (s *part3AnswerOptionService) GetByID(ctx context.Context, id uint) (*model.Part3AnswerOption, error) {

	// 从数据库获取
	entity, err := s.part3AnswerOptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part3answeroption: %w", err)
	}

	return entity, nil
}

// Update 更新Part3AnswerOption
func (s *part3AnswerOptionService) Update(ctx context.Context, id uint, req *UpdatePart3AnswerOptionRequest) (*model.Part3AnswerOption, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.part3AnswerOptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part3answeroption: %w", err)
	}

	// 更新字段
	if req.ConversationId != nil {
		entity.ConversationId = *req.ConversationId
	}
	if req.QuestionNumber != nil {
		entity.QuestionNumber = *req.QuestionNumber
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
	if req.OptionD != nil {
		entity.OptionD = *req.OptionD
	}
	if req.CorrectAnswer != nil {
		entity.CorrectAnswer = *req.CorrectAnswer
	}

	// 保存更新
	if err := s.part3AnswerOptionRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update part3answeroption", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update part3answeroption: %w", err)
	}

	s.logger.Info("Part3AnswerOption updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Part3AnswerOption
func (s *part3AnswerOptionService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.part3AnswerOptionRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get part3answeroption: %w", err)
	}

	// 删除实体
	if err := s.part3AnswerOptionRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete part3answeroption", "id", id, "error", err)
		return fmt.Errorf("failed to delete part3answeroption: %w", err)
	}

	s.logger.Info("Part3AnswerOption deleted successfully", "id", id)
	return nil
}

// List 获取Part3AnswerOption列表
func (s *part3AnswerOptionService) List(ctx context.Context, opts *ListPart3AnswerOptionOptions) ([]*model.Part3AnswerOption, int64, error) {
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
	entities, total, err := s.part3AnswerOptionRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list part3answeroptions", "error", err)
		return nil, 0, fmt.Errorf("failed to list part3answeroptions: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *part3AnswerOptionService) validateCreateRequest(req *CreatePart3AnswerOptionRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *part3AnswerOptionService) validateUpdateRequest(req *UpdatePart3AnswerOptionRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
