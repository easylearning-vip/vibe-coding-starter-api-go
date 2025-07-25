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

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
	logger      logger.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(
	userService service.UserService,
	logger logger.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账户
// @Tags users
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册请求"
// @Success 201 {object} model.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid register request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", "username", req.Username, "email", req.Email, "error", err)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user.ToPublic())
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags users
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	h.logger.Debug("Login handler called")

	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	h.logger.Debug("Login request parsed", "username", req.Username)

	response, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to login user", "username", req.Username, "error", err)
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "login_failed",
			Message: err.Error(),
		})
		return
	}

	h.logger.Debug("Login successful", "username", req.Username)

	c.JSON(http.StatusOK, response)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前用户的资料信息
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", "user_id", userID, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user.ToPublic())
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前用户的资料信息
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateProfileRequest true "更新资料请求"
// @Success 200 {object} model.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update profile request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to update user profile", "user_id", userID, "error", err)
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user.ToPublic())
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid change password request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		h.logger.Error("Failed to change password", "user_id", userID, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "change_password_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password changed successfully",
	})
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表（需要管理员权限）
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// 检查管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Admin access required",
		})
		return
	}

	opts := h.parseListOptions(c)
	users, total, err := h.userService.GetUsers(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to get users", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_users_failed",
			Message: err.Error(),
		})
		return
	}

	// 转换为公开信息
	publicUsers := make([]*model.PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  publicUsers,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户（需要管理员权限）
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 检查管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
		})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete user", "user_id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User deleted successfully",
	})
}

// RegisterRoutes 注册需要认证的路由
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		// 需要认证的路由
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.POST("/change-password", h.ChangePassword)

		// 管理员路由
		users.GET("", h.GetUsers)
		users.DELETE("/:id", h.DeleteUser)
	}
}

// 辅助方法

func (h *UserHandler) getUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

func (h *UserHandler) isAdmin(c *gin.Context) bool {
	if role, exists := c.Get("user_role"); exists {
		return role == "admin"
	}
	return false
}

func (h *UserHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	// 构建过滤器
	filters := make(map[string]interface{})
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
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

// 响应结构体

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ListResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// AuthMiddleware 认证中间件占位符
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现 JWT 认证逻辑
		c.Next()
	}
}
