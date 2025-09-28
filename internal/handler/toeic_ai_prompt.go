package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// ToeicAiPromptHandler ToeicAiPrompt处理器
type ToeicAiPromptHandler struct {
	toeicAiPromptService service.ToeicAiPromptService
	logger               logger.Logger
}

// NewToeicAiPromptHandler 创建ToeicAiPrompt处理器
func NewToeicAiPromptHandler(
	toeicAiPromptService service.ToeicAiPromptService,
	logger logger.Logger,
) *ToeicAiPromptHandler {
	return &ToeicAiPromptHandler{
		toeicAiPromptService: toeicAiPromptService,
		logger:               logger,
	}
}

// RegisterRoutes 注册路由
func (h *ToeicAiPromptHandler) RegisterRoutes(r *gin.RouterGroup) {
	toeicAiPrompts := r.Group("/toeic-ai-prompts")
	{
		toeicAiPrompts.POST("", h.Create)
		toeicAiPrompts.GET("", h.List)
		toeicAiPrompts.GET("/:id", h.GetByID)
		toeicAiPrompts.PUT("/:id", h.Update)
		toeicAiPrompts.DELETE("/:id", h.Delete)
	}
}

// Create 创建ToeicAiPrompt
// @Summary 创建ToeicAiPrompt
// @Description 创建新的TOEIC AI提示词
// @Tags toeic-ai-prompts
// @Accept json
// @Produce json
// @Param request body service.CreateToeicAiPromptRequest true "创建ToeicAiPrompt请求"
// @Success 201 {object} model.ToeicAiPrompt
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/toeic-ai-prompts [post]
func (h *ToeicAiPromptHandler) Create(c *gin.Context) {
	var req service.CreateToeicAiPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create toeic ai prompt request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	entity, err := h.toeicAiPromptService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create toeic ai prompt", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create toeic ai prompt",
		})
		return
	}

	c.JSON(http.StatusCreated, entity)
}

// GetByID 根据ID获取ToeicAiPrompt
// @Summary 获取ToeicAiPrompt详情
// @Description 根据ID获取TOEIC AI提示词详情
// @Tags toeic-ai-prompts
// @Accept json
// @Produce json
// @Param id path int true "ToeicAiPrompt ID"
// @Success 200 {object} model.ToeicAiPrompt
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/toeic-ai-prompts/{id} [get]
func (h *ToeicAiPromptHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid toeic ai prompt ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid toeic ai prompt ID",
		})
		return
	}

	entity, err := h.toeicAiPromptService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get toeic ai prompt", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "ToeicAiPrompt not found",
		})
		return
	}

	c.JSON(http.StatusOK, entity)
}

// Update 更新ToeicAiPrompt
// @Summary 更新ToeicAiPrompt
// @Description 更新TOEIC AI提示词信息
// @Tags toeic-ai-prompts
// @Accept json
// @Produce json
// @Param id path int true "ToeicAiPrompt ID"
// @Param request body service.UpdateToeicAiPromptRequest true "更新ToeicAiPrompt请求"
// @Success 200 {object} model.ToeicAiPrompt
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/toeic-ai-prompts/{id} [put]
func (h *ToeicAiPromptHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid toeic ai prompt ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid toeic ai prompt ID",
		})
		return
	}

	var req service.UpdateToeicAiPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update toeic ai prompt request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	entity, err := h.toeicAiPromptService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update toeic ai prompt", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update toeic ai prompt",
		})
		return
	}

	c.JSON(http.StatusOK, entity)
}

// Delete 删除ToeicAiPrompt
// @Summary 删除ToeicAiPrompt
// @Description 删除TOEIC AI提示词
// @Tags toeic-ai-prompts
// @Accept json
// @Produce json
// @Param id path int true "ToeicAiPrompt ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/toeic-ai-prompts/{id} [delete]
func (h *ToeicAiPromptHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid toeic ai prompt ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid toeic ai prompt ID",
		})
		return
	}

	if err := h.toeicAiPromptService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete toeic ai prompt", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete toeic ai prompt",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "ToeicAiPrompt deleted successfully",
	})
}

// List 获取ToeicAiPrompt列表
// @Summary 获取ToeicAiPrompt列表
// @Description 获取TOEIC AI提示词列表，支持分页和搜索
// @Tags toeic-ai-prompts
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param sort query string false "排序字段"
// @Param order query string false "排序方向" Enums(asc, desc)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/toeic-ai-prompts [get]
func (h *ToeicAiPromptHandler) List(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sort := c.Query("sort")
	order := c.Query("order")
	search := c.Query("search")

	opts := &service.ListToeicAiPromptOptions{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
		Search:   search,
		Filters:  make(map[string]interface{}),
	}

	entities, total, err := h.toeicAiPromptService.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to list toeic ai prompts", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list toeic ai prompts",
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  entities,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}
