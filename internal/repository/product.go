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
		case "sku":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("sku = ?", v)
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
		case "min_stock":
			if v, ok := value.(int); ok {
				query = query.Where("stock_quantity <= ?", v)
			}
			// 在这里添加更多过滤器
		}
	}
	return query
}

// GetBySKU 根据SKU获取产品
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

// GetByCategoryID 根据分类ID获取产品
func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", categoryID)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by category: %w", err)
	}

	// 应用排序和分页
	query = r.applySortingAndPaging(query, opts)

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}

	return entities, total, nil
}

// GetByCategoryWithSubcategories 根据分类ID获取产品（包含子分类）
func (r *productRepository) GetByCategoryWithSubcategories(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	// 获取所有子分类ID
	subCategoryIDs := r.getSubCategoryIDs(ctx, categoryID)
	allCategoryIDs := append([]uint{categoryID}, subCategoryIDs...)

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id IN ?", allCategoryIDs)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by category with subcategories", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by category with subcategories: %w", err)
	}

	// 应用排序和分页
	query = r.applySortingAndPaging(query, opts)

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by category with subcategories", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by category with subcategories: %w", err)
	}

	return entities, total, nil
}

// GetHotSellingProducts 获取热销产品
func (r *productRepository) GetHotSellingProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	var entities []*model.Product

	// 假设热销产品是库存量低且活跃的产品
	// 在实际应用中，这可能基于销售订单数据
	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("is_active = ? AND stock_quantity > 0", true).
		Order("stock_quantity ASC").
		Limit(limit)

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get hot selling products", "error", err)
		return nil, fmt.Errorf("failed to get hot selling products: %w", err)
	}

	return entities, nil
}

// GetByPriceRange 根据价格区间获取产品
func (r *productRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by price range: %w", err)
	}

	// 应用排序和分页
	query = r.applySortingAndPaging(query, opts)

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by price range", "min_price", minPrice, "max_price", maxPrice, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by price range: %w", err)
	}

	return entities, total, nil
}

// SearchProducts 搜索产品
func (r *productRepository) SearchProducts(ctx context.Context, query string, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", 
			"%"+query+"%", "%"+query+"%", "%"+query+"%")

	// 获取总数
	if err := dbQuery.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count search products", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to count search products: %w", err)
	}

	// 应用排序和分页
	dbQuery = r.applySortingAndPaging(dbQuery, opts)

	// 执行查询
	if err := dbQuery.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to search products", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return entities, total, nil
}

// GetLowStockProducts 获取库存不足的产品
func (r *productRepository) GetLowStockProducts(ctx context.Context) ([]*model.Product, error) {
	var entities []*model.Product

	query := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("stock_quantity <= min_stock AND is_active = ?", true).
		Order("stock_quantity ASC")

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get low stock products", "error", err)
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}

	return entities, nil
}

// GetProductsByStatus 根据状态获取产品
func (r *productRepository) GetProductsByStatus(ctx context.Context, isActive bool, opts ListOptions) ([]*model.Product, int64, error) {
	var entities []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("is_active = ?", isActive)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ?", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count products by status", "is_active", isActive, "error", err)
		return nil, 0, fmt.Errorf("failed to count products by status: %w", err)
	}

	// 应用排序和分页
	query = r.applySortingAndPaging(query, opts)

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get products by status", "is_active", isActive, "error", err)
		return nil, 0, fmt.Errorf("failed to get products by status: %w", err)
	}

	return entities, total, nil
}

// getSubCategoryIDs 递归获取所有子分类ID
func (r *productRepository) getSubCategoryIDs(ctx context.Context, parentID uint) []uint {
	var ids []uint
	var categories []model.ProductCategory

	// 获取直接子分类
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&categories).Error; err != nil {
		r.logger.Error("Failed to get subcategories", "parent_id", parentID, "error", err)
		return ids
	}

	for _, category := range categories {
		ids = append(ids, category.ID)
		// 递归获取子分类的子分类
		subIDs := r.getSubCategoryIDs(ctx, category.ID)
		ids = append(ids, subIDs...)
	}

	return ids
}

// applySortingAndPacking 应用排序和分页
func (r *productRepository) applySortingAndPaging(query *gorm.DB, opts ListOptions) *gorm.DB {
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

	return query
}
