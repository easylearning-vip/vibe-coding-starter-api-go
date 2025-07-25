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

// TagRepositoryTestSuite 标签仓储测试套件
type TagRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.TagRepository
	ctx    context.Context
}

// SetupSuite 设置测试套件
func (suite *TagRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建标签仓储
	suite.repo = repository.NewTagRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *TagRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *TagRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
}

// TestCreate 测试创建标签
func (suite *TagRepositoryTestSuite) TestCreate() {
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), tag.ID)
	assert.NotZero(suite.T(), tag.CreatedAt)
	assert.NotZero(suite.T(), tag.UpdatedAt)
}

// TestCreateDuplicateName 测试创建重复名称标签
func (suite *TagRepositoryTestSuite) TestCreateDuplicateName() {
	tag1 := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	tag2 := &model.Tag{
		Name:        "Go", // 相同名称
		Slug:        "golang",
		Description: "Golang",
		Color:       "#007d9c",
	}

	err := suite.repo.Create(suite.ctx, tag1)
	require.NoError(suite.T(), err)

	err = suite.repo.Create(suite.ctx, tag2)
	assert.Error(suite.T(), err)
}

// TestGetByID 测试根据ID获取标签
func (suite *TagRepositoryTestSuite) TestGetByID() {
	// 创建测试标签
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 获取标签
	foundTag, err := suite.repo.GetByID(suite.ctx, tag.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), tag.ID, foundTag.ID)
	assert.Equal(suite.T(), tag.Name, foundTag.Name)
	assert.Equal(suite.T(), tag.Slug, foundTag.Slug)
}

// TestGetByIDNotFound 测试获取不存在的标签
func (suite *TagRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "tag not found")
}

// TestGetBySlug 测试根据slug获取标签
func (suite *TagRepositoryTestSuite) TestGetBySlug() {
	// 创建测试标签
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 根据slug获取标签
	foundTag, err := suite.repo.GetBySlug(suite.ctx, tag.Slug)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), tag.ID, foundTag.ID)
	assert.Equal(suite.T(), tag.Slug, foundTag.Slug)
}

// TestGetBySlugNotFound 测试获取不存在slug的标签
func (suite *TagRepositoryTestSuite) TestGetBySlugNotFound() {
	_, err := suite.repo.GetBySlug(suite.ctx, "nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "tag not found")
}

// TestGetByName 测试根据名称获取标签
func (suite *TagRepositoryTestSuite) TestGetByName() {
	// 创建测试标签
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 根据名称获取标签
	foundTag, err := suite.repo.GetByName(suite.ctx, tag.Name)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), tag.ID, foundTag.ID)
	assert.Equal(suite.T(), tag.Name, foundTag.Name)
}

// TestGetByNames 测试根据名称列表获取标签
func (suite *TagRepositoryTestSuite) TestGetByNames() {
	// 创建测试标签
	tags := []*model.Tag{
		{
			Name:        "Go",
			Slug:        "go",
			Description: "Go programming language",
			Color:       "#00ADD8",
		},
		{
			Name:        "Python",
			Slug:        "python",
			Description: "Python programming language",
			Color:       "#3776ab",
		},
		{
			Name:        "JavaScript",
			Slug:        "javascript",
			Description: "JavaScript programming language",
			Color:       "#f7df1e",
		},
	}

	for _, tag := range tags {
		err := suite.repo.Create(suite.ctx, tag)
		require.NoError(suite.T(), err)
	}

	// 根据名称列表获取标签
	names := []string{"Go", "Python"}
	foundTags, err := suite.repo.GetByNames(suite.ctx, names)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), foundTags, 2)

	// 验证返回的标签
	tagNames := make([]string, len(foundTags))
	for i, tag := range foundTags {
		tagNames[i] = tag.Name
	}
	assert.Contains(suite.T(), tagNames, "Go")
	assert.Contains(suite.T(), tagNames, "Python")
}

// TestUpdate 测试更新标签
func (suite *TagRepositoryTestSuite) TestUpdate() {
	// 创建测试标签
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 更新标签信息
	tag.Description = "Updated description"
	tag.Color = "#007d9c"

	err = suite.repo.Update(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedTag, err := suite.repo.GetByID(suite.ctx, tag.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated description", updatedTag.Description)
	assert.Equal(suite.T(), "#007d9c", updatedTag.Color)
}

// TestDelete 测试删除标签
func (suite *TagRepositoryTestSuite) TestDelete() {
	// 创建测试标签
	tag := &model.Tag{
		Name:        "Go",
		Slug:        "go",
		Description: "Go programming language",
		Color:       "#00ADD8",
	}

	err := suite.repo.Create(suite.ctx, tag)
	require.NoError(suite.T(), err)

	// 删除标签
	err = suite.repo.Delete(suite.ctx, tag.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, tag.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取标签列表
func (suite *TagRepositoryTestSuite) TestList() {
	// 创建多个测试标签
	tags := []*model.Tag{
		{
			Name:        "Go",
			Slug:        "go",
			Description: "Go programming language",
			Color:       "#00ADD8",
		},
		{
			Name:        "Python",
			Slug:        "python",
			Description: "Python programming language",
			Color:       "#3776ab",
		},
		{
			Name:        "JavaScript",
			Slug:        "javascript",
			Description: "JavaScript programming language",
			Color:       "#f7df1e",
		},
	}

	for _, tag := range tags {
		err := suite.repo.Create(suite.ctx, tag)
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

// TestListWithFilters 测试带过滤器的标签列表
func (suite *TagRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	suite.createTestTags()

	// 测试名称过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"name": "Go",
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "Go", result[0].Name)
}

// TestListWithSearch 测试带搜索的标签列表
func (suite *TagRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	suite.createTestTags()

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "go",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestTags 创建测试标签数据
func (suite *TagRepositoryTestSuite) createTestTags() {
	tags := []*model.Tag{
		{
			Name:        "Go",
			Slug:        "go",
			Description: "Go programming language",
			Color:       "#00ADD8",
		},
		{
			Name:        "Python",
			Slug:        "python",
			Description: "Python programming language",
			Color:       "#3776ab",
		},
	}

	for _, tag := range tags {
		err := suite.repo.Create(suite.ctx, tag)
		require.NoError(suite.T(), err)
	}
}

// TestTagRepositoryTestSuite 运行测试套件
func TestTagRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TagRepositoryTestSuite))
}
