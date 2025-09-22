package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// Part4AnswerOptionHandler Part4AnswerOption处理器
type Part4AnswerOptionHandler struct {
	part4AnswerOptionService service.Part4AnswerOptionService
	logger         logger.Logger
}

// NewPart4AnswerOptionHandler 创建Part4AnswerOption处理器
func NewPart4AnswerOptionHandler(
	part4AnswerOptionService service.Part4AnswerOptionService,
	logger logger.Logger,
) *Part4AnswerOptionHandler {
	return &Part4AnswerOptionHandler{
		part4AnswerOptionService: part4AnswerOptionService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *Part4AnswerOptionHandler) RegisterRoutes(r *gin.RouterGroup) {
	part4answeroptions := r.Group("/part4answeroptions")
	{

		part4answeroptions.POST("", h.Create)
		part4answeroptions.GET("", h.List)
		part4answeroptions.GET("/:id", h.GetByID)
		part4answeroptions.PUT("/:id", h.Update)
		part4answeroptions.DELETE("/:id", h.Delete)
	}
}

// Create 创建Part4AnswerOption
// @Summary 创建Part4AnswerOption
// @Description 创建新的Part4AnswerOption
// @Tags part4answeroptions
// @Accept json
// @Produce json
// @Param request body service.CreatePart4AnswerOptionRequest true "创建Part4AnswerOption请求"
// @Success 201 {object} model.Part4AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4answeroptions [post]
func (h *Part4AnswerOptionHandler) Create(c *gin.Context) {
	var req service.CreatePart4AnswerOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create part4answeroption request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part4AnswerOption, err := h.part4AnswerOptionService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create part4answeroption", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, part4AnswerOption)
}

// GetByID 根据ID获取Part4AnswerOption
// @Summary 获取Part4AnswerOption详情
// @Description 根据ID获取Part4AnswerOption详情
// @Tags part4answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part4AnswerOption ID"
// @Success 200 {object} model.Part4AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4answeroptions/{id} [get]
func (h *Part4AnswerOptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4AnswerOption ID",
		})
		return
	}

	part4AnswerOption, err := h.part4AnswerOptionService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get part4answeroption", "part4answeroption_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "part4answeroption_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part4AnswerOption)
}

// Update 更新Part4AnswerOption
// @Summary 更新Part4AnswerOption
// @Description 更新Part4AnswerOption信息
// @Tags part4answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part4AnswerOption ID"
// @Param request body service.UpdatePart4AnswerOptionRequest true "更新Part4AnswerOption请求"
// @Success 200 {object} model.Part4AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4answeroptions/{id} [put]
func (h *Part4AnswerOptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4AnswerOption ID",
		})
		return
	}

	var req service.UpdatePart4AnswerOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update part4answeroption request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part4AnswerOption, err := h.part4AnswerOptionService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update part4answeroption", "part4answeroption_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part4AnswerOption)
}

// Delete 删除Part4AnswerOption
// @Summary 删除Part4AnswerOption
// @Description 删除Part4AnswerOption
// @Tags part4answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part4AnswerOption ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4answeroptions/{id} [delete]
func (h *Part4AnswerOptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4AnswerOption ID",
		})
		return
	}

	if err := h.part4AnswerOptionService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete part4answeroption", "part4answeroption_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Part4AnswerOption deleted successfully",
	})
}

// List 获取Part4AnswerOption列表
// @Summary 获取Part4AnswerOption列表
// @Description 获取Part4AnswerOption列表，支持分页、搜索和过滤
// @Tags part4answeroptions
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4answeroptions [get]
func (h *Part4AnswerOptionHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListPart4AnswerOptionOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	part4answeroptions, total, err := h.part4AnswerOptionService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get part4answeroptions", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_part4answeroptions_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  part4answeroptions,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *Part4AnswerOptionHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreatePart4AnswerOptionRequest 创建Part4AnswerOption请求
type CreatePart4AnswerOptionRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdatePart4AnswerOptionRequest 更新Part4AnswerOption请求
type UpdatePart4AnswerOptionRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
