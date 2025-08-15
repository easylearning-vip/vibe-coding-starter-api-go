package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/test/testutil"
)

// ProductCategoryRepositoryTestSuite ProductCategory仓储测试套件
type ProductCategoryRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.ProductCategoryRepository
	ctx    context.Context
}

// SetupSuite 设置测试套件
func (suite *ProductCategoryRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建仓储实例
	suite.repo = repository.NewProductCategoryRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *ProductCategoryRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
}

// SetupTest 每个测试前的设置
func (suite *ProductCategoryRepositoryTestSuite) SetupTest() {
	// 清理数据
	suite.db.Clean(suite.T())
}

// TestCreate 测试创建ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestCreate() {
	// 准备测试数据
	entity := &model.ProductCategory{
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}

	// 执行创建
	err := suite.repo.Create(suite.ctx, entity)

	// 验证结果
	suite.NoError(err)
	suite.NotZero(entity.ID)
	suite.Equal("Test ProductCategory", entity.Name)
	suite.Equal("Test Description", entity.Description)
	suite.NotZero(entity.CreatedAt)
	suite.NotZero(entity.UpdatedAt)
}

// TestCreateDuplicateName 测试创建重复名称的ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestCreateDuplicateName() {
	// 创建第一个ProductCategory
	entity1 := &model.ProductCategory{
		Name:        "Duplicate Name",
		Description: "First Description",
	}
	err := suite.repo.Create(suite.ctx, entity1)
	suite.NoError(err)

	// 尝试创建重复名称的ProductCategory
	entity2 := &model.ProductCategory{
		Name:        "Duplicate Name",
		Description: "Second Description",
	}
	err = suite.repo.Create(suite.ctx, entity2)
	suite.Error(err)
}

// TestGetByID 测试根据ID获取ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestGetByID() {
	// 创建测试数据
	entity := &model.ProductCategory{
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 根据ID获取
	found, err := suite.repo.GetByID(suite.ctx, entity.ID)

	// 验证结果
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(entity.ID, found.ID)
	suite.Equal(entity.Name, found.Name)
	suite.Equal(entity.Description, found.Description)
}

// TestGetByIDNotFound 测试获取不存在的ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestGetByIDNotFound() {
	// 尝试获取不存在的ProductCategory
	found, err := suite.repo.GetByID(suite.ctx, 999)

	// 验证结果
	suite.Error(err)
	suite.Nil(found)
	suite.Contains(err.Error(), "not found")
}

// TestGetByName 测试根据名称获取ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestGetByName() {
	// 创建测试数据
	entity := &model.ProductCategory{
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 根据名称获取
	found, err := suite.repo.GetByName(suite.ctx, "Test ProductCategory")

	// 验证结果
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(entity.ID, found.ID)
	suite.Equal(entity.Name, found.Name)
}

// TestUpdate 测试更新ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestUpdate() {
	// 创建测试数据
	entity := &model.ProductCategory{
		Name:        "Original Name",
		Description: "Original Description",
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 更新数据
	entity.Name = "Updated Name"
	entity.Description = "Updated Description"
	err = suite.repo.Update(suite.ctx, entity)

	// 验证结果
	suite.NoError(err)

	// 重新获取验证
	updated, err := suite.repo.GetByID(suite.ctx, entity.ID)
	suite.NoError(err)
	suite.Equal("Updated Name", updated.Name)
	suite.Equal("Updated Description", updated.Description)
}

// TestDelete 测试删除ProductCategory
func (suite *ProductCategoryRepositoryTestSuite) TestDelete() {
	// 创建测试数据
	entity := &model.ProductCategory{
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}
	err := suite.repo.Create(suite.ctx, entity)
	suite.NoError(err)

	// 删除数据
	err = suite.repo.Delete(suite.ctx, entity.ID)
	suite.NoError(err)

	// 验证已删除
	found, err := suite.repo.GetByID(suite.ctx, entity.ID)
	suite.Error(err)
	suite.Nil(found)
}

// TestList 测试获取ProductCategory列表
func (suite *ProductCategoryRepositoryTestSuite) TestList() {
	// 创建测试数据
	entities := []*model.ProductCategory{
		{Name: "ProductCategory 1", Description: "Description 1"},
		{Name: "ProductCategory 2", Description: "Description 2"},
		{Name: "ProductCategory 3", Description: "Description 3"},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 获取列表
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 3)
	suite.Equal(int64(3), total)
}

// TestListWithFilters 测试带过滤器的列表查询
func (suite *ProductCategoryRepositoryTestSuite) TestListWithFilters() {
	// 创建测试数据
	entities := []*model.ProductCategory{
		{Name: "Active ProductCategory", Description: "Active Description"},
		{Name: "Inactive ProductCategory", Description: "Inactive Description"},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 使用过滤器查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"name": "Active ProductCategory",
		},
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(1), total)
	suite.Equal("Active ProductCategory", result[0].Name)
}

// TestListWithSearch 测试带搜索的列表查询
func (suite *ProductCategoryRepositoryTestSuite) TestListWithSearch() {
	// 创建测试数据
	entities := []*model.ProductCategory{
		{Name: "Searchable ProductCategory", Description: "Description"},
		{Name: "Another ProductCategory", Description: "Description"},
	}

	for _, entity := range entities {
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 使用搜索查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "Searchable",
	}
	result, total, err := suite.repo.List(suite.ctx, opts)

	// 验证结果
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(1), total)
	suite.Contains(result[0].Name, "Searchable")
}

// TestPagination 测试分页
func (suite *ProductCategoryRepositoryTestSuite) TestPagination() {
	// 创建测试数据
	for i := 1; i <= 5; i++ {
		entity := &model.ProductCategory{
			Name:        fmt.Sprintf("ProductCategory %d", i),
			Description: fmt.Sprintf("Description %d", i),
		}
		err := suite.repo.Create(suite.ctx, entity)
		suite.NoError(err)
	}

	// 第一页
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 2,
	}
	result, total, err := suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(int64(5), total)

	// 第二页
	opts.Page = 2
	result, total, err = suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(int64(5), total)

	// 第三页
	opts.Page = 3
	result, total, err = suite.repo.List(suite.ctx, opts)
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal(int64(5), total)
}

// TestProductCategoryRepositoryTestSuite 运行测试套件
func TestProductCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProductCategoryRepositoryTestSuite))
}
