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

func (m *MockDepartmentRepository) GetByCode(ctx context.Context, code string) (*model.Department, error) {
	args := m.Called(ctx, code)
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

func (m *MockDepartmentRepository) GetChildrenTree(ctx context.Context, parentId uint) ([]*model.Department, error) {
	args := m.Called(ctx, parentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Department), args.Error(1)
}

// MockProductCategoryRepository ProductCategory仓储模拟
type MockProductCategoryRepository struct {
	mock.Mock
}

func (m *MockProductCategoryRepository) Create(ctx context.Context, productCategory *model.ProductCategory) error {
	args := m.Called(ctx, productCategory)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) GetByID(ctx context.Context, id uint) (*model.ProductCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) Update(ctx context.Context, productCategory *model.ProductCategory) error {
	args := m.Called(ctx, productCategory)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.ProductCategory, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.ProductCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductCategoryRepository) GetByName(ctx context.Context, name string) (*model.ProductCategory, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) GetByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductCategoryRepository) GetCategoriesWithProductCount(ctx context.Context, opts repository.ListOptions) ([]*repository.CategoryWithProductCount, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*repository.CategoryWithProductCount), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductCategoryRepository) HasChildren(ctx context.Context, categoryID uint) (bool, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockProductCategoryRepository) HasProducts(ctx context.Context, categoryID uint) (bool, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(bool), args.Error(1)
}

// MockProductRepository Product仓储模拟
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetByName(ctx context.Context, name string) (*model.Product, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*model.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) GetByCategoryID(ctx context.Context, categoryID uint, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, categoryID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetByCategoryIDs(ctx context.Context, categoryIDs []uint, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, categoryIDs, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetByPriceRange(ctx context.Context, minPrice float64, maxPrice float64, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, minPrice, maxPrice, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetLowStockProducts(ctx context.Context, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetHotSellingProducts(ctx context.Context, minSales int, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, minSales, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) SearchByName(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) SearchBySKU(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetActiveProducts(ctx context.Context, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetProductsInStock(ctx context.Context, opts repository.ListOptions) ([]*model.Product, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, productID uint, quantityChange int) error {
	args := m.Called(ctx, productID, quantityChange)
	return args.Error(0)
}

func (m *MockProductRepository) BatchUpdatePrices(ctx context.Context, updates map[uint]float64) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}
