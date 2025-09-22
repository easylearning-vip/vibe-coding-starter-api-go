package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// DifficultyLevelHandler DifficultyLevel处理器
type DifficultyLevelHandler struct {
	difficultyLevelService service.DifficultyLevelService
	logger         logger.Logger
}

// NewDifficultyLevelHandler 创建DifficultyLevel处理器
func NewDifficultyLevelHandler(
	difficultyLevelService service.DifficultyLevelService,
	logger logger.Logger,
) *DifficultyLevelHandler {
	return &DifficultyLevelHandler{
		difficultyLevelService: difficultyLevelService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *DifficultyLevelHandler) RegisterRoutes(r *gin.RouterGroup) {
	difficultylevels := r.Group("/difficultylevels")
	{

		difficultylevels.POST("", h.Create)
		difficultylevels.GET("", h.List)
		difficultylevels.GET("/:id", h.GetByID)
		difficultylevels.PUT("/:id", h.Update)
		difficultylevels.DELETE("/:id", h.Delete)
	}
}

// Create 创建DifficultyLevel
// @Summary 创建DifficultyLevel
// @Description 创建新的DifficultyLevel
// @Tags difficultylevels
// @Accept json
// @Produce json
// @Param request body service.CreateDifficultyLevelRequest true "创建DifficultyLevel请求"
// @Success 201 {object} model.DifficultyLevel
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/difficultylevels [post]
func (h *DifficultyLevelHandler) Create(c *gin.Context) {
	var req service.CreateDifficultyLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create difficultylevel request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	difficultyLevel, err := h.difficultyLevelService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create difficultylevel", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, difficultyLevel)
}

// GetByID 根据ID获取DifficultyLevel
// @Summary 获取DifficultyLevel详情
// @Description 根据ID获取DifficultyLevel详情
// @Tags difficultylevels
// @Accept json
// @Produce json
// @Param id path int true "DifficultyLevel ID"
// @Success 200 {object} model.DifficultyLevel
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/difficultylevels/{id} [get]
func (h *DifficultyLevelHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid DifficultyLevel ID",
		})
		return
	}

	difficultyLevel, err := h.difficultyLevelService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get difficultylevel", "difficultylevel_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "difficultylevel_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, difficultyLevel)
}

// Update 更新DifficultyLevel
// @Summary 更新DifficultyLevel
// @Description 更新DifficultyLevel信息
// @Tags difficultylevels
// @Accept json
// @Produce json
// @Param id path int true "DifficultyLevel ID"
// @Param request body service.UpdateDifficultyLevelRequest true "更新DifficultyLevel请求"
// @Success 200 {object} model.DifficultyLevel
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/difficultylevels/{id} [put]
func (h *DifficultyLevelHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid DifficultyLevel ID",
		})
		return
	}

	var req service.UpdateDifficultyLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update difficultylevel request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	difficultyLevel, err := h.difficultyLevelService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update difficultylevel", "difficultylevel_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, difficultyLevel)
}

// Delete 删除DifficultyLevel
// @Summary 删除DifficultyLevel
// @Description 删除DifficultyLevel
// @Tags difficultylevels
// @Accept json
// @Produce json
// @Param id path int true "DifficultyLevel ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/difficultylevels/{id} [delete]
func (h *DifficultyLevelHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid DifficultyLevel ID",
		})
		return
	}

	if err := h.difficultyLevelService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete difficultylevel", "difficultylevel_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "DifficultyLevel deleted successfully",
	})
}

// List 获取DifficultyLevel列表
// @Summary 获取DifficultyLevel列表
// @Description 获取DifficultyLevel列表，支持分页、搜索和过滤
// @Tags difficultylevels
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/difficultylevels [get]
func (h *DifficultyLevelHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListDifficultyLevelOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	difficultylevels, total, err := h.difficultyLevelService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get difficultylevels", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_difficultylevels_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  difficultylevels,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *DifficultyLevelHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreateDifficultyLevelRequest 创建DifficultyLevel请求
type CreateDifficultyLevelRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateDifficultyLevelRequest 更新DifficultyLevel请求
type UpdateDifficultyLevelRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
