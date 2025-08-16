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

// CommentRepositoryTestSuite 评论仓储测试套件
type CommentRepositoryTestSuite struct {
	suite.Suite
	db          *testutil.TestDatabase
	cache       *testutil.TestCache
	logger      *testutil.TestLogger
	repo        repository.CommentRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	ctx         context.Context
	testUser    *model.User
	testArticle *model.Article
}

// SetupSuite 设置测试套件
func (suite *CommentRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建仓储
	suite.repo = repository.NewCommentRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
	suite.userRepo = repository.NewUserRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
	suite.articleRepo = repository.NewArticleRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *CommentRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *CommentRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
	suite.createTestData()
}

// createTestData 创建测试数据
func (suite *CommentRepositoryTestSuite) createTestData() {
	// 创建测试用户
	suite.testUser = &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Nickname: "Test User",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}
	err := suite.userRepo.Create(suite.ctx, suite.testUser)
	require.NoError(suite.T(), err)

	// 创建测试文章
	suite.testArticle = &model.Article{
		Title:    "Test Article",
		Slug:     "test-article",
		Content:  "This is a test article content",
		Summary:  "Test summary",
		Status:   model.ArticleStatusPublished,
		AuthorID: suite.testUser.ID,
	}
	err = suite.articleRepo.Create(suite.ctx, suite.testArticle)
	require.NoError(suite.T(), err)
}

// TestCreate 测试创建评论
func (suite *CommentRepositoryTestSuite) TestCreate() {
	comment := &model.Comment{
		Content:   "This is a test comment",
		Status:    model.CommentStatusPending,
		ArticleID: suite.testArticle.ID,
		AuthorID:  suite.testUser.ID,
	}

	err := suite.repo.Create(suite.ctx, comment)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), comment.ID)
	assert.NotZero(suite.T(), comment.CreatedAt)
	assert.NotZero(suite.T(), comment.UpdatedAt)
}

// TestGetByID 测试根据ID获取评论
func (suite *CommentRepositoryTestSuite) TestGetByID() {
	// 创建测试评论
	comment := &model.Comment{
		Content:   "This is a test comment",
		Status:    model.CommentStatusPending,
		ArticleID: suite.testArticle.ID,
		AuthorID:  suite.testUser.ID,
	}

	err := suite.repo.Create(suite.ctx, comment)
	require.NoError(suite.T(), err)

	// 获取评论
	foundComment, err := suite.repo.GetByID(suite.ctx, comment.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), comment.ID, foundComment.ID)
	assert.Equal(suite.T(), comment.Content, foundComment.Content)
	assert.Equal(suite.T(), comment.ArticleID, foundComment.ArticleID)
	assert.Equal(suite.T(), comment.AuthorID, foundComment.AuthorID)

	// 验证预加载的关联数据
	assert.NotNil(suite.T(), foundComment.Author)
	assert.Equal(suite.T(), suite.testUser.ID, foundComment.Author.ID)
	assert.NotNil(suite.T(), foundComment.Article)
	assert.Equal(suite.T(), suite.testArticle.ID, foundComment.Article.ID)
}

// TestGetByIDNotFound 测试获取不存在的评论
func (suite *CommentRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "comment not found")
}

// TestUpdate 测试更新评论
func (suite *CommentRepositoryTestSuite) TestUpdate() {
	// 创建测试评论
	comment := &model.Comment{
		Content:   "This is a test comment",
		Status:    model.CommentStatusPending,
		ArticleID: suite.testArticle.ID,
		AuthorID:  suite.testUser.ID,
	}

	err := suite.repo.Create(suite.ctx, comment)
	require.NoError(suite.T(), err)

	// 更新评论信息
	comment.Content = "Updated comment content"
	comment.Status = model.CommentStatusApproved

	err = suite.repo.Update(suite.ctx, comment)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedComment, err := suite.repo.GetByID(suite.ctx, comment.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated comment content", updatedComment.Content)
	assert.Equal(suite.T(), model.CommentStatusApproved, updatedComment.Status)
}

// TestDelete 测试删除评论
func (suite *CommentRepositoryTestSuite) TestDelete() {
	// 创建测试评论
	comment := &model.Comment{
		Content:   "This is a test comment",
		Status:    model.CommentStatusPending,
		ArticleID: suite.testArticle.ID,
		AuthorID:  suite.testUser.ID,
	}

	err := suite.repo.Create(suite.ctx, comment)
	require.NoError(suite.T(), err)

	// 删除评论
	err = suite.repo.Delete(suite.ctx, comment.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, comment.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取评论列表
func (suite *CommentRepositoryTestSuite) TestList() {
	// 创建多个测试评论
	comments := []*model.Comment{
		{
			Content:   "First comment",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
		{
			Content:   "Second comment",
			Status:    model.CommentStatusPending,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
		{
			Content:   "Third comment",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
	}

	for _, comment := range comments {
		err := suite.repo.Create(suite.ctx, comment)
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

// TestGetByArticle 测试根据文章获取评论列表
func (suite *CommentRepositoryTestSuite) TestGetByArticle() {
	// 创建测试评论
	suite.createTestComments()

	// 测试根据文章获取评论
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByArticle(suite.ctx, suite.testArticle.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), result, 3)

	// 验证所有评论都属于指定文章
	for _, comment := range result {
		assert.Equal(suite.T(), suite.testArticle.ID, comment.ArticleID)
	}
}

// TestGetByAuthor 测试根据作者获取评论列表
func (suite *CommentRepositoryTestSuite) TestGetByAuthor() {
	// 创建测试评论
	suite.createTestComments()

	// 测试根据作者获取评论
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByAuthor(suite.ctx, suite.testUser.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), result, 3)

	// 验证所有评论都属于指定作者
	for _, comment := range result {
		assert.Equal(suite.T(), suite.testUser.ID, comment.AuthorID)
	}
}

// TestGetReplies 测试获取评论回复
func (suite *CommentRepositoryTestSuite) TestGetReplies() {
	// 创建父评论
	parentComment := &model.Comment{
		Content:   "Parent comment",
		Status:    model.CommentStatusApproved,
		ArticleID: suite.testArticle.ID,
		AuthorID:  suite.testUser.ID,
	}
	err := suite.repo.Create(suite.ctx, parentComment)
	require.NoError(suite.T(), err)

	// 创建回复评论
	replies := []*model.Comment{
		{
			Content:   "First reply",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
			ParentID:  &parentComment.ID,
		},
		{
			Content:   "Second reply",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
			ParentID:  &parentComment.ID,
		},
	}

	for _, reply := range replies {
		err := suite.repo.Create(suite.ctx, reply)
		require.NoError(suite.T(), err)
	}

	// 测试获取回复
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetReplies(suite.ctx, parentComment.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)

	// 验证所有回复都属于指定父评论
	for _, reply := range result {
		assert.NotNil(suite.T(), reply.ParentID)
		assert.Equal(suite.T(), parentComment.ID, *reply.ParentID)
	}
}

// TestListWithFilters 测试带过滤器的评论列表
func (suite *CommentRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	suite.createTestComments()

	// 测试状态过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"status": model.CommentStatusApproved,
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)

	// 验证所有评论都是已批准状态
	for _, comment := range result {
		assert.Equal(suite.T(), model.CommentStatusApproved, comment.Status)
	}
}

// TestListWithSearch 测试带搜索的评论列表
func (suite *CommentRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	suite.createTestComments()

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "first",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestComments 创建测试评论数据
func (suite *CommentRepositoryTestSuite) createTestComments() {
	comments := []*model.Comment{
		{
			Content:   "First comment",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
		{
			Content:   "Second comment",
			Status:    model.CommentStatusPending,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
		{
			Content:   "Third comment",
			Status:    model.CommentStatusApproved,
			ArticleID: suite.testArticle.ID,
			AuthorID:  suite.testUser.ID,
		},
	}

	for _, comment := range comments {
		err := suite.repo.Create(suite.ctx, comment)
		require.NoError(suite.T(), err)
	}
}

// TestCommentRepositoryTestSuite 运行测试套件
func TestCommentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}
