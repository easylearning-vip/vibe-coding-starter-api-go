package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

func TestUserLoginWithUsername(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 加载配置
	cfg, err := config.LoadConfig("../configs/config.sqlite.yaml")
	require.NoError(t, err)

	// 使用内存数据库进行测试
	cfg.Database.Database = ":memory:"

	// 创建日志实例
	log, err := logger.New(cfg)
	require.NoError(t, err)

	// 创建数据库连接
	db, err := database.New(cfg, log)
	require.NoError(t, err)

	// 执行数据库迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.Comment{},
		&model.File{},
		&model.DictCategory{},
		&model.DictItem{},
	)
	require.NoError(t, err)

	// 创建缓存
	cacheInstance, err := cache.New(cfg, log)
	require.NoError(t, err)

	// 创建仓储
	userRepo := repository.NewUserRepository(db, log)

	// 创建服务
	userService := service.NewUserService(userRepo, cacheInstance, log, cfg)

	// 创建处理器
	userHandler := handler.NewUserHandler(userService, log)

	// 设置路由
	router := gin.New()
	api := router.Group("/api/v1")

	// 注册公开路由
	users := api.Group("/users")
	users.POST("/register", userHandler.Register)
	users.POST("/login", userHandler.Login)

	// 测试用户注册
	t.Run("Register User", func(t *testing.T) {
		registerReq := service.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Nickname: "Test User",
		}

		reqBody, _ := json.Marshal(registerReq)
		req := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "testuser", response["username"])
		assert.Equal(t, "test@example.com", response["email"])
		assert.NotContains(t, response, "password") // 密码不应该返回
	})

	// 测试用户名登录
	t.Run("Login with Username", func(t *testing.T) {
		loginReq := service.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response service.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "testuser", response.User.Username)
		assert.Equal(t, "test@example.com", response.User.Email)
	})

	// 测试错误的用户名
	t.Run("Login with Invalid Username", func(t *testing.T) {
		loginReq := service.LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "login_failed", response["error"])
	})

	// 测试错误的密码
	t.Run("Login with Invalid Password", func(t *testing.T) {
		loginReq := service.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "login_failed", response["error"])
	})

	// 测试用户名唯一性
	t.Run("Register Duplicate Username", func(t *testing.T) {
		registerReq := service.RegisterRequest{
			Username: "testuser", // 重复的用户名
			Email:    "another@example.com",
			Password: "password123",
			Nickname: "Another User",
		}

		reqBody, _ := json.Marshal(registerReq)
		req := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "registration_failed", response["error"])
		assert.Contains(t, response["message"], "already exists")
	})
}
