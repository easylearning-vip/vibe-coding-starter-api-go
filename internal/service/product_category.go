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
	GetCategoryTree(ctx context.Context, parentID uint) ([]*CategoryTreeNode, error)
	GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error)
	GetChildrenByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error)
	UpdateSortOrder(ctx context.Context, categoryID uint, sortOrder int) error
	CanDeleteCategory(ctx context.Context, categoryID uint) (bool, error)
	GetCategoryWithProductCount(ctx context.Context, categoryID uint) (*CategoryWithProductCount, error)
	BatchUpdateSortOrder(ctx context.Context, sortUpdates map[uint]int) error
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
	ParentID uint                   `json:"parent_id"`
}

// CategoryTreeNode 分类树节点结构
type CategoryTreeNode struct {
	*model.ProductCategory
	Children []*CategoryTreeNode `json:"children"`
	Level    int                 `json:"level"`
	Path     string              `json:"path"`
}

// CategoryWithProductCount 包含产品数量的分类信息
type CategoryWithProductCount struct {
	*model.ProductCategory
	ProductCount int64 `json:"product_count"`
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

// validateCreateRequest 验证创建请求
func (s *productCategoryService) validateCreateRequest(req *CreateProductCategoryRequest) error {
	// 验证父分类是否存在
	if req.ParentId > 0 {
		if _, err := s.productCategoryRepo.GetByID(context.Background(), req.ParentId); err != nil {
			return fmt.Errorf("parent category does not exist: %w", err)
		}
	}
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *productCategoryService) validateUpdateRequest(req *UpdateProductCategoryRequest) error {
	// 验证父分类是否存在
	if req.ParentId != nil && *req.ParentId > 0 {
		if _, err := s.productCategoryRepo.GetByID(context.Background(), *req.ParentId); err != nil {
			return fmt.Errorf("parent category does not exist: %w", err)
		}
	}
	return nil
}

// GetCategoryTree 获取分类树结构
func (s *productCategoryService) GetCategoryTree(ctx context.Context, parentID uint) ([]*CategoryTreeNode, error) {
	// 获取所有分类
	categories, _, err := s.productCategoryRepo.List(ctx, repository.ListOptions{
		Sort:  "sort_order",
		Order: "asc",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// 构建分类树
	categoryMap := make(map[uint]*CategoryTreeNode)
	var tree []*CategoryTreeNode

	// 初始化所有节点
	for _, category := range categories {
		node := &CategoryTreeNode{
			ProductCategory: category,
			Children:        make([]*CategoryTreeNode, 0),
			Level:           0,
			Path:            category.Name,
		}
		categoryMap[category.ID] = node
	}

	// 构建树结构
	for _, category := range categories {
		if category.ParentId == parentID {
			// 根节点
			node := categoryMap[category.ID]
			node.Level = 0
			tree = append(tree, node)
		} else {
			// 子节点
			if parent, exists := categoryMap[category.ParentId]; exists {
				node := categoryMap[category.ID]
				node.Level = parent.Level + 1
				node.Path = fmt.Sprintf("%s > %s", parent.Path, category.Name)
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return tree, nil
}

// GetCategoryPath 获取分类路径（面包屑导航）
func (s *productCategoryService) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	var path []*model.ProductCategory
	currentID := categoryID

	// 防止无限循环
	visited := make(map[uint]bool)
	maxDepth := 10

	for currentID > 0 && len(path) < maxDepth {
		if visited[currentID] {
			break
		}
		visited[currentID] = true

		category, err := s.productCategoryRepo.GetByID(ctx, currentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get category %d: %w", currentID, err)
		}

		// 插入到路径开头
		path = append([]*model.ProductCategory{category}, path...)
		currentID = category.ParentId
	}

	return path, nil
}

// GetChildrenByParentID 根据父分类ID查询子分类
func (s *productCategoryService) GetChildrenByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	filters := map[string]interface{}{
		"parent_id": parentID,
	}

	categories, _, err := s.productCategoryRepo.List(ctx, repository.ListOptions{
		Filters: filters,
		Sort:    "sort_order",
		Order:   "asc",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get children categories: %w", err)
	}

	return categories, nil
}

// UpdateSortOrder 更新分类排序
func (s *productCategoryService) UpdateSortOrder(ctx context.Context, categoryID uint, sortOrder int) error {
	category, err := s.productCategoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	category.SortOrder = sortOrder
	if err := s.productCategoryRepo.Update(ctx, category); err != nil {
		return fmt.Errorf("failed to update sort order: %w", err)
	}

	s.logger.Info("Category sort order updated", "id", categoryID, "sort_order", sortOrder)
	return nil
}

// CanDeleteCategory 检查分类是否可删除
func (s *productCategoryService) CanDeleteCategory(ctx context.Context, categoryID uint) (bool, error) {
	// 检查是否有子分类
	children, err := s.GetChildrenByParentID(ctx, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check children: %w", err)
	}
	if len(children) > 0 {
		return false, fmt.Errorf("category has child categories")
	}

	// 检查是否有关联产品
	// 这里需要产品仓储，暂时返回true
	// TODO: 实现产品关联检查
	s.logger.Warn("Product association check not implemented", "category_id", categoryID)
	return true, nil
}

// GetCategoryWithProductCount 获取分类及其产品数量
func (s *productCategoryService) GetCategoryWithProductCount(ctx context.Context, categoryID uint) (*CategoryWithProductCount, error) {
	category, err := s.productCategoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// TODO: 实现产品数量统计
	// 这里暂时返回0，需要产品仓储支持
	return &CategoryWithProductCount{
		ProductCategory: category,
		ProductCount:    0,
	}, nil
}

// BatchUpdateSortOrder 批量更新排序
func (s *productCategoryService) BatchUpdateSortOrder(ctx context.Context, sortUpdates map[uint]int) error {
	for categoryID, sortOrder := range sortUpdates {
		category, err := s.productCategoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			return fmt.Errorf("failed to get category %d: %w", categoryID, err)
		}

		category.SortOrder = sortOrder
		if err := s.productCategoryRepo.Update(ctx, category); err != nil {
			return fmt.Errorf("failed to update category %d sort order: %w", categoryID, err)
		}
	}

	s.logger.Info("Batch sort orders updated", "count", len(sortUpdates))
	return nil
}
