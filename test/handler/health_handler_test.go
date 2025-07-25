package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/test/mocks"
)

// HealthHandlerTestSuite 健康检查处理器测试套件
type HealthHandlerTestSuite struct {
	suite.Suite
	db      *mocks.MockDatabase
	cache   *mocks.MockCache
	logger  *mocks.MockLogger
	handler *handler.HealthHandler
	router  *gin.Engine
}

// SetupSuite 设置测试套件
func (suite *HealthHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	suite.db = new(mocks.MockDatabase)
	suite.cache = new(mocks.MockCache)
	suite.logger = new(mocks.MockLogger)

	// 创建健康检查处理器
	suite.handler = handler.NewHealthHandler(
		suite.db,
		suite.cache,
		suite.logger,
	)

	// 设置路由
	suite.router = gin.New()
	suite.handler.RegisterRoutes(suite.router)
}

// SetupTest 每个测试前的设置
func (suite *HealthHandlerTestSuite) SetupTest() {
	suite.db.ExpectedCalls = nil
	suite.cache.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestHealthCheck 测试健康检查
func (suite *HealthHandlerTestSuite) TestHealthCheck() {
	// Mock 数据库健康检查
	suite.db.On("Health").Return(nil)

	// Mock 缓存健康检查
	suite.cache.On("Health").Return(nil)

	// Mock 日志
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response.Status)
	assert.Equal(suite.T(), "healthy", response.Services["database"].Status)
	assert.Equal(suite.T(), "healthy", response.Services["cache"].Status)
	assert.NotEmpty(suite.T(), response.Timestamp)

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
}

// TestHealthCheckDatabaseError 测试数据库健康检查失败
func (suite *HealthHandlerTestSuite) TestHealthCheckDatabaseError() {
	// Mock 数据库健康检查失败
	suite.db.On("Health").Return(assert.AnError)

	// Mock 缓存健康检查成功
	suite.cache.On("Health").Return(nil)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	var response handler.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "unhealthy", response.Status)
	assert.Equal(suite.T(), "unhealthy", response.Services["database"].Status)
	assert.Equal(suite.T(), "healthy", response.Services["cache"].Status)

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestHealthCheckCacheError 测试缓存健康检查失败
func (suite *HealthHandlerTestSuite) TestHealthCheckCacheError() {
	// Mock 数据库健康检查成功
	suite.db.On("Health").Return(nil)

	// Mock 缓存健康检查失败
	suite.cache.On("Health").Return(assert.AnError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	var response handler.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "unhealthy", response.Status)
	assert.Equal(suite.T(), "healthy", response.Services["database"].Status)
	assert.Equal(suite.T(), "unhealthy", response.Services["cache"].Status)

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestHealthCheckAllServicesError 测试所有服务健康检查失败
func (suite *HealthHandlerTestSuite) TestHealthCheckAllServicesError() {
	// Mock 数据库健康检查失败
	suite.db.On("Health").Return(assert.AnError)

	// Mock 缓存健康检查失败
	suite.cache.On("Health").Return(assert.AnError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return().Twice()
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	var response handler.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "unhealthy", response.Status)
	assert.Equal(suite.T(), "unhealthy", response.Services["database"].Status)
	assert.Equal(suite.T(), "unhealthy", response.Services["cache"].Status)

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestLiveness 测试存活检查
func (suite *HealthHandlerTestSuite) TestLiveness() {
	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["alive"])
	assert.NotEmpty(suite.T(), response["timestamp"])
}

// TestReadiness 测试就绪检查
func (suite *HealthHandlerTestSuite) TestReadiness() {
	// Mock 数据库健康检查
	suite.db.On("Health").Return(nil)

	// Mock 缓存健康检查
	suite.cache.On("Health").Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), response["ready"].(bool))
	assert.NotEmpty(suite.T(), response["timestamp"])

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
}

// TestReadinessNotReady 测试就绪检查失败
func (suite *HealthHandlerTestSuite) TestReadinessNotReady() {
	// Mock 数据库健康检查失败
	suite.db.On("Health").Return(assert.AnError)

	// Mock 缓存健康检查成功
	suite.cache.On("Health").Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.False(suite.T(), response["ready"].(bool))
	assert.NotEmpty(suite.T(), response["timestamp"])

	// 验证mock调用
	suite.db.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
}

// TestHealthHandlerTestSuite 运行测试套件
func TestHealthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthHandlerTestSuite))
}
