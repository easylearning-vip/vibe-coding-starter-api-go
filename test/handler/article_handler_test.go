package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// ArticleHandlerTestSuite 文章处理器测试套件
type ArticleHandlerTestSuite struct {
	suite.Suite
	articleService *mocks.MockArticleService
	logger         *mocks.MockLogger
	handler        *handler.ArticleHandler
	router         *gin.Engine
}

// SetupSuite 设置测试套件
func (suite *ArticleHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	suite.articleService = new(mocks.MockArticleService)
	suite.logger = new(mocks.MockLogger)

	// 创建文章处理器
	suite.handler = handler.NewArticleHandler(
		suite.articleService,
		suite.logger,
	)

	// 设置路由
	suite.router = gin.New()
	api := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(api)
}

// SetupTest 每个测试前的设置
func (suite *ArticleHandlerTestSuite) SetupTest() {
	suite.articleService.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestCreate 测试创建文章
func (suite *ArticleHandlerTestSuite) TestCreate() {
	reqBody := service.CreateArticleRequest{
		Title:      "Test Article",
		Content:    "This is a test article content.",
		Summary:    "Test article summary",
		CategoryID: func() *uint { id := uint(1); return &id }(),
		Status:     model.ArticleStatusDraft,
	}

	article := &model.Article{
		BaseModel:  model.BaseModel{ID: 1},
		Title:      reqBody.Title,
		Content:    reqBody.Content,
		Summary:    reqBody.Summary,
		CategoryID: reqBody.CategoryID,
		Status:     reqBody.Status,
	}

	// Mock 文章服务
	suite.articleService.On("Create", mock.Anything, &reqBody).Return(article, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response model.Article
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.ID, response.ID)
	assert.Equal(suite.T(), article.Title, response.Title)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestCreateInvalidJSON 测试创建文章时JSON格式错误
func (suite *ArticleHandlerTestSuite) TestCreateInvalidJSON() {
	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 创建无效JSON请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBufferString("invalid json"))
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

// TestGetByID 测试根据ID获取文章
func (suite *ArticleHandlerTestSuite) TestGetByID() {
	articleID := uint(1)
	article := &model.Article{
		BaseModel: model.BaseModel{ID: articleID},
		Title:     "Test Article",
		Content:   "Test content",
		Status:    model.ArticleStatusPublished,
		ViewCount: 5,
	}

	// Mock 文章服务
	suite.articleService.On("GetByID", mock.Anything, articleID).Return(article, nil)
	suite.articleService.On("IncrementViewCount", mock.Anything, articleID).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+strconv.Itoa(int(articleID)), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.Article
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.ID, response.ID)
	assert.Equal(suite.T(), article.Title, response.Title)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestGetByIDNotFound 测试获取不存在的文章
func (suite *ArticleHandlerTestSuite) TestGetByIDNotFound() {
	articleID := uint(999)

	// Mock 文章服务返回错误
	suite.articleService.On("GetByID", mock.Anything, articleID).Return(nil, assert.AnError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+strconv.Itoa(int(articleID)), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "article_not_found", response.Error)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestList 测试获取文章列表
func (suite *ArticleHandlerTestSuite) TestList() {
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

	// Mock 文章服务
	suite.articleService.On("List", mock.Anything, mock.AnythingOfType("repository.ListOptions")).Return(articles, int64(2), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.ListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), response.Total)
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 10, response.Size)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestListWithStatusFilter 测试带状态过滤的文章列表
func (suite *ArticleHandlerTestSuite) TestListWithStatusFilter() {
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Published Article",
			Status:    model.ArticleStatusPublished,
		},
	}

	// Mock 文章服务
	suite.articleService.On("List", mock.Anything, mock.MatchedBy(func(opts repository.ListOptions) bool {
		return opts.Filters != nil && opts.Filters["status"] == model.ArticleStatusPublished
	})).Return(articles, int64(1), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles?status=published", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.ListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), response.Total)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestUpdate 测试更新文章
func (suite *ArticleHandlerTestSuite) TestUpdate() {
	articleID := uint(1)
	reqBody := service.UpdateArticleRequest{
		Title:   "Updated Title",
		Content: "Updated content",
		Status:  model.ArticleStatusPublished,
	}

	updatedArticle := &model.Article{
		BaseModel: model.BaseModel{ID: articleID},
		Title:     reqBody.Title,
		Content:   reqBody.Content,
		Status:    reqBody.Status,
	}

	// Mock 文章服务
	suite.articleService.On("Update", mock.Anything, articleID, &reqBody).Return(updatedArticle, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/"+strconv.Itoa(int(articleID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.Article
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedArticle.ID, response.ID)
	assert.Equal(suite.T(), updatedArticle.Title, response.Title)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestDelete 测试删除文章
func (suite *ArticleHandlerTestSuite) TestDelete() {
	articleID := uint(1)

	// Mock 文章服务
	suite.articleService.On("Delete", mock.Anything, articleID).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/"+strconv.Itoa(int(articleID)), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Article deleted successfully", response.Message)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestSearch 测试搜索文章
func (suite *ArticleHandlerTestSuite) TestSearch() {
	query := "Go programming"
	articles := []*model.Article{
		{
			BaseModel: model.BaseModel{ID: 1},
			Title:     "Go Programming Guide",
			Content:   "Learn Go programming",
			Status:    model.ArticleStatusPublished,
		},
	}

	// Mock 文章服务
	suite.articleService.On("Search", mock.Anything, query, mock.AnythingOfType("repository.ListOptions")).Return(articles, int64(1), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/search?q="+url.QueryEscape(query), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.ListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), response.Total)

	// 验证mock调用
	suite.articleService.AssertExpectations(suite.T())
}

// TestArticleHandlerTestSuite 运行测试套件
func TestArticleHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}
