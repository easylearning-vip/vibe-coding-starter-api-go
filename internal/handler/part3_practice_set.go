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
		s.POST("/auto-generate", h.AutoGenerate)
		s.GET("", h.ListSets)
		s.GET("/:id", h.GetSetByID)
		s.GET("/:id/questions", h.ListSetQuestions)
		s.PUT("/:id", h.UpdateSet)
		s.DELETE("/:id", h.DeleteSet)
	}
	it := r.Group("/user/part3-practice-set-items")
	{
		it.POST("", h.CreateItem)
		it.GET("", h.ListItems)
		it.GET("/:id/detail", h.GetItemDetail)
		it.POST("/:id/answer", h.SubmitAnswer)
		it.DELETE("/:id", h.DeleteItem)
	}
}

// RegisterPublicRoutes 仅注册匿名访问的只读接口
func (h *Part3PracticeSetHandler) RegisterPublicRoutes(r *gin.RouterGroup) {
	p := r.Group("/public/toeic/part3")
	{
		p.GET("/sets/:id/questions", h.PublicListSetQuestions)
		p.GET("/items/:id/detail", h.PublicGetItemDetail)
	}
}

// PublicListSetQuestions 匿名：按题集ID查看所有题目信息
func (h *Part3PracticeSetHandler) PublicListSetQuestions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	res, err := h.svc.ListSetQuestions(c.Request.Context(), 0, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}
	// build slim response by mapping question text and options of current question_index
	type q struct {
		QuestionText string `json:"question_text"`
		OptionA      string `json:"option_a"`
		OptionB      string `json:"option_b"`
		OptionC      string `json:"option_c"`
		OptionD      string `json:"option_d"`
	}
	type out struct {
		ItemID   uint `json:"item_id"`
		Question q    `json:"question"`
	}
	outs := make([]out, 0, len(res))
	for _, it := range res {
		// pick question text by index
		var text string
		switch it.QuestionIndex {
		case 1:
			text = it.Conversation.Question1
		case 2:
			text = it.Conversation.Question2
		default:
			text = it.Conversation.Question3
		}
		var oa, ob, oc, od string
		for _, ao := range it.Conversation.AnswerOptions {
			if ao.QuestionNumber == it.QuestionIndex {
				oa, ob, oc, od = ao.OptionA, ao.OptionB, ao.OptionC, ao.OptionD
				break
			}
		}
		outs = append(outs, out{ItemID: it.ItemID, Question: q{QuestionText: text, OptionA: oa, OptionB: ob, OptionC: oc, OptionD: od}})
	}
	c.JSON(http.StatusOK, outs)
}

// PublicGetItemDetail 匿名：按明细ID获取单题详情（含答案选项）
func (h *Part3PracticeSetHandler) PublicGetItemDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid_id", Message: "Invalid ID"})
		return
	}
	res, err := h.svc.GetItemDetail(c.Request.Context(), 0, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}
	// slim response
	type q struct {
		QuestionText string `json:"question_text"`
		OptionA      string `json:"option_a"`
		OptionB      string `json:"option_b"`
		OptionC      string `json:"option_c"`
		OptionD      string `json:"option_d"`
	}
	// pick text and options by index
	var text string
	switch res.QuestionIndex {
	case 1:
		text = res.Conversation.Question1
	case 2:
		text = res.Conversation.Question2
	default:
		text = res.Conversation.Question3
	}
	var oa, ob, oc, od string
	for _, ao := range res.Conversation.AnswerOptions {
		if ao.QuestionNumber == res.QuestionIndex {
			oa, ob, oc, od = ao.OptionA, ao.OptionB, ao.OptionC, ao.OptionD
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{"item_id": res.ItemID, "question": q{QuestionText: text, OptionA: oa, OptionB: ob, OptionC: oc, OptionD: od}})
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

func (h *Part3PracticeSetHandler) AutoGenerate(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: "User not authenticated"})
		return
	}
	var req service.GeneratePracticeSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	set, items, err := h.svc.AutoGenerateSet(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "generate_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"set": set, "items": items})
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

func (h *Part3PracticeSetHandler) ListSetQuestions(c *gin.Context) {
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
	res, err := h.svc.ListSetQuestions(c.Request.Context(), userID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "list_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
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

func (h *Part3PracticeSetHandler) GetItemDetail(c *gin.Context) {
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
	res, err := h.svc.GetItemDetail(c.Request.Context(), userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Part3PracticeSetHandler) SubmitAnswer(c *gin.Context) {
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
	var req service.SubmitPart3AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if err := h.svc.SubmitAnswer(c.Request.Context(), userID, uint(id), &req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "submit_failed", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "submitted"})
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
