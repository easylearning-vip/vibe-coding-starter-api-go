package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// productRepository Product仓储实现
type productRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewProductRepository 创建Product仓储
func NewProductRepository(db database.Database, logger logger.Logger) ProductRepository {
	return &productRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Product
func (r *productRepository) Create(ctx context.Context, entity *model.Product) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Product
func (r *productRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Product not found")
		}
		r.logger.Error("Failed to get Product by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Product: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Product
func (r *productRepository) GetByName(ctx context.Context, name string) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Product not found")
		}
		r.logger.Error("Failed to get Product by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Product: %w", err)
	}
	return &entity, nil
}

// Update 更新Product
func (r *productRepository) Update(ctx context.Context, entity *model.Product) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Product", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Product: %w", err)
	}
	return nil
}

// Delete 删除Product
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Product{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Product", "id", id, "error", err)
		return fmt.Errorf("failed to delete Product: %w", err)
	}
	return nil
}

// List 获取Product列表
func (r *productRepository) List(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count <no value>", "error", err)
		return nil, 0, fmt.Errorf("failed to count <no value>: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to list <no value>", "error", err)
		return nil, 0, fmt.Errorf("failed to list <no value>: %w", err)
	}

	return entities, total, nil
}

// GetBySKU 根据SKU获取产品
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get product by SKU", "sku", sku, "error", err)
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}
	return &entity, nil
}

// SearchProducts 搜索产品
func (r *productRepository) SearchProducts(ctx context.Context, keyword string, filters map[string]interface{}, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 应用过滤器
	query = r.applyFilters(query, filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count search results", "error", err)
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// 应用排序和分页
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order != "" {
			order = opts.Order
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	}

	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to search products", "error", err)
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return entities, total, nil
}

// GetByCategory 根据分类获取产品
func (r *productRepository) GetByCategory(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Where("category_id = ?", categoryID)

	// 应用排序和分页
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order != "" {
			order = opts.Order
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	}

	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by category", "category_id", categoryID, "error", err)
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	return entities, nil
}

// GetByCategoryWithSubCategories 根据分类（包含子分类）获取产品
func (r *productRepository) GetByCategoryWithSubCategories(ctx context.Context, categoryIDs []uint, opts ListOptions) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Where("category_id IN ?", categoryIDs)

	// 应用排序和分页
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order != "" {
			order = opts.Order
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	}

	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by categories", "category_ids", categoryIDs, "error", err)
		return nil, fmt.Errorf("failed to get products by categories: %w", err)
	}

	return entities, nil
}

// GetByPriceRange 根据价格区间获取产品
func (r *productRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts ListOptions) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	// 应用排序和分页
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order != "" {
			order = opts.Order
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	}

	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return entities, nil
}

// GetLowStockProducts 获取低库存产品
func (r *productRepository) GetLowStockProducts(ctx context.Context, threshold int, opts ListOptions) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Where("stock_quantity <= ? AND stock_quantity > 0", threshold)

	// 应用排序和分页
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order != "" {
			order = opts.Order
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("stock_quantity ASC")
	}

	if opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get low stock products", "threshold", threshold, "error", err)
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}

	return entities, nil
}

// GetPopularProducts 获取热销产品（这里简单按创建时间排序，实际应该根据销量）
func (r *productRepository) GetPopularProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Where("is_active = ?", true).Order("created_at DESC").Limit(limit)

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get popular products", "limit", limit, "error", err)
		return nil, fmt.Errorf("failed to get popular products: %w", err)
	}

	return entities, nil
}

// BatchUpdatePrices 批量更新价格
func (r *productRepository) BatchUpdatePrices(ctx context.Context, updates map[uint]map[string]float64) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for productID, prices := range updates {
		if err := tx.Model(&model.Product{}).Where("id = ?", productID).Updates(prices).Error; err != nil {
			tx.Rollback()
			r.logger.Error("Failed to update product price", "product_id", productID, "prices", prices, "error", err)
			return fmt.Errorf("failed to update product price for product %d: %w", productID, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit batch price update", "error", err)
		return fmt.Errorf("failed to commit batch price update: %w", err)
	}

	return nil
}

// UpdateStock 更新库存
func (r *productRepository) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).Update("stock_quantity", quantity).Error; err != nil {
		r.logger.Error("Failed to update stock", "product_id", productID, "quantity", quantity, "error", err)
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

// UpdateStatus 更新产品状态
func (r *productRepository) UpdateStatus(ctx context.Context, productID uint, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productID).Update("is_active", isActive).Error; err != nil {
		r.logger.Error("Failed to update product status", "product_id", productID, "is_active", isActive, "error", err)
		return fmt.Errorf("failed to update product status: %w", err)
	}
	return nil
}

// applyFilters 应用过滤器
func (r *productRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name = ?", v)
			}
		case "category_id":
			if v, ok := value.(uint); ok {
				query = query.Where("category_id = ?", v)
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				query = query.Where("is_active = ?", v)
			}
		case "min_price":
			if v, ok := value.(float64); ok {
				query = query.Where("price >= ?", v)
			}
		case "max_price":
			if v, ok := value.(float64); ok {
				query = query.Where("price <= ?", v)
			}
		case "sku":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("sku = ?", v)
			}
			// 在这里添加更多过滤器
		}
	}
	return query
}
