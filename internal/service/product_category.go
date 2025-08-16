package service

import (
	"context"
	"fmt"


	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"

)

// ProductCategoryService ProductCategory服务接口
type ProductCategoryService interface {
	Create(ctx context.Context, req *CreateProductCategoryRequest) (*model.ProductCategory, error)
	GetByID(ctx context.Context, id uint) (*model.ProductCategory, error)
	Update(ctx context.Context, id uint, req *UpdateProductCategoryRequest) (*model.ProductCategory, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListProductCategoryOptions) ([]*model.ProductCategory, int64, error)
}

// productCategoryService ProductCategory服务实现
type productCategoryService struct {
	productCategoryRepo repository.ProductCategoryRepository

	logger      logger.Logger
}

// NewProductCategoryService 创建ProductCategory服务
func NewProductCategoryService(
	productCategoryRepo repository.ProductCategoryRepository,

	logger logger.Logger,
) ProductCategoryService {
	return &productCategoryService{
		productCategoryRepo: productCategoryRepo,

		logger:      logger,
	}
}

// CreateProductCategoryRequest 创建ProductCategory请求
type CreateProductCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	ParentId uint `json:"parent_id" validate:"required,min=0"`
	SortOrder int `json:"sort_order" validate:"required,min=0"`
	IsActive bool `json:"is_active"`
}

// UpdateProductCategoryRequest 更新ProductCategory请求
type UpdateProductCategoryRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	ParentId *uint `json:"parent_id,omitempty" validate:"omitempty,min=0"`
	SortOrder *int `json:"sort_order,omitempty" validate:"omitempty,min=0"`
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty"`
}

// ListProductCategoryOptions 列表查询选项
type ListProductCategoryOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建ProductCategory
func (s *productCategoryService) Create(ctx context.Context, req *CreateProductCategoryRequest) (*model.ProductCategory, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.ProductCategory{
		Name: req.Name,
		Description: req.Description,
		ParentId: req.ParentId,
		SortOrder: req.SortOrder,
		IsActive: req.IsActive,
	}

	// 保存到数据库
	if err := s.productCategoryRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create productcategory", "error", err)
		return nil, fmt.Errorf("failed to create productcategory: %w", err)
	}



	s.logger.Info("ProductCategory created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取ProductCategory
func (s *productCategoryService) GetByID(ctx context.Context, id uint) (*model.ProductCategory, error) {


	// 从数据库获取
	entity, err := s.productCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get productcategory: %w", err)
	}



	return entity, nil
}

// Update 更新ProductCategory
func (s *productCategoryService) Update(ctx context.Context, id uint, req *UpdateProductCategoryRequest) (*model.ProductCategory, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.productCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get productcategory: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.ParentId != nil {
		entity.ParentId = *req.ParentId
	}
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}

	// 保存更新
	if err := s.productCategoryRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update productcategory", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update productcategory: %w", err)
	}



	s.logger.Info("ProductCategory updated successfully", "id", id)
	return entity, nil
}

// Delete 删除ProductCategory
func (s *productCategoryService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.productCategoryRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get productcategory: %w", err)
	}

	// 删除实体
	if err := s.productCategoryRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete productcategory", "id", id, "error", err)
		return fmt.Errorf("failed to delete productcategory: %w", err)
	}



	s.logger.Info("ProductCategory deleted successfully", "id", id)
	return nil
}

// List 获取ProductCategory列表
func (s *productCategoryService) List(ctx context.Context, opts *ListProductCategoryOptions) ([]*model.ProductCategory, int64, error) {
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
	entities, total, err := s.productCategoryRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list productcategories", "error", err)
		return nil, 0, fmt.Errorf("failed to list productcategories: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *productCategoryService) validateCreateRequest(req *CreateProductCategoryRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *productCategoryService) validateUpdateRequest(req *UpdateProductCategoryRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}


