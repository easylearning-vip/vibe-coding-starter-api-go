package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// Part4TalkService Part4Talk服务接口
type Part4TalkService interface {
	Create(ctx context.Context, req *CreatePart4TalkRequest) (*model.Part4Talk, error)
	GetByID(ctx context.Context, id uint) (*model.Part4Talk, error)
	Update(ctx context.Context, id uint, req *UpdatePart4TalkRequest) (*model.Part4Talk, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListPart4TalkOptions) ([]*model.Part4Talk, int64, error)
}

// part4TalkService Part4Talk服务实现
type part4TalkService struct {
	part4TalkRepo repository.Part4TalkRepository

	logger logger.Logger
}

// NewPart4TalkService 创建Part4Talk服务
func NewPart4TalkService(
	part4TalkRepo repository.Part4TalkRepository,

	logger logger.Logger,
) Part4TalkService {
	return &part4TalkService{
		part4TalkRepo: part4TalkRepo,

		logger: logger,
	}
}

// CreatePart4TalkRequest 创建Part4Talk请求
type CreatePart4TalkRequest struct {
	TestId            int32          `json:"test_id" validate:"required,min=0"`
	TalkNumber        int32          `json:"talk_number" validate:"required,min=0"`
	Title             sql.NullString `json:"title" validate:"required"`
	Content           string         `json:"content" validate:"required,min=1,max=255"`
	Question1         string         `json:"question1" validate:"required,min=1,max=255"`
	Question2         string         `json:"question2" validate:"required,min=1,max=255"`
	Question3         string         `json:"question3" validate:"required,min=1,max=255"`
	ScenarioId        sql.NullInt32  `json:"scenario_id" validate:"required"`
	DifficultyLevelId sql.NullInt32  `json:"difficulty_level_id" validate:"required"`
}

// UpdatePart4TalkRequest 更新Part4Talk请求
type UpdatePart4TalkRequest struct {
	TestId            *int32          `json:"test_id,omitempty" validate:"omitempty,min=0"`
	TalkNumber        *int32          `json:"talk_number,omitempty" validate:"omitempty,min=0"`
	Title             *sql.NullString `json:"title,omitempty" validate:"omitempty"`
	Content           *string         `json:"content,omitempty" validate:"omitempty,min=1,max=255"`
	Question1         *string         `json:"question1,omitempty" validate:"omitempty,min=1,max=255"`
	Question2         *string         `json:"question2,omitempty" validate:"omitempty,min=1,max=255"`
	Question3         *string         `json:"question3,omitempty" validate:"omitempty,min=1,max=255"`
	ScenarioId        *sql.NullInt32  `json:"scenario_id,omitempty" validate:"omitempty"`
	DifficultyLevelId *sql.NullInt32  `json:"difficulty_level_id,omitempty" validate:"omitempty"`
}

// ListPart4TalkOptions 列表查询选项
type ListPart4TalkOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Part4Talk
func (s *part4TalkService) Create(ctx context.Context, req *CreatePart4TalkRequest) (*model.Part4Talk, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Part4Talk{
		TestId:            req.TestId,
		TalkNumber:        req.TalkNumber,
		Title:             req.Title,
		Content:           req.Content,
		Question1:         req.Question1,
		Question2:         req.Question2,
		Question3:         req.Question3,
		ScenarioId:        req.ScenarioId,
		DifficultyLevelId: req.DifficultyLevelId,
	}

	// 保存到数据库
	if err := s.part4TalkRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create part4talk", "error", err)
		return nil, fmt.Errorf("failed to create part4talk: %w", err)
	}

	s.logger.Info("Part4Talk created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Part4Talk
func (s *part4TalkService) GetByID(ctx context.Context, id uint) (*model.Part4Talk, error) {

	// 从数据库获取
	entity, err := s.part4TalkRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part4talk: %w", err)
	}

	return entity, nil
}

// Update 更新Part4Talk
func (s *part4TalkService) Update(ctx context.Context, id uint, req *UpdatePart4TalkRequest) (*model.Part4Talk, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.part4TalkRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get part4talk: %w", err)
	}

	// 更新字段
	if req.TestId != nil {
		entity.TestId = *req.TestId
	}
	if req.TalkNumber != nil {
		entity.TalkNumber = *req.TalkNumber
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
	if err := s.part4TalkRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update part4talk", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update part4talk: %w", err)
	}

	s.logger.Info("Part4Talk updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Part4Talk
func (s *part4TalkService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.part4TalkRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get part4talk: %w", err)
	}

	// 删除实体
	if err := s.part4TalkRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete part4talk", "id", id, "error", err)
		return fmt.Errorf("failed to delete part4talk: %w", err)
	}

	s.logger.Info("Part4Talk deleted successfully", "id", id)
	return nil
}

// List 获取Part4Talk列表
func (s *part4TalkService) List(ctx context.Context, opts *ListPart4TalkOptions) ([]*model.Part4Talk, int64, error) {
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
	entities, total, err := s.part4TalkRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list part4talks", "error", err)
		return nil, 0, fmt.Errorf("failed to list part4talks: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *part4TalkService) validateCreateRequest(req *CreatePart4TalkRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *part4TalkService) validateUpdateRequest(req *UpdatePart4TalkRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
