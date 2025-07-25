package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService service.ArticleService
	logger         logger.Logger
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(
	articleService service.ArticleService,
	logger logger.Logger,
) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		logger:         logger,
	}
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateArticleRequest true "创建文章请求"
// @Success 201 {object} model.Article
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid create article request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// 设置作者ID
	req.AuthorID = userID.(uint)

	article, err := h.articleService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create article", "title", req.Title, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetByID 根据 ID 获取文章
// @Summary 获取文章详情
// @Description 根据 ID 获取文章详情
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.Article
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/{id} [get]
func (h *ArticleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid article ID",
		})
		return
	}

	article, err := h.articleService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get article", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "article_not_found",
			Message: err.Error(),
		})
		return
	}

	// 增加浏览次数
	if err := h.articleService.IncrementViewCount(c.Request.Context(), uint(id)); err != nil {
		h.logger.Warn("Failed to increment view count", "id", id, "error", err)
	}

	c.JSON(http.StatusOK, article)
}

// List 获取文章列表（公共接口，不需要认证）
// @Summary 获取文章列表
// @Description 获取文章列表
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Param status query string false "文章状态" Enums(draft, published, archived)
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)

	articles, total, err := h.articleService.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to get articles", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_articles_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  articles,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// ListUserArticles 获取当前用户的文章列表（需要认证）
// @Summary 获取当前用户的文章列表
// @Description 获取当前登录用户的文章列表
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Param status query string false "文章状态" Enums(draft, published, archived)
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles [get]
func (h *ArticleHandler) ListUserArticles(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	opts := h.parseListOptions(c)
	// 强制设置作者ID为当前用户
	opts.Filters["author_id"] = userID.(uint)

	articles, total, err := h.articleService.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to get user articles list", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  articles,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// ListAllArticles 获取所有文章列表（管理员专用）
// @Summary 获取所有文章列表
// @Description 获取所有文章列表（管理员权限）
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Param status query string false "文章状态" Enums(draft, published, archived)
// @Param author_id query int false "作者ID"
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/articles [get]
func (h *ArticleHandler) ListAllArticles(c *gin.Context) {
	opts := h.parseListOptions(c)
	// 管理员可以查看所有文章，不强制过滤作者

	articles, total, err := h.articleService.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to get all articles list", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  articles,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body service.UpdateArticleRequest true "更新文章请求"
// @Success 200 {object} model.Article
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/{id} [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid article ID",
		})
		return
	}

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update article request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	article, err := h.articleService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update article", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除指定文章
// @Tags articles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid article ID",
		})
		return
	}

	if err := h.articleService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete article", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Article deleted successfully",
	})
}

// Search 搜索文章
// @Summary 搜索文章
// @Description 根据关键词搜索文章
// @Tags articles
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/search [get]
func (h *ArticleHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_query",
			Message: "Search query is required",
		})
		return
	}

	opts := h.parseListOptions(c)
	articles, total, err := h.articleService.Search(c.Request.Context(), query, opts)
	if err != nil {
		h.logger.Error("Failed to search articles", "query", query, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "search_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  articles,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// RegisterRoutes 注册路由
func (h *ArticleHandler) RegisterRoutes(r *gin.RouterGroup) {
	articles := r.Group("/articles")
	{
		// 公共路由（不需要认证）
		articles.GET("", h.List)
		articles.GET("/search", h.Search)
		articles.GET("/:id", h.GetByID)

		// 需要认证的路由（在服务器层面已经处理认证）
		articles.POST("", h.Create)
		articles.PUT("/:id", h.Update)
		articles.DELETE("/:id", h.Delete)
	}
}

// 辅助方法
func (h *ArticleHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	// 构建过滤器
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if authorID := c.Query("author_id"); authorID != "" {
		if id, err := strconv.Atoi(authorID); err == nil {
			filters["author_id"] = id
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			filters["category_id"] = id
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filters["created_after"] = startDate + " 00:00:00"
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["created_before"] = endDate + " 23:59:59"
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
