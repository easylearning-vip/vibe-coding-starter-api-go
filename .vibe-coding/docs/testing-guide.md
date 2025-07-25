# 测试指南

## 测试策略概览

本项目采用多层次测试策略，确保代码质量和系统稳定性：

```mermaid
pyramid
    title 测试金字塔
    
    "E2E Tests" : 5
    "Integration Tests" : 15
    "Unit Tests" : 80
```

### 测试层次
1. **单元测试 (80%)**: 测试单个函数/方法的逻辑
2. **集成测试 (15%)**: 测试模块间的交互
3. **端到端测试 (5%)**: 测试完整的业务流程

## 测试工具链

### 核心测试框架
- **testing**: Go 标准测试包
- **testify/suite**: 测试套件和断言
- **testify/mock**: Mock 对象生成
- **testify/assert**: 断言库

### 测试数据库
- **SQLite**: 内存数据库，快速测试
- **testcontainers**: 容器化测试环境

### 测试工具
- **go test**: 运行测试
- **go test -cover**: 覆盖率分析
- **go test -race**: 竞态条件检测

## 项目测试结构

```
test/
├── config/                 # 测试配置
│   └── test_config.go
├── handler/                # Handler 层测试
│   ├── user_handler_test.go
│   ├── article_handler_test.go
│   └── file_handler_test.go
├── service/                # Service 层测试
│   ├── user_service_test.go
│   ├── article_service_test.go
│   └── file_service_test.go
├── repository/             # Repository 层测试
│   ├── user_repository_test.go
│   ├── article_repository_test.go
│   └── file_repository_test.go
├── mocks/                  # Mock 对象
│   ├── repository_mocks.go
│   ├── service_mocks.go
│   └── cache_mocks.go
├── testutil/               # 测试工具
│   ├── database.go
│   ├── cache.go
│   └── logger.go
└── scripts/                # 测试脚本
    └── run_repository_tests.sh
```

## 单元测试

### 测试套件模式
```go
package service

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/suite"
    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/test/mocks"
)

type UserServiceTestSuite struct {
    suite.Suite
    service     UserService
    mockRepo    *mocks.MockUserRepository
    mockCache   *mocks.MockCache
    mockLogger  *mocks.MockLogger
}

func (s *UserServiceTestSuite) SetupTest() {
    s.mockRepo = &mocks.MockUserRepository{}
    s.mockCache = &mocks.MockCache{}
    s.mockLogger = &mocks.MockLogger{}
    
    s.service = NewUserService(s.mockRepo, s.mockCache, s.mockLogger)
}

func (s *UserServiceTestSuite) TearDownTest() {
    // 清理资源
}

func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

### 测试用例示例
```go
func (s *UserServiceTestSuite) TestRegister_Success() {
    // Given
    req := &RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        Nickname: "Test User",
    }
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Nickname: "Test User",
    }
    
    s.mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, repository.ErrUserNotFound)
    s.mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, repository.ErrUserNotFound)
    s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(expectedUser, nil)
    
    // When
    user, err := s.service.Register(context.Background(), req)
    
    // Then
    s.NoError(err)
    s.Equal("testuser", user.Username)
    s.Equal("test@example.com", user.Email)
    s.Equal("Test User", user.Nickname)
    s.mockRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestRegister_UserExists() {
    // Given
    req := &RegisterRequest{
        Username: "existinguser",
        Email:    "existing@example.com",
        Password: "password123",
    }
    
    existingUser := &model.User{
        ID:       1,
        Username: "existinguser",
        Email:    "existing@example.com",
    }
    
    s.mockRepo.On("GetByUsername", mock.Anything, "existinguser").Return(existingUser, nil)
    
    // When
    user, err := s.service.Register(context.Background(), req)
    
    // Then
    s.Error(err)
    s.Nil(user)
    s.Equal(ErrUserExists, err)
    s.mockRepo.AssertExpectations(s.T())
}
```

### Mock 对象定义
```go
package mocks

import (
    "context"
    "github.com/stretchr/testify/mock"
    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/internal/repository"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
    args := m.Called(ctx, user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    args := m.Called(ctx, username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}
```

## 集成测试

### 数据库集成测试
```go
package repository

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/suite"
    "vibe-coding-starter/internal/model"
    "vibe-coding-starter/test/testutil"
)

type UserRepositoryIntegrationTestSuite struct {
    suite.Suite
    db   database.Database
    repo repository.UserRepository
}

func (s *UserRepositoryIntegrationTestSuite) SetupSuite() {
    // 设置测试数据库
    db, err := testutil.SetupTestDatabase()
    s.Require().NoError(err)
    s.db = db
    
    // 运行迁移
    err = s.db.AutoMigrate(&model.User{})
    s.Require().NoError(err)
    
    // 创建 repository
    logger := testutil.NewTestLogger()
    s.repo = repository.NewUserRepository(s.db, logger)
}

func (s *UserRepositoryIntegrationTestSuite) TearDownSuite() {
    s.db.Close()
}

func (s *UserRepositoryIntegrationTestSuite) SetupTest() {
    // 清理测试数据
    s.db.GetDB().Exec("DELETE FROM users")
}

func (s *UserRepositoryIntegrationTestSuite) TestCreate_Success() {
    // Given
    user := &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        Nickname: "Test User",
    }
    
    // When
    createdUser, err := s.repo.Create(context.Background(), user)
    
    // Then
    s.NoError(err)
    s.NotZero(createdUser.ID)
    s.Equal("testuser", createdUser.Username)
    s.Equal("test@example.com", createdUser.Email)
    s.NotZero(createdUser.CreatedAt)
}

func TestUserRepositoryIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryIntegrationTestSuite))
}
```

### HTTP 集成测试
```go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/suite"
    "vibe-coding-starter/internal/handler"
    "vibe-coding-starter/test/mocks"
)

type UserHandlerIntegrationTestSuite struct {
    suite.Suite
    router      *gin.Engine
    mockService *mocks.MockUserService
    handler     *handler.UserHandler
}

func (s *UserHandlerIntegrationTestSuite) SetupTest() {
    gin.SetMode(gin.TestMode)
    s.router = gin.New()
    
    s.mockService = &mocks.MockUserService{}
    s.handler = handler.NewUserHandler(s.mockService)
    
    // 注册路由
    api := s.router.Group("/api/v1")
    s.handler.RegisterRoutes(api)
}

func (s *UserHandlerIntegrationTestSuite) TestRegister_Success() {
    // Given
    reqBody := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
        "nickname": "Test User",
    }
    
    expectedUser := &model.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Nickname: "Test User",
    }
    
    s.mockService.On("Register", mock.Anything, mock.AnythingOfType("*service.RegisterRequest")).Return(expectedUser, nil)
    
    jsonBody, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    
    // When
    w := httptest.NewRecorder()
    s.router.ServeHTTP(w, req)
    
    // Then
    s.Equal(http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    s.NoError(err)
    s.Equal(float64(201), response["code"])
    s.Equal("user registered successfully", response["message"])
    
    s.mockService.AssertExpectations(s.T())
}
```

## 端到端测试

### 完整业务流程测试
```go
func (s *E2ETestSuite) TestUserRegistrationAndLogin() {
    // 1. 用户注册
    registerReq := map[string]interface{}{
        "username": "e2euser",
        "email":    "e2e@example.com",
        "password": "password123",
        "nickname": "E2E User",
    }
    
    registerResp := s.makeRequest("POST", "/api/v1/users/register", registerReq, "")
    s.Equal(201, registerResp.Code)
    
    // 2. 用户登录
    loginReq := map[string]interface{}{
        "username": "e2euser",
        "password": "password123",
    }
    
    loginResp := s.makeRequest("POST", "/api/v1/users/login", loginReq, "")
    s.Equal(200, loginResp.Code)
    
    var loginData map[string]interface{}
    json.Unmarshal(loginResp.Body.Bytes(), &loginData)
    token := loginData["data"].(map[string]interface{})["token"].(string)
    
    // 3. 获取用户信息
    profileResp := s.makeRequest("GET", "/api/v1/users/profile", nil, token)
    s.Equal(200, profileResp.Code)
    
    var profileData map[string]interface{}
    json.Unmarshal(profileResp.Body.Bytes(), &profileData)
    user := profileData["data"].(map[string]interface{})
    s.Equal("e2euser", user["username"])
    s.Equal("e2e@example.com", user["email"])
}
```

## 测试工具和辅助函数

### 测试数据库设置
```go
package testutil

import (
    "vibe-coding-starter/pkg/database"
    "vibe-coding-starter/internal/config"
)

func SetupTestDatabase() (database.Database, error) {
    cfg := &config.Config{
        Database: config.DatabaseConfig{
            Driver:   "sqlite",
            Database: ":memory:",
        },
    }
    
    logger := NewTestLogger()
    return database.New(cfg, logger)
}
```

### 测试缓存设置
```go
package testutil

import (
    "context"
    "time"
    "vibe-coding-starter/pkg/cache"
)

type TestCache struct {
    data map[string]string
}

func NewTestCache() cache.Cache {
    return &TestCache{
        data: make(map[string]string),
    }
}

func (c *TestCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    c.data[key] = value.(string)
    return nil
}

func (c *TestCache) Get(ctx context.Context, key string) (string, error) {
    if val, exists := c.data[key]; exists {
        return val, nil
    }
    return "", cache.ErrKeyNotFound
}
```

### 测试日志设置
```go
package testutil

import (
    "vibe-coding-starter/pkg/logger"
    "go.uber.org/zap"
    "go.uber.org/zap/zaptest"
)

func NewTestLogger() logger.Logger {
    zapLogger := zaptest.NewLogger(nil)
    return &testLogger{logger: zapLogger}
}

type testLogger struct {
    logger *zap.Logger
}

func (l *testLogger) Info(msg string, fields ...interface{}) {
    l.logger.Sugar().Infow(msg, fields...)
}

func (l *testLogger) Error(msg string, fields ...interface{}) {
    l.logger.Sugar().Errorw(msg, fields...)
}
```

## 测试运行和报告

### 运行测试命令
```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./internal/service/...

# 运行测试并生成覆盖率报告
make test-coverage

# 运行竞态条件检测
go test -race ./...

# 运行基准测试
go test -bench=. ./...

# 运行集成测试
./test/scripts/run_repository_tests.sh
```

### 覆盖率报告
```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率统计
go tool cover -func=coverage.out
```

### 测试报告示例
```
=== RUN   TestUserServiceTestSuite
=== RUN   TestUserServiceTestSuite/TestRegister_Success
=== RUN   TestUserServiceTestSuite/TestRegister_UserExists
=== RUN   TestUserServiceTestSuite/TestLogin_Success
=== RUN   TestUserServiceTestSuite/TestLogin_InvalidPassword
--- PASS: TestUserServiceTestSuite (0.01s)
    --- PASS: TestUserServiceTestSuite/TestRegister_Success (0.00s)
    --- PASS: TestUserServiceTestSuite/TestRegister_UserExists (0.00s)
    --- PASS: TestUserServiceTestSuite/TestLogin_Success (0.00s)
    --- PASS: TestUserServiceTestSuite/TestLogin_InvalidPassword (0.00s)
PASS
coverage: 85.2% of statements
```

## 测试最佳实践

### 1. 测试命名规范
```go
// 格式: Test[FunctionName]_[Scenario]_[ExpectedResult]
func TestUserService_Register_Success(t *testing.T) {}
func TestUserService_Register_UserExists_ReturnsError(t *testing.T) {}
func TestUserService_Login_InvalidPassword_ReturnsError(t *testing.T) {}
```

### 2. 测试数据管理
```go
// 使用工厂函数创建测试数据
func createTestUser() *model.User {
    return &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashedpassword",
        Nickname: "Test User",
    }
}

// 使用表驱动测试
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "test@example.com", true},
        {"invalid email", "invalid-email", false},
        {"empty email", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3. 测试隔离
- 每个测试用例独立运行
- 使用 SetupTest/TearDownTest 清理状态
- 避免测试间的依赖关系

### 4. Mock 使用原则
- 只 Mock 外部依赖
- 验证 Mock 调用
- 保持 Mock 简单

### 5. 测试覆盖率目标
- 单元测试覆盖率 > 80%
- 关键业务逻辑覆盖率 > 90%
- 边界条件和错误处理全覆盖

## CI/CD 集成

### GitHub Actions 配置
```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.23
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

通过遵循这个测试指南，可以确保代码质量，提高系统稳定性，并为持续集成和部署提供可靠的保障。
