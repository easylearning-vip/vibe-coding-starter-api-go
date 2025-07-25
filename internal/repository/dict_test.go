package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/test/testutil"
)

// DictRepositoryTestSuite 数据字典仓储测试套件
type DictRepositoryTestSuite struct {
	suite.Suite
	db           *testutil.TestDatabase
	categoryRepo DictCategoryRepository
	itemRepo     DictItemRepository
	dictRepo     DictRepository
	ctx          context.Context
}

// SetupSuite 设置测试套件
func (suite *DictRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = testutil.NewTestDatabase(suite.T())

	// 创建仓储实例
	testLogger := testutil.NewTestLogger(suite.T())
	logger := testLogger.CreateTestLogger()
	dbAdapter := suite.db.CreateTestDatabase()
	suite.categoryRepo = NewDictCategoryRepository(dbAdapter, logger)
	suite.itemRepo = NewDictItemRepository(dbAdapter, logger)
	suite.dictRepo = NewDictRepository(dbAdapter, logger)

	// 创建表结构
	err := suite.db.GetDB().AutoMigrate(&model.DictCategory{}, &model.DictItem{})
	require.NoError(suite.T(), err)
}

// TearDownSuite 清理测试套件
func (suite *DictRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

// SetupTest 设置每个测试
func (suite *DictRepositoryTestSuite) SetupTest() {
	// 清理数据
	suite.db.GetDB().Exec("DELETE FROM dict_items")
	suite.db.GetDB().Exec("DELETE FROM dict_categories")
}

// TestDictCategoryRepository_Create 测试创建字典分类
func (suite *DictRepositoryTestSuite) TestDictCategoryRepository_Create() {
	category := &model.DictCategory{
		Code:        "test_category",
		Name:        "测试分类",
		Description: "这是一个测试分类",
		SortOrder:   1,
	}

	err := suite.categoryRepo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), category.ID)
	assert.NotZero(suite.T(), category.CreatedAt)
	assert.NotZero(suite.T(), category.UpdatedAt)
}

// TestDictCategoryRepository_GetByCode 测试根据代码获取分类
func (suite *DictRepositoryTestSuite) TestDictCategoryRepository_GetByCode() {
	// 创建测试分类
	category := suite.createTestCategory("test_category", "测试分类")

	// 根据代码获取分类
	result, err := suite.categoryRepo.GetByCode(suite.ctx, "test_category")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), category.ID, result.ID)
	assert.Equal(suite.T(), category.Code, result.Code)
	assert.Equal(suite.T(), category.Name, result.Name)
}

// TestDictCategoryRepository_GetByCode_NotFound 测试获取不存在的分类
func (suite *DictRepositoryTestSuite) TestDictCategoryRepository_GetByCode_NotFound() {
	result, err := suite.categoryRepo.GetByCode(suite.ctx, "nonexistent")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestDictItemRepository_Create 测试创建字典项
func (suite *DictRepositoryTestSuite) TestDictItemRepository_Create() {
	// 先创建分类
	category := suite.createTestCategory("test_category", "测试分类")

	active := true
	item := &model.DictItem{
		CategoryCode: category.Code,
		ItemKey:      "test_key",
		ItemValue:    "测试值",
		Description:  "这是一个测试字典项",
		SortOrder:    1,
		IsActive:     &active,
	}

	err := suite.itemRepo.Create(suite.ctx, item)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), item.ID)
	assert.NotZero(suite.T(), item.CreatedAt)
	assert.NotZero(suite.T(), item.UpdatedAt)
}

// TestDictItemRepository_GetByCategory 测试根据分类获取字典项
func (suite *DictRepositoryTestSuite) TestDictItemRepository_GetByCategory() {
	// 创建测试分类和字典项
	category := suite.createTestCategory("test_category", "测试分类")
	item1 := suite.createTestItem(category.Code, "key1", "值1", true)
	item2 := suite.createTestItem(category.Code, "key2", "值2", true)

	// 获取分类下的所有字典项
	items, err := suite.itemRepo.GetByCategory(suite.ctx, category.Code)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), items, 2)

	// 验证返回的项
	itemKeys := make(map[string]bool)
	for _, item := range items {
		itemKeys[item.ItemKey] = true
		assert.Equal(suite.T(), category.Code, item.CategoryCode)
	}
	assert.True(suite.T(), itemKeys[item1.ItemKey])
	assert.True(suite.T(), itemKeys[item2.ItemKey])
}

// TestDictItemRepository_GetActiveByCategory 测试获取启用的字典项
func (suite *DictRepositoryTestSuite) TestDictItemRepository_GetActiveByCategory() {
	// 创建测试分类和字典项
	category := suite.createTestCategory("test_category", "测试分类")
	suite.createTestItem(category.Code, "key1", "值1", true)  // 启用
	suite.createTestItem(category.Code, "key2", "值2", false) // 禁用
	suite.createTestItem(category.Code, "key3", "值3", true)  // 启用

	// 获取启用的字典项
	items, err := suite.itemRepo.GetActiveByCategory(suite.ctx, category.Code)
	require.NoError(suite.T(), err)

	assert.Len(suite.T(), items, 2) // 只有2个启用的项

	// 验证所有返回的项都是启用的
	for _, item := range items {
		assert.True(suite.T(), item.IsActive != nil && *item.IsActive)
	}
}

// TestDictItemRepository_GetByCategoryAndKey 测试根据分类和键值获取字典项
func (suite *DictRepositoryTestSuite) TestDictItemRepository_GetByCategoryAndKey() {
	// 创建测试分类和字典项
	category := suite.createTestCategory("test_category", "测试分类")
	item := suite.createTestItem(category.Code, "test_key", "测试值", true)

	// 根据分类和键值获取字典项
	result, err := suite.itemRepo.GetByCategoryAndKey(suite.ctx, category.Code, "test_key")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), item.ID, result.ID)
	assert.Equal(suite.T(), item.CategoryCode, result.CategoryCode)
	assert.Equal(suite.T(), item.ItemKey, result.ItemKey)
	assert.Equal(suite.T(), item.ItemValue, result.ItemValue)
}

// TestDictRepository_Integration 测试组合仓储的集成功能
func (suite *DictRepositoryTestSuite) TestDictRepository_Integration() {
	// 创建分类
	categoryReq := &model.DictCategory{
		Code:        "integration_test",
		Name:        "集成测试分类",
		Description: "用于集成测试的分类",
		SortOrder:   1,
	}

	err := suite.dictRepo.CreateCategory(suite.ctx, categoryReq)
	require.NoError(suite.T(), err)

	// 获取分类
	category, err := suite.dictRepo.GetCategoryByCode(suite.ctx, "integration_test")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), categoryReq.Code, category.Code)

	// 创建字典项
	active := true
	itemReq := &model.DictItem{
		CategoryCode: "integration_test",
		ItemKey:      "test_key",
		ItemValue:    "测试值",
		Description:  "集成测试字典项",
		SortOrder:    1,
		IsActive:     &active,
	}

	err = suite.dictRepo.CreateItem(suite.ctx, itemReq)
	require.NoError(suite.T(), err)

	// 获取字典项
	items, err := suite.dictRepo.GetActiveItemsByCategory(suite.ctx, "integration_test")
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), items, 1)
	assert.Equal(suite.T(), itemReq.ItemKey, items[0].ItemKey)
	assert.Equal(suite.T(), itemReq.ItemValue, items[0].ItemValue)
}

// 辅助方法

func (suite *DictRepositoryTestSuite) createTestCategory(code, name string) *model.DictCategory {
	category := &model.DictCategory{
		Code:        code,
		Name:        name,
		Description: "测试分类描述",
		SortOrder:   1,
	}

	err := suite.categoryRepo.Create(suite.ctx, category)
	require.NoError(suite.T(), err)
	return category
}

func (suite *DictRepositoryTestSuite) createTestItem(categoryCode, key, value string, isActive bool) *model.DictItem {
	item := &model.DictItem{
		CategoryCode: categoryCode,
		ItemKey:      key,
		ItemValue:    value,
		Description:  "测试字典项描述",
		SortOrder:    1,
		IsActive:     &isActive,
	}

	err := suite.itemRepo.Create(suite.ctx, item)
	require.NoError(suite.T(), err)
	return item
}

// TestDictRepositoryTestSuite 运行测试套件
func TestDictRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DictRepositoryTestSuite))
}
