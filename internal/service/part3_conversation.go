package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// Part3ConversationService Part3Conversation服务接口
type Part3ConversationService interface {
	Create(ctx context.Context, req *CreatePart3ConversationRequest) (*model.Part3Conversation, error)
	GetByID(ctx context.Context, id uint) (*model.Part3Conversation, error)
	Update(ctx context.Context, id uint, req *UpdatePart3ConversationRequest) (*model.Part3Conversation, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListPart3ConversationOptions) ([]*model.Part3Conversation, int64, error)
}

// part3ConversationService Part3Conversation服务实现
type part3ConversationService struct {
	part3ConversationRepo repository.Part3ConversationRepository

	logger logger.Logger
}

// NewPart3ConversationService 创建Part3Conversation服务
func NewPart3ConversationService(
	part3ConversationRepo repository.Part3ConversationRepository,

	logger logger.Logger,
) Part3ConversationService {
	return &part3ConversationService{
		part3ConversationRepo: part3ConversationRepo,

		logger: logger,
	}
}

// CreatePart3ConversationRequest 创建Part3Conversation请求
type CreatePart3ConversationRequest struct {
	TestId             int32          `json:"test_id" validate:"required,min=0"`
	ConversationNumber int32          `json:"conversation_number" validate:"required,min=0"`
	Title              sql.NullString `json:"title" validate:"required"`
	Content            string         `json:"content" validate:"required,min=1,max=255"`
	Question1          string         `json:"question1" validate:"required,min=1,max=255"`
	Question2          string         `json:"question2" validate:"required,min=1,max=255"`
	Question3          string         `json:"question3" validate:"required,min=1,max=255"`
	ScenarioId         sql.NullInt32  `json:"scenario_id" validate:"required"`
	DifficultyLevelId  sql.NullInt32  `json:"difficulty_level_id" validate:"required"`
}

// UpdatePart3ConversationRequest 更新Part3Conversation请求
type UpdatePart3ConversationRequest struct {
	TestId             *int32          `json:"test_id,omitempty" validate:"omitempty,min=0"`
	ConversationNumber *int32          `json:"conversation_number,omitempty" validate:"omitempty,min=0"`
	Title              *sql.NullString `json:"title,omitempty" validate:"omitempty"`
	Content            *string         `json:"content,omitempty" validate:"omitempty,min=1,max=255"`
	Question1          *string         `json:"question1,omitempty" validate:"omitempty,min=1,max=255"`
	Question2          *string         `json:"question2,omitempty" validate:"omitempty,min=1,max=255"`
	Question3          *string         `json:"question3,omitempty" validate:"omitempty,min=1,max=255"`
	ScenarioId         *sql.NullInt32  `json:"scenario_id,omitempty" validate:"omitempty"`
	DifficultyLevelId  *sql.NullInt32  `json:"difficulty_level_id,omitempty" validate:"omitempty"`
}

// ListPart3ConversationOptions 列表查询选项
type ListPart3ConversationOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Part3Conversation
func (s *part3ConversationService) Create(ctx context.Context, req *CreatePart3ConversationRequest) (*model.Part3Conversation, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Part3Conversation{
		TestId:             req.TestId,
		ConversationNumber: req.ConversationNumber,
		Title:              req.Title,
		Content:            req.Content,
		Question1:          req.Question1,
		Question2:          req.Question2,
		Question3:          req.Question3,
		ScenarioId:         req.ScenarioId,
		DifficultyLevelId:  req.DifficultyLevelId,
	}

	// 保存到数据库
	if err := s.part3ConversationRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create part3conversation", "error", err)
		return nil, fmt.Errorf("failed to create part3conversation: %w", err)
	}

	s.logger.Info("Part3Conversation created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Part3Conversation
func (s *part3ConversationService) GetByID(ctx context.Context, id uint) (*model.Part3Conversation, error) {

	// 从数据库获取
	entity, err := s.part3ConversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part3conversation: %w", err)
	}

	return entity, nil
}

// Update 更新Part3Conversation
func (s *part3ConversationService) Update(ctx context.Context, id uint, req *UpdatePart3ConversationRequest) (*model.Part3Conversation, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.part3ConversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part3conversation: %w", err)
	}

	// 更新字段
	if req.TestId != nil {
		entity.TestId = *req.TestId
	}
	if req.ConversationNumber != nil {
		entity.ConversationNumber = *req.ConversationNumber
	}
	if req.Title != nil {
		entity.Title = *req.Title
	}
	if req.Content != nil {
		entity.Content = *req.Content
	}
	if req.Question1 != nil {
		entity.Question1 = *req.Question1
	}
	if req.Question2 != nil {
		entity.Question2 = *req.Question2
	}
	if req.Question3 != nil {
		entity.Question3 = *req.Question3
	}
	if req.ScenarioId != nil {
		entity.ScenarioId = *req.ScenarioId
	}
	if req.DifficultyLevelId != nil {
		entity.DifficultyLevelId = *req.DifficultyLevelId
	}

	// 保存更新
	if err := s.part3ConversationRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update part3conversation", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update part3conversation: %w", err)
	}

	s.logger.Info("Part3Conversation updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Part3Conversation
func (s *part3ConversationService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.part3ConversationRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get part3conversation: %w", err)
	}

	// 删除实体
	if err := s.part3ConversationRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete part3conversation", "id", id, "error", err)
		return fmt.Errorf("failed to delete part3conversation: %w", err)
	}

	s.logger.Info("Part3Conversation deleted successfully", "id", id)
	return nil
}

// List 获取Part3Conversation列表
func (s *part3ConversationService) List(ctx context.Context, opts *ListPart3ConversationOptions) ([]*model.Part3Conversation, int64, error) {
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
	entities, total, err := s.part3ConversationRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list part3conversations", "error", err)
		return nil, 0, fmt.Errorf("failed to list part3conversations: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *part3ConversationService) validateCreateRequest(req *CreatePart3ConversationRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *part3ConversationService) validateUpdateRequest(req *UpdatePart3ConversationRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
