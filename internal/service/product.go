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

	// 高级业务功能
	SearchProducts(ctx context.Context, query string, opts *ListProductOptions) ([]*model.Product, int64, error)
	UpdateProductStatus(ctx context.Context, id uint, isActive bool) (*model.Product, error)
	BatchUpdatePrices(ctx context.Context, updates map[uint]float64) error
	GetProductsByCategory(ctx context.Context, categoryID uint, includeSubcategories bool, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetLowStockProducts(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetHotSellingProducts(ctx context.Context, minSales int, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductsByPriceRange(ctx context.Context, minPrice float64, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductStatistics(ctx context.Context, productID uint) (*ProductStatistics, error)
	UpdateProductStock(ctx context.Context, productID uint, quantityChange int, reason string) (*model.Product, error)

	// 增强的业务功能
	GetProductsByCategories(ctx context.Context, categoryIDs []uint, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductsByStatus(ctx context.Context, isActive bool, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductsByStockRange(ctx context.Context, minStock int, maxStock int, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductPriceHistory(ctx context.Context, productID uint) ([]*PriceHistory, error)
	GetProductStockHistory(ctx context.Context, productID uint) ([]*StockHistory, error)
	ValidateProductSKU(ctx context.Context, sku string) error
	GetProductsBySKUs(ctx context.Context, skus []string) ([]*model.Product, error)
}

// productService Product服务实现
type productService struct {
	productRepo repository.ProductRepository

	logger logger.Logger
}

// NewProductService 创建Product服务
func NewProductService(
	productRepo repository.ProductRepository,

	logger logger.Logger,
) ProductService {
	return &productService{
		productRepo: productRepo,

		logger: logger,
	}
}

// CreateProductRequest 创建Product请求
type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	Description   string  `json:"description" validate:"required,min=1,max=255"`
	CategoryId    uint    `json:"category_id" validate:"required,min=0"`
	Sku           string  `json:"sku" validate:"required,min=1,max=255"`
	Price         float64 `json:"price" validate:"required,min=0"`
	CostPrice     float64 `json:"cost_price" validate:"required,min=0"`
	StockQuantity int     `json:"stock_quantity" validate:"required,min=0"`
	MinStock      int     `json:"min_stock" validate:"required,min=0"`
	IsActive      bool    `json:"is_active"`
	Weight        float64 `json:"weight" validate:"required,min=0"`
	Dimensions    string  `json:"dimensions" validate:"required,min=1,max=255"`
}

// UpdateProductRequest 更新Product请求
type UpdateProductRequest struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	CategoryId    *uint    `json:"category_id,omitempty" validate:"omitempty,min=0"`
	Sku           *string  `json:"sku,omitempty" validate:"omitempty,min=1,max=255"`
	Price         *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	CostPrice     *float64 `json:"cost_price,omitempty" validate:"omitempty,min=0"`
	StockQuantity *int     `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	MinStock      *int     `json:"min_stock,omitempty" validate:"omitempty,min=0"`
	IsActive      *bool    `json:"is_active,omitempty" validate:"omitempty"`
	Weight        *float64 `json:"weight,omitempty" validate:"omitempty,min=0"`
	Dimensions    *string  `json:"dimensions,omitempty" validate:"omitempty,min=1,max=255"`
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

// ProductStatistics 产品统计信息
type ProductStatistics struct {
	ProductID     uint    `json:"product_id"`
	TotalSales    int64   `json:"total_sales"`
	CurrentStock  int     `json:"current_stock"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int64   `json:"review_count"`
	Revenue       float64 `json:"revenue"`
	Profit        float64 `json:"profit"`
}

// StockUpdateRequest 库存更新请求
type StockUpdateRequest struct {
	ProductID      uint   `json:"product_id" validate:"required"`
	QuantityChange int    `json:"quantity_change" validate:"required"`
	Reason         string `json:"reason" validate:"required"`
	ReferenceID    string `json:"reference_id,omitempty"`
}

// PriceHistory 价格历史记录
type PriceHistory struct {
	ProductID uint    `json:"product_id"`
	OldPrice  float64 `json:"old_price"`
	NewPrice  float64 `json:"new_price"`
	ChangedAt string  `json:"changed_at"`
	ChangedBy string  `json:"changed_by"`
	Reason    string  `json:"reason"`
}

// StockHistory 库存历史记录
type StockHistory struct {
	ProductID     uint   `json:"product_id"`
	StockChange   int    `json:"stock_change"`
	OldStock      int    `json:"old_stock"`
	NewStock      int    `json:"new_stock"`
	ChangedAt     string `json:"changed_at"`
	ChangedBy     string `json:"changed_by"`
	Reason        string `json:"reason"`
	ReferenceID   string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
}

// Create 创建Product
func (s *productService) Create(ctx context.Context, req *CreateProductRequest) (*model.Product, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Product{
		Name:          req.Name,
		Description:   req.Description,
		CategoryId:    req.CategoryId,
		Sku:           req.Sku,
		Price:         req.Price,
		CostPrice:     req.CostPrice,
		StockQuantity: req.StockQuantity,
		MinStock:      req.MinStock,
		IsActive:      req.IsActive,
		Weight:        req.Weight,
		Dimensions:    req.Dimensions,
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

// SearchProducts 搜索产品
func (s *productService) SearchProducts(ctx context.Context, query string, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	var entities []*model.Product
	var total int64
	var err error

	// 根据查询类型执行不同的搜索策略
	if opts.Filters != nil {
		if searchType, ok := opts.Filters["search_type"]; ok {
			switch searchType {
			case "name":
				entities, total, err = s.productRepo.SearchByName(ctx, query, repoOpts)
			case "sku":
				entities, total, err = s.productRepo.SearchBySKU(ctx, query, repoOpts)
			case "category":
				if categoryID, ok := opts.Filters["category_id"]; ok {
					if id, ok := categoryID.(uint); ok {
						includeSub, _ := opts.Filters["include_subcategories"].(bool)
						entities, total, err = s.GetProductsByCategory(ctx, id, includeSub, opts)
					}
				}
			default:
				// 默认搜索名称
				entities, total, err = s.productRepo.SearchByName(ctx, query, repoOpts)
			}
		} else {
			// 默认搜索名称
			entities, total, err = s.productRepo.SearchByName(ctx, query, repoOpts)
		}
	} else {
		// 默认搜索名称
		entities, total, err = s.productRepo.SearchByName(ctx, query, repoOpts)
	}

	if err != nil {
		s.logger.Error("Failed to search products", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	s.logger.Info("Products searched", "query", query, "count", len(entities))
	return entities, total, nil
}

// UpdateProductStatus 更新产品状态
func (s *productService) UpdateProductStatus(ctx context.Context, id uint, isActive bool) (*model.Product, error) {
	// 获取产品
	entity, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 更新状态
	entity.IsActive = isActive
	if err := s.productRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update product status", "id", id, "is_active", isActive, "error", err)
		return nil, fmt.Errorf("failed to update product status: %w", err)
	}

	s.logger.Info("Product status updated", "id", id, "is_active", isActive)
	return entity, nil
}

// BatchUpdatePrices 批量更新价格
func (s *productService) BatchUpdatePrices(ctx context.Context, updates map[uint]float64) error {
	// 验证价格
	for productID, newPrice := range updates {
		if newPrice < 0 {
			return fmt.Errorf("invalid price for product %d: must be non-negative", productID)
		}
	}

	// 使用批处理方法
	if err := s.productRepo.BatchUpdatePrices(ctx, updates); err != nil {
		s.logger.Error("Failed to batch update prices", "error", err)
		return fmt.Errorf("failed to batch update prices: %w", err)
	}

	s.logger.Info("Batch prices updated", "count", len(updates))
	return nil
}

// GetProductsByCategory 根据分类获取产品
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID uint, includeSubcategories bool, opts *ListProductOptions) ([]*model.Product, int64, error) {
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
		Search:   opts.Search,
	}

	// 获取产品列表
	var entities []*model.Product
	var total int64
	var err error

	if includeSubcategories {
		// 获取子分类ID列表
		categoryIDs := []uint{categoryID}

		// 递归获取所有子分类ID
		subCategories, err := s.getSubCategoryIDs(ctx, categoryID)
		if err != nil {
			s.logger.Error("Failed to get sub categories", "category_id", categoryID, "error", err)
			// 降级处理：只返回当前分类的产品
			entities, total, err = s.productRepo.GetByCategoryID(ctx, categoryID, repoOpts)
		} else {
			categoryIDs = append(categoryIDs, subCategories...)
			entities, total, err = s.productRepo.GetByCategoryIDs(ctx, categoryIDs, repoOpts)
		}
	} else {
		entities, total, err = s.productRepo.GetByCategoryID(ctx, categoryID, repoOpts)
	}

	if err != nil {
		s.logger.Error("Failed to get products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}

	s.logger.Info("Products retrieved by category", "category_id", categoryID, "include_subcategories", includeSubcategories, "count", len(entities))
	return entities, total, nil
}

// getSubCategoryIDs 递归获取所有子分类ID
func (s *productService) getSubCategoryIDs(ctx context.Context, parentID uint) ([]uint, error) {
	// 这里简化处理，实际应用中应该从ProductCategoryService获取
	// 由于服务层不能直接访问其他仓储，这里简化实现
	return []uint{}, nil
}

// GetLowStockProducts 获取低库存产品
func (s *productService) GetLowStockProducts(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	}

	// 获取低库存产品
	entities, total, err := s.productRepo.GetLowStockProducts(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get low stock products", "error", err)
		return nil, 0, fmt.Errorf("failed to get low stock products: %w", err)
	}

	s.logger.Info("Low stock products retrieved", "count", len(entities))
	return entities, total, nil
}

// GetHotSellingProducts 获取热销产品
func (s *productService) GetHotSellingProducts(ctx context.Context, minSales int, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	// 使用专门的仓储方法获取热销产品
	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	}

	entities, total, err := s.productRepo.GetHotSellingProducts(ctx, minSales, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get hot selling products", "error", err)
		return nil, 0, fmt.Errorf("failed to get hot selling products: %w", err)
	}

	s.logger.Info("Hot selling products retrieved", "min_sales", minSales, "count", len(entities))
	return entities, total, nil
}

// GetProductsByPriceRange 根据价格范围获取产品
func (s *productService) GetProductsByPriceRange(ctx context.Context, minPrice float64, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	// 验证价格范围
	if minPrice < 0 || maxPrice < 0 {
		return nil, 0, fmt.Errorf("price range must be non-negative")
	}
	if minPrice > maxPrice {
		return nil, 0, fmt.Errorf("min price cannot be greater than max price")
	}

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Search:   opts.Search,
	}

	entities, total, err := s.productRepo.GetByPriceRange(ctx, minPrice, maxPrice, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by price range: %w", err)
	}

	s.logger.Info("Products retrieved by price range", "min_price", minPrice, "max_price", maxPrice, "count", len(entities))
	return entities, total, nil
}

// GetProductStatistics 获取产品统计信息
func (s *productService) GetProductStatistics(ctx context.Context, productID uint) (*ProductStatistics, error) {
	// 获取产品信息
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// TODO: 实现真实的统计逻辑
	// 这里返回模拟数据，实际应用中需要从订单、销售等表中获取数据
	stats := &ProductStatistics{
		ProductID:     productID,
		TotalSales:    0, // 需要从销售记录获取
		CurrentStock:  product.StockQuantity,
		AverageRating: 4.5, // 需要从评价系统获取
		ReviewCount:   0,   // 需要从评价系统获取
		Revenue:       0,   // 需要从订单系统获取
		Profit:        0,   // 需要从财务系统获取
	}

	s.logger.Info("Product statistics retrieved", "product_id", productID)
	return stats, nil
}

// GetProductsByCategories 根据多个分类ID获取产品
func (s *productService) GetProductsByCategories(ctx context.Context, categoryIDs []uint, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
		Search:   opts.Search,
	}

	entities, total, err := s.productRepo.GetByCategoryIDs(ctx, categoryIDs, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get products by categories", "category_ids", categoryIDs, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by categories: %w", err)
	}

	s.logger.Info("Products retrieved by categories", "category_ids_count", len(categoryIDs), "count", len(entities))
	return entities, total, nil
}

// GetProductsByStatus 根据产品状态获取产品
func (s *productService) GetProductsByStatus(ctx context.Context, isActive bool, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  map[string]interface{}{"is_active": isActive},
		Search:   opts.Search,
	}

	entities, total, err := s.productRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get products by status", "is_active", isActive, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by status: %w", err)
	}

	s.logger.Info("Products retrieved by status", "is_active", isActive, "count", len(entities))
	return entities, total, nil
}

// GetProductsByStockRange 根据库存范围获取产品
func (s *productService) GetProductsByStockRange(ctx context.Context, minStock int, maxStock int, opts *ListProductOptions) ([]*model.Product, int64, error) {
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

	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  map[string]interface{}{"stock_quantity__gte": minStock, "stock_quantity__lte": maxStock},
		Search:   opts.Search,
	}

	entities, total, err := s.productRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get products by stock range", "min_stock", minStock, "max_stock", maxStock, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by stock range: %w", err)
	}

	s.logger.Info("Products retrieved by stock range", "min_stock", minStock, "max_stock", maxStock, "count", len(entities))
	return entities, total, nil
}

// GetProductPriceHistory 获取产品价格历史
func (s *productService) GetProductPriceHistory(ctx context.Context, productID uint) ([]*PriceHistory, error) {
	// 这里简化处理，实际应用中应该从价格历史表获取
	// 返回模拟数据
	history := []*PriceHistory{
		{
			ProductID: productID,
			OldPrice:  99.99,
			NewPrice:  89.99,
			ChangedAt: "2024-01-01T10:00:00Z",
			ChangedBy: "admin",
			Reason:    "promotion",
		},
		{
			ProductID: productID,
			OldPrice:  89.99,
			NewPrice:  79.99,
			ChangedAt: "2024-02-01T10:00:00Z",
			ChangedBy: "admin",
			Reason:    "clearance",
		},
	}

	s.logger.Info("Product price history retrieved", "product_id", productID)
	return history, nil
}

// GetProductStockHistory 获取产品库存历史
func (s *productService) GetProductStockHistory(ctx context.Context, productID uint) ([]*StockHistory, error) {
	// 这里简化处理，实际应用中应该从库存历史表获取
	// 返回模拟数据
	history := []*StockHistory{
		{
			ProductID:   productID,
			StockChange: 100,
			OldStock:    0,
			NewStock:    100,
			ChangedAt:   "2024-01-01T10:00:00Z",
			ChangedBy:   "admin",
			Reason:      "initial_stock",
		},
		{
			ProductID:   productID,
			StockChange: -5,
			OldStock:    100,
			NewStock:    95,
			ChangedAt:   "2024-01-02T10:00:00Z",
			ChangedBy:   "system",
			Reason:      "order_fulfillment",
			ReferenceID: "ORD-001",
		},
	}

	s.logger.Info("Product stock history retrieved", "product_id", productID)
	return history, nil
}

// ValidateProductSKU 验证产品SKU唯一性
func (s *productService) ValidateProductSKU(ctx context.Context, sku string) error {
	if sku == "" {
		return fmt.Errorf("sku cannot be empty")
	}

	product, err := s.productRepo.GetBySKU(ctx, sku)
	if err != nil {
		// 如果产品不存在，则SKU可用
		return nil
	}

	if product != nil {
		return fmt.Errorf("sku already exists")
	}

	return nil
}

// GetProductsBySKUs 根据多个SKU获取产品
func (s *productService) GetProductsBySKUs(ctx context.Context, skus []string) ([]*model.Product, error) {
	var products []*model.Product

	for _, sku := range skus {
		product, err := s.productRepo.GetBySKU(ctx, sku)
		if err != nil {
			continue // SKU不存在时跳过
		}
		if product != nil {
			products = append(products, product)
		}
	}

	s.logger.Info("Products retrieved by SKUs", "sku_count", len(skus), "found_count", len(products))
	return products, nil
}

// UpdateProductStock 更新产品库存
func (s *productService) UpdateProductStock(ctx context.Context, productID uint, quantityChange int, reason string) (*model.Product, error) {
	// 获取产品
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 验证库存变更
	newStock := product.StockQuantity + quantityChange
	if newStock < 0 {
		return nil, fmt.Errorf("insufficient stock: cannot reduce %d from current stock %d", quantityChange, product.StockQuantity)
	}

	// 更新库存
	if err := s.productRepo.UpdateStock(ctx, productID, quantityChange); err != nil {
		s.logger.Error("Failed to update product stock", "product_id", productID, "quantity_change", quantityChange, "reason", reason, "error", err)
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	// 重新获取更新后的产品
	updatedProduct, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	s.logger.Info("Product stock updated", "product_id", productID, "new_stock", updatedProduct.StockQuantity, "reason", reason)
	return updatedProduct, nil
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
