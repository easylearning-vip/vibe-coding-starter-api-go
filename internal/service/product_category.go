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

	// 层级分类管理功能
	GetCategoryTree(ctx context.Context) ([]*CategoryTreeNode, error)
	GetChildren(ctx context.Context, parentID uint) ([]*model.ProductCategory, error)
	GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error)
	ValidateParentChild(ctx context.Context, parentID, childID uint) error
	BatchUpdateSortOrder(ctx context.Context, updates []SortOrderUpdate) error
	CanDelete(ctx context.Context, categoryID uint) (bool, string, error)
}

// productCategoryService ProductCategory服务实现
type productCategoryService struct {
	productCategoryRepo repository.ProductCategoryRepository

	logger logger.Logger
}

// NewProductCategoryService 创建ProductCategory服务
func NewProductCategoryService(
	productCategoryRepo repository.ProductCategoryRepository,

	logger logger.Logger,
) ProductCategoryService {
	return &productCategoryService{
		productCategoryRepo: productCategoryRepo,

		logger: logger,
	}
}

// CreateProductCategoryRequest 创建ProductCategory请求
type CreateProductCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	ParentId    uint   `json:"parent_id" validate:"required,min=0"`
	SortOrder   int    `json:"sort_order" validate:"required,min=0"`
	IsActive    bool   `json:"is_active"`
}

// UpdateProductCategoryRequest 更新ProductCategory请求
type UpdateProductCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	ParentId    *uint   `json:"parent_id,omitempty" validate:"omitempty,min=0"`
	SortOrder   *int    `json:"sort_order,omitempty" validate:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active,omitempty" validate:"omitempty"`
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

// CategoryTreeNode 分类树节点
type CategoryTreeNode struct {
	*model.ProductCategory
	Children []*CategoryTreeNode `json:"children"`
}

// SortOrderUpdate 排序更新
type SortOrderUpdate struct {
	ID        uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

// Create 创建ProductCategory
func (s *productCategoryService) Create(ctx context.Context, req *CreateProductCategoryRequest) (*model.ProductCategory, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.ProductCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentId:    req.ParentId,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
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

// GetCategoryTree 获取分类树
func (s *productCategoryService) GetCategoryTree(ctx context.Context) ([]*CategoryTreeNode, error) {
	// 获取所有分类
	categories, err := s.productCategoryRepo.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to get all categories for tree", "error", err)
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}

	// 构建分类映射
	categoryMap := make(map[uint]*CategoryTreeNode)
	var rootNodes []*CategoryTreeNode

	// 第一遍：创建所有节点
	for _, category := range categories {
		node := &CategoryTreeNode{
			ProductCategory: category,
			Children:        make([]*CategoryTreeNode, 0),
		}
		categoryMap[category.ID] = node
	}

	// 第二遍：建立父子关系
	for _, category := range categories {
		node := categoryMap[category.ID]
		if category.ParentId == 0 {
			// 根节点
			rootNodes = append(rootNodes, node)
		} else {
			// 子节点
			if parent, exists := categoryMap[category.ParentId]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return rootNodes, nil
}

// GetChildren 获取子分类
func (s *productCategoryService) GetChildren(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	children, err := s.productCategoryRepo.GetByParentID(ctx, parentID)
	if err != nil {
		s.logger.Error("Failed to get children categories", "parent_id", parentID, "error", err)
		return nil, fmt.Errorf("failed to get children categories: %w", err)
	}

	return children, nil
}

// GetCategoryPath 获取分类路径
func (s *productCategoryService) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	path, err := s.productCategoryRepo.GetCategoryPath(ctx, categoryID)
	if err != nil {
		s.logger.Error("Failed to get category path", "category_id", categoryID, "error", err)
		return nil, fmt.Errorf("failed to get category path: %w", err)
	}

	return path, nil
}

// ValidateParentChild 验证父子关系
func (s *productCategoryService) ValidateParentChild(ctx context.Context, parentID, childID uint) error {
	if parentID == childID {
		return fmt.Errorf("category cannot be its own parent")
	}

	// 检查是否会形成循环引用
	path, err := s.productCategoryRepo.GetCategoryPath(ctx, parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent path: %w", err)
	}

	for _, category := range path {
		if category.ID == childID {
			return fmt.Errorf("circular reference detected: category %d is already an ancestor of category %d", childID, parentID)
		}
	}

	return nil
}

// BatchUpdateSortOrder 批量更新排序
func (s *productCategoryService) BatchUpdateSortOrder(ctx context.Context, updates []SortOrderUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	// 转换为map格式
	updateMap := make(map[uint]int)
	for _, update := range updates {
		updateMap[update.ID] = update.SortOrder
	}

	err := s.productCategoryRepo.BatchUpdateSortOrder(ctx, updateMap)
	if err != nil {
		s.logger.Error("Failed to batch update sort order", "error", err)
		return fmt.Errorf("failed to batch update sort order: %w", err)
	}

	return nil
}

// CanDelete 检查是否可以删除分类
func (s *productCategoryService) CanDelete(ctx context.Context, categoryID uint) (bool, string, error) {
	// 检查是否有子分类
	hasChildren, err := s.productCategoryRepo.HasChildren(ctx, categoryID)
	if err != nil {
		s.logger.Error("Failed to check if category has children", "category_id", categoryID, "error", err)
		return false, "", fmt.Errorf("failed to check children: %w", err)
	}

	if hasChildren {
		return false, "分类下存在子分类，无法删除", nil
	}

	// 检查是否有关联的产品
	productCount, err := s.productCategoryRepo.CountProductsByCategory(ctx, categoryID)
	if err != nil {
		s.logger.Error("Failed to count products in category", "category_id", categoryID, "error", err)
		return false, "", fmt.Errorf("failed to count products: %w", err)
	}

	if productCount > 0 {
		return false, fmt.Sprintf("分类下存在 %d 个产品，无法删除", productCount), nil
	}

	return true, "", nil
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
