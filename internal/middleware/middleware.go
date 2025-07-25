package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// Middleware 中间件管理器
type Middleware struct {
	config     *config.Config
	logger     logger.Logger
	cache      cache.Cache
	auth       *AuthMiddleware
	permission *PermissionMiddleware
	rateLimit  *RateLimitMiddleware
	logging    *LoggingMiddleware
	cors       *CORSMiddleware
	security   *SecurityMiddleware
}

// NewMiddleware 创建中间件管理器
func NewMiddleware(
	config *config.Config,
	logger logger.Logger,
	cache cache.Cache,
) *Middleware {
	return &Middleware{
		config:     config,
		logger:     logger,
		cache:      cache,
		auth:       NewAuthMiddleware(config, cache, logger),
		permission: NewPermissionMiddleware(config, cache, logger),
		rateLimit:  NewRateLimitMiddleware(config, cache, logger),
		logging:    NewLoggingMiddleware(config, logger),
		cors:       NewCORSMiddleware(config, logger),
		security:   NewSecurityMiddleware(config, logger),
	}
}

// Auth 获取认证中间件
func (m *Middleware) Auth() *AuthMiddleware {
	return m.auth
}

// Permission 获取权限中间件
func (m *Middleware) Permission() *PermissionMiddleware {
	return m.permission
}

// RateLimit 获取限流中间件
func (m *Middleware) RateLimit() *RateLimitMiddleware {
	return m.rateLimit
}

// Logging 获取日志中间件
func (m *Middleware) Logging() *LoggingMiddleware {
	return m.logging
}

// CORS 获取 CORS 中间件
func (m *Middleware) CORS() *CORSMiddleware {
	return m.cors
}

// Security 获取安全中间件
func (m *Middleware) Security() *SecurityMiddleware {
	return m.security
}

// SetupGlobalMiddleware 设置全局中间件
func (m *Middleware) SetupGlobalMiddleware(engine *gin.Engine) {
	// 恢复中间件（必须在最前面）
	engine.Use(gin.Recovery())

	// 请求 ID 和日志中间件
	engine.Use(m.logging.StructuredLogging())

	// 安全中间件
	engine.Use(m.security.SecurityHeaders())

	// CORS 中间件
	if m.config.Server.Mode == "debug" {
		engine.Use(m.cors.DevCORS())
	} else {
		engine.Use(m.cors.CORS())
	}

	// 全局限流
	engine.Use(m.rateLimit.IPRateLimit(100, 200)) // 每分钟 100 次请求

	// 安全检查
	engine.Use(m.security.RequestSizeLimit(10 * 1024 * 1024)) // 10MB 限制
	engine.Use(m.logging.SecurityLogging())
}

// SetupAPIMiddleware 设置 API 中间件
func (m *Middleware) SetupAPIMiddleware() []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{
		// API 专用限流
		m.rateLimit.UserRateLimit(60, 120), // 每分钟 60 次请求
		
		// 错误日志
		m.logging.ErrorLogging(),
	}

	return middlewares
}

// SetupAuthMiddleware 设置认证中间件组
func (m *Middleware) SetupAuthMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.RequireAuth(),
	}
}

// SetupAdminMiddleware 设置管理员中间件组
func (m *Middleware) SetupAdminMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.RequireAuth(),
		m.auth.RequireRole("admin"),
		m.rateLimit.AdminRateLimit(),
	}
}

// SetupPublicMiddleware 设置公共接口中间件组
func (m *Middleware) SetupPublicMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.OptionalAuth(),
		m.rateLimit.IPRateLimit(30, 60), // 更严格的限流
	}
}

// 便捷方法

// RequireAuth 需要认证
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return m.auth.RequireAuth()
}

// RequireRole 需要角色
func (m *Middleware) RequireRole(roles ...string) gin.HandlerFunc {
	return m.auth.RequireRole(roles...)
}

// RequirePermission 需要权限
func (m *Middleware) RequirePermission(permission Permission) gin.HandlerFunc {
	return m.permission.RequirePermissions(permission)
}

// RequireOwnership 需要所有权
func (m *Middleware) RequireOwnership(resourceType string) gin.HandlerFunc {
	return m.permission.RequireOwnership(resourceType)
}

// IPRateLimit IP 限流
func (m *Middleware) IPRateLimit(rate, burst int) gin.HandlerFunc {
	return m.rateLimit.IPRateLimit(rate, burst)
}

// UserRateLimit 用户限流
func (m *Middleware) UserRateLimit(rate, burst int) gin.HandlerFunc {
	return m.rateLimit.UserRateLimit(rate, burst)
}

// LoginRateLimit 登录限流
func (m *Middleware) LoginRateLimit() gin.HandlerFunc {
	return m.rateLimit.LoginRateLimit()
}

// RegisterRateLimit 注册限流
func (m *Middleware) RegisterRateLimit() gin.HandlerFunc {
	return m.rateLimit.RegisterRateLimit()
}

// UploadRateLimit 上传限流
func (m *Middleware) UploadRateLimit() gin.HandlerFunc {
	return m.rateLimit.UploadRateLimit()
}

// NoCache 禁用缓存
func (m *Middleware) NoCache() gin.HandlerFunc {
	return m.security.NoCache()
}

// RequestTimeout 请求超时
func (m *Middleware) RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return m.security.Timeout(timeout)
}

// IPWhitelist IP 白名单
func (m *Middleware) IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return m.security.IPWhitelist(allowedIPs)
}

// BasicAuth 基础认证
func (m *Middleware) BasicAuth(username, password string) gin.HandlerFunc {
	return m.security.BasicAuth(username, password)
}

// 预定义的中间件组合

// PublicAPI 公共 API 中间件组合
func (m *Middleware) PublicAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.OptionalAuth(),
		m.rateLimit.IPRateLimit(30, 60),
	}
}

// ProtectedAPI 受保护的 API 中间件组合
func (m *Middleware) ProtectedAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.RequireAuth(),
		m.rateLimit.UserRateLimit(60, 120),
	}
}

// AdminAPI 管理员 API 中间件组合
func (m *Middleware) AdminAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.RequireAuth(),
		m.auth.RequireRole("admin"),
		m.rateLimit.AdminRateLimit(),
	}
}

// FileUploadAPI 文件上传 API 中间件组合
func (m *Middleware) FileUploadAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.auth.RequireAuth(),
		m.rateLimit.UploadRateLimit(),
		m.security.RequestSizeLimit(50 * 1024 * 1024), // 50MB
	}
}

// AuthAPI 认证相关 API 中间件组合
func (m *Middleware) AuthAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.rateLimit.LoginRateLimit(),
		m.security.NoCache(),
	}
}

// HealthCheckAPI 健康检查 API 中间件组合
func (m *Middleware) HealthCheckAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.rateLimit.IPRateLimit(10, 20),
		m.security.NoCache(),
	}
}

// WebhookAPI Webhook API 中间件组合
func (m *Middleware) WebhookAPI(allowedIPs []string) []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{
		m.rateLimit.IPRateLimit(100, 200),
		m.security.NoCache(),
	}
	
	if len(allowedIPs) > 0 {
		middlewares = append(middlewares, m.security.IPWhitelist(allowedIPs))
	}
	
	return middlewares
}

// DevelopmentAPI 开发环境 API 中间件组合
func (m *Middleware) DevelopmentAPI() []gin.HandlerFunc {
	if m.config.Server.Mode != "debug" {
		return m.ProtectedAPI()
	}
	
	return []gin.HandlerFunc{
		m.auth.OptionalAuth(),
		m.rateLimit.IPRateLimit(1000, 2000), // 开发环境更宽松的限流
	}
}

// 工具方法

// GenerateToken 生成 JWT token
func (m *Middleware) GenerateToken(user interface{}) (string, error) {
	// 这里需要类型断言，实际使用时需要传入正确的用户类型
	return "", nil // 占位符实现
}

// ValidateToken 验证 JWT token
func (m *Middleware) ValidateToken(token string) (interface{}, error) {
	// 占位符实现
	return nil, nil
}

// RevokeToken 撤销 JWT token
func (m *Middleware) RevokeToken(token string) error {
	return m.auth.RevokeToken(token)
}

// RefreshToken 刷新 JWT token
func (m *Middleware) RefreshToken(token string) (string, error) {
	return m.auth.RefreshToken(token)
}

// ClearUserPermissions 清除用户权限缓存
func (m *Middleware) ClearUserPermissions(userID uint) error {
	return m.permission.ClearUserPermissions(userID)
}

// GetRateLimitStatus 获取限流状态
func (m *Middleware) GetRateLimitStatus(key string) (int, int, time.Time, error) {
	// 使用默认配置
	config := RateLimitConfig{
		Rate:   60,
		Burst:  120,
		Window: time.Minute,
	}
	return m.rateLimit.GetRateLimitStatus(key, config)
}

// ClearRateLimit 清除限流记录
func (m *Middleware) ClearRateLimit(key string) error {
	return m.rateLimit.ClearRateLimit(key)
}
