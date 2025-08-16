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

		// 层级分类管理API
		productcategories.GET("/tree", h.GetCategoryTree)
		productcategories.GET("/:id/path", h.GetCategoryPath)
		productcategories.GET("/:id/children", h.GetChildrenByParentID)
		productcategories.GET("/:id/product-count", h.GetCategoryWithProductCount)
		productcategories.PUT("/:id/sort-order", h.UpdateSortOrder)
		productcategories.POST("/batch-sort", h.BatchUpdateSortOrder)
		productcategories.GET("/:id/can-delete", h.CanDeleteCategory)
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

// GetCategoryTree 获取分类树结构
// @Summary 获取分类树结构
// @Description 获取指定父分类下的分类树结构
// @Tags productcategories
// @Accept json
// @Produce json
// @Param parent_id query int false "父分类ID" default(0)
// @Success 200 {array} service.CategoryTreeNode
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/tree [get]
func (h *ProductCategoryHandler) GetCategoryTree(c *gin.Context) {
	parentID, _ := strconv.ParseUint(c.DefaultQuery("parent_id", "0"), 10, 32)

	tree, err := h.productCategoryService.GetCategoryTree(c.Request.Context(), uint(parentID))
	if err != nil {
		h.logger.Error("Failed to get category tree", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_tree_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetCategoryPath 获取分类路径
// @Summary 获取分类路径（面包屑导航）
// @Description 获取指定分类的完整路径（从根到当前分类）
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {array} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id}/path [get]
func (h *ProductCategoryHandler) GetCategoryPath(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid category ID",
		})
		return
	}

	path, err := h.productCategoryService.GetCategoryPath(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get category path", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_path_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, path)
}

// GetChildrenByParentID 获取子分类
// @Summary 获取指定父分类的子分类列表
// @Description 获取指定父分类下的所有子分类
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "父分类ID"
// @Success 200 {array} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id}/children [get]
func (h *ProductCategoryHandler) GetChildrenByParentID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid category ID",
		})
		return
	}

	children, err := h.productCategoryService.GetChildrenByParentID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get children categories", "parent_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_children_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, children)
}

// UpdateSortOrder 更新分类排序
// @Summary 更新分类排序
// @Description 更新指定分类的排序值
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param request body SortOrderRequest true "排序更新请求"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id}/sort-order [put]
func (h *ProductCategoryHandler) UpdateSortOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid category ID",
		})
		return
	}

	var req SortOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid sort order request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.productCategoryService.UpdateSortOrder(c.Request.Context(), uint(id), req.SortOrder); err != nil {
		h.logger.Error("Failed to update sort order", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Sort order updated successfully",
	})
}

// BatchUpdateSortOrder 批量更新排序
// @Summary 批量更新分类排序
// @Description 批量更新多个分类的排序值
// @Tags productcategories
// @Accept json
// @Produce json
// @Param request body BatchSortRequest true "批量排序请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/batch-sort [post]
func (h *ProductCategoryHandler) BatchUpdateSortOrder(c *gin.Context) {
	var req BatchSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid batch sort request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.productCategoryService.BatchUpdateSortOrder(c.Request.Context(), req.SortUpdates); err != nil {
		h.logger.Error("Failed to batch update sort orders", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "batch_update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Sort orders updated successfully",
	})
}

// GetCategoryWithProductCount 获取分类及其产品数量
// @Summary 获取分类及其产品数量
// @Description 获取指定分类的信息及其关联的产品数量
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} service.CategoryWithProductCount
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id}/product-count [get]
func (h *ProductCategoryHandler) GetCategoryWithProductCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid category ID",
		})
		return
	}

	category, err := h.productCategoryService.GetCategoryWithProductCount(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get category with product count", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CanDeleteCategory 检查分类是否可删除
// @Summary 检查分类是否可删除
// @Description 检查指定分类是否可以删除（无子分类且无关联产品）
// @Tags productcategories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/productcategories/{id}/can-delete [get]
func (h *ProductCategoryHandler) CanDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid category ID",
		})
		return
	}

	canDelete, err := h.productCategoryService.CanDeleteCategory(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to check if category can be deleted", "category_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "check_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"can_delete": canDelete})
}

// 辅助方法

func (h *ProductCategoryHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	parentID, _ := strconv.ParseUint(c.DefaultQuery("parent_id", "0"), 10, 32)

	// 构建过滤器
	filters := make(map[string]interface{})
	if parentID > 0 {
		filters["parent_id"] = uint(parentID)
	}

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
	ParentId    uint   `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// UpdateProductCategoryRequest 更新ProductCategory请求
type UpdateProductCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	ParentId    *uint   `json:"parent_id,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// SortOrderRequest 排序更新请求
type SortOrderRequest struct {
	SortOrder int `json:"sort_order" validate:"required,min=0"`
}

// BatchSortRequest 批量排序请求
type BatchSortRequest struct {
	SortUpdates map[uint]int `json:"sort_updates" validate:"required"`
}
