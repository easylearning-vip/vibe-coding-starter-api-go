package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// DifficultyLevelService DifficultyLevel服务接口
type DifficultyLevelService interface {
	Create(ctx context.Context, req *CreateDifficultyLevelRequest) (*model.DifficultyLevel, error)
	GetByID(ctx context.Context, id uint) (*model.DifficultyLevel, error)
	Update(ctx context.Context, id uint, req *UpdateDifficultyLevelRequest) (*model.DifficultyLevel, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListDifficultyLevelOptions) ([]*model.DifficultyLevel, int64, error)
}

// difficultyLevelService DifficultyLevel服务实现
type difficultyLevelService struct {
	difficultyLevelRepo repository.DifficultyLevelRepository

	logger logger.Logger
}

// NewDifficultyLevelService 创建DifficultyLevel服务
func NewDifficultyLevelService(
	difficultyLevelRepo repository.DifficultyLevelRepository,

	logger logger.Logger,
) DifficultyLevelService {
	return &difficultyLevelService{
		difficultyLevelRepo: difficultyLevelRepo,

		logger: logger,
	}
}

// CreateDifficultyLevelRequest 创建DifficultyLevel请求
type CreateDifficultyLevelRequest struct {
	Name            string         `json:"name" validate:"required,min=1,max=255"`
	Description     string         `json:"description" validate:"required,min=1,max=255"`
	Characteristics sql.NullString `json:"characteristics" validate:"required"`
}

// UpdateDifficultyLevelRequest 更新DifficultyLevel请求
type UpdateDifficultyLevelRequest struct {
	Name            *string         `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description     *string         `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Characteristics *sql.NullString `json:"characteristics,omitempty" validate:"omitempty"`
}

// ListDifficultyLevelOptions 列表查询选项
type ListDifficultyLevelOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建DifficultyLevel
func (s *difficultyLevelService) Create(ctx context.Context, req *CreateDifficultyLevelRequest) (*model.DifficultyLevel, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.DifficultyLevel{
		Name:            req.Name,
		Description:     req.Description,
		Characteristics: req.Characteristics,
	}

	// 保存到数据库
	if err := s.difficultyLevelRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create difficultylevel", "error", err)
		return nil, fmt.Errorf("failed to create difficultylevel: %w", err)
	}

	s.logger.Info("DifficultyLevel created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取DifficultyLevel
func (s *difficultyLevelService) GetByID(ctx context.Context, id uint) (*model.DifficultyLevel, error) {

	// 从数据库获取
	entity, err := s.difficultyLevelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get difficultylevel: %w", err)
	}

	return entity, nil
}

// Update 更新DifficultyLevel
func (s *difficultyLevelService) Update(ctx context.Context, id uint, req *UpdateDifficultyLevelRequest) (*model.DifficultyLevel, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.difficultyLevelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get difficultylevel: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.Characteristics != nil {
		entity.Characteristics = *req.Characteristics
	}

	// 保存更新
	if err := s.difficultyLevelRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update difficultylevel", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update difficultylevel: %w", err)
	}

	s.logger.Info("DifficultyLevel updated successfully", "id", id)
	return entity, nil
}

// Delete 删除DifficultyLevel
func (s *difficultyLevelService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.difficultyLevelRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get difficultylevel: %w", err)
	}

	// 删除实体
	if err := s.difficultyLevelRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete difficultylevel", "id", id, "error", err)
		return fmt.Errorf("failed to delete difficultylevel: %w", err)
	}

	s.logger.Info("DifficultyLevel deleted successfully", "id", id)
	return nil
}

// List 获取DifficultyLevel列表
func (s *difficultyLevelService) List(ctx context.Context, opts *ListDifficultyLevelOptions) ([]*model.DifficultyLevel, int64, error) {
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
	entities, total, err := s.difficultyLevelRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list difficultylevels", "error", err)
		return nil, 0, fmt.Errorf("failed to list difficultylevels: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *difficultyLevelService) validateCreateRequest(req *CreateDifficultyLevelRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *difficultyLevelService) validateUpdateRequest(req *UpdateDifficultyLevelRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
