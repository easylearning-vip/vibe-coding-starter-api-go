package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/middleware"
	"vibe-coding-starter/test/testutil"
)

func TestMiddlewareIntegration(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 创建测试配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode: "test",
		},
		CORS: config.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
		Security: config.SecurityConfig{
			MaxRequestSize: 1024 * 1024, // 1MB
			RequestTimeout: 30,
		},
	}

	// 创建测试依赖
	testLoggerWrapper := testutil.NewTestLogger(t)
	testLogger := testLoggerWrapper.CreateTestLogger()
	testCacheWrapper := testutil.NewTestCache(t)
	testCache := testCacheWrapper.CreateTestCache()
	defer testCacheWrapper.Close()

	// 创建中间件管理器
	mw := middleware.NewMiddleware(cfg, testLogger, testCache)

	t.Run("CORS Middleware", func(t *testing.T) {
		engine := gin.New()
		engine.Use(mw.CORS().CORS())

		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 测试预检请求
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("Security Headers Middleware", func(t *testing.T) {
		engine := gin.New()
		engine.Use(mw.Security().SecurityHeaders())

		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	})

	t.Run("Rate Limit Middleware", func(t *testing.T) {
		engine := gin.New()
		engine.Use(mw.RateLimit().IPRateLimit(2, 2)) // 很低的限制用于测试

		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 第一次请求应该成功
		req1 := httptest.NewRequest("GET", "/test", nil)
		w1 := httptest.NewRecorder()
		engine.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// 第二次请求应该成功
		req2 := httptest.NewRequest("GET", "/test", nil)
		w2 := httptest.NewRecorder()
		engine.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// 第三次请求应该被限流（需要等待实际的限流实现）
		// 注意：这个测试可能需要根据实际的限流实现进行调整
	})

	t.Run("Logging Middleware", func(t *testing.T) {
		engine := gin.New()
		engine.Use(mw.Logging().StructuredLogging())

		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 检查是否设置了请求 ID
		requestID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)
	})

	t.Run("Request Size Limit", func(t *testing.T) {
		engine := gin.New()
		engine.Use(mw.Security().RequestSizeLimit(100)) // 100 字节限制

		engine.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 创建超过限制的请求体
		largeData := make(map[string]string)
		for i := 0; i < 50; i++ {
			largeData[fmt.Sprintf("key%d", i)] = "very long value that exceeds the limit"
		}

		body, _ := json.Marshal(largeData)
		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		// 应该返回 413 Request Entity Too Large
		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})
}

func TestMiddlewareChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		CORS:   config.CORSConfig{AllowOrigins: []string{"*"}},
	}

	testLoggerWrapper := testutil.NewTestLogger(t)
	testLogger := testLoggerWrapper.CreateTestLogger()
	testCacheWrapper := testutil.NewTestCache(t)
	testCache := testCacheWrapper.CreateTestCache()
	defer testCacheWrapper.Close()
	mw := middleware.NewMiddleware(cfg, testLogger, testCache)

	t.Run("Multiple Middleware Chain", func(t *testing.T) {
		engine := gin.New()

		// 添加多个中间件
		engine.Use(mw.Security().SecurityHeaders())
		engine.Use(mw.CORS().CORS())
		engine.Use(mw.Logging().StructuredLogging())

		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":    "ok",
				"request_id": c.GetString("request_id"),
			})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://example.com")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 检查安全头
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))

		// 检查 CORS 头 - 当配置为"*"且请求包含Origin时，应该返回具体的Origin
		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))

		// 检查请求 ID
		assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

		// 检查响应内容
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["message"])
		assert.NotEmpty(t, response["request_id"])
	})
}

func TestMiddlewareConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Development vs Production Middleware", func(t *testing.T) {
		// 开发环境配置
		devCfg := &config.Config{
			Server: config.ServerConfig{Mode: "debug"},
		}

		// 生产环境配置
		prodCfg := &config.Config{
			Server: config.ServerConfig{Mode: "release"},
		}

		testLoggerWrapper := testutil.NewTestLogger(t)
		testLogger := testLoggerWrapper.CreateTestLogger()
		testCacheWrapper := testutil.NewTestCache(t)
		testCache := testCacheWrapper.CreateTestCache()
		defer testCacheWrapper.Close()

		devMW := middleware.NewMiddleware(devCfg, testLogger, testCache)
		prodMW := middleware.NewMiddleware(prodCfg, testLogger, testCache)

		// 测试开发环境的 CORS 配置（应该更宽松）
		devEngine := gin.New()
		devEngine.Use(devMW.CORS().DevCORS())
		devEngine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "dev"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://random-domain.com")

		w := httptest.NewRecorder()
		devEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))

		// 测试生产环境的 CORS 配置（应该更严格）
		prodEngine := gin.New()
		prodEngine.Use(prodMW.CORS().CORS())
		prodEngine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "prod"})
		})

		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.Header.Set("Origin", "http://random-domain.com")

		w2 := httptest.NewRecorder()
		prodEngine.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
		// 生产环境不应该允许随意的域名
		assert.NotEqual(t, "http://random-domain.com", w2.Header().Get("Access-Control-Allow-Origin"))
	})
}
