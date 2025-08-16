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
	SearchProducts(ctx context.Context, req *SearchProductRequest) ([]*model.Product, int64, error)
	GetProductsByCategory(ctx context.Context, categoryID uint, includeSubCategories bool, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetLowStockProducts(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error)
	GetPopularProducts(ctx context.Context, limit int) ([]*model.Product, error)

	// 价格管理
	BatchUpdatePrices(ctx context.Context, updates []PriceUpdate) error
	UpdateProductStatus(ctx context.Context, productID uint, isActive bool) error

	// 库存管理
	UpdateStock(ctx context.Context, productID uint, quantity int, operation StockOperation) error
	CheckStockAvailability(ctx context.Context, productID uint, requiredQuantity int) (bool, error)
	GetStockAlert(ctx context.Context) ([]*model.Product, error)
}

// productService Product服务实现
type productService struct {
	productRepo         repository.ProductRepository
	productCategoryRepo repository.ProductCategoryRepository
	logger              logger.Logger
}

// NewProductService 创建Product服务
func NewProductService(
	productRepo repository.ProductRepository,
	productCategoryRepo repository.ProductCategoryRepository,
	logger logger.Logger,
) ProductService {
	return &productService{
		productRepo:         productRepo,
		productCategoryRepo: productCategoryRepo,
		logger:              logger,
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

// SearchProductRequest 产品搜索请求
type SearchProductRequest struct {
	Keyword    string   `json:"keyword"`
	CategoryID *uint    `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
	SKU        string   `json:"sku,omitempty"`
}

// PriceUpdate 价格更新
type PriceUpdate struct {
	ProductID uint    `json:"product_id"`
	Price     float64 `json:"price"`
	CostPrice float64 `json:"cost_price"`
}

// StockOperation 库存操作类型
type StockOperation string

const (
	StockOperationAdd      StockOperation = "add"      // 增加库存
	StockOperationSubtract StockOperation = "subtract" // 减少库存
	StockOperationSet      StockOperation = "set"      // 设置库存
)

// SearchProducts 搜索产品
func (s *productService) SearchProducts(ctx context.Context, req *SearchProductRequest) ([]*model.Product, int64, error) {
	// 构建过滤器
	filters := make(map[string]interface{})

	if req.CategoryID != nil {
		filters["category_id"] = *req.CategoryID
	}
	if req.MinPrice != nil {
		filters["min_price"] = *req.MinPrice
	}
	if req.MaxPrice != nil {
		filters["max_price"] = *req.MaxPrice
	}
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}
	if req.SKU != "" {
		filters["sku"] = req.SKU
	}

	// 构建查询选项
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 50, // 默认分页
		Sort:     "created_at",
		Order:    "desc",
	}

	products, total, err := s.productRepo.SearchProducts(ctx, req.Keyword, filters, opts)
	if err != nil {
		s.logger.Error("Failed to search products", "error", err)
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}

// GetProductsByCategory 根据分类获取产品
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID uint, includeSubCategories bool, opts *ListProductOptions) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	var err error

	// 转换查询选项
	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Search:   opts.Search,
		Filters:  opts.Filters,
	}

	if includeSubCategories {
		// 获取所有子分类ID
		categoryIDs := []uint{categoryID}
		children, err := s.productCategoryRepo.GetByParentID(ctx, categoryID)
		if err != nil {
			s.logger.Error("Failed to get child categories", "category_id", categoryID, "error", err)
			return nil, 0, fmt.Errorf("failed to get child categories: %w", err)
		}

		for _, child := range children {
			categoryIDs = append(categoryIDs, child.ID)
		}

		products, err = s.productRepo.GetByCategoryWithSubCategories(ctx, categoryIDs, repoOpts)
		if err != nil {
			return nil, 0, err
		}
		total = int64(len(products)) // 简化处理，实际应该单独计算总数
	} else {
		products, err = s.productRepo.GetByCategory(ctx, categoryID, repoOpts)
		if err != nil {
			return nil, 0, err
		}
		total = int64(len(products)) // 简化处理，实际应该单独计算总数
	}

	return products, total, nil
}

// GetProductsByPriceRange 根据价格区间获取产品
func (s *productService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts *ListProductOptions) ([]*model.Product, int64, error) {
	// 转换查询选项
	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Search:   opts.Search,
		Filters:  opts.Filters,
	}

	products, err := s.productRepo.GetByPriceRange(ctx, minPrice, maxPrice, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return products, int64(len(products)), nil
}

// GetLowStockProducts 获取低库存产品
func (s *productService) GetLowStockProducts(ctx context.Context, opts *ListProductOptions) ([]*model.Product, int64, error) {
	// 转换查询选项
	repoOpts := repository.ListOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Search:   opts.Search,
		Filters:  opts.Filters,
	}

	// 使用最小库存作为阈值，这里设为10
	threshold := 10
	products, err := s.productRepo.GetLowStockProducts(ctx, threshold, repoOpts)
	if err != nil {
		s.logger.Error("Failed to get low stock products", "threshold", threshold, "error", err)
		return nil, 0, fmt.Errorf("failed to get low stock products: %w", err)
	}

	return products, int64(len(products)), nil
}

// GetPopularProducts 获取热销产品
func (s *productService) GetPopularProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	products, err := s.productRepo.GetPopularProducts(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to get popular products", "limit", limit, "error", err)
		return nil, fmt.Errorf("failed to get popular products: %w", err)
	}

	return products, nil
}

// BatchUpdatePrices 批量更新价格
func (s *productService) BatchUpdatePrices(ctx context.Context, updates []PriceUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	// 转换为repository需要的格式
	repoUpdates := make(map[uint]map[string]float64)
	for _, update := range updates {
		repoUpdates[update.ProductID] = map[string]float64{
			"price":      update.Price,
			"cost_price": update.CostPrice,
		}
	}

	err := s.productRepo.BatchUpdatePrices(ctx, repoUpdates)
	if err != nil {
		s.logger.Error("Failed to batch update prices", "error", err)
		return fmt.Errorf("failed to batch update prices: %w", err)
	}

	return nil
}

// UpdateProductStatus 更新产品状态
func (s *productService) UpdateProductStatus(ctx context.Context, productID uint, isActive bool) error {
	err := s.productRepo.UpdateStatus(ctx, productID, isActive)
	if err != nil {
		s.logger.Error("Failed to update product status", "product_id", productID, "is_active", isActive, "error", err)
		return fmt.Errorf("failed to update product status: %w", err)
	}

	return nil
}

// UpdateStock 更新库存
func (s *productService) UpdateStock(ctx context.Context, productID uint, quantity int, operation StockOperation) error {
	// 获取当前产品信息
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.logger.Error("Failed to get product for stock update", "product_id", productID, "error", err)
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	var newQuantity int
	switch operation {
	case StockOperationAdd:
		newQuantity = product.StockQuantity + quantity
	case StockOperationSubtract:
		newQuantity = product.StockQuantity - quantity
		if newQuantity < 0 {
			return fmt.Errorf("insufficient stock: current %d, requested %d", product.StockQuantity, quantity)
		}
	case StockOperationSet:
		newQuantity = quantity
		if newQuantity < 0 {
			return fmt.Errorf("stock quantity cannot be negative")
		}
	default:
		return fmt.Errorf("invalid stock operation: %s", operation)
	}

	err = s.productRepo.UpdateStock(ctx, productID, newQuantity)
	if err != nil {
		s.logger.Error("Failed to update stock", "product_id", productID, "new_quantity", newQuantity, "error", err)
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

// CheckStockAvailability 检查库存可用性
func (s *productService) CheckStockAvailability(ctx context.Context, productID uint, requiredQuantity int) (bool, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		s.logger.Error("Failed to get product for stock check", "product_id", productID, "error", err)
		return false, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return false, fmt.Errorf("product not found")
	}

	return product.StockQuantity >= requiredQuantity, nil
}

// GetStockAlert 获取库存警报
func (s *productService) GetStockAlert(ctx context.Context) ([]*model.Product, error) {
	// 使用默认的查询选项
	opts := &ListProductOptions{
		Page:     1,
		PageSize: 100,
		Sort:     "stock_quantity",
		Order:    "asc",
		Filters:  make(map[string]interface{}),
	}

	products, _, err := s.GetLowStockProducts(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to get stock alert", "error", err)
		return nil, fmt.Errorf("failed to get stock alert: %w", err)
	}

	// 过滤出需要警报的产品（库存低于最小库存）
	var alertProducts []*model.Product
	for _, product := range products {
		if product.StockQuantity <= product.MinStock {
			alertProducts = append(alertProducts, product)
		}
	}

	return alertProducts, nil
}
