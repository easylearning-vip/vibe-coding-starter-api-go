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
	logger         logger.Logger
}

// NewProductCategoryHandler 创建ProductCategory处理器
func NewProductCategoryHandler(
	productCategoryService service.ProductCategoryService,
	logger logger.Logger,
) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		productCategoryService: productCategoryService,
		logger:         logger,
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
