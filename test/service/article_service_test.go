package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// ArticleServiceTestSuite 文章服务测试套件
type ArticleServiceTestSuite struct {
	suite.Suite
	articleRepo *mocks.MockArticleRepository
	cache       *mocks.MockCache
	logger      *mocks.MockLogger
	service     service.ArticleService
	ctx         context.Context
}

// SetupSuite 设置测试套件
func (suite *ArticleServiceTestSuite) SetupSuite() {
	suite.articleRepo = new(mocks.MockArticleRepository)
	suite.cache = new(mocks.MockCache)
	suite.logger = new(mocks.MockLogger)
	suite.ctx = context.Background()

	// 创建文章服务
	suite.service = service.NewArticleService(
		suite.articleRepo,
		suite.cache,
		suite.logger,
	)
}

// SetupTest 每个测试前的设置
func (suite *ArticleServiceTestSuite) SetupTest() {
	// 重置所有mock
	suite.articleRepo.ExpectedCalls = nil
	suite.cache.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestCreate 测试创建文章
func (suite *ArticleServiceTestSuite) TestCreate() {
	req := &service.CreateArticleRequest{
		Title:      "Test Article",
		Content:    "This is a test article content.",
		Summary:    "Test article summary",
		CategoryID: func() *uint { id := uint(1); return &id }(),
		Status:     model.ArticleStatusDraft,
	}

	// Mock 创建文章
	suite.articleRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Article")).Return(nil).Run(func(args mock.Arguments) {
		article := args.Get(1).(*model.Article)
		article.ID = 1
	})

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行创建
	article, err := suite.service.Create(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), article)
	assert.Equal(suite.T(), req.Title, article.Title)
	assert.Equal(suite.T(), req.Content, article.Content)
	assert.Equal(suite.T(), req.Summary, article.Summary)
	assert.Equal(suite.T(), req.Status, article.Status)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestCreateWithDefaultStatus 测试创建文章时使用默认状态
func (suite *ArticleServiceTestSuite) TestCreateWithDefaultStatus() {
	req := &service.CreateArticleRequest{
		Title:   "Test Article",
		Content: "This is a test article content.",
		Summary: "Test article summary",
		// Status 为空，应该使用默认值
	}

	// Mock 创建文章
	suite.articleRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Article")).Return(nil).Run(func(args mock.Arguments) {
		article := args.Get(1).(*model.Article)
		article.ID = 1
		// 验证默认状态被设置
		assert.Equal(suite.T(), model.ArticleStatusDraft, article.Status)
	})

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行创建
	article, err := suite.service.Create(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), article)
	assert.Equal(suite.T(), model.ArticleStatusDraft, article.Status)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestCreateError 测试创建文章失败
func (suite *ArticleServiceTestSuite) TestCreateError() {
	req := &service.CreateArticleRequest{
		Title:   "Test Article",
		Content: "This is a test article content.",
		Summary: "Test article summary",
		Status:  model.ArticleStatusDraft,
	}

	expectedError := errors.New("database error")

	// Mock 创建文章失败
	suite.articleRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Article")).Return(expectedError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行创建
	article, err := suite.service.Create(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), article)
	assert.Contains(suite.T(), err.Error(), "failed to create article")

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetByID 测试根据ID获取文章
func (suite *ArticleServiceTestSuite) TestGetByID() {
	articleID := uint(1)
	article := &model.Article{
		BaseModel: model.BaseModel{ID: articleID},
		Title:     "Test Article",
		Content:   "Test content",
		Status:    model.ArticleStatusPublished,
	}

	// Mock 获取文章
	suite.articleRepo.On("GetByID", suite.ctx, articleID).Return(article, nil)

	// 执行获取
	result, err := suite.service.GetByID(suite.ctx, articleID)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), article.ID, result.ID)
	assert.Equal(suite.T(), article.Title, result.Title)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestGetByIDNotFound 测试获取不存在的文章
func (suite *ArticleServiceTestSuite) TestGetByIDNotFound() {
	articleID := uint(999)

	// Mock 文章不存在
	suite.articleRepo.On("GetByID", suite.ctx, articleID).Return(nil, gorm.ErrRecordNotFound)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行获取
	result, err := suite.service.GetByID(suite.ctx, articleID)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetBySlug 测试根据slug获取文章
func (suite *ArticleServiceTestSuite) TestGetBySlug() {
	slug := "test-article"
	article := &model.Article{
		BaseModel: model.BaseModel{ID: 1},
		Title:     "Test Article",
		Slug:      slug,
		Content:   "Test content",
		Status:    model.ArticleStatusPublished,
	}

	// Mock 获取文章
	suite.articleRepo.On("GetBySlug", suite.ctx, slug).Return(article, nil)

	// 执行获取
	result, err := suite.service.GetBySlug(suite.ctx, slug)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), article.ID, result.ID)
	assert.Equal(suite.T(), article.Slug, result.Slug)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestUpdate 测试更新文章
func (suite *ArticleServiceTestSuite) TestUpdate() {
	articleID := uint(1)
	existingArticle := &model.Article{
		BaseModel: model.BaseModel{ID: articleID},
		Title:     "Old Title",
		Content:   "Old content",
		Status:    model.ArticleStatusDraft,
	}

	req := &service.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New content",
		Status:  model.ArticleStatusPublished,
	}

	// Mock 获取现有文章
	suite.articleRepo.On("GetByID", suite.ctx, articleID).Return(existingArticle, nil)

	// Mock 更新文章
	suite.articleRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.Article")).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行更新
	result, err := suite.service.Update(suite.ctx, articleID, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Title, result.Title)
	assert.Equal(suite.T(), req.Content, result.Content)
	assert.Equal(suite.T(), req.Status, result.Status)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestDelete 测试删除文章
func (suite *ArticleServiceTestSuite) TestDelete() {
	articleID := uint(1)
	article := &model.Article{
		BaseModel: model.BaseModel{ID: articleID},
		Title:     "Test Article",
		Content:   "Test content",
		Status:    model.ArticleStatusPublished,
	}

	// Mock 获取文章
	suite.articleRepo.On("GetByID", suite.ctx, articleID).Return(article, nil)

	// Mock 删除文章
	suite.articleRepo.On("Delete", suite.ctx, articleID).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 执行删除
	err := suite.service.Delete(suite.ctx, articleID)

	// 验证结果
	require.NoError(suite.T(), err)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestList 测试获取文章列表
func (suite *ArticleServiceTestSuite) TestList() {
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Article 1",
			Status:    model.ArticleStatusPublished,
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Title:     "Article 2",
			Status:    model.ArticleStatusDraft,
		},
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 获取文章列表
	suite.articleRepo.On("List", suite.ctx, opts).Return(articles, int64(2), nil)

	// 执行获取列表
	result, total, err := suite.service.List(suite.ctx, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int64(2), total)
	assert.Equal(suite.T(), articles[0].ID, result[0].ID)
	assert.Equal(suite.T(), articles[1].ID, result[1].ID)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestGetPublished 测试获取已发布文章列表
func (suite *ArticleServiceTestSuite) TestGetPublished() {
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Published Article 1",
			Status:    model.ArticleStatusPublished,
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Title:     "Published Article 2",
			Status:    model.ArticleStatusPublished,
		},
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 获取已发布文章列表
	suite.articleRepo.On("GetPublished", suite.ctx, opts).Return(articles, int64(2), nil)

	// 执行获取已发布文章列表
	result, total, err := suite.service.GetPublished(suite.ctx, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int64(2), total)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestSearch 测试搜索文章
func (suite *ArticleServiceTestSuite) TestSearch() {
	query := "Go programming"
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Go Programming Guide",
			Content:   "Learn Go programming",
			Status:    model.ArticleStatusPublished,
		},
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 搜索文章
	suite.articleRepo.On("Search", suite.ctx, query, opts).Return(articles, int64(1), nil)

	// 执行搜索
	result, total, err := suite.service.Search(suite.ctx, query, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), int64(1), total)
	assert.Equal(suite.T(), articles[0].ID, result[0].ID)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestSearchEmptyQuery 测试空查询搜索
func (suite *ArticleServiceTestSuite) TestSearchEmptyQuery() {
	query := ""
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Article 1",
			Status:    model.ArticleStatusPublished,
		},
	}

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 获取文章列表（空查询时调用List）
	suite.articleRepo.On("List", suite.ctx, opts).Return(articles, int64(1), nil)

	// 执行搜索
	result, total, err := suite.service.Search(suite.ctx, query, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), int64(1), total)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestIncrementViewCount 测试增加浏览次数
func (suite *ArticleServiceTestSuite) TestIncrementViewCount() {
	articleID := uint(1)

	// Mock 增加浏览次数
	suite.articleRepo.On("IncrementViewCount", suite.ctx, articleID).Return(nil)

	// 执行增加浏览次数
	err := suite.service.IncrementViewCount(suite.ctx, articleID)

	// 验证结果
	require.NoError(suite.T(), err)

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
}

// TestIncrementViewCountError 测试增加浏览次数失败
func (suite *ArticleServiceTestSuite) TestIncrementViewCountError() {
	articleID := uint(1)
	expectedError := errors.New("database error")

	// Mock 增加浏览次数失败
	suite.articleRepo.On("IncrementViewCount", suite.ctx, articleID).Return(expectedError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行增加浏览次数
	err := suite.service.IncrementViewCount(suite.ctx, articleID)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to increment view count")

	// 验证mock调用
	suite.articleRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestArticleServiceTestSuite 运行测试套件
func TestArticleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleServiceTestSuite))
}
