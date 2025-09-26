package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "vibe-coding-starter/internal/service"
    "vibe-coding-starter/pkg/logger"
)

// Part3PracticeItemHandler handles user Part3 practice records
type Part3PracticeItemHandler struct {
    svc    service.Part3PracticeItemService
    logger logger.Logger
}

func NewPart3PracticeItemHandler(svc service.Part3PracticeItemService, logger logger.Logger) *Part3PracticeItemHandler {
    return &Part3PracticeItemHandler{svc: svc, logger: logger}
}

// RegisterRoutes registers protected user routes
func (h *Part3PracticeItemHandler) RegisterRoutes(r *gin.RouterGroup) {
    grp := r.Group("/user/part3-practice-items")
    {
        grp.POST("", h.Create)
        grp.GET("", h.List)
        grp.GET("/stats", h.Stats)
        grp.GET("/:id", h.GetByID)
        grp.PUT("/:id", h.Update)
        grp.DELETE("/:id", h.Delete)
    }
}

func (h *Part3PracticeItemHandler) getUserID(c *gin.Context) uint {
    if v, ok := c.Get("user_id"); ok {
        if id, ok := v.(uint); ok { return id }
    }
    return 0
}

func (h *Part3PracticeItemHandler) Create(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    var req service.CreatePart3PracticeItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()}); return
    }

    item, err := h.svc.Create(c.Request.Context(), userID, &req)
    if err != nil { c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "create_failed", Message: err.Error()}); return }
    c.JSON(http.StatusCreated, item)
}

func (h *Part3PracticeItemHandler) GetByID(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"}); return }

    item, err := h.svc.GetByID(c.Request.Context(), userID, uint(id))
    if err != nil { c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()}); return }
    c.JSON(http.StatusOK, item)
}

func (h *Part3PracticeItemHandler) Update(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"}); return }

    var req service.UpdatePart3PracticeItemRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()}); return }

    item, err := h.svc.Update(c.Request.Context(), userID, uint(id), &req)
    if err != nil { c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "update_failed", Message: err.Error()}); return }
    c.JSON(http.StatusOK, item)
}

func (h *Part3PracticeItemHandler) Delete(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil { c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"}); return }

    if err := h.svc.Delete(c.Request.Context(), userID, uint(id)); err != nil { c.JSON(http.StatusNotFound, ErrorResponse{Error: "delete_failed", Message: err.Error()}); return }
    c.JSON(http.StatusOK, SuccessResponse{Message: "Deleted"})
}

func (h *Part3PracticeItemHandler) List(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    filters := map[string]interface{}{}
    if v := c.Query("scenario_id"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["scenario_id"] = x } }
    if v := c.Query("difficulty_level_id"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["difficulty_level_id"] = x } }
    if v := c.Query("conversation_id"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["conversation_id"] = uint(x) } }
    if v := c.Query("question_index"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["question_index"] = x } }
    if v := c.Query("is_correct"); v != "" { filters["is_correct"] = (v == "true" || v == "1") }

    items, total, err := h.svc.List(c.Request.Context(), userID, &service.ListPart3PracticeItemOptions{
        Page: page, PageSize: pageSize, Sort: c.DefaultQuery("sort", "created_at"), Order: c.DefaultQuery("order", "desc"), Filters: filters,
    })
    if err != nil { c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "list_failed", Message: err.Error()}); return }

    c.JSON(http.StatusOK, ListResponse{Data: items, Total: total, Page: page, Size: pageSize})
}

func (h *Part3PracticeItemHandler) Stats(c *gin.Context) {
    userID := h.getUserID(c)
    if userID == 0 { c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"}); return }

    filters := map[string]interface{}{}
    if v := c.Query("scenario_id"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["scenario_id"] = x } }
    if v := c.Query("difficulty_level_id"); v != "" { if x, err := strconv.Atoi(v); err == nil { filters["difficulty_level_id"] = x } }

    stats, err := h.svc.Stats(c.Request.Context(), userID, filters)
    if err != nil { c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "stats_failed", Message: err.Error()}); return }
    c.JSON(http.StatusOK, stats)
}

