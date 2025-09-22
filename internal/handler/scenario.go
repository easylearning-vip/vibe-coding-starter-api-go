package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// ScenarioHandler Scenario处理器
type ScenarioHandler struct {
	scenarioService service.ScenarioService
	logger         logger.Logger
}

// NewScenarioHandler 创建Scenario处理器
func NewScenarioHandler(
	scenarioService service.ScenarioService,
	logger logger.Logger,
) *ScenarioHandler {
	return &ScenarioHandler{
		scenarioService: scenarioService,
		logger:         logger,
	}
}

// RegisterRoutes 注册路由
func (h *ScenarioHandler) RegisterRoutes(r *gin.RouterGroup) {
	scenarios := r.Group("/scenarios")
	{

		scenarios.POST("", h.Create)
		scenarios.GET("", h.List)
		scenarios.GET("/:id", h.GetByID)
		scenarios.PUT("/:id", h.Update)
		scenarios.DELETE("/:id", h.Delete)
	}
}

// Create 创建Scenario
// @Summary 创建Scenario
// @Description 创建新的Scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param request body service.CreateScenarioRequest true "创建Scenario请求"
// @Success 201 {object} model.Scenario
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/scenarios [post]
func (h *ScenarioHandler) Create(c *gin.Context) {
	var req service.CreateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create scenario request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	scenario, err := h.scenarioService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create scenario", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, scenario)
}

// GetByID 根据ID获取Scenario
// @Summary 获取Scenario详情
// @Description 根据ID获取Scenario详情
// @Tags scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Success 200 {object} model.Scenario
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/scenarios/{id} [get]
func (h *ScenarioHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Scenario ID",
		})
		return
	}

	scenario, err := h.scenarioService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get scenario", "scenario_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "scenario_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Update 更新Scenario
// @Summary 更新Scenario
// @Description 更新Scenario信息
// @Tags scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Param request body service.UpdateScenarioRequest true "更新Scenario请求"
// @Success 200 {object} model.Scenario
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/scenarios/{id} [put]
func (h *ScenarioHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Scenario ID",
		})
		return
	}

	var req service.UpdateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update scenario request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	scenario, err := h.scenarioService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update scenario", "scenario_id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Delete 删除Scenario
// @Summary 删除Scenario
// @Description 删除Scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param id path int true "Scenario ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/scenarios/{id} [delete]
func (h *ScenarioHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid Scenario ID",
		})
		return
	}

	if err := h.scenarioService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete scenario", "scenario_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Scenario deleted successfully",
	})
}

// List 获取Scenario列表
// @Summary 获取Scenario列表
// @Description 获取Scenario列表，支持分页、搜索和过滤
// @Tags scenarios
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/scenarios [get]
func (h *ScenarioHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	serviceOpts := &service.ListScenarioOptions{
		Page:     opts.Page,
		PageSize: opts.PageSize,
		Search:   opts.Search,
		Sort:     opts.Sort,
		Order:    opts.Order,
		Filters:  opts.Filters,
	}

	scenarios, total, err := h.scenarioService.List(c.Request.Context(), serviceOpts)
	if err != nil {
		h.logger.Error("Failed to get scenarios", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_scenarios_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  scenarios,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// 辅助方法

func (h *ScenarioHandler) parseListOptions(c *gin.Context) repository.ListOptions {
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

// CreateScenarioRequest 创建Scenario请求
type CreateScenarioRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateScenarioRequest 更新Scenario请求
type UpdateScenarioRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}
