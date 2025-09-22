package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
)

// MockUserRepository 用户仓储模拟
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.User, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockArticleRepository 文章仓储模拟
type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *model.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleRepository) Update(ctx context.Context, article *model.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleRepository) GetByAuthor(ctx context.Context, authorID uint, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, authorID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetByCategory(ctx context.Context, categoryID uint, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, categoryID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetByTag(ctx context.Context, tagID uint, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, tagID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetPublished(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) IncrementViewCount(ctx context.Context, articleID uint) error {
	args := m.Called(ctx, articleID)
	return args.Error(0)
}

// MockFileRepository 文件仓储模拟
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *model.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(ctx context.Context, id uint) (*model.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileRepository) Update(ctx context.Context, file *model.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.File, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.File), args.Get(1).(int64), args.Error(2)
}

func (m *MockFileRepository) GetByHash(ctx context.Context, hash string) (*model.File, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileRepository) GetByOwner(ctx context.Context, ownerID uint, opts repository.ListOptions) ([]*model.File, int64, error) {
	args := m.Called(ctx, ownerID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.File), args.Get(1).(int64), args.Error(2)
}

// MockDictRepository 数据字典仓储模拟
type MockDictRepository struct {
	mock.Mock
}

func (m *MockDictRepository) GetCategoryByCode(ctx context.Context, code string) (*model.DictCategory, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictCategory), args.Error(1)
}

func (m *MockDictRepository) CreateCategory(ctx context.Context, category *model.DictCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockDictRepository) GetItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	args := m.Called(ctx, categoryCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DictItem), args.Error(1)
}

func (m *MockDictRepository) GetActiveItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	args := m.Called(ctx, categoryCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DictItem), args.Error(1)
}

func (m *MockDictRepository) GetItemByID(ctx context.Context, id uint) (*model.DictItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictItem), args.Error(1)
}

func (m *MockDictRepository) CreateItem(ctx context.Context, item *model.DictItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockDictRepository) UpdateItem(ctx context.Context, item *model.DictItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockDictRepository) DeleteItem(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDictRepository) GetAllCategories(ctx context.Context) ([]*model.DictCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DictCategory), args.Error(1)
}

func (m *MockDictRepository) DeleteCategory(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockDepartmentRepository Department仓储模拟
type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) Create(ctx context.Context, department *model.Department) error {
	args := m.Called(ctx, department)
	return args.Error(0)
}

func (m *MockDepartmentRepository) GetByID(ctx context.Context, id uint) (*model.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	args := m.Called(ctx, department)
	return args.Error(0)
}

func (m *MockDepartmentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Department, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Department), args.Get(1).(int64), args.Error(2)
}

func (m *MockDepartmentRepository) GetByName(ctx context.Context, name string) (*model.Department, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetByParentId(ctx context.Context, parentId uint) ([]*model.Department, error) {
	args := m.Called(ctx, parentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetByCode(ctx context.Context, code string) (*model.Department, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetChildrenTree(ctx context.Context, parentId uint) ([]*model.Department, error) {
	args := m.Called(ctx, parentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

// MockDifficultyLevelRepository DifficultyLevel仓储模拟
type MockDifficultyLevelRepository struct {
	mock.Mock
}

func (m *MockDifficultyLevelRepository) Create(ctx context.Context, difficultyLevel *model.DifficultyLevel) error {
	args := m.Called(ctx, difficultyLevel)
	return args.Error(0)
}

func (m *MockDifficultyLevelRepository) GetByID(ctx context.Context, id uint) (*model.DifficultyLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DifficultyLevel), args.Error(1)
}

func (m *MockDifficultyLevelRepository) Update(ctx context.Context, difficultyLevel *model.DifficultyLevel) error {
	args := m.Called(ctx, difficultyLevel)
	return args.Error(0)
}

func (m *MockDifficultyLevelRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDifficultyLevelRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.DifficultyLevel, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.DifficultyLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockDifficultyLevelRepository) GetByName(ctx context.Context, name string) (*model.DifficultyLevel, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DifficultyLevel), args.Error(1)
}

// MockScenarioRepository Scenario仓储模拟
type MockScenarioRepository struct {
	mock.Mock
}

func (m *MockScenarioRepository) Create(ctx context.Context, scenario *model.Scenario) error {
	args := m.Called(ctx, scenario)
	return args.Error(0)
}

func (m *MockScenarioRepository) GetByID(ctx context.Context, id uint) (*model.Scenario, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Scenario), args.Error(1)
}

func (m *MockScenarioRepository) Update(ctx context.Context, scenario *model.Scenario) error {
	args := m.Called(ctx, scenario)
	return args.Error(0)
}

func (m *MockScenarioRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockScenarioRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Scenario, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Scenario), args.Get(1).(int64), args.Error(2)
}

func (m *MockScenarioRepository) GetByName(ctx context.Context, name string) (*model.Scenario, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Scenario), args.Error(1)
}

// MockTestRepository Test仓储模拟
type MockTestRepository struct {
	mock.Mock
}

func (m *MockTestRepository) Create(ctx context.Context, test *model.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *MockTestRepository) GetByID(ctx context.Context, id uint) (*model.Test, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Test), args.Error(1)
}

func (m *MockTestRepository) Update(ctx context.Context, test *model.Test) error {
	args := m.Called(ctx, test)
	return args.Error(0)
}

func (m *MockTestRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTestRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Test, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Test), args.Get(1).(int64), args.Error(2)
}

func (m *MockTestRepository) GetByName(ctx context.Context, name string) (*model.Test, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Test), args.Error(1)
}

// MockPart2QuestionRepository Part2Question仓储模拟
type MockPart2QuestionRepository struct {
	mock.Mock
}

func (m *MockPart2QuestionRepository) Create(ctx context.Context, part2Question *model.Part2Question) error {
	args := m.Called(ctx, part2Question)
	return args.Error(0)
}

func (m *MockPart2QuestionRepository) GetByID(ctx context.Context, id uint) (*model.Part2Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part2Question), args.Error(1)
}

func (m *MockPart2QuestionRepository) Update(ctx context.Context, part2Question *model.Part2Question) error {
	args := m.Called(ctx, part2Question)
	return args.Error(0)
}

func (m *MockPart2QuestionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPart2QuestionRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Part2Question, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Part2Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockPart2QuestionRepository) GetByName(ctx context.Context, name string) (*model.Part2Question, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part2Question), args.Error(1)
}

// MockPart3ConversationRepository Part3Conversation仓储模拟
type MockPart3ConversationRepository struct {
	mock.Mock
}

func (m *MockPart3ConversationRepository) Create(ctx context.Context, part3Conversation *model.Part3Conversation) error {
	args := m.Called(ctx, part3Conversation)
	return args.Error(0)
}

func (m *MockPart3ConversationRepository) GetByID(ctx context.Context, id uint) (*model.Part3Conversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part3Conversation), args.Error(1)
}

func (m *MockPart3ConversationRepository) Update(ctx context.Context, part3Conversation *model.Part3Conversation) error {
	args := m.Called(ctx, part3Conversation)
	return args.Error(0)
}

func (m *MockPart3ConversationRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPart3ConversationRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Part3Conversation, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Part3Conversation), args.Get(1).(int64), args.Error(2)
}

func (m *MockPart3ConversationRepository) GetByName(ctx context.Context, name string) (*model.Part3Conversation, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part3Conversation), args.Error(1)
}

// MockPart3AnswerOptionRepository Part3AnswerOption仓储模拟
type MockPart3AnswerOptionRepository struct {
	mock.Mock
}

func (m *MockPart3AnswerOptionRepository) Create(ctx context.Context, part3AnswerOption *model.Part3AnswerOption) error {
	args := m.Called(ctx, part3AnswerOption)
	return args.Error(0)
}

func (m *MockPart3AnswerOptionRepository) GetByID(ctx context.Context, id uint) (*model.Part3AnswerOption, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part3AnswerOption), args.Error(1)
}

func (m *MockPart3AnswerOptionRepository) Update(ctx context.Context, part3AnswerOption *model.Part3AnswerOption) error {
	args := m.Called(ctx, part3AnswerOption)
	return args.Error(0)
}

func (m *MockPart3AnswerOptionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPart3AnswerOptionRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Part3AnswerOption, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Part3AnswerOption), args.Get(1).(int64), args.Error(2)
}

func (m *MockPart3AnswerOptionRepository) GetByName(ctx context.Context, name string) (*model.Part3AnswerOption, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part3AnswerOption), args.Error(1)
}

// MockPart4TalkRepository Part4Talk仓储模拟
type MockPart4TalkRepository struct {
	mock.Mock
}

func (m *MockPart4TalkRepository) Create(ctx context.Context, part4Talk *model.Part4Talk) error {
	args := m.Called(ctx, part4Talk)
	return args.Error(0)
}

func (m *MockPart4TalkRepository) GetByID(ctx context.Context, id uint) (*model.Part4Talk, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part4Talk), args.Error(1)
}

func (m *MockPart4TalkRepository) Update(ctx context.Context, part4Talk *model.Part4Talk) error {
	args := m.Called(ctx, part4Talk)
	return args.Error(0)
}

func (m *MockPart4TalkRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPart4TalkRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Part4Talk, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Part4Talk), args.Get(1).(int64), args.Error(2)
}

func (m *MockPart4TalkRepository) GetByName(ctx context.Context, name string) (*model.Part4Talk, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part4Talk), args.Error(1)
}

// MockPart4AnswerOptionRepository Part4AnswerOption仓储模拟
type MockPart4AnswerOptionRepository struct {
	mock.Mock
}

func (m *MockPart4AnswerOptionRepository) Create(ctx context.Context, part4AnswerOption *model.Part4AnswerOption) error {
	args := m.Called(ctx, part4AnswerOption)
	return args.Error(0)
}

func (m *MockPart4AnswerOptionRepository) GetByID(ctx context.Context, id uint) (*model.Part4AnswerOption, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part4AnswerOption), args.Error(1)
}

func (m *MockPart4AnswerOptionRepository) Update(ctx context.Context, part4AnswerOption *model.Part4AnswerOption) error {
	args := m.Called(ctx, part4AnswerOption)
	return args.Error(0)
}

func (m *MockPart4AnswerOptionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPart4AnswerOptionRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Part4AnswerOption, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Part4AnswerOption), args.Get(1).(int64), args.Error(2)
}

func (m *MockPart4AnswerOptionRepository) GetByName(ctx context.Context, name string) (*model.Part4AnswerOption, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Part4AnswerOption), args.Error(1)
}
