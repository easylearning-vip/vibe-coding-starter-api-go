package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// Part3AnswerOptionHandler Part3AnswerOption处理器
type Part3AnswerOptionHandler struct {
	part3AnswerOptionService service.Part3AnswerOptionService
	logger         logger.Logger
}

// NewPart3AnswerOptionHandler 创建Part3AnswerOption处理器
func NewPart3AnswerOptionHandler(
	part3AnswerOptionService service.Part3AnswerOptionService,
	logger logger.Logger,
) *Part3AnswerOptionHandler {
	return &Part3AnswerOptionHandler{
		part3AnswerOptionService: part3AnswerOptionService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *Part3AnswerOptionHandler) RegisterRoutes(r *gin.RouterGroup) {
	part3answeroptions := r.Group("/part3answeroptions")
	{

		part3answeroptions.POST("", h.Create)
		part3answeroptions.GET("", h.List)
		part3answeroptions.GET("/:id", h.GetByID)
		part3answeroptions.PUT("/:id", h.Update)
		part3answeroptions.DELETE("/:id", h.Delete)
	}
}

// Create 创建Part3AnswerOption
// @Summary 创建Part3AnswerOption
// @Description 创建新的Part3AnswerOption
// @Tags part3answeroptions
// @Accept json
// @Produce json
// @Param request body service.CreatePart3AnswerOptionRequest true "创建Part3AnswerOption请求"
// @Success 201 {object} model.Part3AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3answeroptions [post]
func (h *Part3AnswerOptionHandler) Create(c *gin.Context) {
	var req service.CreatePart3AnswerOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create part3answeroption request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part3AnswerOption, err := h.part3AnswerOptionService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create part3answeroption", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, part3AnswerOption)
}

// GetByID 根据ID获取Part3AnswerOption
// @Summary 获取Part3AnswerOption详情
// @Description 根据ID获取Part3AnswerOption详情
// @Tags part3answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part3AnswerOption ID"
// @Success 200 {object} model.Part3AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3answeroptions/{id} [get]
func (h *Part3AnswerOptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3AnswerOption ID",
		})
		return
	}

	part3AnswerOption, err := h.part3AnswerOptionService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get part3answeroption", "part3answeroption_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "part3answeroption_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part3AnswerOption)
}

// Update 更新Part3AnswerOption
// @Summary 更新Part3AnswerOption
// @Description 更新Part3AnswerOption信息
// @Tags part3answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part3AnswerOption ID"
// @Param request body service.UpdatePart3AnswerOptionRequest true "更新Part3AnswerOption请求"
// @Success 200 {object} model.Part3AnswerOption
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3answeroptions/{id} [put]
func (h *Part3AnswerOptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3AnswerOption ID",
		})
		return
	}

	var req service.UpdatePart3AnswerOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update part3answeroption request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part3AnswerOption, err := h.part3AnswerOptionService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update part3answeroption", "part3answeroption_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part3AnswerOption)
}

// Delete 删除Part3AnswerOption
// @Summary 删除Part3AnswerOption
// @Description 删除Part3AnswerOption
// @Tags part3answeroptions
// @Accept json
// @Produce json
// @Param id path int true "Part3AnswerOption ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3answeroptions/{id} [delete]
func (h *Part3AnswerOptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3AnswerOption ID",
		})
		return
	}

	if err := h.part3AnswerOptionService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete part3answeroption", "part3answeroption_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Part3AnswerOption deleted successfully",
	})
}

// List 获取Part3AnswerOption列表
// @Summary 获取Part3AnswerOption列表
// @Description 获取Part3AnswerOption列表，支持分页、搜索和过滤
// @Tags part3answeroptions
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3answeroptions [get]
func (h *Part3AnswerOptionHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListPart3AnswerOptionOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	part3answeroptions, total, err := h.part3AnswerOptionService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get part3answeroptions", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_part3answeroptions_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  part3answeroptions,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *Part3AnswerOptionHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreatePart3AnswerOptionRequest 创建Part3AnswerOption请求
type CreatePart3AnswerOptionRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdatePart3AnswerOptionRequest 更新Part3AnswerOption请求
type UpdatePart3AnswerOptionRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
