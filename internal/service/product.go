package service

import (
	"context"
	"fmt"


	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"

)

// ProductService Product服务接口
type ProductService interface {
	Create(ctx context.Context, req *CreateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	Update(ctx context.Context, id uint, req *UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)
	GetByCategoryID(ctx context.Context, categoryID uint, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetByCategoryWithSubcategories(ctx context.Context, categoryID uint, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetHotSellingProducts(ctx context.Context, limit int) ([]*model.Product, error)
	GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error)
	SearchProducts(ctx context.Context, query string, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetLowStockProducts(ctx context.Context) ([]*model.Product, error)
	GetProductsByStatus(ctx context.Context, isActive bool, opts *ListProductOptions) ([]*model.Product, int64, error)
	BatchUpdatePrices(ctx context.Context, updates []*UpdatePriceRequest) error
	BatchUpdateStatus(ctx context.Context, updates []*UpdateStatusRequest) error
}

// productService Product服务实现
type productService struct {
	productRepo repository.ProductRepository

	logger      logger.Logger
}

// NewProductService 创建Product服务
func NewProductService(
	productRepo repository.ProductRepository,

	logger logger.Logger,
) ProductService {
	return &productService{
		productRepo: productRepo,

		logger:      logger,
	}
}

// CreateProductRequest 创建Product请求
type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	CategoryId uint `json:"category_id" validate:"required,min=0"`
	Sku string `json:"sku" validate:"required,min=1,max=255"`
	Price float64 `json:"price" validate:"required,min=0"`
	CostPrice float64 `json:"cost_price" validate:"required,min=0"`
	StockQuantity int `json:"stock_quantity" validate:"required,min=0"`
	MinStock int `json:"min_stock" validate:"required,min=0"`
	IsActive bool `json:"is_active"`
	Weight float64 `json:"weight" validate:"required,min=0"`
	Dimensions string `json:"dimensions" validate:"required,min=1,max=255"`
}

// UpdateProductRequest 更新Product请求
type UpdateProductRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	CategoryId *uint `json:"category_id,omitempty" validate:"omitempty,min=0"`
	Sku *string `json:"sku,omitempty" validate:"omitempty,min=1,max=255"`
	Price *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	CostPrice *float64 `json:"cost_price,omitempty" validate:"omitempty,min=0"`
	StockQuantity *int `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	MinStock *int `json:"min_stock,omitempty" validate:"omitempty,min=0"`
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty"`
	Weight *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Dimensions *string `json:"dimensions,omitempty" validate:"omitempty,min=1,max=255"`
}

// ListProductOptions 列表查询选项
type ListProductOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Product
func (s *productService) Create(ctx context.Context, req *CreateProductRequest) (*model.Product, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Product{
		Name: req.Name,
		Description: req.Description,
		CategoryId: req.CategoryId,
		Sku: req.Sku,
		Price: req.Price,
		CostPrice: req.CostPrice,
		StockQuantity: req.StockQuantity,
		MinStock: req.MinStock,
		IsActive: req.IsActive,
		Weight: req.Weight,
		Dimensions: req.Dimensions,
	}

	// 保存到数据库
	if err := s.productRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create product", "error", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}



	s.logger.Info("Product created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Product
func (s *productService) GetByID(ctx context.Context, id uint) (*model.Product, error) {


	// 从数据库获取
	entity, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}



	return entity, nil
}

// Update 更新Product
func (s *productService) Update(ctx context.Context, id uint, req *UpdateProductRequest) (*model.Product, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.CategoryId != nil {
		entity.CategoryId = *req.CategoryId
	}
	if req.Sku != nil {
		entity.Sku = *req.Sku
	}
	if req.Price != nil {
		entity.Price = *req.Price
	}
	if req.CostPrice != nil {
		entity.CostPrice = *req.CostPrice
	}
	if req.StockQuantity != nil {
		entity.StockQuantity = *req.StockQuantity
	}
	if req.MinStock != nil {
		entity.MinStock = *req.MinStock
	}
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	if req.Weight != nil {
		entity.Weight = *req.Weight
	}
	if req.Dimensions != nil {
		entity.Dimensions = *req.Dimensions
	}

	// 保存更新
	if err := s.productRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update product", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}



	s.logger.Info("Product updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Product
func (s *productService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.productRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// 删除实体
	if err := s.productRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete product", "id", id, "error", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}



	s.logger.Info("Product deleted successfully", "id", id)
	return nil
}

// List 获取Product列表
func (s *productService) List(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error) {
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
	entities, total, err := s.productRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list products", "error", err)
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *productService) validateCreateRequest(req *CreateProductRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *productService) validateUpdateRequest(req *UpdateProductRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// UpdatePriceRequest 更新价格请求
type UpdatePriceRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Price     float64 `json:"price" validate:"required,min=0"`
	CostPrice *float64 `json:"cost_price,omitempty"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	ProductID uint  `json:"product_id" validate:"required"`
	IsActive  bool  `json:"is_active"`
}

// GetBySKU 根据SKU获取产品
func (s *productService) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	product, err := s.productRepo.GetBySKU(ctx, sku)
	if err != nil {
		s.logger.Error("Failed to get product by SKU", "sku", sku, "error", err)
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return product, nil
}

// GetByCategoryID 根据分类ID获取产品
func (s *productService) GetByCategoryID(ctx context.Context, categoryID uint, opts *ListProductOptions) ([]*model.Product, int64, error) {
	products, total, err := s.productRepo.GetByCategoryID(ctx, categoryID, repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	})
	if err != nil {
		s.logger.Error("Failed to get products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}
	return products, total, nil
}

// GetByCategoryWithSubcategories 根据分类ID获取产品（包含子分类）
func (s *productService) GetByCategoryWithSubcategories(ctx context.Context, categoryID uint, opts *ListProductOptions) ([]*model.Product, int64, error) {
	products, total, err := s.productRepo.GetByCategoryWithSubcategories(ctx, categoryID, repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	})
	if err != nil {
		s.logger.Error("Failed to get products by category with subcategories", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category with subcategories: %w", err)
	}
	return products, total, nil
}

// GetHotSellingProducts 获取热销产品
func (s *productService) GetHotSellingProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 10 // 默认限制
	}
	
	products, err := s.productRepo.GetHotSellingProducts(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to get hot selling products", "error", err)
		return nil, fmt.Errorf("failed to get hot selling products: %w", err)
	}
	return products, nil
}

// GetByPriceRange 根据价格区间获取产品
func (s *productService) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error) {
	if minPrice < 0 {
		minPrice = 0
	}
	if maxPrice <= 0 {
		maxPrice = 999999.99 // 设置一个合理的最大值
	}
	if minPrice > maxPrice {
		minPrice, maxPrice = maxPrice, minPrice // 交换最小值和最大值
	}
	
	products, total, err := s.productRepo.GetByPriceRange(ctx, minPrice, maxPrice, repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	})
	if err != nil {
		s.logger.Error("Failed to get products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by price range: %w", err)
	}
	return products, total, nil
}

// SearchProducts 搜索产品
func (s *productService) SearchProducts(ctx context.Context, query string, opts *ListProductOptions) ([]*model.Product, int64, error) {
	if query == "" {
		return s.List(ctx, opts) // 如果搜索查询为空，返回所有产品
	}
	
	products, total, err := s.productRepo.SearchProducts(ctx, query, repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	})
	if err != nil {
		s.logger.Error("Failed to search products", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	return products, total, nil
}

// GetLowStockProducts 获取库存不足的产品
func (s *productService) GetLowStockProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepo.GetLowStockProducts(ctx)
	if err != nil {
		s.logger.Error("Failed to get low stock products", "error", err)
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	return products, nil
}

// GetProductsByStatus 根据状态获取产品
func (s *productService) GetProductsByStatus(ctx context.Context, isActive bool, opts *ListProductOptions) ([]*model.Product, int64, error) {
	products, total, err := s.productRepo.GetProductsByStatus(ctx, isActive, repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	})
	if err != nil {
		s.logger.Error("Failed to get products by status", "is_active", isActive, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by status: %w", err)
	}
	return products, total, nil
}

// BatchUpdatePrices 批量更新价格
func (s *productService) BatchUpdatePrices(ctx context.Context, updates []*UpdatePriceRequest) error {
	for _, update := range updates {
		product, err := s.productRepo.GetByID(ctx, update.ProductID)
		if err != nil {
			s.logger.Error("Failed to get product for price update", "product_id", update.ProductID, "error", err)
			return fmt.Errorf("failed to get product %d: %w", update.ProductID, err)
		}
		
		// 记录价格变更历史（这里可以扩展为单独的价格历史表）
		product.Price = update.Price
		if update.CostPrice != nil {
			product.CostPrice = *update.CostPrice
		}
		
		if err := s.productRepo.Update(ctx, product); err != nil {
			s.logger.Error("Failed to update product price", "product_id", update.ProductID, "error", err)
			return fmt.Errorf("failed to update product price: %w", err)
		}
	}
	
	s.logger.Info("Batch updated product prices successfully", "count", len(updates))
	return nil
}

// BatchUpdateStatus 批量更新状态
func (s *productService) BatchUpdateStatus(ctx context.Context, updates []*UpdateStatusRequest) error {
	for _, update := range updates {
		product, err := s.productRepo.GetByID(ctx, update.ProductID)
		if err != nil {
			s.logger.Error("Failed to get product for status update", "product_id", update.ProductID, "error", err)
			return fmt.Errorf("failed to get product %d: %w", update.ProductID, err)
		}
		
		product.IsActive = update.IsActive
		
		if err := s.productRepo.Update(ctx, product); err != nil {
			s.logger.Error("Failed to update product status", "product_id", update.ProductID, "error", err)
			return fmt.Errorf("failed to update product status: %w", err)
		}
	}
	
	s.logger.Info("Batch updated product status successfully", "count", len(updates))
	return nil
}


