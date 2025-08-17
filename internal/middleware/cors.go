package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/logger"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string      `json:"allow_origins"`     // 允许的源
	AllowMethods     []string      `json:"allow_methods"`     // 允许的方法
	AllowHeaders     []string      `json:"allow_headers"`     // 允许的请求头
	ExposeHeaders    []string      `json:"expose_headers"`    // 暴露的响应头
	AllowCredentials bool          `json:"allow_credentials"` // 是否允许凭证
	MaxAge           time.Duration `json:"max_age"`           // 预检请求缓存时间
}

// CORSMiddleware CORS 中间件
type CORSMiddleware struct {
	config     *config.Config
	logger     logger.Logger
	corsConfig CORSConfig
}

// NewCORSMiddleware 创建 CORS 中间件
func NewCORSMiddleware(
	config *config.Config,
	logger logger.Logger,
) *CORSMiddleware {
	// 默认 CORS 配置
	corsConfig := CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
			"https://yourdomain.com",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-API-Key",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// 从配置文件覆盖默认配置
	if config.CORS.AllowOrigins != nil {
		corsConfig.AllowOrigins = config.CORS.AllowOrigins
	}
	if config.CORS.AllowMethods != nil {
		corsConfig.AllowMethods = config.CORS.AllowMethods
	}
	if config.CORS.AllowHeaders != nil {
		corsConfig.AllowHeaders = config.CORS.AllowHeaders
	}
	if config.CORS.ExposeHeaders != nil {
		corsConfig.ExposeHeaders = config.CORS.ExposeHeaders
	}
	if config.CORS.MaxAge > 0 {
		corsConfig.MaxAge = time.Duration(config.CORS.MaxAge) * time.Second
	}

	return &CORSMiddleware{
		config:     config,
		logger:     logger,
		corsConfig: corsConfig,
	}
}

// CORS 标准 CORS 中间件
func (m *CORSMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该源
		if m.isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if m.allowsAllOrigins() {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// 设置允许的方法
		if len(m.corsConfig.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(m.corsConfig.AllowMethods, ", "))
		}

		// 设置允许的请求头
		if len(m.corsConfig.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(m.corsConfig.AllowHeaders, ", "))
		}

		// 设置暴露的响应头
		if len(m.corsConfig.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(m.corsConfig.ExposeHeaders, ", "))
		}

		// 设置是否允许凭证
		if m.corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			m.handlePreflightRequest(c)
			return
		}

		c.Next()
	}
}

// DynamicCORS 动态 CORS 中间件（根据请求动态配置）
func (m *CORSMiddleware) DynamicCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 根据不同的 API 路径设置不同的 CORS 策略
		corsConfig := m.getCORSConfigForPath(c.Request.URL.Path)

		// 检查是否允许该源
		if m.isOriginAllowedForConfig(origin, corsConfig) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置其他 CORS 头
		m.setCORSHeaders(c, corsConfig)

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			m.handlePreflightRequestWithConfig(c, corsConfig)
			return
		}

		c.Next()
	}
}

// StrictCORS 严格的 CORS 中间件（用于生产环境）
func (m *CORSMiddleware) StrictCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 严格检查源
		if !m.isOriginAllowed(origin) {
			m.logger.Warn("CORS: Origin not allowed",
				"origin", origin,
				"path", c.Request.URL.Path,
				"method", c.Request.Method)

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "cors_forbidden",
				"message": "Origin not allowed",
			})
			c.Abort()
			return
		}

		// 设置 CORS 头
		c.Header("Access-Control-Allow-Origin", origin)
		m.setCORSHeaders(c, m.corsConfig)

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			m.handlePreflightRequest(c)
			return
		}

		c.Next()
	}
}

// DevCORS 开发环境 CORS 中间件（宽松配置）
func (m *CORSMiddleware) DevCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开发环境允许所有源
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 检查源是否被允许
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range m.corsConfig.AllowOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// 支持通配符匹配
		if m.matchesWildcard(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

// isOriginAllowedForConfig 检查源是否被特定配置允许
func (m *CORSMiddleware) isOriginAllowedForConfig(origin string, config CORSConfig) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range config.AllowOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		if m.matchesWildcard(allowedOrigin, origin) {
			return true
		}
	}

	return false
}

// allowsAllOrigins 检查是否允许所有源
func (m *CORSMiddleware) allowsAllOrigins() bool {
	for _, origin := range m.corsConfig.AllowOrigins {
		if origin == "*" {
			return true
		}
	}
	return false
}

// matchesWildcard 通配符匹配
func (m *CORSMiddleware) matchesWildcard(pattern, origin string) bool {
	// 简单的通配符匹配，支持 *.example.com 格式
	if strings.HasPrefix(pattern, "*.") {
		domain := pattern[2:]
		return strings.HasSuffix(origin, "."+domain) || origin == domain
	}
	return false
}

// handlePreflightRequest 处理预检请求
func (m *CORSMiddleware) handlePreflightRequest(c *gin.Context) {
	// 设置预检请求缓存时间
	c.Header("Access-Control-Max-Age", strconv.Itoa(int(m.corsConfig.MaxAge.Seconds())))

	// 检查请求的方法是否被允许
	requestMethod := c.Request.Header.Get("Access-Control-Request-Method")
	if requestMethod != "" && !m.isMethodAllowed(requestMethod) {
		m.logger.Warn("CORS: Method not allowed in preflight",
			"method", requestMethod,
			"origin", c.Request.Header.Get("Origin"))
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	// 检查请求的头是否被允许
	requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
	if requestHeaders != "" && !m.areHeadersAllowed(requestHeaders) {
		m.logger.Warn("CORS: Headers not allowed in preflight",
			"headers", requestHeaders,
			"origin", c.Request.Header.Get("Origin"))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// handlePreflightRequestWithConfig 使用特定配置处理预检请求
func (m *CORSMiddleware) handlePreflightRequestWithConfig(c *gin.Context, config CORSConfig) {
	c.Header("Access-Control-Max-Age", strconv.Itoa(int(config.MaxAge.Seconds())))

	requestMethod := c.Request.Header.Get("Access-Control-Request-Method")
	if requestMethod != "" && !m.isMethodAllowedForConfig(requestMethod, config) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
	if requestHeaders != "" && !m.areHeadersAllowedForConfig(requestHeaders, config) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// isMethodAllowed 检查方法是否被允许
func (m *CORSMiddleware) isMethodAllowed(method string) bool {
	for _, allowedMethod := range m.corsConfig.AllowMethods {
		if allowedMethod == method {
			return true
		}
	}
	return false
}

// isMethodAllowedForConfig 检查方法是否被特定配置允许
func (m *CORSMiddleware) isMethodAllowedForConfig(method string, config CORSConfig) bool {
	for _, allowedMethod := range config.AllowMethods {
		if allowedMethod == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed 检查请求头是否被允许
func (m *CORSMiddleware) areHeadersAllowed(headers string) bool {
	requestHeaders := strings.Split(headers, ",")
	for _, header := range requestHeaders {
		header = strings.TrimSpace(header)
		if !m.isHeaderAllowed(header) {
			return false
		}
	}
	return true
}

// areHeadersAllowedForConfig 检查请求头是否被特定配置允许
func (m *CORSMiddleware) areHeadersAllowedForConfig(headers string, config CORSConfig) bool {
	requestHeaders := strings.Split(headers, ",")
	for _, header := range requestHeaders {
		header = strings.TrimSpace(header)
		if !m.isHeaderAllowedForConfig(header, config) {
			return false
		}
	}
	return true
}

// isHeaderAllowed 检查单个请求头是否被允许
func (m *CORSMiddleware) isHeaderAllowed(header string) bool {
	header = strings.ToLower(header)
	for _, allowedHeader := range m.corsConfig.AllowHeaders {
		if strings.ToLower(allowedHeader) == header {
			return true
		}
	}
	return false
}

// isHeaderAllowedForConfig 检查单个请求头是否被特定配置允许
func (m *CORSMiddleware) isHeaderAllowedForConfig(header string, config CORSConfig) bool {
	header = strings.ToLower(header)
	for _, allowedHeader := range config.AllowHeaders {
		if strings.ToLower(allowedHeader) == header {
			return true
		}
	}
	return false
}

// setCORSHeaders 设置 CORS 响应头
func (m *CORSMiddleware) setCORSHeaders(c *gin.Context, config CORSConfig) {
	if len(config.AllowMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
	}

	if len(config.AllowHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
	}

	if len(config.ExposeHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
	}

	if config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}

// getCORSConfigForPath 根据路径获取 CORS 配置
func (m *CORSMiddleware) getCORSConfigForPath(path string) CORSConfig {
	// 根据不同的 API 路径返回不同的 CORS 配置
	switch {
	case strings.HasPrefix(path, "/api/v1/admin"):
		// 管理员 API 使用更严格的 CORS 配置
		return CORSConfig{
			AllowOrigins:     []string{"https://admin.yourdomain.com"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           time.Hour,
		}
	case strings.HasPrefix(path, "/api/v1/public"):
		// 公共 API 使用宽松的 CORS 配置
		return CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: false,
			MaxAge:           24 * time.Hour,
		}
	default:
		// 默认配置
		return m.corsConfig
	}
}
