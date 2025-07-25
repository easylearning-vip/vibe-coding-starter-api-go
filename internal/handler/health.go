package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db     database.Database
	cache  cache.Cache
	logger logger.Logger
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(
	db database.Database,
	cache cache.Cache,
	logger logger.Logger,
) *HealthHandler {
	return &HealthHandler{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Services  map[string]ServiceInfo `json:"services"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Health 健康检查
// @Summary 健康检查
// @Description 检查应用程序和依赖服务的健康状态
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0", // 可以从构建时注入
		Services:  make(map[string]ServiceInfo),
	}

	// 检查数据库
	dbStatus := h.checkDatabase()
	response.Services["database"] = dbStatus
	if dbStatus.Status != "healthy" {
		response.Status = "unhealthy"
	}

	// 检查缓存
	cacheStatus := h.checkCache()
	response.Services["cache"] = cacheStatus
	if cacheStatus.Status != "healthy" {
		response.Status = "unhealthy"
	}

	// 根据整体状态设置 HTTP 状态码
	statusCode := http.StatusOK
	if response.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	h.logger.Debug("Health check completed", "status", response.Status)
	c.JSON(statusCode, response)
}

// Readiness 就绪检查
// @Summary 就绪检查
// @Description 检查应用程序是否准备好接收流量
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// 检查关键依赖是否就绪
	ready := true
	checks := make(map[string]interface{})

	// 检查数据库连接
	if err := h.db.Health(); err != nil {
		ready = false
		checks["database"] = map[string]interface{}{
			"status": "not_ready",
			"error":  err.Error(),
		}
	} else {
		checks["database"] = map[string]interface{}{
			"status": "ready",
		}
	}

	// 检查缓存连接
	if err := h.cache.Health(); err != nil {
		ready = false
		checks["cache"] = map[string]interface{}{
			"status": "not_ready",
			"error":  err.Error(),
		}
	} else {
		checks["cache"] = map[string]interface{}{
			"status": "ready",
		}
	}

	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now(),
		"checks":    checks,
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Liveness 存活检查
// @Summary 存活检查
// @Description 检查应用程序是否存活
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	// 简单的存活检查，只要应用程序能响应就认为是存活的
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes 注册路由
func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/ready", h.Readiness)
	r.GET("/live", h.Liveness)
}

// checkDatabase 检查数据库状态
func (h *HealthHandler) checkDatabase() ServiceInfo {
	if err := h.db.Health(); err != nil {
		h.logger.Error("Database health check failed", "error", err)
		return ServiceInfo{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}
	return ServiceInfo{Status: "healthy"}
}

// checkCache 检查缓存状态
func (h *HealthHandler) checkCache() ServiceInfo {
	if err := h.cache.Health(); err != nil {
		h.logger.Error("Cache health check failed", "error", err)
		return ServiceInfo{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}
	return ServiceInfo{Status: "healthy"}
}

// startTime 应用启动时间
var startTime = time.Now()

// Ready 就绪检查别名
func (h *HealthHandler) Ready(c *gin.Context) {
	h.Readiness(c)
}

// Live 存活检查别名
func (h *HealthHandler) Live(c *gin.Context) {
	h.Liveness(c)
}
