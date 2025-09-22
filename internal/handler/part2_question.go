package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// Part2QuestionHandler Part2Question处理器
type Part2QuestionHandler struct {
	part2QuestionService service.Part2QuestionService
	logger         logger.Logger
}

// NewPart2QuestionHandler 创建Part2Question处理器
func NewPart2QuestionHandler(
	part2QuestionService service.Part2QuestionService,
	logger logger.Logger,
) *Part2QuestionHandler {
	return &Part2QuestionHandler{
		part2QuestionService: part2QuestionService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *Part2QuestionHandler) RegisterRoutes(r *gin.RouterGroup) {
	part2questions := r.Group("/part2questions")
	{

		part2questions.POST("", h.Create)
		part2questions.GET("", h.List)
		part2questions.GET("/:id", h.GetByID)
		part2questions.PUT("/:id", h.Update)
		part2questions.DELETE("/:id", h.Delete)
	}
}

// Create 创建Part2Question
// @Summary 创建Part2Question
// @Description 创建新的Part2Question
// @Tags part2questions
// @Accept json
// @Produce json
// @Param request body service.CreatePart2QuestionRequest true "创建Part2Question请求"
// @Success 201 {object} model.Part2Question
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part2questions [post]
func (h *Part2QuestionHandler) Create(c *gin.Context) {
	var req service.CreatePart2QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create part2question request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part2Question, err := h.part2QuestionService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create part2question", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, part2Question)
}

// GetByID 根据ID获取Part2Question
// @Summary 获取Part2Question详情
// @Description 根据ID获取Part2Question详情
// @Tags part2questions
// @Accept json
// @Produce json
// @Param id path int true "Part2Question ID"
// @Success 200 {object} model.Part2Question
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part2questions/{id} [get]
func (h *Part2QuestionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part2Question ID",
		})
		return
	}

	part2Question, err := h.part2QuestionService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get part2question", "part2question_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "part2question_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part2Question)
}

// Update 更新Part2Question
// @Summary 更新Part2Question
// @Description 更新Part2Question信息
// @Tags part2questions
// @Accept json
// @Produce json
// @Param id path int true "Part2Question ID"
// @Param request body service.UpdatePart2QuestionRequest true "更新Part2Question请求"
// @Success 200 {object} model.Part2Question
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part2questions/{id} [put]
func (h *Part2QuestionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part2Question ID",
		})
		return
	}

	var req service.UpdatePart2QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update part2question request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	part2Question, err := h.part2QuestionService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update part2question", "part2question_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, part2Question)
}

// Delete 删除Part2Question
// @Summary 删除Part2Question
// @Description 删除Part2Question
// @Tags part2questions
// @Accept json
// @Produce json
// @Param id path int true "Part2Question ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part2questions/{id} [delete]
func (h *Part2QuestionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Part2Question ID",
		})
		return
	}

	if err := h.part2QuestionService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete part2question", "part2question_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Part2Question deleted successfully",
	})
}

// List 获取Part2Question列表
// @Summary 获取Part2Question列表
// @Description 获取Part2Question列表，支持分页、搜索和过滤
// @Tags part2questions
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/part2questions [get]
func (h *Part2QuestionHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListPart2QuestionOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	part2questions, total, err := h.part2QuestionService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get part2questions", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_part2questions_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  part2questions,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *Part2QuestionHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreatePart2QuestionRequest 创建Part2Question请求
type CreatePart2QuestionRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdatePart2QuestionRequest 更新Part2Question请求
type UpdatePart2QuestionRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
