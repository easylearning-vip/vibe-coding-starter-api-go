package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// Part4TalkHandler Part4Talk处理器
type Part4TalkHandler struct {
	part4TalkService service.Part4TalkService
	logger           logger.Logger
}

// NewPart4TalkHandler 创建Part4Talk处理器
func NewPart4TalkHandler(
	part4TalkService service.Part4TalkService,
	logger logger.Logger,
) *Part4TalkHandler {
	return &Part4TalkHandler{
		part4TalkService: part4TalkService,
		logger:           logger,
	}
}

// RegisterRoutes 注册路由
func (h *Part4TalkHandler) RegisterRoutes(r *gin.RouterGroup) {
	part4talks := r.Group("/part4talks")
	{

		part4talks.POST("", h.Create)
		part4talks.GET("", h.List)
		part4talks.GET("/:id", h.GetByID)
		part4talks.PUT("/:id", h.Update)
		part4talks.DELETE("/:id", h.Delete)
	}
}

// Create 创建Part4Talk
// @Summary 创建Part4Talk
// @Description 创建新的Part4Talk
// @Tags part4talks
// @Accept json
// @Produce json
// @Param request body service.CreatePart4TalkRequest true "创建Part4Talk请求"
// @Success 201 {object} model.Part4Talk
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4talks [post]
func (h *Part4TalkHandler) Create(c *gin.Context) {
	var req service.CreatePart4TalkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create part4talk request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part4Talk, err := h.part4TalkService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create part4talk", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, part4Talk)
}

// GetByID 根据ID获取Part4Talk
// @Summary 获取Part4Talk详情
// @Description 根据ID获取Part4Talk详情
// @Tags part4talks
// @Accept json
// @Produce json
// @Param id path int true "Part4Talk ID"
// @Success 200 {object} model.Part4Talk
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4talks/{id} [get]
func (h *Part4TalkHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4Talk ID",
		})
		return
	}

	part4Talk, err := h.part4TalkService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get part4talk", "part4talk_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "part4talk_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part4Talk)
}

// Update 更新Part4Talk
// @Summary 更新Part4Talk
// @Description 更新Part4Talk信息
// @Tags part4talks
// @Accept json
// @Produce json
// @Param id path int true "Part4Talk ID"
// @Param request body service.UpdatePart4TalkRequest true "更新Part4Talk请求"
// @Success 200 {object} model.Part4Talk
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4talks/{id} [put]
func (h *Part4TalkHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4Talk ID",
		})
		return
	}

	var req service.UpdatePart4TalkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update part4talk request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part4Talk, err := h.part4TalkService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update part4talk", "part4talk_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part4Talk)
}

// Delete 删除Part4Talk
// @Summary 删除Part4Talk
// @Description 删除Part4Talk
// @Tags part4talks
// @Accept json
// @Produce json
// @Param id path int true "Part4Talk ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4talks/{id} [delete]
func (h *Part4TalkHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part4Talk ID",
		})
		return
	}

	if err := h.part4TalkService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete part4talk", "part4talk_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Part4Talk deleted successfully",
	})
}

// List 获取Part4Talk列表
// @Summary 获取Part4Talk列表
// @Description 获取Part4Talk列表，支持分页、搜索和过滤。如果包含include_answer_options=true参数，则预加载答案选项
// @Tags part4talks
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Param include_answer_options query bool false "是否包含答案选项" default(false)
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part4talks [get]
func (h *Part4TalkHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListPart4TalkOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	// 检查是否需要预加载答案选项
	includeAnswerOptions := c.Query("include_answer_options") == "true"

	var part4talks []*model.Part4Talk
	var total int64
	var err error

	if includeAnswerOptions {
		part4talks, total, err = h.part4TalkService.ListWithAnswerOptions(c.Request.Context(), serviceOpts)
	} else {
		part4talks, total, err = h.part4TalkService.List(c.Request.Context(), serviceOpts)
	}

	if err != nil {
		h.logger.Error("Failed to get part4talks", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_part4talks_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  part4talks,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *Part4TalkHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	// 构建过滤器
	filters := make(map[string]interface{})

	// 添加 test_id 过滤器
	if testID := c.Query("test_id"); testID != "" {
		if id, err := strconv.Atoi(testID); err == nil {
			filters["test_id"] = id
		}
	}

	// 添加 scenario_id 过滤器
	if scenarioID := c.Query("scenario_id"); scenarioID != "" {
		if id, err := strconv.Atoi(scenarioID); err == nil {
			filters["scenario_id"] = id
		}
	}

	// 添加 difficulty_level_id 过滤器
	if difficultyLevelID := c.Query("difficulty_level_id"); difficultyLevelID != "" {
		if id, err := strconv.Atoi(difficultyLevelID); err == nil {
			filters["difficulty_level_id"] = id
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

// 请求结构体

// CreatePart4TalkRequest 创建Part4Talk请求
type CreatePart4TalkRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdatePart4TalkRequest 更新Part4Talk请求
type UpdatePart4TalkRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
