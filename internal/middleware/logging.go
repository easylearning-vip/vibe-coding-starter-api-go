package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/logger"
)

// LoggingConfig 日志配置
type LoggingConfig struct {
	SkipPaths       []string `json:"skip_paths"`        // 跳过记录的路径
	SkipMethods     []string `json:"skip_methods"`      // 跳过记录的方法
	LogRequestBody  bool     `json:"log_request_body"`  // 是否记录请求体
	LogResponseBody bool     `json:"log_response_body"` // 是否记录响应体
	MaxBodySize     int      `json:"max_body_size"`     // 最大记录的请求/响应体大小
	SensitiveFields []string `json:"sensitive_fields"`  // 敏感字段列表
}

// LoggingMiddleware 请求日志中间件
type LoggingMiddleware struct {
	config    *config.Config
	logger    logger.Logger
	logConfig LoggingConfig
}

// NewLoggingMiddleware 创建请求日志中间件
func NewLoggingMiddleware(
	config *config.Config,
	logger logger.Logger,
) *LoggingMiddleware {
	logConfig := LoggingConfig{
		SkipPaths: []string{
			"/health",
			"/ready",
			"/live",
			"/metrics",
			"/favicon.ico",
		},
		SkipMethods: []string{
			"OPTIONS",
		},
		LogRequestBody:  true,
		LogResponseBody: false,     // 默认不记录响应体，避免日志过大
		MaxBodySize:     1024 * 10, // 10KB
		SensitiveFields: []string{
			"password",
			"token",
			"secret",
			"key",
			"authorization",
			"cookie",
		},
	}

	return &LoggingMiddleware{
		config:    config,
		logger:    logger,
		logConfig: logConfig,
	}
}

// RequestLogging 请求日志中间件
func (m *LoggingMiddleware) RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过日志记录
		if m.shouldSkip(c) {
			c.Next()
			return
		}

		// 生成请求 ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if m.logConfig.LogRequestBody && c.Request.Body != nil {
			requestBody = m.readRequestBody(c)
		}

		// 创建响应写入器包装器
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			logResponse:    m.logConfig.LogResponseBody,
		}
		c.Writer = responseWriter

		// 记录请求信息
		m.logRequest(c, requestID, requestBody)

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 记录响应信息
		m.logResponse(c, requestID, duration, responseWriter)
	}
}

// StructuredLogging 结构化日志中间件
func (m *LoggingMiddleware) StructuredLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 检查是否跳过
		if m.shouldSkip(c) {
			c.Next()
			return
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 构建日志字段
		logFields := []interface{}{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"content_length", c.Request.ContentLength,
			"response_size", c.Writer.Size(),
		}

		// 添加用户信息（只在值存在且不为nil时添加）
		if userID != nil {
			logFields = append(logFields, "user_id", userID)
		}
		if username != nil {
			logFields = append(logFields, "username", username)
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// 根据状态码选择日志级别
		switch {
		case c.Writer.Status() >= 500:
			m.logger.Error("HTTP Request", logFields...)
		case c.Writer.Status() >= 400:
			m.logger.Warn("HTTP Request", logFields...)
		default:
			m.logger.Info("HTTP Request", logFields...)
		}
	}
}

// AccessLogging 访问日志中间件（类似 Apache/Nginx 格式）
func (m *LoggingMiddleware) AccessLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 检查是否跳过
		if m.shouldSkip(c) {
			c.Next()
			return
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 生成访问日志（Common Log Format 扩展）
		logLine := fmt.Sprintf("%s - %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %dms",
			c.ClientIP(),
			m.getUsername(c),
			startTime.Format("02/Jan/2006:15:04:05 -0700"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			c.Writer.Size(),
			c.Request.Referer(),
			c.Request.UserAgent(),
			duration.Milliseconds(),
		)

		m.logger.Info("Access Log", "log", logLine)
	}
}

// ErrorLogging 错误日志中间件
func (m *LoggingMiddleware) ErrorLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录错误
		if len(c.Errors) > 0 {
			requestID, _ := c.Get("request_id")
			userID, _ := c.Get("user_id")

			for _, err := range c.Errors {
				m.logger.Error("Request Error",
					"request_id", requestID,
					"user_id", userID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"error", err.Error(),
					"type", err.Type,
				)
			}
		}
	}
}

// SecurityLogging 安全日志中间件
func (m *LoggingMiddleware) SecurityLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检测可疑请求
		suspicious := m.detectSuspiciousRequest(c)
		if suspicious {
			m.logSecurityEvent(c, "suspicious_request")
		}

		c.Next()

		// 记录认证失败
		if c.Writer.Status() == http.StatusUnauthorized {
			m.logSecurityEvent(c, "authentication_failed")
		}

		// 记录权限拒绝
		if c.Writer.Status() == http.StatusForbidden {
			m.logSecurityEvent(c, "authorization_failed")
		}
	}
}

// shouldSkip 检查是否应该跳过日志记录
func (m *LoggingMiddleware) shouldSkip(c *gin.Context) bool {
	path := c.Request.URL.Path
	method := c.Request.Method

	// 检查跳过的路径
	for _, skipPath := range m.logConfig.SkipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// 检查跳过的方法
	for _, skipMethod := range m.logConfig.SkipMethods {
		if method == skipMethod {
			return true
		}
	}

	return false
}

// readRequestBody 读取请求体
func (m *LoggingMiddleware) readRequestBody(c *gin.Context) []byte {
	if c.Request.Body == nil {
		return nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		m.logger.Error("Failed to read request body", "error", err)
		return nil
	}

	// 恢复请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 限制记录的大小
	if len(body) > m.logConfig.MaxBodySize {
		body = body[:m.logConfig.MaxBodySize]
	}

	return body
}

// logRequest 记录请求信息
func (m *LoggingMiddleware) logRequest(c *gin.Context, requestID string, requestBody []byte) {
	logFields := []interface{}{
		"request_id", requestID,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"query", c.Request.URL.RawQuery,
		"client_ip", c.ClientIP(),
		"user_agent", c.Request.UserAgent(),
		"content_length", c.Request.ContentLength,
		"content_type", c.Request.Header.Get("Content-Type"),
	}

	// 添加请求头（过滤敏感信息）
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if !m.isSensitiveField(key) {
			headers[key] = strings.Join(values, ", ")
		}
	}
	logFields = append(logFields, "headers", headers)

	// 添加请求体
	if requestBody != nil && len(requestBody) > 0 {
		if m.isJSONContent(c.Request.Header.Get("Content-Type")) {
			// 过滤敏感字段
			filteredBody := m.filterSensitiveFields(requestBody)
			logFields = append(logFields, "request_body", string(filteredBody))
		} else {
			logFields = append(logFields, "request_body_size", len(requestBody))
		}
	}

	m.logger.Info("HTTP Request Started", logFields...)
}

// logResponse 记录响应信息
func (m *LoggingMiddleware) logResponse(c *gin.Context, requestID string, duration time.Duration, rw *responseWriter) {
	logFields := []interface{}{
		"request_id", requestID,
		"status", c.Writer.Status(),
		"duration_ms", duration.Milliseconds(),
		"response_size", c.Writer.Size(),
	}

	// 添加用户信息
	if userID, exists := c.Get("user_id"); exists {
		logFields = append(logFields, "user_id", userID)
	}

	// 添加响应体
	if m.logConfig.LogResponseBody && rw.body.Len() > 0 {
		responseBody := rw.body.Bytes()
		if len(responseBody) > m.logConfig.MaxBodySize {
			responseBody = responseBody[:m.logConfig.MaxBodySize]
		}
		logFields = append(logFields, "response_body", string(responseBody))
	}

	// 根据状态码选择日志级别
	switch {
	case c.Writer.Status() >= 500:
		m.logger.Error("HTTP Request Completed", logFields...)
	case c.Writer.Status() >= 400:
		m.logger.Warn("HTTP Request Completed", logFields...)
	default:
		m.logger.Info("HTTP Request Completed", logFields...)
	}
}

// logSecurityEvent 记录安全事件
func (m *LoggingMiddleware) logSecurityEvent(c *gin.Context, eventType string) {
	requestID, _ := c.Get("request_id")
	userID, _ := c.Get("user_id")

	m.logger.Warn("Security Event",
		"event_type", eventType,
		"request_id", requestID,
		"user_id", userID,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"client_ip", c.ClientIP(),
		"user_agent", c.Request.UserAgent(),
	)
}

// detectSuspiciousRequest 检测可疑请求
func (m *LoggingMiddleware) detectSuspiciousRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	userAgent := c.Request.UserAgent()

	// 检测 SQL 注入尝试
	if strings.Contains(path, "union") || strings.Contains(path, "select") ||
		strings.Contains(path, "drop") || strings.Contains(path, "insert") {
		return true
	}

	// 检测 XSS 尝试
	if strings.Contains(path, "<script") || strings.Contains(path, "javascript:") {
		return true
	}

	// 检测路径遍历尝试
	if strings.Contains(path, "../") || strings.Contains(path, "..\\") {
		return true
	}

	// 检测可疑的 User-Agent
	if userAgent == "" || strings.Contains(userAgent, "sqlmap") ||
		strings.Contains(userAgent, "nmap") {
		return true
	}

	return false
}

// getUsername 获取用户名
func (m *LoggingMiddleware) getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return "-"
}

// isSensitiveField 检查是否为敏感字段
func (m *LoggingMiddleware) isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	for _, sensitive := range m.logConfig.SensitiveFields {
		if strings.Contains(field, sensitive) {
			return true
		}
	}
	return false
}

// isJSONContent 检查是否为 JSON 内容
func (m *LoggingMiddleware) isJSONContent(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}

// filterSensitiveFields 过滤敏感字段
func (m *LoggingMiddleware) filterSensitiveFields(body []byte) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}

	// 递归过滤敏感字段
	m.filterMapSensitiveFields(data)

	filtered, err := json.Marshal(data)
	if err != nil {
		return body
	}

	return filtered
}

// filterMapSensitiveFields 递归过滤 map 中的敏感字段
func (m *LoggingMiddleware) filterMapSensitiveFields(data map[string]interface{}) {
	for key, value := range data {
		if m.isSensitiveField(key) {
			data[key] = "***FILTERED***"
		} else if subMap, ok := value.(map[string]interface{}); ok {
			m.filterMapSensitiveFields(subMap)
		}
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	gin.ResponseWriter
	body        *bytes.Buffer
	logResponse bool
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.logResponse {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	if w.logResponse {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}
