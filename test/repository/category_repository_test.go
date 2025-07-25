package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/test/testutil"
)

// CategoryRepositoryTestSuite 分类仓储测试套件
type CategoryRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.CategoryRepository
	ctx    context.Context
}

// SetupSuite 设置测试套件
func (suite *CategoryRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建分类仓储
	suite.repo = repository.NewCategoryRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *CategoryRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *CategoryRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
}

// TestCreate 测试创建分类
func (suite *CategoryRepositoryTestSuite) TestCreate() {
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
		Icon:        "tech-icon",
		SortOrder:   1,
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), category.ID)
	assert.NotZero(suite.T(), category.CreatedAt)
	assert.NotZero(suite.T(), category.UpdatedAt)
}

// TestCreateDuplicateName 测试创建重复名称分类
func (suite *CategoryRepositoryTestSuite) TestCreateDuplicateName() {
	category1 := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
	}

	category2 := &model.Category{
		Name:        "Technology", // 相同名称
		Slug:        "tech",
		Description: "Tech articles",
		Color:       "#28a745",
	}

	err := suite.repo.Create(suite.ctx, category1)
	require.NoError(suite.T(), err)

	err = suite.repo.Create(suite.ctx, category2)
	assert.Error(suite.T(), err)
}

// TestGetByID 测试根据ID获取分类
func (suite *CategoryRepositoryTestSuite) TestGetByID() {
	// 创建测试分类
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
		SortOrder:   1,
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 获取分类
	foundCategory, err := suite.repo.GetByID(suite.ctx, category.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.ID, foundCategory.ID)
	assert.Equal(suite.T(), category.Name, foundCategory.Name)
	assert.Equal(suite.T(), category.Slug, foundCategory.Slug)
}

// TestGetByIDNotFound 测试获取不存在的分类
func (suite *CategoryRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "category not found")
}

// TestGetBySlug 测试根据slug获取分类
func (suite *CategoryRepositoryTestSuite) TestGetBySlug() {
	// 创建测试分类
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 根据slug获取分类
	foundCategory, err := suite.repo.GetBySlug(suite.ctx, category.Slug)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.ID, foundCategory.ID)
	assert.Equal(suite.T(), category.Slug, foundCategory.Slug)
}

// TestGetBySlugNotFound 测试获取不存在slug的分类
func (suite *CategoryRepositoryTestSuite) TestGetBySlugNotFound() {
	_, err := suite.repo.GetBySlug(suite.ctx, "nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "category not found")
}

// TestGetByName 测试根据名称获取分类
func (suite *CategoryRepositoryTestSuite) TestGetByName() {
	// 创建测试分类
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 根据名称获取分类
	foundCategory, err := suite.repo.GetByName(suite.ctx, category.Name)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.ID, foundCategory.ID)
	assert.Equal(suite.T(), category.Name, foundCategory.Name)
}

// TestUpdate 测试更新分类
func (suite *CategoryRepositoryTestSuite) TestUpdate() {
	// 创建测试分类
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
		SortOrder:   1,
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 更新分类信息
	category.Description = "Updated description"
	category.Color = "#28a745"
	category.SortOrder = 2

	err = suite.repo.Update(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedCategory, err := suite.repo.GetByID(suite.ctx, category.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated description", updatedCategory.Description)
	assert.Equal(suite.T(), "#28a745", updatedCategory.Color)
	assert.Equal(suite.T(), 2, updatedCategory.SortOrder)
}

// TestDelete 测试删除分类
func (suite *CategoryRepositoryTestSuite) TestDelete() {
	// 创建测试分类
	category := &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology related articles",
		Color:       "#007bff",
	}

	err := suite.repo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)

	// 删除分类
	err = suite.repo.Delete(suite.ctx, category.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, category.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取分类列表
func (suite *CategoryRepositoryTestSuite) TestList() {
	// 创建多个测试分类
	categories := []*model.Category{
		{
			Name:        "Technology",
			Slug:        "technology",
			Description: "Technology articles",
			Color:       "#007bff",
			SortOrder:   1,
		},
		{
			Name:        "Programming",
			Slug:        "programming",
			Description: "Programming tutorials",
			Color:       "#28a745",
			SortOrder:   2,
		},
		{
			Name:        "Design",
			Slug:        "design",
			Description: "Design resources",
			Color:       "#ffc107",
			SortOrder:   3,
		},
	}

	for _, category := range categories {
		err := suite.repo.Create(suite.ctx, category)
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

// TestListWithFilters 测试带过滤器的分类列表
func (suite *CategoryRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	suite.createTestCategories()

	// 测试名称过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"name": "Technology",
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Technology", result[0].Name)
}

// TestListWithSearch 测试带搜索的分类列表
func (suite *CategoryRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	suite.createTestCategories()

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "tech",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestCategories 创建测试分类数据
func (suite *CategoryRepositoryTestSuite) createTestCategories() {
	categories := []*model.Category{
		{
			Name:        "Technology",
			Slug:        "technology",
			Description: "Technology articles",
			Color:       "#007bff",
			SortOrder:   1,
		},
		{
			Name:        "Programming",
			Slug:        "programming",
			Description: "Programming tutorials",
			Color:       "#28a745",
			SortOrder:   2,
		},
	}

	for _, category := range categories {
		err := suite.repo.Create(suite.ctx, category)
		require.NoError(suite.T(), err)
	}
}

// TestCategoryRepositoryTestSuite 运行测试套件
func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}
