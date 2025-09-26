package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

type Part3PracticeSetHandler struct {
	svc    service.Part3PracticeSetService
	logger logger.Logger
}

func NewPart3PracticeSetHandler(svc service.Part3PracticeSetService, logger logger.Logger) *Part3PracticeSetHandler {
	return &Part3PracticeSetHandler{svc: svc, logger: logger}
}

func (h *Part3PracticeSetHandler) RegisterRoutes(r *gin.RouterGroup) {
	s := r.Group("/user/part3-practice-sets")
	{
		s.POST("", h.CreateSet)
		s.GET("", h.ListSets)
		s.GET("/:id", h.GetSetByID)
		s.PUT("/:id", h.UpdateSet)
		s.DELETE("/:id", h.DeleteSet)
	}
	it := r.Group("/user/part3-practice-set-items")
	{
		it.POST("", h.CreateItem)
		it.GET("", h.ListItems)
		it.DELETE("/:id", h.DeleteItem)
	}
}

func (h *Part3PracticeSetHandler) getUserID(c *gin.Context) uint {
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

func (h *Part3PracticeSetHandler) CreateSet(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	var req service.CreatePracticeSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	e, err := h.svc.CreateSet(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "create_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}
func (h *Part3PracticeSetHandler) GetSetByID(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	e, err := h.svc.GetSetByID(c.Request.Context(), userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}
func (h *Part3PracticeSetHandler) UpdateSet(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	var req service.UpdatePracticeSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	e, err := h.svc.UpdateSet(c.Request.Context(), userID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "update_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}
func (h *Part3PracticeSetHandler) DeleteSet(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	if err := h.svc.DeleteSet(c.Request.Context(), userID, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "delete_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "Deleted"})
}
func (h *Part3PracticeSetHandler) ListSets(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	filters := map[string]interface{}{}
	if v := c.Query("scenario_id"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			filters["scenario_id"] = x
		}
	}
	if v := c.Query("difficulty_level_id"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			filters["difficulty_level_id"] = x
		}
	}
	items, total, err := h.svc.ListSets(c.Request.Context(), userID, &service.ListPracticeSetOptions{Page: page, PageSize: pageSize, Sort: c.DefaultQuery("sort", "created_at"), Order: c.DefaultQuery("order", "desc"), Filters: filters})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ListResponse{Data: items, Total: total, Page: page, Size: pageSize})
}

func (h *Part3PracticeSetHandler) CreateItem(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	var req service.CreatePart3SetItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	e, err := h.svc.CreateItem(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "create_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}
func (h *Part3PracticeSetHandler) ListItems(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	setIDStr := c.Query("set_id")
	if setIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: "set_id is required"})
		return
	}
	setID64, err := strconv.ParseUint(setIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_request", Message: "invalid set_id"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.ListItems(c.Request.Context(), userID, uint(setID64), &service.ListPracticeSetItemOptions{Page: page, PageSize: pageSize, Sort: c.DefaultQuery("sort", "order_index"), Order: c.DefaultQuery("order", "asc"), Filters: map[string]interface{}{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ListResponse{Data: items, Total: total, Page: page, Size: pageSize})
}
func (h *Part3PracticeSetHandler) DeleteItem(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	if err := h.svc.DeleteItem(c.Request.Context(), userID, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "delete_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "Deleted"})
}
