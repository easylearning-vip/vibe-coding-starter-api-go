package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// DepartmentHandler Department处理器
// @Description Department部门管理处理器，提供完整的部门CRUD操作和树形结构管理功能
// @Tags departments
// @Security BearerAuth
// @Accept json
// @Produce json
// @BasePath /api/v1
// @Schemes http https
// @Host localhost:8081
// @Contact.name API Support
// @Contact.email support@example.com
// @License.name MIT
// @License.url https://opensource.org/licenses/MIT
// @ExternalDocs.description 更多API文档
// @ExternalDocs.url https://github.com/your-org/vibe-coding-starter-api-go
// @TermsOfService http://example.com/terms/
type DepartmentHandler struct {
	departmentService service.DepartmentService
	logger            logger.Logger
}

// NewDepartmentHandler 创建Department处理器
func NewDepartmentHandler(
	departmentService service.DepartmentService,
	logger logger.Logger,
) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
		logger:            logger,
	}
}

// RegisterRoutes 注册路由
func (h *DepartmentHandler) RegisterRoutes(r *gin.RouterGroup) {
	departments := r.Group("/departments")
	{
		departments.POST("", h.Create)
		departments.GET("", h.List)
		departments.GET("/tree", h.GetTree)
		departments.GET("/:id", h.GetByID)
		departments.GET("/:id/children", h.GetChildren)
		departments.GET("/:id/path", h.GetPath)
		departments.PUT("/:id", h.Update)
		departments.PUT("/:id/move", h.Move)
		departments.DELETE("/:id", h.Delete)
	}
}

// Create 创建Department
// @Summary 创建Department
// @Description 创建新的Department部门信息，支持树形结构管理，可指定父部门
// @Tags departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param request body service.CreateDepartmentRequest true "创建Department请求体"
// @Success 201 {object} model.Department "创建成功的部门信息"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "父部门不存在"
// @Failure 409 {object} ErrorResponse "部门已存在"
// @Failure 422 {object} ErrorResponse "数据验证失败"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/departments [post]
// @Example {json} 请求示例（根部门）:
// {
//   "name": "技术部",
//   "code": "TECH",
//   "description": "负责公司技术开发和维护",
//   "parent_id": 0,
//   "sort": 1,
//   "status": "active",
//   "manager_id": 1
// }
// @Example {json} 请求示例（子部门）:
// {
//   "name": "前端开发组",
//   "code": "FE-DEV",
//   "description": "负责前端应用开发",
//   "parent_id": 1,
//   "sort": 1,
//   "status": "active",
//   "manager_id": 2
// }
// @Example {json} 成功响应示例:
// {
//   "id": 2,
//   "name": "前端开发组",
//   "code": "FE-DEV",
//   "description": "负责前端应用开发",
//   "parent_id": 1,
//   "sort": 1,
//   "status": "active",
//   "manager_id": 2,
//   "parent": null,
//   "children": [],
//   "path": "1",
//   "level": 2,
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-15T10:30:00Z"
// }
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req service.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create department request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	department, err := h.departmentService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create department", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, department)
}

// GetByID 根据ID获取Department
// @Summary 获取Department详情
// @Description 根据ID获取Department详细信息，包含完整的树形结构和关联信息
// @Tags departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Department ID" min(1) example(1)
// @Success 200 {object} model.Department "部门详细信息"
// @Failure 400 {object} ErrorResponse "无效的ID格式"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "部门不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/departments/{id} [get]
// @Example {json} 成功响应示例:
// {
//   "id": 1,
//   "name": "技术部",
//   "code": "TECH",
//   "description": "负责公司技术开发和维护",
//   "parent_id": 0,
//   "sort": 1,
//   "status": "active",
//   "manager_id": 1,
//   "parent": null,
//   "children": [
//     {
//       "id": 2,
//       "name": "前端开发组",
//       "code": "FE-DEV",
//       "description": "负责前端应用开发",
//       "parent_id": 1,
//       "sort": 1,
//       "status": "active",
//       "manager_id": 2
//     }
//   ],
//   "path": "",
//   "level": 1,
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-15T10:30:00Z"
// }
func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	department, err := h.departmentService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get department", "department_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "department_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, department)
}

// Update 更新Department
// @Summary 更新Department
// @Description 更新Department部门信息，支持部分字段更新，可修改部门层级关系
// @Tags departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Department ID" min(1) example(1)
// @Param request body service.UpdateDepartmentRequest true "更新Department请求体"
// @Success 200 {object} model.Department "更新后的部门信息"
// @Failure 400 {object} ErrorResponse "无效的ID格式或请求参数"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "部门不存在或父部门不存在"
// @Failure 409 {object} ErrorResponse "部门代码已存在"
// @Failure 422 {object} ErrorResponse "数据验证失败"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/departments/{id} [put]
// @Example {json} 请求示例（部分更新）:
// {
//   "name": "技术研发部",
//   "description": "负责公司技术研发和创新",
//   "manager_id": 3
// }
// @Example {json} 成功响应示例:
// {
//   "id": 1,
//   "name": "技术研发部",
//   "code": "TECH",
//   "description": "负责公司技术研发和创新",
//   "parent_id": 0,
//   "sort": 1,
//   "status": "active",
//   "manager_id": 3,
//   "parent": null,
//   "children": [],
//   "path": "",
//   "level": 1,
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-16T14:20:00Z"
// }
func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	var req service.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update department request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	department, err := h.departmentService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update department", "department_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, department)
}

// Delete 删除Department
// @Summary 删除Department
// @Description 删除指定的Department部门信息，删除后无法恢复，如果存在子部门将无法删除
// @Tags departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Department ID" min(1) example(1)
// @Success 200 {object} SuccessResponse "删除成功"
// @Failure 400 {object} ErrorResponse "无效的ID格式"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "部门不存在"
// @Failure 409 {object} ErrorResponse "存在子部门，无法删除"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/departments/{id} [delete]
// @Example {json} 成功响应示例:
// {
//   "message": "Department deleted successfully"
// }
// @Example {json} 409错误响应示例:
// {
//   "error": "delete_failed",
//   "message": "Cannot delete department with child departments"
// }
func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	if err := h.departmentService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete department", "department_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department deleted successfully",
	})
}

// List 获取Department列表
// @Summary 获取Department列表
// @Description 获取Department列表，支持分页、搜索和过滤
// @Tags departments
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments [get]
func (h *DepartmentHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListDepartmentOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	departments, total, err := h.departmentService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get departments", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_departments_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  departments,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *DepartmentHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	// 构建过滤器
	filters := make(map[string]interface{})
	// 在这里添加特定的过滤器逻辑

	return repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Sort:     c.DefaultQuery("sort", "created_at"),
		Order:    c.DefaultQuery("order", "desc"),
		Filters:  filters,
	}
}

// 请求结构体

// CreateDepartmentRequest 创建Department请求
type CreateDepartmentRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateDepartmentRequest 更新Department请求
type UpdateDepartmentRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// GetTree 获取部门树结构
// @Summary 获取部门树结构
// @Description 获取完整的部门树结构
// @Tags departments
// @Accept json
// @Produce json
// @Success 200 {object} []model.Department
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/tree [get]
func (h *DepartmentHandler) GetTree(c *gin.Context) {
	tree, err := h.departmentService.GetTree(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get department tree", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_tree_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetChildren 获取子部门
// @Summary 获取子部门
// @Description 获取指定部门的直接子部门
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} []model.Department
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{id}/children [get]
func (h *DepartmentHandler) GetChildren(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	children, err := h.departmentService.GetChildren(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get department children", "department_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_children_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, children)
}

// GetPath 获取部门路径
// @Summary 获取部门路径
// @Description 获取从根部门到指定部门的完整路径
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} []model.Department
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{id}/path [get]
func (h *DepartmentHandler) GetPath(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	path, err := h.departmentService.GetPath(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get department path", "department_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_path_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, path)
}

// Move 移动部门
// @Summary 移动部门
// @Description 将部门移动到新的父部门下
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body MoveDepartmentRequest true "移动部门请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{id}/move [put]
func (h *DepartmentHandler) Move(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Department ID",
		})
		return
	}

	var req MoveDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid move department request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.departmentService.Move(c.Request.Context(), uint(id), req.NewParentId); err != nil {
		h.logger.Error("Failed to move department", "department_id", id, "new_parent_id", req.NewParentId, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "move_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Department moved successfully",
	})
}

// MoveDepartmentRequest 移动部门请求
type MoveDepartmentRequest struct {
	NewParentId uint `json:"new_parent_id" validate:"min=0"`
}
