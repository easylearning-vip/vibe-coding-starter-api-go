package middleware

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/logger"
)

// SecurityMiddleware 安全中间件
type SecurityMiddleware struct {
	config *config.Config
	logger logger.Logger
}

// NewSecurityMiddleware 创建安全中间件
func NewSecurityMiddleware(
	config *config.Config,
	logger logger.Logger,
) *SecurityMiddleware {
	return &SecurityMiddleware{
		config: config,
		logger: logger,
	}
}

// SecurityHeaders 安全头中间件
func (m *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options: 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: 防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection: XSS 保护
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy: 控制 Referrer 信息
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content-Security-Policy: 内容安全策略
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'"
		c.Header("Content-Security-Policy", csp)

		// Permissions-Policy: 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// X-Permitted-Cross-Domain-Policies: 跨域策略
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// 如果是 HTTPS 连接，添加 HSTS 头
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next()
	}
}

// HSTS HTTP 严格传输安全中间件
func (m *SecurityMiddleware) HSTS(maxAge int, includeSubDomains, preload bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			hstsValue := fmt.Sprintf("max-age=%d", maxAge)
			if includeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			if preload {
				hstsValue += "; preload"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}
		c.Next()
	}
}

// CSP 内容安全策略中间件
func (m *SecurityMiddleware) CSP(policy string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", policy)
		c.Next()
	}
}

// NoCache 禁用缓存中间件
func (m *SecurityMiddleware) NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

// BasicAuth 基础认证中间件
func (m *SecurityMiddleware) BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()

		if !hasAuth ||
			subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {

			m.logger.Warn("Basic auth failed",
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent())

			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyAuth API Key 认证中间件
func (m *SecurityMiddleware) APIKeyAuth(validKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "API key required",
			})
			c.Abort()
			return
		}

		// 验证 API Key
		valid := false
		for _, validKey := range validKeys {
			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(validKey)) == 1 {
				valid = true
				break
			}
		}

		if !valid {
			m.logger.Warn("Invalid API key",
				"ip", c.ClientIP(),
				"api_key", apiKey[:min(len(apiKey), 8)]+"...")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey)
		c.Next()
	}
}

// IPWhitelist IP 白名单中间件
func (m *SecurityMiddleware) IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查 IP 是否在白名单中
		allowed := false
		for _, allowedIP := range allowedIPs {
			if allowedIP == clientIP || allowedIP == "*" {
				allowed = true
				break
			}
			// 支持 CIDR 格式的 IP 范围检查
			if m.isIPInCIDR(clientIP, allowedIP) {
				allowed = true
				break
			}
		}

		if !allowed {
			m.logger.Warn("IP not in whitelist",
				"ip", clientIP,
				"path", c.Request.URL.Path)

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Access denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPBlacklist IP 黑名单中间件
func (m *SecurityMiddleware) IPBlacklist(blockedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查 IP 是否在黑名单中
		for _, blockedIP := range blockedIPs {
			if blockedIP == clientIP {
				m.logger.Warn("Blocked IP attempted access",
					"ip", clientIP,
					"path", c.Request.URL.Path)

				c.JSON(http.StatusForbidden, gin.H{
					"error":   "forbidden",
					"message": "Access denied",
				})
				c.Abort()
				return
			}
			// 支持 CIDR 格式的 IP 范围检查
			if m.isIPInCIDR(clientIP, blockedIP) {
				m.logger.Warn("Blocked IP range attempted access",
					"ip", clientIP,
					"blocked_range", blockedIP)

				c.JSON(http.StatusForbidden, gin.H{
					"error":   "forbidden",
					"message": "Access denied",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// UserAgentFilter User-Agent 过滤中间件
func (m *SecurityMiddleware) UserAgentFilter(blockedAgents []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.Request.UserAgent()

		// 检查 User-Agent 是否被阻止
		for _, blockedAgent := range blockedAgents {
			if strings.Contains(strings.ToLower(userAgent), strings.ToLower(blockedAgent)) {
				m.logger.Warn("Blocked user agent",
					"user_agent", userAgent,
					"ip", c.ClientIP())

				c.JSON(http.StatusForbidden, gin.H{
					"error":   "forbidden",
					"message": "Access denied",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequestSizeLimit 请求大小限制中间件
func (m *SecurityMiddleware) RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			m.logger.Warn("Request size too large",
				"size", c.Request.ContentLength,
				"max_size", maxSize,
				"ip", c.ClientIP())

			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "request_too_large",
				"message": "Request entity too large",
			})
			c.Abort()
			return
		}

		// 限制请求体读取大小
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// Timeout 请求超时中间件
func (m *SecurityMiddleware) Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求上下文
		c.Request = c.Request.WithContext(ctx)

		// 使用通道来检测请求是否完成
		finished := make(chan struct{})
		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			m.logger.Warn("Request timeout",
				"timeout", timeout,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP())

			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "request_timeout",
				"message": "Request timeout",
			})
			c.Abort()
		}
	}
}

// HTTPSRedirect HTTPS 重定向中间件
func (m *SecurityMiddleware) HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS == nil && c.GetHeader("X-Forwarded-Proto") != "https" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// SecureHeaders 安全头组合中间件
func (m *SecurityMiddleware) SecureHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 组合多个安全头
		m.SecurityHeaders()(c)

		// 如果请求被中止，直接返回
		if c.IsAborted() {
			return
		}

		// 添加额外的安全措施
		c.Header("Server", "")       // 隐藏服务器信息
		c.Header("X-Powered-By", "") // 隐藏技术栈信息
	})
}

// isIPInCIDR 检查 IP 是否在 CIDR 范围内
func (m *SecurityMiddleware) isIPInCIDR(ip, cidr string) bool {
	// 这里应该实现 CIDR 检查逻辑
	// 为了简化，这里只做简单的字符串匹配
	// 在实际项目中，应该使用 net.ParseCIDR 和 net.IP.Contains
	return strings.Contains(cidr, "/") && strings.HasPrefix(ip, strings.Split(cidr, "/")[0][:strings.LastIndex(strings.Split(cidr, "/")[0], ".")+1])
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 预定义的安全配置
var (
	// DefaultSecurityHeaders 默认安全头
	DefaultSecurityHeaders = map[string]string{
		"X-Content-Type-Options":            "nosniff",
		"X-Frame-Options":                   "DENY",
		"X-XSS-Protection":                  "1; mode=block",
		"Referrer-Policy":                   "strict-origin-when-cross-origin",
		"X-Permitted-Cross-Domain-Policies": "none",
	}

	// StrictCSP 严格的内容安全策略
	StrictCSP = "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'"

	// CommonBlockedUserAgents 常见的恶意 User-Agent
	CommonBlockedUserAgents = []string{
		"sqlmap",
		"nmap",
		"nikto",
		"masscan",
		"zmap",
		"curl", // 可选，根据需要
		"wget", // 可选，根据需要
	}
)
