package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// DictHandler 数据字典处理器
type DictHandler struct {
	dictService service.DictService
	logger      logger.Logger
}

// NewDictHandler 创建数据字典处理器
func NewDictHandler(
	dictService service.DictService,
	logger logger.Logger,
) *DictHandler {
	return &DictHandler{
		dictService: dictService,
		logger:      logger,
	}
}

// GetCategories 获取所有字典分类
// @Summary 获取所有字典分类
// @Description 获取系统中所有的数据字典分类
// @Tags dict
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]model.DictCategory} "获取成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/categories [get]
func (h *DictHandler) GetCategories(c *gin.Context) {
	categories, err := h.dictService.GetDictCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get dict categories", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_categories_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    categories,
	})
}

// GetItemsByCategory 根据分类获取字典项
// @Summary 根据分类获取字典项
// @Description 根据分类代码获取该分类下的所有字典项
// @Tags dict
// @Accept json
// @Produce json
// @Param category path string true "分类代码"
// @Success 200 {object} Response{data=[]model.DictItem} "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/items/{category} [get]
func (h *DictHandler) GetItemsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		h.logger.Error("Category parameter is required")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_parameter",
			Message: "Category parameter is required",
		})
		return
	}

	items, err := h.dictService.GetDictItems(c.Request.Context(), category)
	if err != nil {
		h.logger.Error("Failed to get dict items", "category", category, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_items_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    items,
	})
}

// GetItemByKey 获取特定字典项
// @Summary 获取特定字典项
// @Description 根据分类代码和项键值获取特定的字典项
// @Tags dict
// @Accept json
// @Produce json
// @Param category path string true "分类代码"
// @Param key path string true "项键值"
// @Success 200 {object} Response{data=model.DictItem} "获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 404 {object} ErrorResponse "字典项不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/item/{category}/{key} [get]
func (h *DictHandler) GetItemByKey(c *gin.Context) {
	category := c.Param("category")
	key := c.Param("key")

	if category == "" || key == "" {
		h.logger.Error("Category and key parameters are required")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_parameters",
			Message: "Category and key parameters are required",
		})
		return
	}

	item, err := h.dictService.GetDictItemByKey(c.Request.Context(), category, key)
	if err != nil {
		h.logger.Error("Failed to get dict item", "category", category, "key", key, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "item_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    item,
	})
}

// CreateCategory 创建字典分类
// @Summary 创建字典分类
// @Description 创建新的数据字典分类
// @Tags dict
// @Accept json
// @Produce json
// @Param request body service.CreateCategoryRequest true "分类创建数据"
// @Success 201 {object} Response{data=model.DictCategory} "创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 409 {object} ErrorResponse "分类已存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/categories [post]
func (h *DictHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	category, err := h.dictService.CreateDictCategory(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create dict category", "code", req.Code, "error", err)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "Category created successfully",
		Data:    category,
	})
}

// DeleteCategory 删除字典分类
// @Summary 删除字典分类
// @Description 删除指定ID的字典分类
// @Tags dict
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} Response{data=map[string]string} "删除成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 404 {object} ErrorResponse "分类不存在"
// @Failure 409 {object} ErrorResponse "分类下还有字典项，无法删除"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/categories/{id} [delete]
func (h *DictHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a positive integer",
		})
		return
	}

	if err := h.dictService.DeleteDictCategory(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete dict category", "id", id, "error", err)
		// 检查是否是因为分类下还有字典项
		if strings.Contains(err.Error(), "contains") && strings.Contains(err.Error(), "items") {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "category_has_items",
				Message: err.Error(),
			})
		} else {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "delete_failed",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Dict category deleted successfully",
		Data:    map[string]string{"message": "Dict category deleted successfully"},
	})
}

// CreateItem 创建字典项
// @Summary 创建字典项
// @Description 创建新的数据字典项
// @Tags dict
// @Accept json
// @Produce json
// @Param request body service.CreateItemRequest true "字典项创建数据"
// @Success 201 {object} Response{data=model.DictItem} "创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 409 {object} ErrorResponse "字典项已存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/items [post]
func (h *DictHandler) CreateItem(c *gin.Context) {
	var req service.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	item, err := h.dictService.CreateDictItem(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create dict item", "category_code", req.CategoryCode, "item_key", req.ItemKey, "error", err)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "Dict item created successfully",
		Data:    item,
	})
}

// UpdateItem 更新字典项
// @Summary 更新字典项
// @Description 更新指定ID的字典项
// @Tags dict
// @Accept json
// @Produce json
// @Param id path int true "字典项ID"
// @Param request body service.UpdateItemRequest true "字典项更新数据"
// @Success 200 {object} Response{data=model.DictItem} "更新成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 404 {object} ErrorResponse "字典项不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/items/{id} [put]
func (h *DictHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a positive integer",
		})
		return
	}

	var req service.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	item, err := h.dictService.UpdateDictItem(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update dict item", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Dict item updated successfully",
		Data:    item,
	})
}

// DeleteItem 删除字典项
// @Summary 删除字典项
// @Description 删除指定ID的字典项
// @Tags dict
// @Accept json
// @Produce json
// @Param id path int true "字典项ID"
// @Success 204 "删除成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 404 {object} ErrorResponse "字典项不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/items/{id} [delete]
func (h *DictHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a positive integer",
		})
		return
	}

	if err := h.dictService.DeleteDictItem(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete dict item", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Dict item deleted successfully",
		Data:    map[string]string{"message": "Dict item deleted successfully"},
	})
}

// InitDefaultData 初始化默认数据
// @Summary 初始化默认数据
// @Description 初始化系统默认的数据字典数据
// @Tags dict
// @Accept json
// @Produce json
// @Success 200 {object} Response "初始化成功"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /api/v1/dict/init [post]
func (h *DictHandler) InitDefaultData(c *gin.Context) {
	if err := h.dictService.InitDefaultDictData(c.Request.Context()); err != nil {
		h.logger.Error("Failed to initialize default dict data", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "init_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Default dictionary data initialized successfully",
		Data:    nil,
	})
}

// RegisterRoutes 注册路由
func (h *DictHandler) RegisterRoutes(r *gin.RouterGroup) {
	dict := r.Group("/dict")
	{
		// 获取分类和字典项
		dict.GET("/categories", h.GetCategories)
		dict.GET("/items/:category", h.GetItemsByCategory)
		dict.GET("/item/:category/:key", h.GetItemByKey)

		// 创建分类和字典项
		dict.POST("/categories", h.CreateCategory)
		dict.POST("/items", h.CreateItem)

		// 删除分类
		dict.DELETE("/categories/:id", h.DeleteCategory)

		// 更新和删除字典项
		dict.PUT("/items/:id", h.UpdateItem)
		dict.DELETE("/items/:id", h.DeleteItem)

		// 初始化默认数据
		dict.POST("/init", h.InitDefaultData)
	}
}

// 响应结构体

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
