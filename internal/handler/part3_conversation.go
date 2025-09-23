package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// Part3ConversationHandler Part3Conversation处理器
type Part3ConversationHandler struct {
	part3ConversationService service.Part3ConversationService
	logger                   logger.Logger
}

// NewPart3ConversationHandler 创建Part3Conversation处理器
func NewPart3ConversationHandler(
	part3ConversationService service.Part3ConversationService,
	logger logger.Logger,
) *Part3ConversationHandler {
	return &Part3ConversationHandler{
		part3ConversationService: part3ConversationService,
		logger:                   logger,
	}
}

// RegisterRoutes 注册路由
func (h *Part3ConversationHandler) RegisterRoutes(r *gin.RouterGroup) {
	part3conversations := r.Group("/part3conversations")
	{

		part3conversations.POST("", h.Create)
		part3conversations.GET("", h.List)
		part3conversations.GET("/:id", h.GetByID)
		part3conversations.PUT("/:id", h.Update)
		part3conversations.DELETE("/:id", h.Delete)
	}
}

// Create 创建Part3Conversation
// @Summary 创建Part3Conversation
// @Description 创建新的Part3Conversation
// @Tags part3conversations
// @Accept json
// @Produce json
// @Param request body service.CreatePart3ConversationRequest true "创建Part3Conversation请求"
// @Success 201 {object} model.Part3Conversation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/part3conversations [post]
func (h *Part3ConversationHandler) Create(c *gin.Context) {
	var req service.CreatePart3ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create part3conversation request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part3Conversation, err := h.part3ConversationService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create part3conversation", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, part3Conversation)
}

// GetByID 根据ID获取Part3Conversation
// @Summary 获取Part3Conversation详情
// @Description 根据ID获取Part3Conversation详情
// @Tags part3conversations
// @Accept json
// @Produce json
// @Param id path int true "Part3Conversation ID"
// @Success 200 {object} model.Part3Conversation
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3conversations/{id} [get]
func (h *Part3ConversationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3Conversation ID",
		})
		return
	}

	part3Conversation, err := h.part3ConversationService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get part3conversation", "part3conversation_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "part3conversation_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part3Conversation)
}

// Update 更新Part3Conversation
// @Summary 更新Part3Conversation
// @Description 更新Part3Conversation信息
// @Tags part3conversations
// @Accept json
// @Produce json
// @Param id path int true "Part3Conversation ID"
// @Param request body service.UpdatePart3ConversationRequest true "更新Part3Conversation请求"
// @Success 200 {object} model.Part3Conversation
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/part3conversations/{id} [put]
func (h *Part3ConversationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3Conversation ID",
		})
		return
	}

	var req service.UpdatePart3ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update part3conversation request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part3Conversation, err := h.part3ConversationService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update part3conversation", "part3conversation_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part3Conversation)
}

// Delete 删除Part3Conversation
// @Summary 删除Part3Conversation
// @Description 删除Part3Conversation
// @Tags part3conversations
// @Accept json
// @Produce json
// @Param id path int true "Part3Conversation ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3conversations/{id} [delete]
func (h *Part3ConversationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part3Conversation ID",
		})
		return
	}

	if err := h.part3ConversationService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete part3conversation", "part3conversation_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Part3Conversation deleted successfully",
	})
}

// List 获取Part3Conversation列表
// @Summary 获取Part3Conversation列表
// @Description 获取Part3Conversation列表，支持分页、搜索和过滤
// @Tags part3conversations
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part3conversations [get]
func (h *Part3ConversationHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListPart3ConversationOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	part3conversations, total, err := h.part3ConversationService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get part3conversations", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_part3conversations_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  part3conversations,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *Part3ConversationHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	// 构建过滤器
	filters := make(map[string]interface{})

	// 添加测试ID过滤器
	if testId := c.Query("test_id"); testId != "" {
		if id, err := strconv.ParseInt(testId, 10, 32); err == nil {
			filters["test_id"] = int32(id)
		}
	}

	// 添加场景ID过滤器
	if scenarioId := c.Query("scenario_id"); scenarioId != "" {
		if id, err := strconv.ParseInt(scenarioId, 10, 32); err == nil {
			filters["scenario_id"] = int32(id)
		}
	}

	// 添加难度级别ID过滤器
	if difficultyId := c.Query("difficulty_level_id"); difficultyId != "" {
		if id, err := strconv.ParseInt(difficultyId, 10, 32); err == nil {
			filters["difficulty_level_id"] = int32(id)
		}
	}

	// 添加对话编号过滤器
	if conversationNumber := c.Query("conversation_number"); conversationNumber != "" {
		if num, err := strconv.ParseInt(conversationNumber, 10, 32); err == nil {
			filters["conversation_number"] = int32(num)
		}
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

// 请求结构体 - 这些结构体已经在service层定义，这里不需要重复定义
