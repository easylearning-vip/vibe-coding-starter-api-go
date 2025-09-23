package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// ScenarioService Scenario服务接口
type ScenarioService interface {
	Create(ctx context.Context, req *CreateScenarioRequest) (*model.Scenario, error)
	GetByID(ctx context.Context, id uint) (*model.Scenario, error)
	Update(ctx context.Context, id uint, req *UpdateScenarioRequest) (*model.Scenario, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListScenarioOptions) ([]*model.Scenario, int64, error)
}

// scenarioService Scenario服务实现
type scenarioService struct {
	scenarioRepo repository.ScenarioRepository

	logger logger.Logger
}

// NewScenarioService 创建Scenario服务
func NewScenarioService(
	scenarioRepo repository.ScenarioRepository,

	logger logger.Logger,
) ScenarioService {
	return &scenarioService{
		scenarioRepo: scenarioRepo,

		logger: logger,
	}
}

// CreateScenarioRequest 创建Scenario请求
type CreateScenarioRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"required,min=1,max=255"`
	Examples    *string `json:"examples,omitempty"`
}

// UpdateScenarioRequest 更新Scenario请求
type UpdateScenarioRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Examples    *string `json:"examples,omitempty"`
}

// ListScenarioOptions 列表查询选项
type ListScenarioOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Scenario
func (s *scenarioService) Create(ctx context.Context, req *CreateScenarioRequest) (*model.Scenario, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Scenario{
		Name:        req.Name,
		Description: req.Description,
	}

	// 处理可选的 Examples 字段
	if req.Examples != nil {
		entity.Examples = sql.NullString{
			String: *req.Examples,
			Valid:  true,
		}
	} else {
		entity.Examples = sql.NullString{
			Valid: false,
		}
	}

	// 保存到数据库
	if err := s.scenarioRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create scenario", "error", err)
		return nil, fmt.Errorf("failed to create scenario: %w", err)
	}

	s.logger.Info("Scenario created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Scenario
func (s *scenarioService) GetByID(ctx context.Context, id uint) (*model.Scenario, error) {

	// 从数据库获取
	entity, err := s.scenarioRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scenario: %w", err)
	}

	return entity, nil
}

// Update 更新Scenario
func (s *scenarioService) Update(ctx context.Context, id uint, req *UpdateScenarioRequest) (*model.Scenario, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.scenarioRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scenario: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.Examples != nil {
		entity.Examples = sql.NullString{
			String: *req.Examples,
			Valid:  true,
		}
	}

	// 保存更新
	if err := s.scenarioRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update scenario", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update scenario: %w", err)
	}

	s.logger.Info("Scenario updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Scenario
func (s *scenarioService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.scenarioRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get scenario: %w", err)
	}

	// 删除实体
	if err := s.scenarioRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete scenario", "id", id, "error", err)
		return fmt.Errorf("failed to delete scenario: %w", err)
	}

	s.logger.Info("Scenario deleted successfully", "id", id)
	return nil
}

// List 获取Scenario列表
func (s *scenarioService) List(ctx context.Context, opts *ListScenarioOptions) ([]*model.Scenario, int64, error) {
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
	entities, total, err := s.scenarioRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list scenarios", "error", err)
		return nil, 0, fmt.Errorf("failed to list scenarios: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *scenarioService) validateCreateRequest(req *CreateScenarioRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *scenarioService) validateUpdateRequest(req *UpdateScenarioRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
