package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/logger"
)

// DepartmentService Department服务接口
type DepartmentService interface {
	Create(ctx context.Context, req *CreateDepartmentRequest) (*model.Department, error)
	GetByID(ctx context.Context, id uint) (*model.Department, error)
	Update(ctx context.Context, id uint, req *UpdateDepartmentRequest) (*model.Department, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts *ListDepartmentOptions) ([]*model.Department, int64, error)
	GetTree(ctx context.Context) ([]*model.Department, error)
	GetChildren(ctx context.Context, parentId uint) ([]*model.Department, error)
	GetPath(ctx context.Context, id uint) ([]*model.Department, error)
	Move(ctx context.Context, id uint, newParentId uint) error
}

// departmentService Department服务实现
type departmentService struct {
	departmentRepo repository.DepartmentRepository

	logger      logger.Logger
}

// NewDepartmentService 创建Department服务
func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,

	logger logger.Logger,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,

		logger:      logger,
	}
}

// CreateDepartmentRequest 创建Department请求
type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Code string `json:"code" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	ParentId uint `json:"parent_id" validate:"required,min=0"`
	Sort int `json:"sort" validate:"required,min=0"`
	Status string `json:"status" validate:"required,min=1,max=255"`
	ManagerId uint `json:"manager_id" validate:"required,min=0"`
}

// UpdateDepartmentRequest 更新Department请求
type UpdateDepartmentRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Code *string `json:"code,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	ParentId *uint `json:"parent_id,omitempty" validate:"omitempty,min=0"`
	Sort *int `json:"sort,omitempty" validate:"omitempty,min=0"`
	Status *string `json:"status,omitempty" validate:"omitempty,min=1,max=255"`
	ManagerId *uint `json:"manager_id,omitempty" validate:"omitempty,min=0"`
}

// ListDepartmentOptions 列表查询选项
type ListDepartmentOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// Create 创建Department
func (s *departmentService) Create(ctx context.Context, req *CreateDepartmentRequest) (*model.Department, error) {
	// 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 创建模型
	entity := &model.Department{
		Name: req.Name,
		Code: req.Code,
		Description: req.Description,
		ParentId: req.ParentId,
		Sort: req.Sort,
		Status: req.Status,
		ManagerId: req.ManagerId,
	}

	// 保存到数据库
	if err := s.departmentRepo.Create(ctx, entity); err != nil {
		s.logger.Error("Failed to create department", "error", err)
		return nil, fmt.Errorf("failed to create department: %w", err)
	}



	s.logger.Info("Department created successfully", "id", entity.ID)
	return entity, nil
}

// GetByID 根据ID获取Department
func (s *departmentService) GetByID(ctx context.Context, id uint) (*model.Department, error) {


	// 从数据库获取
	entity, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}



	return entity, nil
}

// Update 更新Department
func (s *departmentService) Update(ctx context.Context, id uint, req *UpdateDepartmentRequest) (*model.Department, error) {
	// 验证请求
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有实体
	entity, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		entity.Name = *req.Name
	}
	if req.Code != nil {
		entity.Code = *req.Code
	}
	if req.Description != nil {
		entity.Description = *req.Description
	}
	if req.ParentId != nil {
		entity.ParentId = *req.ParentId
	}
	if req.Sort != nil {
		entity.Sort = *req.Sort
	}
	if req.Status != nil {
		entity.Status = *req.Status
	}
	if req.ManagerId != nil {
		entity.ManagerId = *req.ManagerId
	}

	// 保存更新
	if err := s.departmentRepo.Update(ctx, entity); err != nil {
		s.logger.Error("Failed to update department", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update department: %w", err)
	}



	s.logger.Info("Department updated successfully", "id", id)
	return entity, nil
}

// Delete 删除Department
func (s *departmentService) Delete(ctx context.Context, id uint) error {
	// 检查实体是否存在
	if _, err := s.departmentRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("failed to get department: %w", err)
	}

	// 删除实体
	if err := s.departmentRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete department", "id", id, "error", err)
		return fmt.Errorf("failed to delete department: %w", err)
	}



	s.logger.Info("Department deleted successfully", "id", id)
	return nil
}

// List 获取Department列表
func (s *departmentService) List(ctx context.Context, opts *ListDepartmentOptions) ([]*model.Department, int64, error) {
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
	entities, total, err := s.departmentRepo.List(ctx, repoOpts)
	if err != nil {
		s.logger.Error("Failed to list departments", "error", err)
		return nil, 0, fmt.Errorf("failed to list departments: %w", err)
	}

	return entities, total, nil
}

// validateCreateRequest 验证创建请求
func (s *departmentService) validateCreateRequest(req *CreateDepartmentRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *departmentService) validateUpdateRequest(req *UpdateDepartmentRequest) error {
	// 使用 validate 标签进行验证
	// 这里可以添加自定义验证逻辑
	return nil
}

// GetTree 获取部门树结构
func (s *departmentService) GetTree(ctx context.Context) ([]*model.Department, error) {
	// 获取所有根部门（ParentId = 0）
	rootDepts, err := s.departmentRepo.GetByParentId(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get root departments: %w", err)
	}

	// 递归构建树结构
	for _, dept := range rootDepts {
		if err := s.buildTree(ctx, dept); err != nil {
			return nil, fmt.Errorf("failed to build tree: %w", err)
		}
	}

	return rootDepts, nil
}

// GetChildren 获取子部门
func (s *departmentService) GetChildren(ctx context.Context, parentId uint) ([]*model.Department, error) {
	children, err := s.departmentRepo.GetByParentId(ctx, parentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}
	return children, nil
}

// GetPath 获取部门路径
func (s *departmentService) GetPath(ctx context.Context, id uint) ([]*model.Department, error) {
	dept, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	if dept.Path == "" {
		return []*model.Department{dept}, nil
	}

	// 解析路径获取ID列表
	pathIds := s.parsePath(dept.Path)
	var pathDepts []*model.Department

	for _, pathId := range pathIds {
		pathDept, err := s.departmentRepo.GetByID(ctx, pathId)
		if err != nil {
			return nil, fmt.Errorf("failed to get path department: %w", err)
		}
		pathDepts = append(pathDepts, pathDept)
	}

	pathDepts = append(pathDepts, dept)
	return pathDepts, nil
}

// Move 移动部门
func (s *departmentService) Move(ctx context.Context, id uint, newParentId uint) error {
	// 检查是否移动到自己的子部门
	if newParentId != 0 {
		pathDepts, err := s.GetPath(ctx, newParentId)
		if err != nil {
			return fmt.Errorf("failed to check parent path: %w", err)
		}

		for _, pathDept := range pathDepts {
			if pathDept.ID == id {
				return fmt.Errorf("cannot move department to its own child")
			}
		}
	}

	// 获取当前部门
	dept, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get department: %w", err)
	}

	// 更新父部门ID
	dept.ParentId = newParentId

	// 重新计算路径和层级
	if newParentId == 0 {
		dept.Path = ""
		dept.Level = 1
	} else {
		parent, err := s.departmentRepo.GetByID(ctx, newParentId)
		if err != nil {
			return fmt.Errorf("failed to get parent department: %w", err)
		}
		dept.Level = parent.Level + 1
		if parent.Path == "" {
			dept.Path = fmt.Sprintf("%d", parent.ID)
		} else {
			dept.Path = fmt.Sprintf("%s,%d", parent.Path, parent.ID)
		}
	}

	// 保存更新
	if err := s.departmentRepo.Update(ctx, dept); err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	// 递归更新所有子部门的路径和层级
	if err := s.updateChildrenPaths(ctx, id); err != nil {
		return fmt.Errorf("failed to update children paths: %w", err)
	}

	s.logger.Info("Department moved successfully", "id", id, "new_parent_id", newParentId)
	return nil
}

// buildTree 递归构建树结构
func (s *departmentService) buildTree(ctx context.Context, parent *model.Department) error {
	children, err := s.departmentRepo.GetByParentId(ctx, parent.ID)
	if err != nil {
		return err
	}

	parent.Children = make([]model.Department, len(children))
	for i, child := range children {
		parent.Children[i] = *child
		if err := s.buildTree(ctx, &parent.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

// parsePath 解析路径字符串
func (s *departmentService) parsePath(path string) []uint {
	if path == "" {
		return []uint{}
	}

	var ids []uint
	for _, idStr := range strings.Split(path, ",") {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	return ids
}

// updateChildrenPaths 递归更新子部门路径
func (s *departmentService) updateChildrenPaths(ctx context.Context, parentId uint) error {
	children, err := s.departmentRepo.GetByParentId(ctx, parentId)
	if err != nil {
		return err
	}

	for _, child := range children {
		// 更新当前子部门的路径
		parent, err := s.departmentRepo.GetByID(ctx, parentId)
		if err != nil {
			return err
		}

		child.Level = parent.Level + 1
		if parent.Path == "" {
			child.Path = fmt.Sprintf("%d", parent.ID)
		} else {
			child.Path = fmt.Sprintf("%s,%d", parent.Path, parent.ID)
		}

		if err := s.departmentRepo.Update(ctx, child); err != nil {
			return err
		}

		// 递归更新子部门的子部门
		if err := s.updateChildrenPaths(ctx, child.ID); err != nil {
			return err
		}
	}

	return nil
}


