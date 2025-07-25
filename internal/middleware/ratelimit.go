package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Rate     int           `json:"rate"`   // 每秒允许的请求数
	Burst    int           `json:"burst"`  // 突发请求数
	Window   time.Duration `json:"window"` // 时间窗口
	KeyFunc  KeyFunc       `json:"-"`      // 生成限流键的函数
	SkipFunc SkipFunc      `json:"-"`      // 跳过限流的函数
}

// KeyFunc 生成限流键的函数类型
type KeyFunc func(*gin.Context) string

// SkipFunc 跳过限流的函数类型
type SkipFunc func(*gin.Context) bool

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	config *config.Config
	cache  cache.Cache
	logger logger.Logger
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(
	config *config.Config,
	cache cache.Cache,
	logger logger.Logger,
) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		config: config,
		cache:  cache,
		logger: logger,
	}
}

// RateLimit 通用限流中间件
func (m *RateLimitMiddleware) RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过限流
		if config.SkipFunc != nil && config.SkipFunc(c) {
			c.Next()
			return
		}

		// 生成限流键
		key := config.KeyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		// 检查限流
		allowed, remaining, resetTime, err := m.checkRateLimit(key, config)
		if err != nil {
			m.logger.Error("Rate limit check failed", "error", err, "key", key)
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			m.logger.Warn("Rate limit exceeded",
				"key", key,
				"path", c.Request.URL.Path,
				"method", c.Request.Method)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimit 基于 IP 的限流
func (m *RateLimitMiddleware) IPRateLimit(rate, burst int) gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   rate,
		Burst:  burst,
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			return "ip:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// UserRateLimit 基于用户的限流
func (m *RateLimitMiddleware) UserRateLimit(rate, burst int) gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   rate,
		Burst:  burst,
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("user:%v", userID)
			}
			return "ip:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// APIKeyRateLimit 基于 API Key 的限流
func (m *RateLimitMiddleware) APIKeyRateLimit(rate, burst int) gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   rate,
		Burst:  burst,
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey != "" {
				return "api_key:" + apiKey
			}
			return "ip:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// EndpointRateLimit 基于端点的限流
func (m *RateLimitMiddleware) EndpointRateLimit(rate, burst int) gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   rate,
		Burst:  burst,
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			endpoint := c.Request.Method + ":" + c.FullPath()
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("endpoint:%s:user:%v", endpoint, userID)
			}
			return fmt.Sprintf("endpoint:%s:ip:%s", endpoint, c.ClientIP())
		},
	}
	return m.RateLimit(config)
}

// LoginRateLimit 登录接口专用限流
func (m *RateLimitMiddleware) LoginRateLimit() gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   5,  // 每分钟 5 次
		Burst:  10, // 突发 10 次
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			return "login:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// RegisterRateLimit 注册接口专用限流
func (m *RateLimitMiddleware) RegisterRateLimit() gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   2, // 每分钟 2 次
		Burst:  5, // 突发 5 次
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			return "register:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// UploadRateLimit 文件上传限流
func (m *RateLimitMiddleware) UploadRateLimit() gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   10, // 每分钟 10 次
		Burst:  20, // 突发 20 次
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("upload:user:%v", userID)
			}
			return "upload:ip:" + c.ClientIP()
		},
	}
	return m.RateLimit(config)
}

// AdminRateLimit 管理员接口限流（更宽松）
func (m *RateLimitMiddleware) AdminRateLimit() gin.HandlerFunc {
	config := RateLimitConfig{
		Rate:   100, // 每分钟 100 次
		Burst:  200, // 突发 200 次
		Window: time.Minute,
		KeyFunc: func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return fmt.Sprintf("admin:user:%v", userID)
			}
			return "admin:ip:" + c.ClientIP()
		},
		SkipFunc: func(c *gin.Context) bool {
			// 跳过超级管理员的限流
			role, exists := c.Get("user_role")
			return exists && role == "super_admin"
		},
	}
	return m.RateLimit(config)
}

// checkRateLimit 检查限流状态
func (m *RateLimitMiddleware) checkRateLimit(key string, config RateLimitConfig) (bool, int, time.Time, error) {
	ctx := context.Background()
	now := time.Now()
	window := config.Window

	// 使用滑动窗口算法
	windowStart := now.Truncate(window)
	windowKey := fmt.Sprintf("rate_limit:%s:%d", key, windowStart.Unix())

	// 获取当前窗口的请求计数
	countStr, err := m.cache.Get(ctx, windowKey)
	if err != nil {
		// 如果键不存在，说明是新窗口
		countStr = "0"
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 0
	}

	// 检查是否超过限制
	if count >= config.Rate {
		resetTime := windowStart.Add(window)
		return false, 0, resetTime, nil
	}

	// 增加计数
	newCount := count + 1
	err = m.cache.Set(ctx, windowKey, strconv.Itoa(newCount), window)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	remaining := config.Rate - newCount
	if remaining < 0 {
		remaining = 0
	}

	resetTime := windowStart.Add(window)
	return true, remaining, resetTime, nil
}

// TokenBucketRateLimit 令牌桶限流实现
func (m *RateLimitMiddleware) TokenBucketRateLimit(rateLimit rate.Limit, burst int, keyFunc KeyFunc) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		// 获取或创建限流器
		limiter, exists := limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rateLimit, burst)
			limiters[key] = limiter
		}

		if !limiter.Allow() {
			m.logger.Warn("Token bucket rate limit exceeded",
				"key", key,
				"path", c.Request.URL.Path)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SlidingWindowRateLimit 滑动窗口限流
func (m *RateLimitMiddleware) SlidingWindowRateLimit(limit int, window time.Duration, keyFunc KeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		ctx := context.Background()
		windowKey := fmt.Sprintf("sliding_window:%s", key)

		// 简化的滑动窗口实现：使用计数器
		countStr, err := m.cache.Get(ctx, windowKey)
		var count int64 = 0
		if err == nil && countStr != "" {
			if parsedCount, parseErr := strconv.ParseInt(countStr, 10, 64); parseErr == nil {
				count = parsedCount
			}
		}

		if int(count) >= limit {
			m.logger.Warn("Sliding window rate limit exceeded",
				"key", key,
				"count", count,
				"limit", limit)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		// 增加计数器
		newCount := count + 1
		if err := m.cache.Set(ctx, windowKey, strconv.FormatInt(newCount, 10), window); err != nil {
			m.logger.Error("Failed to update sliding window count", "error", err)
		}

		c.Next()
	}
}

// GetRateLimitStatus 获取限流状态
func (m *RateLimitMiddleware) GetRateLimitStatus(key string, config RateLimitConfig) (int, int, time.Time, error) {
	ctx := context.Background()
	now := time.Now()
	window := config.Window

	windowStart := now.Truncate(window)
	windowKey := fmt.Sprintf("rate_limit:%s:%d", key, windowStart.Unix())

	countStr, err := m.cache.Get(ctx, windowKey)
	if err != nil {
		countStr = "0"
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 0
	}

	remaining := config.Rate - count
	if remaining < 0 {
		remaining = 0
	}

	resetTime := windowStart.Add(window)
	return count, remaining, resetTime, nil
}

// ClearRateLimit 清除限流记录
func (m *RateLimitMiddleware) ClearRateLimit(key string) error {
	// 简化实现：由于缓存接口不支持模式匹配，暂时返回 nil
	// 在实际使用中，可以考虑扩展缓存接口或使用具体的 Redis 客户端
	m.logger.Info("ClearRateLimit called", "key", key)
	return nil
}

// 预定义的键生成函数
var (
	// IPKeyFunc 基于 IP 生成键
	IPKeyFunc = func(c *gin.Context) string {
		return "ip:" + c.ClientIP()
	}

	// UserKeyFunc 基于用户 ID 生成键
	UserKeyFunc = func(c *gin.Context) string {
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		return "ip:" + c.ClientIP()
	}

	// EndpointKeyFunc 基于端点生成键
	EndpointKeyFunc = func(c *gin.Context) string {
		return fmt.Sprintf("endpoint:%s:%s", c.Request.Method, c.FullPath())
	}
)
