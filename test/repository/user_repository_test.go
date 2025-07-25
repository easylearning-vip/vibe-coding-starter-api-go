package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/test/testutil"
)

// UserRepositoryTestSuite 用户仓储测试套件
type UserRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.UserRepository
	ctx    context.Context
}

// SetupSuite 设置测试套件
func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建用户仓储
	suite.repo = repository.NewUserRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
}

// TestCreate 测试创建用户
func (suite *UserRepositoryTestSuite) TestCreate() {
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)
	assert.NotZero(suite.T(), user.CreatedAt)
	assert.NotZero(suite.T(), user.UpdatedAt)
}

// TestCreateDuplicateEmail 测试创建重复邮箱用户
func (suite *UserRepositoryTestSuite) TestCreateDuplicateEmail() {
	user1 := &model.User{
		Username: "testuser1",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User 1",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	user2 := &model.User{
		Username: "testuser2",
		Email:    "test@example.com", // 相同邮箱
		Password: "password123",
		Nickname: "Test User 2",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user1)
	require.NoError(suite.T(), err)

	err = suite.repo.Create(suite.ctx, user2)
	assert.Error(suite.T(), err)
}

// TestGetByID 测试根据ID获取用户
func (suite *UserRepositoryTestSuite) TestGetByID() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 获取用户
	foundUser, err := suite.repo.GetByID(suite.ctx, user.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), user.Username, foundUser.Username)
	assert.Equal(suite.T(), user.Email, foundUser.Email)
}

// TestGetByIDNotFound 测试获取不存在的用户
func (suite *UserRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// TestGetByEmail 测试根据邮箱获取用户
func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 根据邮箱获取用户
	foundUser, err := suite.repo.GetByEmail(suite.ctx, user.Email)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), user.Email, foundUser.Email)
}

// TestGetByEmailNotFound 测试获取不存在邮箱的用户
func (suite *UserRepositoryTestSuite) TestGetByEmailNotFound() {
	_, err := suite.repo.GetByEmail(suite.ctx, "nonexistent@example.com")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// TestGetByUsername 测试根据用户名获取用户
func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 根据用户名获取用户
	foundUser, err := suite.repo.GetByUsername(suite.ctx, user.Username)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), user.Username, foundUser.Username)
}

// TestUpdate 测试更新用户
func (suite *UserRepositoryTestSuite) TestUpdate() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 更新用户信息
	user.Nickname = "Updated User"
	user.Status = model.UserStatusInactive

	err = suite.repo.Update(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedUser, err := suite.repo.GetByID(suite.ctx, user.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated User", updatedUser.Nickname)
	assert.Equal(suite.T(), model.UserStatusInactive, updatedUser.Status)
}

// TestDelete 测试删除用户
func (suite *UserRepositoryTestSuite) TestDelete() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 删除用户
	err = suite.repo.Delete(suite.ctx, user.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, user.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取用户列表
func (suite *UserRepositoryTestSuite) TestList() {
	// 创建多个测试用户
	users := []*model.User{
		{
			Username: "user1",
			Email:    "user1@example.com",
			Password: "password123",
			Nickname: "User 1",
			Role:     model.UserRoleUser,
			Status:   model.UserStatusActive,
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Password: "password123",
			Nickname: "User 2",
			Role:     model.UserRoleAdmin,
			Status:   model.UserStatusActive,
		},
		{
			Username: "user3",
			Email:    "user3@example.com",
			Password: "password123",
			Nickname: "User 3",
			Role:     model.UserRoleUser,
			Status:   model.UserStatusInactive,
		},
	}

	for _, user := range users {
		err := suite.repo.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
	}

	// 测试基本列表查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), result, 3)
}

// TestListWithFilters 测试带过滤器的用户列表
func (suite *UserRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	suite.createTestUsers()

	// 测试角色过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"role": model.UserRoleAdmin,
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), model.UserRoleAdmin, result[0].Role)
}

// TestListWithSearch 测试带搜索的用户列表
func (suite *UserRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	suite.createTestUsers()

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "admin",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestUsers 创建测试用户数据
func (suite *UserRepositoryTestSuite) createTestUsers() {
	users := []*model.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "password123",
			Nickname: "Administrator",
			Role:     model.UserRoleAdmin,
			Status:   model.UserStatusActive,
		},
		{
			Username: "user1",
			Email:    "user1@example.com",
			Password: "password123",
			Nickname: "Regular User",
			Role:     model.UserRoleUser,
			Status:   model.UserStatusActive,
		},
	}

	for _, user := range users {
		err := suite.repo.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
	}
}

// TestUpdateLastLogin 测试更新最后登录时间
func (suite *UserRepositoryTestSuite) TestUpdateLastLogin() {
	// 创建测试用户
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// 更新最后登录时间
	err = suite.repo.UpdateLastLogin(suite.ctx, user.ID)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedUser, err := suite.repo.GetByID(suite.ctx, user.ID)
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser.LastLogin)
	assert.WithinDuration(suite.T(), time.Now(), *updatedUser.LastLogin, time.Minute)
}

// TestPagination 测试分页功能
func (suite *UserRepositoryTestSuite) TestPagination() {
	// 创建10个测试用户
	for i := 0; i < 10; i++ {
		user := &model.User{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password123",
			Nickname: fmt.Sprintf("User %d", i),
			Role:     model.UserRoleUser,
			Status:   model.UserStatusActive,
		}
		err := suite.repo.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
	}

	// 测试第一页
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 3,
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(10), total)
	assert.Len(suite.T(), result, 3)

	// 测试第二页
	opts.Page = 2
	result, total, err = suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(10), total)
	assert.Len(suite.T(), result, 3)
}

// TestUserRepositoryTestSuite 运行测试套件
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
