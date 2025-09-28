package service

import (
	"context"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// ToeicAiPromptService ToeicAiPrompt服务接口
type ToeicAiPromptService interface {
	Create(ctx context.Context, req *CreateToeicAiPromptRequest) (*model.ToeicAiPrompt, error)
	GetByID(ctx context.Context, id uint) (*model.ToeicAiPrompt, error)
	Update(ctx context.Context, id uint, req *UpdateToeicAiPromptRequest) (*model.ToeicAiPrompt, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListToeicAiPromptOptions) ([]*model.ToeicAiPrompt, int64, error)
}

// toeicAiPromptService ToeicAiPrompt服务实现
type toeicAiPromptService struct {
	toeicAiPromptRepo repository.ToeicAiPromptRepository
	logger            logger.Logger
}

// NewToeicAiPromptService 创建ToeicAiPrompt服务
func NewToeicAiPromptService(
	toeicAiPromptRepo repository.ToeicAiPromptRepository,
	logger logger.Logger,
) ToeicAiPromptService {
	return &toeicAiPromptService{
		toeicAiPromptRepo: toeicAiPromptRepo,
		logger:            logger,
	}
}

// CreateToeicAiPromptRequest 创建ToeicAiPrompt请求
type CreateToeicAiPromptRequest struct {
	Part2Prompt string `json:"part2_prompt" validate:"required"`
	Part3Prompt string `json:"part3_prompt" validate:"required"`
	Part4Prompt string `json:"part4_prompt" validate:"required"`
}

// UpdateToeicAiPromptRequest 更新ToeicAiPrompt请求
type UpdateToeicAiPromptRequest struct {
	Part2Prompt string `json:"part2_prompt" validate:"required"`
	Part3Prompt string `json:"part3_prompt" validate:"required"`
	Part4Prompt string `json:"part4_prompt" validate:"required"`
}

// ListToeicAiPromptOptions 列表查询选项
type ListToeicAiPromptOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建ToeicAiPrompt
func (s *toeicAiPromptService) Create(ctx context.Context, req *CreateToeicAiPromptRequest) (*model.ToeicAiPrompt, error) {
	entity := &model.ToeicAiPrompt{
		Part2Prompt: req.Part2Prompt,
		Part3Prompt: req.Part3Prompt,
		Part4Prompt: req.Part4Prompt,
	}

	if err := s.toeicAiPromptRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create toeic ai prompt", "error", err)
		return nil, fmt.Errorf("failed to create toeic ai prompt: %w", err)
	}

	s.logger.Info("ToeicAiPrompt created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取ToeicAiPrompt
func (s *toeicAiPromptService) GetByID(ctx context.Context, id uint) (*model.ToeicAiPrompt, error) {
	entity, err := s.toeicAiPromptRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get toeic ai prompt by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get toeic ai prompt: %w", err)
	}

	return entity, nil
}

// Update 更新ToeicAiPrompt
func (s *toeicAiPromptService) Update(ctx context.Context, id uint, req *UpdateToeicAiPromptRequest) (*model.ToeicAiPrompt, error) {
	// 先获取现有实体
	entity, err := s.toeicAiPromptRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get toeic ai prompt for update", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get toeic ai prompt: %w", err)
	}

	// 更新字段
	entity.Part2Prompt = req.Part2Prompt
	entity.Part3Prompt = req.Part3Prompt
	entity.Part4Prompt = req.Part4Prompt

	if err := s.toeicAiPromptRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update toeic ai prompt", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update toeic ai prompt: %w", err)
	}

	s.logger.Info("ToeicAiPrompt updated successfully", "id", id)
	return entity, nil
}

// Delete 删除ToeicAiPrompt
func (s *toeicAiPromptService) Delete(ctx context.Context, id uint) error {
	// 先检查实体是否存在
	_, err := s.toeicAiPromptRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get toeic ai prompt for delete", "id", id, "error", err)
		return fmt.Errorf("failed to get toeic ai prompt: %w", err)
	}

	if err := s.toeicAiPromptRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete toeic ai prompt", "id", id, "error", err)
		return fmt.Errorf("failed to delete toeic ai prompt: %w", err)
	}

	s.logger.Info("ToeicAiPrompt deleted successfully", "id", id)
	return nil
}

// List 获取ToeicAiPrompt列表
func (s *toeicAiPromptService) List(ctx context.Context, opts *ListToeicAiPromptOptions) ([]*model.ToeicAiPrompt, int64, error) {
	if opts == nil {
		opts = &ListToeicAiPromptOptions{}
	}

	// 设置默认值
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	}

	entities, total, err := s.toeicAiPromptRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list toeic ai prompts", "error", err)
		return nil, 0, fmt.Errorf("failed to list toeic ai prompts: %w", err)
	}

	return entities, total, nil
}
