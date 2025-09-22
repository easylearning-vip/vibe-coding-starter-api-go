package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/test/testutil"

	"github.com/stretchr/testify/suite"
)

// TestRepositoryTestSuite Test仓储测试套件
type TestRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.TestRepository
	ctx    context.Context
}

// SetupSuite 设置测试套件
func (suite *TestRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建仓储实例
	suite.repo = repository.NewTestRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *TestRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
}

// SetupTest 每个测试前的设置
func (suite *TestRepositoryTestSuite) SetupTest() {
	// 清理数据
	suite.db.Clean(suite.T())
}

// TestCreate 测试创建Test
func (suite *TestRepositoryTestSuite) TestCreate() {
	// 准备测试数据
	entity := &model.Test{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}

	// 执行创建
	err := suite.repo.Create(suite.ctx, entity)

	// 验证结果
	suite.NoError(err)
	suite.NotZero(entity.ID)
	suite.Equal("Test Test", entity.Name)
	suite.Equal("Test Description", entity.Description.String)
	suite.True(entity.Description.Valid)
	suite.NotZero(entity.CreatedAt)
	suite.NotZero(entity.UpdatedAt)
}

// TestCreateDuplicateName 测试创建重复名称的Test
func (suite *TestRepositoryTestSuite) TestCreateDuplicateName() {
	// 创建第一个Test
	entity1 := &model.Test{
		Name:        "Duplicate Name",
		Description: sql.NullString{String: "First Description", Valid: true},
	}
	err := suite.repo.Create(suite.ctx, entity1)
	suite.NoError(err)

	// 尝试创建重复名称的Test
	entity2 := &model.Test{
		Name:        "Duplicate Name",
		Description: sql.NullString{String: "Second Description", Valid: true},
	}
	err = suite.repo.Create(suite.ctx, entity2)
	suite.NoError(err) // Test model doesn't have unique constraint on name
}

// TestGetByID 测试根据ID获取Test
func (suite *TestRepositoryTestSuite) TestGetByID() {
	// 创建测试数据
	entity := &model.Test{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 根据ID获取
	found, err := suite.repo.GetByID(suite.ctx, entity.ID)

	// 验证结果
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(entity.ID, found.ID)
	suite.Equal(entity.Name, found.Name)
	suite.Equal(entity.Description, found.Description)
}

// TestGetByIDNotFound 测试获取不存在的Test
func (suite *TestRepositoryTestSuite) TestGetByIDNotFound() {
	// 尝试获取不存在的Test
	found, err := suite.repo.GetByID(suite.ctx, 999)

	// 验证结果
	suite.Error(err)
	suite.Nil(found)
	suite.Contains(err.Error(), "not found")
}

// TestGetByName 测试根据名称获取Test
func (suite *TestRepositoryTestSuite) TestGetByName() {
	// 创建测试数据
	entity := &model.Test{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 根据名称获取
	found, err := suite.repo.GetByName(suite.ctx, "Test Test")

	// 验证结果
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(entity.ID, found.ID)
	suite.Equal(entity.Name, found.Name)
}

// TestUpdate 测试更新Test
func (suite *TestRepositoryTestSuite) TestUpdate() {
	// 创建测试数据
	entity := &model.Test{
		Name:        "Original Name",
		Description: sql.NullString{String: "Original Description", Valid: true},
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 更新数据
	entity.Name = "Updated Name"
	entity.Description = sql.NullString{String: "Updated Description", Valid: true}
	err = suite.repo.Update(suite.ctx, entity)

	// 验证结果
	suite.NoError(err)

	// 重新获取验证
	updated, err := suite.repo.GetByID(suite.ctx, entity.ID)
	suite.NoError(err)
	suite.Equal("Updated Name", updated.Name)
	suite.Equal("Updated Description", updated.Description.String)
	suite.True(updated.Description.Valid)
}

// TestDelete 测试删除Test
func (suite *TestRepositoryTestSuite) TestDelete() {
	// 创建测试数据
	entity := &model.Test{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 删除数据
	err = suite.repo.Delete(suite.ctx, entity.ID)
	suite.NoError(err)

	// 验证已删除
	found, err := suite.repo.GetByID(suite.ctx, entity.ID)
	suite.Error(err)
	suite.Nil(found)
}

// TestList 测试获取Test列表
func (suite *TestRepositoryTestSuite) TestList() {
	// 创建测试数据
	entities := []*model.Test{
		{Name: "Test 1", Description: sql.NullString{String: "Description 1", Valid: true}},
		{Name: "Test 2", Description: sql.NullString{String: "Description 2", Valid: true}},
		{Name: "Test 3", Description: sql.NullString{String: "Description 3", Valid: true}},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 获取列表
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 3)
	suite.Equal(int64(3), total)
}

// TestListWithFilters 测试带过滤器的列表查询
func (suite *TestRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	entities := []*model.Test{
		{Name: "Active Test", Description: sql.NullString{String: "Active Description", Valid: true}},
		{Name: "Inactive Test", Description: sql.NullString{String: "Inactive Description", Valid: true}},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 使用过滤器查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"name": "Active Test",
		},
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(1), total)
	suite.Equal("Active Test", result[0].Name)
}

// TestListWithSearch 测试带搜索的列表查询
func (suite *TestRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	entities := []*model.Test{
		{Name: "Searchable Test", Description: sql.NullString{String: "Description", Valid: true}},
		{Name: "Another Test", Description: sql.NullString{String: "Description", Valid: true}},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 使用搜索查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "Searchable",
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(1), total)
	suite.Contains(result[0].Name, "Searchable")
}

// TestPagination 测试分页
func (suite *TestRepositoryTestSuite) TestPagination() {
	// 创建测试数据
	for i := 1; i <= 5; i++ {
		entity := &model.Test{
			Name:        fmt.Sprintf("Test %d", i),
			Description: sql.NullString{String: fmt.Sprintf("Description %d", i), Valid: true},
		}
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 第一页
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 2,
	}
	result, total, err := suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(int64(5), total)

	// 第二页
	opts.Page = 2
	result, total, err = suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(int64(5), total)

	// 第三页
	opts.Page = 3
	result, total, err = suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(5), total)
}

// TestTestRepositoryTestSuite 运行测试套件
func TestTestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TestRepositoryTestSuite))
}
