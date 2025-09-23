package service

import (
	"context"
	"database/sql"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// TestService Test服务接口
type TestService interface {
	Create(ctx context.Context, req *CreateTestRequest) (*model.Test, error)
	GetByID(ctx context.Context, id uint) (*model.Test, error)
	Update(ctx context.Context, id uint, req *UpdateTestRequest) (*model.Test, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListTestOptions) ([]*model.Test, int64, error)
}

// testService Test服务实现
type testService struct {
	testRepo repository.TestRepository

	logger logger.Logger
}

// NewTestService 创建Test服务
func NewTestService(
	testRepo repository.TestRepository,

	logger logger.Logger,
) TestService {
	return &testService{
		testRepo: testRepo,

		logger: logger,
	}
}

// CreateTestRequest 创建Test请求
type CreateTestRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty"`
}

// UpdateTestRequest 更新Test请求
type UpdateTestRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty"`
}

// ListTestOptions 列表查询选项
type ListTestOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Test
func (s *testService) Create(ctx context.Context, req *CreateTestRequest) (*model.Test, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Test{
		Name: req.Name,
	}

	// 处理可选的 Description 字段
	if req.Description != nil {
		entity.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	} else {
		entity.Description = sql.NullString{
			Valid: false,
		}
	}

	// 保存到数据库
	if err := s.testRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create test", "error", err)
		return nil, fmt.Errorf("failed to create test: %w", err)
	}

	s.logger.Info("Test created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Test
func (s *testService) GetByID(ctx context.Context, id uint) (*model.Test, error) {

	// 从数据库获取
	entity, err := s.testRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	return entity, nil
}

// Update 更新Test
func (s *testService) Update(ctx context.Context, id uint, req *UpdateTestRequest) (*model.Test, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.testRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get test: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	// 保存更新
	if err := s.testRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update test", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update test: %w", err)
	}

	s.logger.Info("Test updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Test
func (s *testService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.testRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get test: %w", err)
	}

	// 删除实体
	if err := s.testRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete test", "id", id, "error", err)
		return fmt.Errorf("failed to delete test: %w", err)
	}

	s.logger.Info("Test deleted successfully", "id", id)
	return nil
}

// List 获取Test列表
func (s *testService) List(ctx context.Context, opts *ListTestOptions) ([]*model.Test, int64, error) {
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
	entities, total, err := s.testRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list tests", "error", err)
		return nil, 0, fmt.Errorf("failed to list tests: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *testService) validateCreateRequest(req *CreateTestRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *testService) validateUpdateRequest(req *UpdateTestRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}
