package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// UserServiceTestSuite 用户服务测试套件
type UserServiceTestSuite struct {
	suite.Suite
	userRepo *mocks.MockUserRepository
	cache    *mocks.MockCache
	logger   *mocks.MockLogger
	config   *config.Config
	service  service.UserService
	ctx      context.Context
}

// SetupSuite 设置测试套件
func (suite *UserServiceTestSuite) SetupSuite() {
	suite.userRepo = new(mocks.MockUserRepository)
	suite.cache = new(mocks.MockCache)
	suite.logger = new(mocks.MockLogger)
	suite.ctx = context.Background()

	// 创建测试配置
	suite.config = &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			Issuer:     "test-issuer",
			Expiration: 86400, // 24 hours in seconds
		},
	}

	// 创建用户服务
	suite.service = service.NewUserService(
		suite.userRepo,
		suite.cache,
		suite.logger,
		suite.config,
	)
}

// SetupTest 每个测试前的设置
func (suite *UserServiceTestSuite) SetupTest() {
	// 重置所有mock
	suite.userRepo.ExpectedCalls = nil
	suite.cache.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestRegister 测试用户注册
func (suite *UserServiceTestSuite) TestRegister() {
	req := &service.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	// Mock 检查邮箱是否存在
	suite.userRepo.On("GetByEmail", suite.ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock 检查用户名是否存在
	suite.userRepo.On("GetByUsername", suite.ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)

	// Mock 创建用户
	suite.userRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*model.User)
		user.ID = 1
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	})

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行注册
	user, err := suite.service.Register(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), req.Username, user.Username)
	assert.Equal(suite.T(), req.Email, user.Email)
	assert.Equal(suite.T(), req.Nickname, user.Nickname)
	assert.Equal(suite.T(), model.UserRoleUser, user.Role)
	assert.Equal(suite.T(), model.UserStatusActive, user.Status)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestRegisterDuplicateEmail 测试注册重复邮箱
func (suite *UserServiceTestSuite) TestRegisterDuplicateEmail() {
	req := &service.RegisterRequest{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: 1},
		Username:  "existing",
		Email:     "existing@example.com",
	}

	// Mock 检查邮箱已存在
	suite.userRepo.On("GetByEmail", suite.ctx, req.Email).Return(existingUser, nil)

	// 执行注册
	user, err := suite.service.Register(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "already exists")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
}

// TestRegisterDuplicateUsername 测试注册重复用户名
func (suite *UserServiceTestSuite) TestRegisterDuplicateUsername() {
	req := &service.RegisterRequest{
		Username: "existing",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
	}

	existingUser := &model.User{
		BaseModel: model.BaseModel{ID: 1},
		Username:  "existing",
		Email:     "existing@example.com",
	}

	// Mock 检查邮箱不存在
	suite.userRepo.On("GetByEmail", suite.ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Mock 检查用户名已存在
	suite.userRepo.On("GetByUsername", suite.ctx, req.Username).Return(existingUser, nil)

	// 执行注册
	user, err := suite.service.Register(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "already exists")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
}

// TestLogin 测试用户登录
func (suite *UserServiceTestSuite) TestLogin() {
	req := &service.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	user := &model.User{
		BaseModel: model.BaseModel{ID: 1},
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash of "password"
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	// Mock 获取用户
	suite.userRepo.On("GetByUsername", suite.ctx, req.Username).Return(user, nil)

	// Mock 更新最后登录时间
	suite.userRepo.On("UpdateLastLogin", suite.ctx, user.ID).Return(nil)

	// Mock 日志 - 匹配展开后的参数
	suite.logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	suite.logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	suite.logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	// 执行登录
	response, err := suite.service.Login(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.NotEmpty(suite.T(), response.Token)
	assert.NotNil(suite.T(), response.User)
	assert.Equal(suite.T(), user.ID, response.User.ID)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestLoginInvalidUsername 测试登录无效用户名
func (suite *UserServiceTestSuite) TestLoginInvalidUsername() {
	req := &service.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	// Mock 用户不存在
	suite.userRepo.On("GetByUsername", suite.ctx, req.Username).Return(nil, gorm.ErrRecordNotFound)

	// Mock 日志
	suite.logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	suite.logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	// 执行登录
	response, err := suite.service.Login(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "invalid username or password")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestLoginInactiveUser 测试登录非活跃用户
func (suite *UserServiceTestSuite) TestLoginInactiveUser() {
	req := &service.LoginRequest{
		Username: "inactive",
		Password: "password123",
	}

	user := &model.User{
		BaseModel: model.BaseModel{ID: 1},
		Username:  "inactive",
		Email:     "inactive@example.com",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
		Role:      model.UserRoleUser,
		Status:    model.UserStatusInactive,
	}

	// Mock 获取用户
	suite.userRepo.On("GetByUsername", suite.ctx, req.Username).Return(user, nil)

	// Mock 日志
	suite.logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	suite.logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	suite.logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	// 执行登录
	response, err := suite.service.Login(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "not active")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetProfile 测试获取用户资料
func (suite *UserServiceTestSuite) TestGetProfile() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Nickname:  "Test User",
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// 执行获取资料
	result, err := suite.service.GetProfile(suite.ctx, userID)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), user.ID, result.ID)
	assert.Equal(suite.T(), user.Username, result.Username)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
}

// TestGetProfileNotFound 测试获取不存在用户的资料
func (suite *UserServiceTestSuite) TestGetProfileNotFound() {
	userID := uint(999)

	// Mock 用户不存在
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(nil, gorm.ErrRecordNotFound)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行获取资料
	result, err := suite.service.GetProfile(suite.ctx, userID)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestUpdateProfile 测试更新用户资料
func (suite *UserServiceTestSuite) TestUpdateProfile() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Nickname:  "Old Nickname",
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	req := &service.UpdateProfileRequest{
		Nickname: "New Nickname",
		Avatar:   "new-avatar.jpg",
	}

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Mock 更新用户
	suite.userRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 执行更新资料
	result, err := suite.service.UpdateProfile(suite.ctx, userID, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Nickname, result.Nickname)
	assert.Equal(suite.T(), req.Avatar, result.Avatar)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestChangePassword 测试修改密码
func (suite *UserServiceTestSuite) TestChangePassword() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash of "password"
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	req := &service.ChangePasswordRequest{
		OldPassword: "password",
		NewPassword: "newpassword123",
	}

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Mock 更新用户
	suite.userRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.User")).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 执行修改密码
	err := suite.service.ChangePassword(suite.ctx, userID, req)

	// 验证结果
	require.NoError(suite.T(), err)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestChangePasswordWrongOldPassword 测试修改密码时旧密码错误
func (suite *UserServiceTestSuite) TestChangePasswordWrongOldPassword() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // bcrypt hash of "password"
		Role:      model.UserRoleUser,
		Status:    model.UserStatusActive,
	}

	req := &service.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Mock 日志
	suite.logger.On("Warn", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("uint")).Return()

	// 执行修改密码
	err := suite.service.ChangePassword(suite.ctx, userID, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid old password")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetUsers 测试获取用户列表
func (suite *UserServiceTestSuite) TestGetUsers() {
	users := []*model.User{
		{
			BaseModel: model.BaseModel{ID: 1},
			Username:  "user1",
			Email:     "user1@example.com",
			Role:      model.UserRoleUser,
			Status:    model.UserStatusActive,
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Username:  "user2",
			Email:     "user2@example.com",
			Role:      model.UserRoleAdmin,
			Status:    model.UserStatusActive,
		},
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 获取用户列表
	suite.userRepo.On("List", suite.ctx, opts).Return(users, int64(2), nil)

	// 执行获取用户列表
	result, total, err := suite.service.GetUsers(suite.ctx, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int64(2), total)
	assert.Equal(suite.T(), users[0].ID, result[0].ID)
	assert.Equal(suite.T(), users[1].ID, result[1].ID)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
}

// TestDeleteUser 测试删除用户
func (suite *UserServiceTestSuite) TestDeleteUser() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
	}

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Mock 删除用户
	suite.userRepo.On("Delete", suite.ctx, userID).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 执行删除用户
	err := suite.service.DeleteUser(suite.ctx, userID)

	// 验证结果
	require.NoError(suite.T(), err)

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestDeleteUserError 测试删除用户失败
func (suite *UserServiceTestSuite) TestDeleteUserError() {
	userID := uint(1)
	user := &model.User{
		BaseModel: model.BaseModel{ID: userID},
		Username:  "testuser",
		Email:     "test@example.com",
	}
	expectedError := errors.New("database error")

	// Mock 获取用户
	suite.userRepo.On("GetByID", suite.ctx, userID).Return(user, nil)

	// Mock 删除用户失败
	suite.userRepo.On("Delete", suite.ctx, userID).Return(expectedError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行删除用户
	err := suite.service.DeleteUser(suite.ctx, userID)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to delete user")

	// 验证mock调用
	suite.userRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestUserServiceTestSuite 运行测试套件
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
