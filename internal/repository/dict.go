package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// dictCategoryRepository 数据字典分类仓储实现
type dictCategoryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDictCategoryRepository 创建数据字典分类仓储
func NewDictCategoryRepository(db database.Database, logger logger.Logger) DictCategoryRepository {
	return &dictCategoryRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建字典分类
func (r *dictCategoryRepository) Create(ctx context.Context, category *model.DictCategory) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		r.logger.Error("Failed to create dict category", "error", err)
		return fmt.Errorf("failed to create dict category: %w", err)
	}
	return nil
}

// GetByID 根据ID获取字典分类
func (r *dictCategoryRepository) GetByID(ctx context.Context, id uint) (*model.DictCategory, error) {
	var category model.DictCategory
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dict category not found with id %d", id)
		}
		r.logger.Error("Failed to get dict category by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get dict category: %w", err)
	}
	return &category, nil
}

// Update 更新字典分类
func (r *dictCategoryRepository) Update(ctx context.Context, category *model.DictCategory) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		r.logger.Error("Failed to update dict category", "id", category.ID, "error", err)
		return fmt.Errorf("failed to update dict category: %w", err)
	}
	return nil
}

// Delete 删除字典分类
func (r *dictCategoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DictCategory{}, id).Error; err != nil {
		r.logger.Error("Failed to delete dict category", "id", id, "error", err)
		return fmt.Errorf("failed to delete dict category: %w", err)
	}
	return nil
}

// List 获取字典分类列表
func (r *dictCategoryRepository) List(ctx context.Context, opts ListOptions) ([]*model.DictCategory, int64, error) {
	var categories []*model.DictCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DictCategory{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("code LIKE ? OR name LIKE ? OR description LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count dict categories", "error", err)
		return nil, 0, fmt.Errorf("failed to count dict categories: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("sort_order ASC, created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&categories).Error; err != nil {
		r.logger.Error("Failed to list dict categories", "error", err)
		return nil, 0, fmt.Errorf("failed to list dict categories: %w", err)
	}

	return categories, total, nil
}

// GetByCode 根据代码获取字典分类
func (r *dictCategoryRepository) GetByCode(ctx context.Context, code string) (*model.DictCategory, error) {
	var category model.DictCategory
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dict category not found with code %s", code)
		}
		r.logger.Error("Failed to get dict category by code", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get dict category: %w", err)
	}
	return &category, nil
}

// applyFilters 应用过滤器
func (r *dictCategoryRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "code":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("code = ?", v)
			}
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name LIKE ?", "%"+v+"%")
			}
		}
	}

	return query
}

// dictItemRepository 数据字典项仓储实现
type dictItemRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDictItemRepository 创建数据字典项仓储
func NewDictItemRepository(db database.Database, logger logger.Logger) DictItemRepository {
	return &dictItemRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建字典项
func (r *dictItemRepository) Create(ctx context.Context, item *model.DictItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		r.logger.Error("Failed to create dict item", "error", err)
		return fmt.Errorf("failed to create dict item: %w", err)
	}
	return nil
}

// GetByID 根据ID获取字典项
func (r *dictItemRepository) GetByID(ctx context.Context, id uint) (*model.DictItem, error) {
	var item model.DictItem
	if err := r.db.WithContext(ctx).Preload("Category").First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dict item not found with id %d", id)
		}
		r.logger.Error("Failed to get dict item by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get dict item: %w", err)
	}
	return &item, nil
}

// Update 更新字典项
func (r *dictItemRepository) Update(ctx context.Context, item *model.DictItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		r.logger.Error("Failed to update dict item", "id", item.ID, "error", err)
		return fmt.Errorf("failed to update dict item: %w", err)
	}
	return nil
}

// Delete 删除字典项
func (r *dictItemRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.DictItem{}, id).Error; err != nil {
		r.logger.Error("Failed to delete dict item", "id", id, "error", err)
		return fmt.Errorf("failed to delete dict item: %w", err)
	}
	return nil
}

// List 获取字典项列表
func (r *dictItemRepository) List(ctx context.Context, opts ListOptions) ([]*model.DictItem, int64, error) {
	var items []*model.DictItem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.DictItem{}).Preload("Category")

	// 应用过滤器
	query = r.applyItemFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("category_code LIKE ? OR item_key LIKE ? OR item_value LIKE ? OR description LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count dict items", "error", err)
		return nil, 0, fmt.Errorf("failed to count dict items: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("category_code ASC, sort_order ASC, created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&items).Error; err != nil {
		r.logger.Error("Failed to list dict items", "error", err)
		return nil, 0, fmt.Errorf("failed to list dict items: %w", err)
	}

	return items, total, nil
}

// GetByCategory 根据分类获取字典项
func (r *dictItemRepository) GetByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	var items []*model.DictItem
	if err := r.db.WithContext(ctx).Where("category_code = ?", categoryCode).
		Order("sort_order ASC, created_at DESC").Find(&items).Error; err != nil {
		r.logger.Error("Failed to get dict items by category", "category_code", categoryCode, "error", err)
		return nil, fmt.Errorf("failed to get dict items: %w", err)
	}
	return items, nil
}

// GetActiveByCategory 根据分类获取启用的字典项
func (r *dictItemRepository) GetActiveByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	var items []*model.DictItem
	if err := r.db.WithContext(ctx).Where("category_code = ? AND is_active = ?", categoryCode, true).
		Order("sort_order ASC, created_at DESC").Find(&items).Error; err != nil {
		r.logger.Error("Failed to get active dict items by category", "category_code", categoryCode, "error", err)
		return nil, fmt.Errorf("failed to get active dict items: %w", err)
	}
	return items, nil
}

// GetByCategoryAndKey 根据分类和键值获取字典项
func (r *dictItemRepository) GetByCategoryAndKey(ctx context.Context, categoryCode, itemKey string) (*model.DictItem, error) {
	var item model.DictItem
	if err := r.db.WithContext(ctx).Where("category_code = ? AND item_key = ?", categoryCode, itemKey).
		First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dict item not found with category %s and key %s", categoryCode, itemKey)
		}
		r.logger.Error("Failed to get dict item by category and key", "category_code", categoryCode, "item_key", itemKey, "error", err)
		return nil, fmt.Errorf("failed to get dict item: %w", err)
	}
	return &item, nil
}

// applyItemFilters 应用字典项过滤器
func (r *dictItemRepository) applyItemFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "category_code":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("category_code = ?", v)
			}
		case "item_key":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("item_key = ?", v)
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				query = query.Where("is_active = ?", v)
			}
		}
	}

	return query
}

// dictRepository 数据字典组合仓储实现
type dictRepository struct {
	categoryRepo DictCategoryRepository
	itemRepo     DictItemRepository
}

// NewDictRepository 创建数据字典组合仓储
func NewDictRepository(db database.Database, logger logger.Logger) DictRepository {
	return &dictRepository{
		categoryRepo: NewDictCategoryRepository(db, logger),
		itemRepo:     NewDictItemRepository(db, logger),
	}
}

// GetAllCategories 获取所有字典分类
func (r *dictRepository) GetAllCategories(ctx context.Context) ([]*model.DictCategory, error) {
	categories, _, err := r.categoryRepo.List(ctx, ListOptions{
		Page:     1,
		PageSize: 1000, // 获取所有分类，假设不会超过1000个
		Sort:     "sort_order",
		Order:    "asc",
	})
	return categories, err
}

// GetCategoryByCode 根据分类代码获取分类
func (r *dictRepository) GetCategoryByCode(ctx context.Context, code string) (*model.DictCategory, error) {
	return r.categoryRepo.GetByCode(ctx, code)
}

// CreateCategory 创建字典分类
func (r *dictRepository) CreateCategory(ctx context.Context, category *model.DictCategory) error {
	return r.categoryRepo.Create(ctx, category)
}

// DeleteCategory 删除字典分类
func (r *dictRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.categoryRepo.Delete(ctx, id)
}

// GetItemsByCategory 根据分类获取所有字典项
func (r *dictRepository) GetItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	return r.itemRepo.GetByCategory(ctx, categoryCode)
}

// GetActiveItemsByCategory 根据分类获取启用的字典项
func (r *dictRepository) GetActiveItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	return r.itemRepo.GetActiveByCategory(ctx, categoryCode)
}

// GetItemByID 根据ID获取字典项
func (r *dictRepository) GetItemByID(ctx context.Context, id uint) (*model.DictItem, error) {
	return r.itemRepo.GetByID(ctx, id)
}

// CreateItem 创建字典项
func (r *dictRepository) CreateItem(ctx context.Context, item *model.DictItem) error {
	return r.itemRepo.Create(ctx, item)
}

// UpdateItem 更新字典项
func (r *dictRepository) UpdateItem(ctx context.Context, item *model.DictItem) error {
	return r.itemRepo.Update(ctx, item)
}

// DeleteItem 删除字典项
func (r *dictRepository) DeleteItem(ctx context.Context, id uint) error {
	return r.itemRepo.Delete(ctx, id)
}
