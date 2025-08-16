package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// ProductHandler Product处理器
// @Description Product管理处理器，提供产品的CRUD操作和高级业务功能
// @Tags products
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
type ProductHandler struct {
	productService service.ProductService
	logger         logger.Logger
}

// NewProductHandler 创建Product处理器
func NewProductHandler(
	productService service.ProductService,
	logger logger.Logger,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	{
		products.POST("", h.Create)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}

// Create 创建Product
// @Summary 创建Product
// @Description 创建新的Product产品信息，包含完整的产品属性如价格、库存、SKU等
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param request body service.CreateProductRequest true "创建Product请求体"
// @Success 201 {object} model.Product "创建成功的产品信息"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 409 {object} ErrorResponse "产品已存在"
// @Failure 422 {object} ErrorResponse "数据验证失败"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/products [post]
// @Example {json} 请求示例:
// {
//   "name": "iPhone 15 Pro",
//   "description": "最新款iPhone，配备A17芯片",
//   "category_id": 1,
//   "sku": "IPHONE15PRO-128GB-NAT",
//   "price": 7999.00,
//   "cost_price": 6500.00,
//   "stock_quantity": 100,
//   "min_stock": 10,
//   "is_active": true,
//   "weight": 0.187,
//   "dimensions": "146.6×70.6×8.25mm"
// }
// @Example {json} 成功响应示例:
// {
//   "id": 1,
//   "name": "iPhone 15 Pro",
//   "description": "最新款iPhone，配备A17芯片",
//   "category_id": 1,
//   "sku": "IPHONE15PRO-128GB-NAT",
//   "price": 7999.00,
//   "cost_price": 6500.00,
//   "stock_quantity": 100,
//   "min_stock": 10,
//   "is_active": true,
//   "weight": 0.187,
//   "dimensions": "146.6×70.6×8.25mm",
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-15T10:30:00Z"
// }
// @Example {json} 错误响应示例:
// {
//   "error": "validation_error",
//   "message": "SKU already exists"
// }
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create product request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetByID 根据ID获取Product
// @Summary 获取Product详情
// @Description 根据ID获取Product详细信息，包括所有产品属性和关联信息
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Product ID" min(1) example(1)
// @Success 200 {object} model.Product "产品详细信息"
// @Failure 400 {object} ErrorResponse "无效的ID格式"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "产品不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/products/{id} [get]
// @Example {json} 成功响应示例:
// {
//   "id": 1,
//   "name": "iPhone 15 Pro",
//   "description": "最新款iPhone，配备A17芯片",
//   "category_id": 1,
//   "sku": "IPHONE15PRO-128GB-NAT",
//   "price": 7999.00,
//   "cost_price": 6500.00,
//   "stock_quantity": 100,
//   "min_stock": 10,
//   "is_active": true,
//   "weight": 0.187,
//   "dimensions": "146.6×70.6×8.25mm",
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-15T10:30:00Z"
// }
// @Example {json} 404错误响应示例:
// {
//   "error": "product_not_found",
//   "message": "Product with ID 999 not found"
// }
// @Example {json} 400错误响应示例:
// {
//   "error": "invalid_id",
//   "message": "Invalid Product ID"
// }
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Product ID",
		})
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get product", "product_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "product_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update 更新Product
// @Summary 更新Product
// @Description 更新Product信息，支持部分字段更新，未提供的字段将保持不变
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Product ID" min(1) example(1)
// @Param request body service.UpdateProductRequest true "更新Product请求体"
// @Success 200 {object} model.Product "更新后的产品信息"
// @Failure 400 {object} ErrorResponse "无效的ID格式或请求参数"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "产品不存在"
// @Failure 409 {object} ErrorResponse "SKU已存在"
// @Failure 422 {object} ErrorResponse "数据验证失败"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/products/{id} [put]
// @Example {json} 请求示例（部分更新）:
// {
//   "price": 7499.00,
//   "stock_quantity": 150,
//   "is_active": false
// }
// @Example {json} 成功响应示例:
// {
//   "id": 1,
//   "name": "iPhone 15 Pro",
//   "description": "最新款iPhone，配备A17芯片",
//   "category_id": 1,
//   "sku": "IPHONE15PRO-128GB-NAT",
//   "price": 7499.00,
//   "cost_price": 6500.00,
//   "stock_quantity": 150,
//   "min_stock": 10,
//   "is_active": false,
//   "weight": 0.187,
//   "dimensions": "146.6×70.6×8.25mm",
//   "created_at": "2024-01-15T10:30:00Z",
//   "updated_at": "2024-01-16T14:20:00Z"
// }
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Product ID",
		})
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update product request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update product", "product_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Delete 删除Product
// @Summary 删除Product
// @Description 删除指定的Product产品信息，删除后无法恢复
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param id path int true "Product ID" min(1) example(1)
// @Success 200 {object} SuccessResponse "删除成功"
// @Failure 400 {object} ErrorResponse "无效的ID格式"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 404 {object} ErrorResponse "产品不存在"
// @Failure 409 {object} ErrorResponse "产品有关联数据，无法删除"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/products/{id} [delete]
// @Example {json} 成功响应示例:
// {
//   "message": "Product deleted successfully"
// }
// @Example {json} 404错误响应示例:
// {
//   "error": "delete_failed",
//   "message": "Product with ID 999 not found"
// }
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Product ID",
		})
		return
	}

	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete product", "product_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// List 获取Product列表
// @Summary 获取Product列表
// @Description 获取Product产品列表，支持分页、搜索、过滤和排序功能
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" default(Bearer <your-token>)
// @Param page query int false "页码" default(1) minimum(1)
// @Param page_size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param search query string false "搜索关键词（支持产品名称、描述、SKU模糊搜索）"
// @Param sort query string false "排序字段" default(created_at) Enums(id, name, price, stock_quantity, created_at, updated_at)
// @Param order query string false "排序方式" default(desc) Enums(asc, desc)
// @Param category_id query int false "分类ID过滤"
// @Param is_active query bool false "激活状态过滤"
// @Param min_price query number false "最低价格过滤"
// @Param max_price query number false "最高价格过滤"
// @Param min_stock query int false "最小库存过滤"
// @Param max_stock query int false "最大库存过滤"
// @Success 200 {object} ListResponse "包含分页信息的产品列表"
// @Success 206 {object} ListResponse "部分内容（分页数据）"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "未授权访问"
// @Failure 422 {object} ErrorResponse "过滤参数格式错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/products [get]
// @Example {json} 成功响应示例:
// {
//   "data": [
//     {
//       "id": 1,
//       "name": "iPhone 15 Pro",
//       "description": "最新款iPhone，配备A17芯片",
//       "category_id": 1,
//       "sku": "IPHONE15PRO-128GB-NAT",
//       "price": 7999.00,
//       "cost_price": 6500.00,
//       "stock_quantity": 100,
//       "min_stock": 10,
//       "is_active": true,
//       "weight": 0.187,
//       "dimensions": "146.6×70.6×8.25mm",
//       "created_at": "2024-01-15T10:30:00Z",
//       "updated_at": "2024-01-15T10:30:00Z"
//     }
//   ],
//   "total": 25,
//   "page": 1,
//   "size": 10
// }
// @Example {json} 带搜索参数请求:
// GET /api/v1/products?search=iPhone&page=1&page_size=5&sort=price&order=asc
func (h *ProductHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListProductOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	products, total, err := h.productService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get products", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_products_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  products,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *ProductHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreateProductRequest 创建Product请求
type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateProductRequest 更新Product请求
type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
