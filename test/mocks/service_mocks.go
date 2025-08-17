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

// MockDepartmentService Department服务模拟
type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) Create(ctx context.Context, req *service.CreateDepartmentRequest) (*model.Department, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentService) GetByID(ctx context.Context, id uint) (*model.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentService) Update(ctx context.Context, id uint, req *service.UpdateDepartmentRequest) (*model.Department, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentService) List(ctx context.Context, opts *service.ListDepartmentOptions) ([]*model.Department, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Department), args.Get(1).(int64), args.Error(2)
}

func (m *MockDepartmentService) GetTree(ctx context.Context) ([]*model.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

func (m *MockDepartmentService) GetChildren(ctx context.Context, parentId uint) ([]*model.Department, error) {
	args := m.Called(ctx, parentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

func (m *MockDepartmentService) GetPath(ctx context.Context, id uint) ([]*model.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

func (m *MockDepartmentService) Move(ctx context.Context, id uint, newParentId uint) error {
	args := m.Called(ctx, id, newParentId)
	return args.Error(0)
}

// MockProductCategoryService ProductCategory服务模拟
type MockProductCategoryService struct {
	mock.Mock
}

func (m *MockProductCategoryService) Create(ctx context.Context, req *service.CreateProductCategoryRequest) (*model.ProductCategory, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryService) GetByID(ctx context.Context, id uint) (*model.ProductCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryService) Update(ctx context.Context, id uint, req *service.UpdateProductCategoryRequest) (*model.ProductCategory, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductCategoryService) List(ctx context.Context, opts *service.ListProductCategoryOptions) ([]*model.ProductCategory, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.ProductCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductCategoryService) GetCategoryTree(ctx context.Context) ([]*service.CategoryTreeNode, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.CategoryTreeNode), args.Error(1)
}

func (m *MockProductCategoryService) GetChildren(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryService) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryService) ValidateParentChild(ctx context.Context, parentID, childID uint) error {
	args := m.Called(ctx, parentID, childID)
	return args.Error(0)
}

func (m *MockProductCategoryService) BatchUpdateSortOrder(ctx context.Context, updates []service.SortOrderUpdate) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockProductCategoryService) CanDelete(ctx context.Context, categoryID uint) (bool, string, error) {
	args := m.Called(ctx, categoryID)
	return args.Bool(0), args.String(1), args.Error(2)
}

// MockProductService Product服务模拟
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(ctx context.Context, req *service.CreateProductRequest) (*model.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) Update(ctx context.Context, id uint, req *service.UpdateProductRequest) (*model.Product, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) List(ctx context.Context, opts *service.ListProductOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) SearchProducts(ctx context.Context, req *service.SearchProductRequest) ([]*model.Product, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductsByCategory(ctx context.Context, categoryID uint, includeSubCategories bool, opts *service.ListProductOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, categoryID, includeSubCategories, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64, opts *service.ListProductOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, minPrice, maxPrice, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetLowStockProducts(ctx context.Context, opts *service.ListProductOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetPopularProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Product), args.Error(1)
}

func (m *MockProductService) BatchUpdatePrices(ctx context.Context, updates []service.PriceUpdate) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockProductService) UpdateProductStatus(ctx context.Context, productID uint, isActive bool) error {
	args := m.Called(ctx, productID, isActive)
	return args.Error(0)
}

func (m *MockProductService) UpdateStock(ctx context.Context, productID uint, quantity int, operation service.StockOperation) error {
	args := m.Called(ctx, productID, quantity, operation)
	return args.Error(0)
}

func (m *MockProductService) CheckStockAvailability(ctx context.Context, productID uint, requiredQuantity int) (bool, error) {
	args := m.Called(ctx, productID, requiredQuantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductService) GetStockAlert(ctx context.Context) ([]*model.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Product), args.Error(1)
}
