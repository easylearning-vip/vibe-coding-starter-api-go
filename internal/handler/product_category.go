package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// ProductCategoryHandler ProductCategory处理器
type ProductCategoryHandler struct {
	productCategoryService service.ProductCategoryService
	logger                 logger.Logger
}

// NewProductCategoryHandler 创建ProductCategory处理器
func NewProductCategoryHandler(
	productCategoryService service.ProductCategoryService,
	logger logger.Logger,
) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		productCategoryService: productCategoryService,
		logger:                 logger,
	}
}

// RegisterRoutes 注册路由
func (h *ProductCategoryHandler) RegisterRoutes(r *gin.RouterGroup) {
	productcategories := r.Group("/productcategories")
	{
		productcategories.POST("", h.Create)
		productcategories.GET("", h.List)
		productcategories.GET("/:id", h.GetByID)
		productcategories.PUT("/:id", h.Update)
		productcategories.DELETE("/:id", h.Delete)

		// 层级分类管理路由
		productcategories.GET("/tree", h.GetCategoryTree)
		productcategories.GET("/:id/children", h.GetChildren)
		productcategories.GET("/:id/path", h.GetCategoryPath)
		productcategories.POST("/batch-sort", h.BatchUpdateSortOrder)
		productcategories.GET("/:id/can-delete", h.CanDelete)
	}
}

// Create 创建ProductCategory
// @Summary 创建ProductCategory
// @Description 创建新的ProductCategory
// @Tags productcategories
// @Accept json
// @Produce json
// @Param request body service.CreateProductCategoryRequest true "创建ProductCategory请求"
// @Success 201 {object} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories [post]
func (h *ProductCategoryHandler) Create(c *gin.Context) {
	var req service.CreateProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create productcategory request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	productCategory, err := h.productCategoryService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create productcategory", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, productCategory)
}

// GetByID 根据ID获取ProductCategory
// @Summary 获取ProductCategory详情
// @Description 根据ID获取ProductCategory详情
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "ProductCategory ID"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id} [get]
func (h *ProductCategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid ProductCategory ID",
		})
		return
	}

	productCategory, err := h.productCategoryService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get productcategory", "productcategory_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "productcategory_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, productCategory)
}

// Update 更新ProductCategory
// @Summary 更新ProductCategory
// @Description 更新ProductCategory信息
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "ProductCategory ID"
// @Param request body service.UpdateProductCategoryRequest true "更新ProductCategory请求"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id} [put]
func (h *ProductCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid ProductCategory ID",
		})
		return
	}

	var req service.UpdateProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update productcategory request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	productCategory, err := h.productCategoryService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update productcategory", "productcategory_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, productCategory)
}

// Delete 删除ProductCategory
// @Summary 删除ProductCategory
// @Description 删除ProductCategory
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "ProductCategory ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id} [delete]
func (h *ProductCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid ProductCategory ID",
		})
		return
	}

	if err := h.productCategoryService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete productcategory", "productcategory_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "ProductCategory deleted successfully",
	})
}

// List 获取ProductCategory列表
// @Summary 获取ProductCategory列表
// @Description 获取ProductCategory列表，支持分页、搜索和过滤
// @Tags productcategories
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories [get]
func (h *ProductCategoryHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListProductCategoryOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	productcategories, total, err := h.productCategoryService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get productcategories", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_productcategories_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  productcategories,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *ProductCategoryHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreateProductCategoryRequest 创建ProductCategory请求
type CreateProductCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateProductCategoryRequest 更新ProductCategory请求
type UpdateProductCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// GetCategoryTree 获取分类树
// @Summary 获取分类树
// @Description 获取完整的分类树结构
// @Tags productcategories
// @Produce json
// @Success 200 {object} Response{data=[]service.CategoryTreeNode} "分类树"
// @Failure 500 {object} Response "服务器错误"
// @Router /api/v1/productcategories/tree [get]
func (h *ProductCategoryHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.productCategoryService.GetCategoryTree(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get category tree", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取分类树失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取分类树成功",
		"data":    tree,
	})
}

// GetChildren 获取子分类
// @Summary 获取子分类
// @Description 获取指定分类的子分类列表
// @Tags productcategories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} Response{data=[]model.ProductCategory} "子分类列表"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /api/v1/productcategories/{id}/children [get]
func (h *ProductCategoryHandler) GetChildren(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的分类ID",
			"data":    nil,
		})
		return
	}

	children, err := h.productCategoryService.GetChildren(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get children categories", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取子分类失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取子分类成功",
		"data":    children,
	})
}

// GetCategoryPath 获取分类路径
// @Summary 获取分类路径
// @Description 获取从根分类到指定分类的完整路径
// @Tags productcategories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} Response{data=[]model.ProductCategory} "分类路径"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /api/v1/productcategories/{id}/path [get]
func (h *ProductCategoryHandler) GetCategoryPath(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的分类ID",
			"data":    nil,
		})
		return
	}

	path, err := h.productCategoryService.GetCategoryPath(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get category path", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取分类路径失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取分类路径成功",
		"data":    path,
	})
}

// BatchUpdateSortOrderRequest 批量更新排序请求
type BatchUpdateSortOrderRequest struct {
	Updates []service.SortOrderUpdate `json:"updates" validate:"required"`
}

// BatchUpdateSortOrder 批量更新排序
// @Summary 批量更新排序
// @Description 批量更新分类的排序顺序
// @Tags productcategories
// @Accept json
// @Produce json
// @Param request body BatchUpdateSortOrderRequest true "批量更新排序请求"
// @Success 200 {object} Response "更新成功"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /api/v1/productcategories/batch-sort [post]
func (h *ProductCategoryHandler) BatchUpdateSortOrder(c *gin.Context) {
	var req BatchUpdateSortOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	err := h.productCategoryService.BatchUpdateSortOrder(c.Request.Context(), req.Updates)
	if err != nil {
		h.logger.Error("Failed to batch update sort order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量更新排序失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量更新排序成功",
		"data":    nil,
	})
}

// CanDeleteResponse 删除检查响应
type CanDeleteResponse struct {
	CanDelete bool   `json:"can_delete"`
	Reason    string `json:"reason,omitempty"`
}

// CanDelete 检查是否可以删除
// @Summary 检查是否可以删除
// @Description 检查指定分类是否可以删除
// @Tags productcategories
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} Response{data=CanDeleteResponse} "删除检查结果"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /api/v1/productcategories/{id}/can-delete [get]
func (h *ProductCategoryHandler) CanDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的分类ID",
			"data":    nil,
		})
		return
	}

	canDelete, reason, err := h.productCategoryService.CanDelete(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to check if category can be deleted", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查删除条件失败",
			"data":    nil,
		})
		return
	}

	response := CanDeleteResponse{
		CanDelete: canDelete,
		Reason:    reason,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "检查删除条件成功",
		"data":    response,
	})
}
