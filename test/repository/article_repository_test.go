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

// ArticleRepositoryTestSuite 文章仓储测试套件
type ArticleRepositoryTestSuite struct {
	suite.Suite
	db       *testutil.TestDatabase
	cache    *testutil.TestCache
	logger   *testutil.TestLogger
	repo     repository.ArticleRepository
	ctx      context.Context
	author   *model.User
	category *model.Category
	tags     []*model.Tag
}

// SetupSuite 设置测试套件
func (suite *ArticleRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建文章仓储
	suite.repo = repository.NewArticleRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *ArticleRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *ArticleRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
	suite.createTestData()
}

// createTestData 创建测试数据
func (suite *ArticleRepositoryTestSuite) createTestData() {
	// 创建作者
	suite.author = &model.User{
		Username: "author",
		Email:    "author@example.com",
		Password: "password123",
		Nickname: "Article Author",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}
	err := suite.db.GetDB().Create(suite.author).Error
	require.NoError(suite.T(), err)

	// 创建分类
	suite.category = &model.Category{
		Name:        "Technology",
		Slug:        "technology",
		Description: "Technology articles",
	}
	err = suite.db.GetDB().Create(suite.category).Error
	require.NoError(suite.T(), err)

	// 创建标签
	suite.tags = []*model.Tag{
		{
			Name:  "Go",
			Slug:  "go",
			Color: "#00ADD8",
		},
		{
			Name:  "Testing",
			Slug:  "testing",
			Color: "#FF6B6B",
		},
	}
	for _, tag := range suite.tags {
		err = suite.db.GetDB().Create(tag).Error
		require.NoError(suite.T(), err)
	}
}

// TestCreate 测试创建文章
func (suite *ArticleRepositoryTestSuite) TestCreate() {
	article := &model.Article{
		Title:      "Test Article",
		Slug:       "test-article",
		Content:    "This is a test article content.",
		Summary:    "Test article summary",
		Status:     model.ArticleStatusDraft,
		AuthorID:   suite.author.ID,
		CategoryID: &suite.category.ID,
	}

	err := suite.repo.Create(suite.ctx, article)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), article.ID)
	assert.NotZero(suite.T(), article.CreatedAt)
	assert.NotZero(suite.T(), article.UpdatedAt)
}

// TestCreateDuplicateSlug 测试创建重复slug的文章
func (suite *ArticleRepositoryTestSuite) TestCreateDuplicateSlug() {
	article1 := &model.Article{
		Title:    "Test Article 1",
		Slug:     "test-article",
		Content:  "Content 1",
		Status:   model.ArticleStatusDraft,
		AuthorID: suite.author.ID,
	}

	article2 := &model.Article{
		Title:    "Test Article 2",
		Slug:     "test-article", // 相同slug
		Content:  "Content 2",
		Status:   model.ArticleStatusDraft,
		AuthorID: suite.author.ID,
	}

	err := suite.repo.Create(suite.ctx, article1)
	require.NoError(suite.T(), err)

	err = suite.repo.Create(suite.ctx, article2)
	assert.Error(suite.T(), err)
}

// TestGetByID 测试根据ID获取文章
func (suite *ArticleRepositoryTestSuite) TestGetByID() {
	// 创建测试文章
	article := suite.createTestArticle("Test Article", "test-article")

	// 获取文章
	foundArticle, err := suite.repo.GetByID(suite.ctx, article.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.ID, foundArticle.ID)
	assert.Equal(suite.T(), article.Title, foundArticle.Title)
	assert.Equal(suite.T(), article.Slug, foundArticle.Slug)
	assert.NotNil(suite.T(), foundArticle.Author)
	assert.Equal(suite.T(), suite.author.ID, foundArticle.Author.ID)
}

// TestGetByIDNotFound 测试获取不存在的文章
func (suite *ArticleRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "article not found")
}

// TestGetBySlug 测试根据slug获取文章
func (suite *ArticleRepositoryTestSuite) TestGetBySlug() {
	// 创建测试文章
	article := suite.createTestArticle("Test Article", "test-article")

	// 根据slug获取文章
	foundArticle, err := suite.repo.GetBySlug(suite.ctx, article.Slug)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.ID, foundArticle.ID)
	assert.Equal(suite.T(), article.Slug, foundArticle.Slug)
}

// TestGetBySlugNotFound 测试获取不存在slug的文章
func (suite *ArticleRepositoryTestSuite) TestGetBySlugNotFound() {
	_, err := suite.repo.GetBySlug(suite.ctx, "nonexistent-slug")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "article not found")
}

// TestUpdate 测试更新文章
func (suite *ArticleRepositoryTestSuite) TestUpdate() {
	// 创建测试文章
	article := suite.createTestArticle("Test Article", "test-article")

	// 更新文章信息
	article.Title = "Updated Article"
	article.Content = "Updated content"
	article.Status = model.ArticleStatusPublished

	err := suite.repo.Update(suite.ctx, article)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedArticle, err := suite.repo.GetByID(suite.ctx, article.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Article", updatedArticle.Title)
	assert.Equal(suite.T(), "Updated content", updatedArticle.Content)
	assert.Equal(suite.T(), model.ArticleStatusPublished, updatedArticle.Status)
}

// TestDelete 测试删除文章
func (suite *ArticleRepositoryTestSuite) TestDelete() {
	// 创建测试文章
	article := suite.createTestArticle("Test Article", "test-article")

	// 删除文章
	err := suite.repo.Delete(suite.ctx, article.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, article.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取文章列表
func (suite *ArticleRepositoryTestSuite) TestList() {
	// 创建多个测试文章
	articles := []struct {
		title  string
		slug   string
		status string
	}{
		{"Article 1", "article-1", model.ArticleStatusDraft},
		{"Article 2", "article-2", model.ArticleStatusPublished},
		{"Article 3", "article-3", model.ArticleStatusArchived},
	}

	for _, a := range articles {
		article := suite.createTestArticle(a.title, a.slug)
		article.Status = a.status
		err := suite.repo.Update(suite.ctx, article)
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

// TestGetByAuthor 测试根据作者获取文章列表
func (suite *ArticleRepositoryTestSuite) TestGetByAuthor() {
	// 创建测试文章
	suite.createTestArticle("Article 1", "article-1")
	suite.createTestArticle("Article 2", "article-2")

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByAuthor(suite.ctx, suite.author.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)

	// 验证所有文章都属于指定作者
	for _, article := range result {
		assert.Equal(suite.T(), suite.author.ID, article.AuthorID)
	}
}

// TestGetPublished 测试获取已发布文章列表
func (suite *ArticleRepositoryTestSuite) TestGetPublished() {
	// 创建不同状态的文章
	articles := []struct {
		title  string
		slug   string
		status string
	}{
		{"Draft Article", "draft-article", model.ArticleStatusDraft},
		{"Published Article 1", "published-article-1", model.ArticleStatusPublished},
		{"Published Article 2", "published-article-2", model.ArticleStatusPublished},
		{"Archived Article", "archived-article", model.ArticleStatusArchived},
	}

	for _, a := range articles {
		article := suite.createTestArticle(a.title, a.slug)
		article.Status = a.status
		if a.status == model.ArticleStatusPublished {
			now := time.Now()
			article.PublishedAt = &now
		}
		err := suite.repo.Update(suite.ctx, article)
		require.NoError(suite.T(), err)
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetPublished(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)

	// 验证所有文章都是已发布状态
	for _, article := range result {
		assert.Equal(suite.T(), model.ArticleStatusPublished, article.Status)
	}
}

// createTestArticle 创建测试文章
func (suite *ArticleRepositoryTestSuite) createTestArticle(title, slug string) *model.Article {
	article := &model.Article{
		Title:      title,
		Slug:       slug,
		Content:    fmt.Sprintf("Content for %s", title),
		Summary:    fmt.Sprintf("Summary for %s", title),
		Status:     model.ArticleStatusDraft,
		AuthorID:   suite.author.ID,
		CategoryID: &suite.category.ID,
	}

	err := suite.repo.Create(suite.ctx, article)
	require.NoError(suite.T(), err)
	return article
}

// TestSearch 测试搜索文章
func (suite *ArticleRepositoryTestSuite) TestSearch() {
	// 创建测试文章
	articles := []struct {
		title   string
		slug    string
		content string
	}{
		{"Go Programming", "go-programming", "Learn Go programming language"},
		{"JavaScript Basics", "javascript-basics", "Basic JavaScript concepts"},
		{"Go Testing", "go-testing", "Testing in Go programming"},
	}

	for _, a := range articles {
		article := &model.Article{
			Title:    a.title,
			Slug:     a.slug,
			Content:  a.content,
			Status:   model.ArticleStatusPublished,
			AuthorID: suite.author.ID,
		}
		err := suite.repo.Create(suite.ctx, article)
		require.NoError(suite.T(), err)
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// 搜索包含"Go"的文章
	result, total, err := suite.repo.Search(suite.ctx, "Go", opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)
}

// TestIncrementViewCount 测试增加浏览次数
func (suite *ArticleRepositoryTestSuite) TestIncrementViewCount() {
	// 创建测试文章
	article := suite.createTestArticle("Test Article", "test-article")
	initialViewCount := article.ViewCount

	// 增加浏览次数
	err := suite.repo.IncrementViewCount(suite.ctx, article.ID)
	require.NoError(suite.T(), err)

	// 验证浏览次数增加
	updatedArticle, err := suite.repo.GetByID(suite.ctx, article.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), initialViewCount+1, updatedArticle.ViewCount)
}

// TestGetByCategory 测试根据分类获取文章
func (suite *ArticleRepositoryTestSuite) TestGetByCategory() {
	// 创建测试文章
	suite.createTestArticle("Article 1", "article-1")
	suite.createTestArticle("Article 2", "article-2")

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByCategory(suite.ctx, suite.category.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)

	// 验证所有文章都属于指定分类
	for _, article := range result {
		assert.Equal(suite.T(), suite.category.ID, *article.CategoryID)
	}
}

// TestGetByTag 测试根据标签获取文章
func (suite *ArticleRepositoryTestSuite) TestGetByTag() {
	// 创建带标签的文章
	article := suite.createTestArticle("Tagged Article", "tagged-article")

	// 关联标签
	err := suite.db.GetDB().Model(article).Association("Tags").Append(suite.tags[0])
	require.NoError(suite.T(), err)

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByTag(suite.ctx, suite.tags[0].ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), article.ID, result[0].ID)
}

// TestListWithFilters 测试带过滤器的文章列表
func (suite *ArticleRepositoryTestSuite) TestListWithFilters() {
	// 创建不同状态的文章
	suite.createTestArticleWithStatus("Draft Article", "draft-article", model.ArticleStatusDraft)
	suite.createTestArticleWithStatus("Published Article", "published-article", model.ArticleStatusPublished)

	// 测试状态过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"status": model.ArticleStatusPublished,
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), model.ArticleStatusPublished, result[0].Status)
}

// TestListWithSearch 测试带搜索的文章列表
func (suite *ArticleRepositoryTestSuite) TestListWithSearch() {
	// 创建测试文章
	suite.createTestArticle("Go Programming Guide", "go-programming-guide")
	suite.createTestArticle("JavaScript Tutorial", "javascript-tutorial")

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "Go",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestArticleWithStatus 创建指定状态的测试文章
func (suite *ArticleRepositoryTestSuite) createTestArticleWithStatus(title, slug, status string) *model.Article {
	article := &model.Article{
		Title:      title,
		Slug:       slug,
		Content:    fmt.Sprintf("Content for %s", title),
		Summary:    fmt.Sprintf("Summary for %s", title),
		Status:     status,
		AuthorID:   suite.author.ID,
		CategoryID: &suite.category.ID,
	}

	if status == model.ArticleStatusPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	err := suite.repo.Create(suite.ctx, article)
	require.NoError(suite.T(), err)
	return article
}

// TestArticleRepositoryTestSuite 运行测试套件
func TestArticleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleRepositoryTestSuite))
}
