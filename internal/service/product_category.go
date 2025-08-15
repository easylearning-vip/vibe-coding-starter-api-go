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
	GetCategoryTree(ctx context.Context) ([]*model.ProductCategory, error)
	GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error)
	GetByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error)
	BatchUpdateSortOrder(ctx context.Context, updates []*UpdateSortOrderRequest) error
	ValidateParentChildRelationship(ctx context.Context, parentID, childID uint) error
	CanDeleteCategory(ctx context.Context, categoryID uint) (bool, error)
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

	// 验证父子关系
	if req.ParentId > 0 {
		if err := s.ValidateParentChildRelationship(ctx, req.ParentId, 0); err != nil {
			return nil, fmt.Errorf("invalid parent-child relationship: %w", err)
		}
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

	// 验证父子关系
	if req.ParentId != nil {
		if err := s.ValidateParentChildRelationship(ctx, *req.ParentId, id); err != nil {
			return nil, fmt.Errorf("invalid parent-child relationship: %w", err)
		}
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

	// 检查是否可删除
	canDelete, err := s.CanDeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if category can be deleted: %w", err)
	}
	
	if !canDelete {
		return fmt.Errorf("category cannot be deleted: it has subcategories or associated products")
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

// UpdateSortOrderRequest 更新排序请求
type UpdateSortOrderRequest struct {
	CategoryID uint `json:"category_id" validate:"required"`
	SortOrder  int  `json:"sort_order" validate:"required"`
}

// GetCategoryTree 获取分类树结构
func (s *productCategoryService) GetCategoryTree(ctx context.Context) ([]*model.ProductCategory, error) {
	tree, err := s.productCategoryRepo.GetCategoryTree(ctx)
	if err != nil {
		s.logger.Error("Failed to get category tree", "error", err)
		return nil, fmt.Errorf("failed to get category tree: %w", err)
	}
	return tree, nil
}

// GetCategoryPath 获取分类路径（面包屑导航）
func (s *productCategoryService) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	path, err := s.productCategoryRepo.GetCategoryPath(ctx, categoryID)
	if err != nil {
		s.logger.Error("Failed to get category path", "category_id", categoryID, "error", err)
		return nil, fmt.Errorf("failed to get category path: %w", err)
	}
	return path, nil
}

// GetByParentID 根据父分类ID获取子分类
func (s *productCategoryService) GetByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	categories, err := s.productCategoryRepo.GetByParentID(ctx, parentID)
	if err != nil {
		s.logger.Error("Failed to get categories by parent ID", "parent_id", parentID, "error", err)
		return nil, fmt.Errorf("failed to get categories by parent ID: %w", err)
	}
	return categories, nil
}

// BatchUpdateSortOrder 批量更新排序
func (s *productCategoryService) BatchUpdateSortOrder(ctx context.Context, updates []*UpdateSortOrderRequest) error {
	for _, update := range updates {
		category, err := s.productCategoryRepo.GetByID(ctx, update.CategoryID)
		if err != nil {
			s.logger.Error("Failed to get category for sort update", "category_id", update.CategoryID, "error", err)
			return fmt.Errorf("failed to get category %d: %w", update.CategoryID, err)
		}
		
		category.SortOrder = update.SortOrder
		if err := s.productCategoryRepo.Update(ctx, category); err != nil {
			s.logger.Error("Failed to update category sort order", "category_id", update.CategoryID, "error", err)
			return fmt.Errorf("failed to update category sort order: %w", err)
		}
	}
	
	s.logger.Info("Batch updated sort order successfully", "count", len(updates))
	return nil
}

// ValidateParentChildRelationship 验证父子关系
func (s *productCategoryService) ValidateParentChildRelationship(ctx context.Context, parentID, childID uint) error {
	if parentID == 0 {
		return nil // 根分类是允许的
	}
	
	if parentID == childID {
		return fmt.Errorf("category cannot be its own parent")
	}
	
	// 检查是否会形成循环引用
	path, err := s.productCategoryRepo.GetCategoryPath(ctx, parentID)
	if err != nil {
		return fmt.Errorf("failed to validate parent-child relationship: %w", err)
	}
	
	for _, category := range path {
		if category.ID == childID {
			return fmt.Errorf("circular reference detected: category %d cannot be parent of category %d", childID, parentID)
		}
	}
	
	return nil
}

// CanDeleteCategory 检查分类是否可删除
func (s *productCategoryService) CanDeleteCategory(ctx context.Context, categoryID uint) (bool, error) {
	// 检查是否有子分类
	hasChildren, err := s.productCategoryRepo.HasChildren(ctx, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check children categories: %w", err)
	}
	
	if hasChildren {
		return false, nil
	}
	
	// 检查是否有关联产品
	productCount, err := s.productCategoryRepo.CountProductsByCategory(ctx, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to count products: %w", err)
	}
	
	if productCount > 0 {
		return false, nil
	}
	
	return true, nil
}


