package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// productCategoryRepository ProductCategory仓储实现
type productCategoryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewProductCategoryRepository 创建ProductCategory仓储
func NewProductCategoryRepository(db database.Database, logger logger.Logger) ProductCategoryRepository {
	return &productCategoryRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建ProductCategory
func (r *productCategoryRepository) Create(ctx context.Context, entity *model.ProductCategory) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取ProductCategory
func (r *productCategoryRepository) GetByID(ctx context.Context, id uint) (*model.ProductCategory, error) {
	var entity model.ProductCategory
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ProductCategory not found")
		}
		r.logger.Error("Failed to get ProductCategory by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get ProductCategory: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取ProductCategory
func (r *productCategoryRepository) GetByName(ctx context.Context, name string) (*model.ProductCategory, error) {
	var entity model.ProductCategory
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ProductCategory not found")
		}
		r.logger.Error("Failed to get ProductCategory by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get ProductCategory: %w", err)
	}
	return &entity, nil
}

// Update 更新ProductCategory
func (r *productCategoryRepository) Update(ctx context.Context, entity *model.ProductCategory) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update ProductCategory", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update ProductCategory: %w", err)
	}
	return nil
}

// Delete 删除ProductCategory
func (r *productCategoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.ProductCategory{}, id).Error; err != nil {
		r.logger.Error("Failed to delete ProductCategory", "id", id, "error", err)
		return fmt.Errorf("failed to delete ProductCategory: %w", err)
	}
	return nil
}

// List 获取ProductCategory列表
func (r *productCategoryRepository) List(ctx context.Context, opts ListOptions) ([]*model.ProductCategory, int64, error) {
	var entities []*model.ProductCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ProductCategory{})

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
func (r *productCategoryRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name = ?", v)
			}
		case "parent_id":
				if v, ok := value.(uint); ok {
					query = query.Where("parent_id = ?", v)
				}
		case "is_active":
				if v, ok := value.(bool); ok {
					query = query.Where("is_active = ?", v)
				}
			// 在这里添加更多过滤器
		}
	}
	return query
}

// GetByParentID 根据父分类ID获取子分类
func (r *productCategoryRepository) GetByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	var categories []*model.ProductCategory
	
	query := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort_order ASC, created_at ASC")
	if err := query.Find(&categories).Error; err != nil {
		r.logger.Error("Failed to get categories by parent ID", "parent_id", parentID, "error", err)
		return nil, fmt.Errorf("failed to get categories by parent ID: %w", err)
	}
	
	return categories, nil
}

// GetCategoryPath 获取分类路径（面包屑导航）
func (r *productCategoryRepository) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	var path []*model.ProductCategory
	currentID := categoryID
	
	for currentID > 0 {
		var category model.ProductCategory
		if err := r.db.WithContext(ctx).First(&category, currentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			r.logger.Error("Failed to get category for path", "category_id", currentID, "error", err)
			return nil, fmt.Errorf("failed to get category for path: %w", err)
		}
		
		path = append([]*model.ProductCategory{&category}, path...)
		currentID = category.ParentId
	}
	
	return path, nil
}

// GetCategoryTree 获取完整的分类树结构
func (r *productCategoryRepository) GetCategoryTree(ctx context.Context) ([]*model.ProductCategory, error) {
	var allCategories []*model.ProductCategory
	
	// 获取所有分类并按排序字段排序
	if err := r.db.WithContext(ctx).Order("parent_id ASC, sort_order ASC, created_at ASC").Find(&allCategories).Error; err != nil {
		r.logger.Error("Failed to get all categories for tree", "error", err)
		return nil, fmt.Errorf("failed to get all categories for tree: %w", err)
	}
	
	// 构建树结构
	return r.buildCategoryTree(allCategories, 0), nil
}

// buildCategoryTree 递归构建分类树
func (r *productCategoryRepository) buildCategoryTree(categories []*model.ProductCategory, parentID uint) []*model.ProductCategory {
	var tree []*model.ProductCategory
	
	for _, category := range categories {
		if category.ParentId == parentID {
			children := r.buildCategoryTree(categories, category.ID)
			if len(children) > 0 {
				category.Children = children
			}
			tree = append(tree, category)
		}
	}
	
	return tree
}

// CountProductsByCategory 统计分类下的产品数量
func (r *productCategoryRepository) CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	
	// 假设有一个Product模型，并且Product模型有CategoryID字段
	// 这里需要根据实际的Product模型结构调整
	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("category_id = ?", categoryID)
	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("Failed to count products by category", "category_id", categoryID, "error", err)
		return 0, fmt.Errorf("failed to count products by category: %w", err)
	}
	
	return count, nil
}

// HasChildren 检查分类是否有子分类
func (r *productCategoryRepository) HasChildren(ctx context.Context, categoryID uint) (bool, error) {
	var count int64
	
	query := r.db.WithContext(ctx).Model(&model.ProductCategory{}).Where("parent_id = ?", categoryID)
	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("Failed to count children categories", "category_id", categoryID, "error", err)
		return false, fmt.Errorf("failed to count children categories: %w", err)
	}
	
	return count > 0, nil
}
