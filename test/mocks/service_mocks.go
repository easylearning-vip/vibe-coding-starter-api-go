package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
)

// MockUserService 用户服务模拟
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *service.RegisterRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.LoginResponse), args.Error(1)
}

func (m *MockUserService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(ctx context.Context, userID uint, req *service.UpdateProfileRequest) (*model.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID uint, req *service.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) GetUsers(ctx context.Context, opts repository.ListOptions) ([]*model.User, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockArticleService 文章服务模拟
type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) Create(ctx context.Context, req *service.CreateArticleRequest) (*model.Article, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) Update(ctx context.Context, id uint, req *service.UpdateArticleRequest) (*model.Article, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleService) List(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleService) GetPublished(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Article, int64, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleService) IncrementViewCount(ctx context.Context, articleID uint) error {
	args := m.Called(ctx, articleID)
	return args.Error(0)
}

// MockFileService 文件服务模拟
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) Upload(ctx context.Context, req *service.UploadRequest) (*model.File, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileService) GetByID(ctx context.Context, id uint) (*model.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileService) Download(ctx context.Context, id uint) (*service.DownloadResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DownloadResponse), args.Error(1)
}

func (m *MockFileService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileService) List(ctx context.Context, opts repository.ListOptions) ([]*model.File, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.File), args.Get(1).(int64), args.Error(2)
}

func (m *MockFileService) GetByOwner(ctx context.Context, ownerID uint, opts repository.ListOptions) ([]*model.File, int64, error) {
	args := m.Called(ctx, ownerID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.File), args.Get(1).(int64), args.Error(2)
}

// MockDictService 数据字典服务模拟
type MockDictService struct {
	mock.Mock
}

func (m *MockDictService) GetDictCategories(ctx context.Context) ([]*model.DictCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DictCategory), args.Error(1)
}

func (m *MockDictService) GetDictItems(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	args := m.Called(ctx, categoryCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DictItem), args.Error(1)
}

func (m *MockDictService) GetDictItemByKey(ctx context.Context, categoryCode, itemKey string) (*model.DictItem, error) {
	args := m.Called(ctx, categoryCode, itemKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictItem), args.Error(1)
}

func (m *MockDictService) CreateDictCategory(ctx context.Context, req *service.CreateCategoryRequest) (*model.DictCategory, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictCategory), args.Error(1)
}

func (m *MockDictService) DeleteDictCategory(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDictService) CreateDictItem(ctx context.Context, req *service.CreateItemRequest) (*model.DictItem, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictItem), args.Error(1)
}

func (m *MockDictService) UpdateDictItem(ctx context.Context, id uint, req *service.UpdateItemRequest) (*model.DictItem, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DictItem), args.Error(1)
}

func (m *MockDictService) DeleteDictItem(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDictService) InitDefaultDictData(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDictService) ClearDefaultDictData(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
