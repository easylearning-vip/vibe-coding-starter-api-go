package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// departmentRepository Department仓储实现
type departmentRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewDepartmentRepository 创建Department仓储
func NewDepartmentRepository(db database.Database, logger logger.Logger) DepartmentRepository {
	return &departmentRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建Department
func (r *departmentRepository) Create(ctx context.Context, entity *model.Department) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create <no value>", "error", err)
		return fmt.Errorf("failed to create <no value>: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Department
func (r *departmentRepository) GetByID(ctx context.Context, id uint) (*model.Department, error) {
	var entity model.Department
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Department not found")
		}
		r.logger.Error("Failed to get Department by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get Department: %w", err)
	}
	return &entity, nil
}

// GetByName 根据名称获取Department
func (r *departmentRepository) GetByName(ctx context.Context, name string) (*model.Department, error) {
	var entity model.Department
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Department not found")
		}
		r.logger.Error("Failed to get Department by name", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get Department: %w", err)
	}
	return &entity, nil
}

// Update 更新Department
func (r *departmentRepository) Update(ctx context.Context, entity *model.Department) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update Department", "id", entity.ID, "error", err)
		return fmt.Errorf("failed to update Department: %w", err)
	}
	return nil
}

// Delete 删除Department
func (r *departmentRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Department{}, id).Error; err != nil {
		r.logger.Error("Failed to delete Department", "id", id, "error", err)
		return fmt.Errorf("failed to delete Department: %w", err)
	}
	return nil
}

// List 获取Department列表
func (r *departmentRepository) List(ctx context.Context, opts ListOptions) ([]*model.Department, int64, error) {
	var entities []*model.Department
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Department{})

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
func (r *departmentRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("name = ?", v)
			}
			// 在这里添加更多过滤器
		}
	}
	return query
}

// GetByParentId 根据父部门ID获取子部门
func (r *departmentRepository) GetByParentId(ctx context.Context, parentId uint) ([]*model.Department, error) {
	var entities []*model.Department
	query := r.db.WithContext(ctx).Where("parent_id = ?", parentId).Order("sort ASC, created_at ASC")

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get departments by parent_id", "parent_id", parentId, "error", err)
		return nil, fmt.Errorf("failed to get departments by parent_id: %w", err)
	}

	return entities, nil
}

// GetByCode 根据代码获取部门
func (r *departmentRepository) GetByCode(ctx context.Context, code string) (*model.Department, error) {
	var entity model.Department
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Department not found")
		}
		r.logger.Error("Failed to get Department by code", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get Department: %w", err)
	}
	return &entity, nil
}

// GetChildrenTree 获取子部门树结构
func (r *departmentRepository) GetChildrenTree(ctx context.Context, parentId uint) ([]*model.Department, error) {
	var entities []*model.Department
	query := r.db.WithContext(ctx).Where("parent_id = ?", parentId).Order("sort ASC, created_at ASC")

	if err := query.Preload("Children").Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get children tree", "parent_id", parentId, "error", err)
		return nil, fmt.Errorf("failed to get children tree: %w", err)
	}

	return entities, nil
}
