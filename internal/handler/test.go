package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// TestHandler Test处理器
type TestHandler struct {
	testService service.TestService
	logger         logger.Logger
}

// NewTestHandler 创建Test处理器
func NewTestHandler(
	testService service.TestService,
	logger logger.Logger,
) *TestHandler {
	return &TestHandler{
		testService: testService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *TestHandler) RegisterRoutes(r *gin.RouterGroup) {
	tests := r.Group("/tests")
	{

		tests.POST("", h.Create)
		tests.GET("", h.List)
		tests.GET("/:id", h.GetByID)
		tests.PUT("/:id", h.Update)
		tests.DELETE("/:id", h.Delete)
	}
}

// Create 创建Test
// @Summary 创建Test
// @Description 创建新的Test
// @Tags tests
// @Accept json
// @Produce json
// @Param request body service.CreateTestRequest true "创建Test请求"
// @Success 201 {object} model.Test
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tests [post]
func (h *TestHandler) Create(c *gin.Context) {
	var req service.CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create test request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	test, err := h.testService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create test", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, test)
}

// GetByID 根据ID获取Test
// @Summary 获取Test详情
// @Description 根据ID获取Test详情
// @Tags tests
// @Accept json
// @Produce json
// @Param id path int true "Test ID"
// @Success 200 {object} model.Test
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tests/{id} [get]
func (h *TestHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Test ID",
		})
		return
	}

	test, err := h.testService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get test", "test_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "test_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, test)
}

// Update 更新Test
// @Summary 更新Test
// @Description 更新Test信息
// @Tags tests
// @Accept json
// @Produce json
// @Param id path int true "Test ID"
// @Param request body service.UpdateTestRequest true "更新Test请求"
// @Success 200 {object} model.Test
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tests/{id} [put]
func (h *TestHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Test ID",
		})
		return
	}

	var req service.UpdateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update test request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	test, err := h.testService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update test", "test_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, test)
}

// Delete 删除Test
// @Summary 删除Test
// @Description 删除Test
// @Tags tests
// @Accept json
// @Produce json
// @Param id path int true "Test ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tests/{id} [delete]
func (h *TestHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Test ID",
		})
		return
	}

	if err := h.testService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete test", "test_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Test deleted successfully",
	})
}

// List 获取Test列表
// @Summary 获取Test列表
// @Description 获取Test列表，支持分页、搜索和过滤
// @Tags tests
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tests [get]
func (h *TestHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListTestOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	tests, total, err := h.testService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get tests", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_tests_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  tests,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *TestHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreateTestRequest 创建Test请求
type CreateTestRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateTestRequest 更新Test请求
type UpdateTestRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
