package repository

import (
	"context"
	"fmt"
	"strings"

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

// GetBySKU 根据SKU获取Product
func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	var entity model.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Product not found")
		}
		r.logger.Error("Failed to get Product by SKU", "sku", sku, "error", err)
		return nil, fmt.Errorf("failed to get Product: %w", err)
	}
	return &entity, nil
}

// GetByCategoryID 根据分类ID获取Product
func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", categoryID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by category: %w", err)
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
		r.logger.Error("Failed to get products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}

	return entities, total, nil
}

// GetByCategoryIDs 根据多个分类ID获取产品（包括子分类）
func (r *productRepository) GetByCategoryIDs(ctx context.Context, categoryIDs []uint, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id IN ?", categoryIDs)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by category IDs", "category_ids", categoryIDs, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by category IDs: %w", err)
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
		r.logger.Error("Failed to get products by category IDs", "category_ids", categoryIDs, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category IDs: %w", err)
	}

	return entities, total, nil
}

// GetHotSellingProducts 获取热销产品（按销量排序）
func (r *productRepository) GetHotSellingProducts(ctx context.Context, minSales int, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	// 这里简化处理，实际应用中应该从订单表获取销量数据
	// 暂时按创建时间倒序，假设新创建的产品更受欢迎
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("is_active = ? AND stock_quantity > ?", true, minSales)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count hot selling products", "error", err)
		return nil, 0, fmt.Errorf("failed to count hot selling products: %w", err)
	}

	// 应用排序 - 按价格降序（模拟热销）
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("price DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get hot selling products", "error", err)
		return nil, 0, fmt.Errorf("failed to get hot selling products: %w", err)
	}

	return entities, total, nil
}

// GetByPriceRange 根据价格范围获取Product
func (r *productRepository) GetByPriceRange(ctx context.Context, minPrice float64, maxPrice float64, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("price >= ? AND price <= ?", minPrice, maxPrice)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by price range", "error", err)
		return nil, 0, fmt.Errorf("failed to count products by price range: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("price ASC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by price range", "error", err)
		return nil, 0, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return entities, total, nil
}

// GetLowStockProducts 获取低库存产品
func (r *productRepository) GetLowStockProducts(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("stock_quantity < min_stock")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count low stock products", "error", err)
		return nil, 0, fmt.Errorf("failed to count low stock products: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("stock_quantity ASC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get low stock products", "error", err)
		return nil, 0, fmt.Errorf("failed to get low stock products: %w", err)
	}

	return entities, total, nil
}

// SearchByName 根据名称搜索Product
func (r *productRepository) SearchByName(ctx context.Context, query string, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	searchQuery := "%" + query + "%"
	queryDB := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("name LIKE ?", searchQuery)

	// 获取总数
	if err := queryDB.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by name", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by name: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		queryDB = queryDB.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		queryDB = queryDB.Order("created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		queryDB = queryDB.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := queryDB.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to search products by name", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search products by name: %w", err)
	}

	return entities, total, nil
}

// SearchBySKU 根据SKU搜索Product
func (r *productRepository) SearchBySKU(ctx context.Context, query string, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	searchQuery := "%" + query + "%"
	queryDB := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("sku LIKE ?", searchQuery)

	// 获取总数
	if err := queryDB.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by SKU", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by SKU: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		queryDB = queryDB.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		queryDB = queryDB.Order("created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		queryDB = queryDB.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := queryDB.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to search products by SKU", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search products by SKU: %w", err)
	}

	return entities, total, nil
}

// GetActiveProducts 获取激活的产品
func (r *productRepository) GetActiveProducts(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("is_active = ?", true)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count active products", "error", err)
		return nil, 0, fmt.Errorf("failed to count active products: %w", err)
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
		r.logger.Error("Failed to get active products", "error", err)
		return nil, 0, fmt.Errorf("failed to get active products: %w", err)
	}

	return entities, total, nil
}

// GetProductsInStock 获取有库存的产品
func (r *productRepository) GetProductsInStock(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("stock_quantity > 0")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products in stock", "error", err)
		return nil, 0, fmt.Errorf("failed to count products in stock: %w", err)
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
		r.logger.Error("Failed to get products in stock", "error", err)
		return nil, 0, fmt.Errorf("failed to get products in stock: %w", err)
	}

	return entities, total, nil
}

// UpdateStock 更新产品库存
func (r *productRepository) UpdateStock(ctx context.Context, productID uint, quantityChange int) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantityChange))

	if result.Error != nil {
		r.logger.Error("Failed to update product stock", "product_id", productID, "quantity_change", quantityChange, "error", result.Error)
		return fmt.Errorf("failed to update product stock: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	r.logger.Info("Product stock updated", "product_id", productID, "quantity_change", quantityChange)
	return nil
}

// BatchUpdatePrices 批量更新产品价格
func (r *productRepository) BatchUpdatePrices(ctx context.Context, updates map[uint]float64) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for productID, newPrice := range updates {
		result := tx.Model(&model.Product{}).
			Where("id = ?", productID).
			Update("price", newPrice)

		if result.Error != nil {
			tx.Rollback()
			r.logger.Error("Failed to update product price", "product_id", productID, "price", newPrice, "error", result.Error)
			return fmt.Errorf("failed to update price for product %d: %w", productID, result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("product %d not found", productID)
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit batch price update", "error", err)
		return fmt.Errorf("failed to commit batch price update: %w", err)
	}

	r.logger.Info("Batch prices updated", "count", len(updates))
	return nil
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

// applyFilters 应用过滤器
func (r *productRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		// 处理特殊操作符
		if strings.HasSuffix(key, "__gt") {
			field := strings.TrimSuffix(key, "__gt")
			query = query.Where(field+" > ?", value)
		} else if strings.HasSuffix(key, "__gte") {
			field := strings.TrimSuffix(key, "__gte")
			query = query.Where(field+" >= ?", value)
		} else if strings.HasSuffix(key, "__lt") {
			field := strings.TrimSuffix(key, "__lt")
			query = query.Where(field+" < ?", value)
		} else if strings.HasSuffix(key, "__lte") {
			field := strings.TrimSuffix(key, "__lte")
			query = query.Where(field+" <= ?", value)
		} else if strings.HasSuffix(key, "__ne") {
			field := strings.TrimSuffix(key, "__ne")
			query = query.Where(field+" != ?", value)
		} else {
			switch key {
			case "name":
				if v, ok := value.(string); ok && v != "" {
					query = query.Where("name = ?", v)
				}
			case "sku":
				if v, ok := value.(string); ok && v != "" {
					query = query.Where("sku = ?", v)
				}
			case "category_id":
				if v, ok := value.(uint); ok {
					query = query.Where("category_id = ?", v)
				}
			case "is_active":
				if v, ok := value.(bool); ok {
					query = query.Where("is_active = ?", v)
				}
			case "price":
				if v, ok := value.(float64); ok {
					query = query.Where("price = ?", v)
				}
			case "stock_quantity":
				if v, ok := value.(int); ok {
					query = query.Where("stock_quantity = ?", v)
				}
				// 在这里添加更多过滤器
			}
		}
	}
	return query
}
