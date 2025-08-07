package handler

import (
	"bytes"
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
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// UserHandlerTestSuite 用户处理器测试套件
type UserHandlerTestSuite struct {
	suite.Suite
	userService *mocks.MockUserService
	logger      *mocks.MockLogger
	handler     *handler.UserHandler
	router      *gin.Engine
}

// SetupSuite 设置测试套件
func (suite *UserHandlerTestSuite) SetupSuite() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	suite.userService = new(mocks.MockUserService)
	suite.logger = new(mocks.MockLogger)

	// 创建用户处理器
	suite.handler = handler.NewUserHandler(
		suite.userService,
		suite.logger,
	)

	// 设置路由
	suite.router = gin.New()
	api := suite.router.Group("/api/v1")

	// 注册公共路由（不需要认证）
	users := api.Group("/users")
	users.POST("/register", suite.handler.Register)
	users.POST("/login", suite.handler.Login)

	// 注册需要认证的路由
	suite.handler.RegisterRoutes(api)
}

// SetupTest 每个测试前的设置
func (suite *UserHandlerTestSuite) SetupTest() {
	// 重置所有mock
	suite.userService.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestRegister 测试用户注册
func (suite *UserHandlerTestSuite) TestRegister() {
	reqBody := service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	user := &model.User{
		BaseModel: model.BaseModel{ID: 1},
		Username:  reqBody.Username,
		Email:     reqBody.Email,
		Nickname:  reqBody.Nickname,
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	// Mock 用户服务
	suite.userService.On("Register", mock.Anything, &reqBody).Return(user, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, response.ID)
	assert.Equal(suite.T(), user.Username, response.Username)
	assert.Equal(suite.T(), user.Email, response.Email)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
}

// TestRegisterInvalidJSON 测试注册时JSON格式错误
func (suite *UserHandlerTestSuite) TestRegisterInvalidJSON() {
	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 创建无效JSON请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "validation_error", response.Error)

	// 验证mock调用
	suite.logger.AssertExpectations(suite.T())
}

// TestLogin 测试用户登录
func (suite *UserHandlerTestSuite) TestLogin() {
	reqBody := service.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResponse := &service.LoginResponse{
		User: &model.PublicUser{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     model.UserRoleUser,
			Status:   model.UserStatusActive,
		},
		Token: "jwt-token-here",
	}

	// Mock 日志
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).Return()

	// Mock 用户服务
	suite.userService.On("Login", mock.Anything, &reqBody).Return(loginResponse, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response service.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), loginResponse.Token, response.Token)
	assert.Equal(suite.T(), loginResponse.User.ID, response.User.ID)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
}

// TestLoginServiceError 测试登录时服务错误
func (suite *UserHandlerTestSuite) TestLoginServiceError() {
	reqBody := service.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Mock 日志
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()
	suite.logger.On("Debug", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything).Return()
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	// Mock 用户服务返回错误
	suite.userService.On("Login", mock.Anything, &reqBody).Return(nil, assert.AnError)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "login_failed", response.Error)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetProfile 测试获取用户资料
func (suite *UserHandlerTestSuite) TestGetProfile() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Nickname:  "Test User",
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	// Mock 用户服务
	suite.userService.On("GetProfile", mock.Anything, userID).Return(user, nil)

	// 创建请求（需要模拟认证中间件设置的用户ID）
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)

	// 直接调用处理器方法
	suite.handler.GetProfile(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, response.ID)
	assert.Equal(suite.T(), user.Username, response.Username)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
}

// TestUpdateProfile 测试更新用户资料
func (suite *UserHandlerTestSuite) TestUpdateProfile() {
	userID := uint(1)
	reqBody := service.UpdateProfileRequest{
		Nickname: "Updated Nickname",
		Avatar:   "new-avatar.jpg",
	}

	updatedUser := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Nickname:  reqBody.Nickname,
		Avatar:    reqBody.Avatar,
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	// Mock 用户服务
	suite.userService.On("UpdateProfile", mock.Anything, userID, &reqBody).Return(updatedUser, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)

	// 直接调用处理器方法
	suite.handler.UpdateProfile(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedUser.ID, response.ID)
	assert.Equal(suite.T(), updatedUser.Nickname, response.Nickname)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
}

// TestChangePassword 测试修改密码
func (suite *UserHandlerTestSuite) TestChangePassword() {
	userID := uint(1)
	reqBody := service.ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}

	// Mock 用户服务
	suite.userService.On("ChangePassword", mock.Anything, userID, &reqBody).Return(nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)

	// 直接调用处理器方法
	suite.handler.ChangePassword(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Password changed successfully", response.Message)

	// 验证mock调用
	suite.userService.AssertExpectations(suite.T())
}

// TestUserHandlerTestSuite 运行测试套件
func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
