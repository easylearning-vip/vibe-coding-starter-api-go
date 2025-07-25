package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// boolPtr 创建bool指针的辅助函数
func boolPtr(b bool) *bool {
	return &b
}

// DictServiceTestSuite 数据字典服务测试套件
type DictServiceTestSuite struct {
	suite.Suite
	dictRepo *mocks.MockDictRepository
	cache    *mocks.MockCache
	logger   *mocks.MockLogger
	service  service.DictService
	ctx      context.Context
}

// SetupTest 设置每个测试
func (suite *DictServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.dictRepo = &mocks.MockDictRepository{}
	suite.cache = &mocks.MockCache{}
	suite.logger = &mocks.MockLogger{}

	suite.service = service.NewDictService(
		suite.dictRepo,
		suite.cache,
		suite.logger,
	)
}

// TestGetDictItems_Success 测试成功获取字典项
func (suite *DictServiceTestSuite) TestGetDictItems_Success() {
	categoryCode := "test_category"
	expectedItems := []*model.DictItem{
		{
			BaseModel:    model.BaseModel{ID: 1},
			CategoryCode: categoryCode,
			ItemKey:      "key1",
			ItemValue:    "值1",
			IsActive:     boolPtr(true),
		},
		{
			BaseModel:    model.BaseModel{ID: 2},
			CategoryCode: categoryCode,
			ItemKey:      "key2",
			ItemValue:    "值2",
			IsActive:     boolPtr(true),
		},
	}

	// Mock 缓存未命中
	suite.cache.On("Get", suite.ctx, "dict_items:test_category").Return("", errors.New("cache miss"))

	// Mock 数据库查询
	suite.dictRepo.On("GetActiveItemsByCategory", suite.ctx, categoryCode).Return(expectedItems, nil)

	// Mock 缓存设置
	suite.cache.On("Set", suite.ctx, "dict_items:test_category", mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Mock 日志
	suite.logger.On("Debug", "Dict items retrieved from database", "category_code", categoryCode, "count", 2)

	// 执行测试
	result, err := suite.service.GetDictItems(suite.ctx, categoryCode)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), expectedItems[0].ItemKey, result[0].ItemKey)
	assert.Equal(suite.T(), expectedItems[1].ItemKey, result[1].ItemKey)

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetDictItems_CacheHit 测试缓存命中
func (suite *DictServiceTestSuite) TestGetDictItems_CacheHit() {
	categoryCode := "test_category"
	cachedData := `[{"id":1,"category_code":"test_category","item_key":"key1","item_value":"值1","is_active":true}]`

	// Mock 缓存命中
	suite.cache.On("Get", suite.ctx, "dict_items:test_category").Return(cachedData, nil)

	// Mock 日志
	suite.logger.On("Debug", "Dict items retrieved from cache", "category_code", categoryCode, "count", 1)

	// 执行测试
	result, err := suite.service.GetDictItems(suite.ctx, categoryCode)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "key1", result[0].ItemKey)
	assert.Equal(suite.T(), "值1", result[0].ItemValue)

	// 验证mock调用
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
	// 不应该调用数据库
	suite.dictRepo.AssertNotCalled(suite.T(), "GetActiveItemsByCategory")
}

// TestGetDictItemByKey_Success 测试成功获取特定字典项
func (suite *DictServiceTestSuite) TestGetDictItemByKey_Success() {
	categoryCode := "test_category"
	itemKey := "test_key"
	allItems := []*model.DictItem{
		{
			BaseModel:    model.BaseModel{ID: 1},
			CategoryCode: categoryCode,
			ItemKey:      itemKey,
			ItemValue:    "测试值",
			IsActive:     boolPtr(true),
		},
		{
			BaseModel:    model.BaseModel{ID: 2},
			CategoryCode: categoryCode,
			ItemKey:      "other_key",
			ItemValue:    "其他值",
			IsActive:     boolPtr(true),
		},
	}

	// Mock 缓存未命中
	suite.cache.On("Get", suite.ctx, "dict_item:test_category:test_key").Return("", errors.New("cache miss"))

	// Mock 数据库查询
	suite.dictRepo.On("GetItemsByCategory", suite.ctx, categoryCode).Return(allItems, nil)

	// Mock 缓存设置
	suite.cache.On("Set", suite.ctx, "dict_item:test_category:test_key", mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

	// Mock 日志
	suite.logger.On("Debug", "Dict item found", "category_code", categoryCode, "item_key", itemKey)

	// 执行测试
	result, err := suite.service.GetDictItemByKey(suite.ctx, categoryCode, itemKey)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), itemKey, result.ItemKey)
	assert.Equal(suite.T(), "测试值", result.ItemValue)

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetDictItemByKey_NotFound 测试获取不存在的字典项
func (suite *DictServiceTestSuite) TestGetDictItemByKey_NotFound() {
	categoryCode := "test_category"
	itemKey := "nonexistent_key"
	allItems := []*model.DictItem{
		{
			BaseModel:    model.BaseModel{ID: 1},
			CategoryCode: categoryCode,
			ItemKey:      "other_key",
			ItemValue:    "其他值",
			IsActive:     boolPtr(true),
		},
	}

	// Mock 缓存未命中
	suite.cache.On("Get", suite.ctx, "dict_item:test_category:nonexistent_key").Return("", errors.New("cache miss"))

	// Mock 数据库查询
	suite.dictRepo.On("GetItemsByCategory", suite.ctx, categoryCode).Return(allItems, nil)

	// Mock 日志
	suite.logger.On("Error", "Failed to get dict items for key lookup", "category_code", categoryCode, "item_key", itemKey, "error", mock.AnythingOfType("*errors.errorString"))

	// 执行测试
	result, err := suite.service.GetDictItemByKey(suite.ctx, categoryCode, itemKey)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "not found")

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
}

// TestCreateDictCategory_Success 测试成功创建字典分类
func (suite *DictServiceTestSuite) TestCreateDictCategory_Success() {
	req := &service.CreateCategoryRequest{
		Code:        "new_category",
		Name:        "新分类",
		Description: "新分类描述",
		SortOrder:   1,
	}

	// Mock 检查分类不存在
	suite.dictRepo.On("GetCategoryByCode", suite.ctx, req.Code).Return(nil, errors.New("not found"))

	// Mock 创建分类
	suite.dictRepo.On("CreateCategory", suite.ctx, mock.AnythingOfType("*model.DictCategory")).Return(nil)

	// Mock 日志
	suite.logger.On("Info", "Dict category created successfully", "code", req.Code, "name", req.Name)

	// 执行测试
	result, err := suite.service.CreateDictCategory(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Code, result.Code)
	assert.Equal(suite.T(), req.Name, result.Name)

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestCreateDictCategory_AlreadyExists 测试创建已存在的分类
func (suite *DictServiceTestSuite) TestCreateDictCategory_AlreadyExists() {
	req := &service.CreateCategoryRequest{
		Code:        "existing_category",
		Name:        "已存在分类",
		Description: "已存在分类描述",
		SortOrder:   1,
	}

	existingCategory := &model.DictCategory{
		BaseModel: model.BaseModel{ID: 1},
		Code:      req.Code,
		Name:      "已存在的分类",
	}

	// Mock 检查分类已存在
	suite.dictRepo.On("GetCategoryByCode", suite.ctx, req.Code).Return(existingCategory, nil)

	// 执行测试
	result, err := suite.service.CreateDictCategory(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "already exists")

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
}

// TestCreateDictItem_Success 测试成功创建字典项
func (suite *DictServiceTestSuite) TestCreateDictItem_Success() {
	req := &service.CreateItemRequest{
		CategoryCode: "test_category",
		ItemKey:      "new_key",
		ItemValue:    "新值",
		Description:  "新字典项描述",
		SortOrder:    1,
		IsActive:     boolPtr(true),
	}

	category := &model.DictCategory{
		BaseModel: model.BaseModel{ID: 1},
		Code:      req.CategoryCode,
		Name:      "测试分类",
	}

	// Mock 检查分类存在
	suite.dictRepo.On("GetCategoryByCode", suite.ctx, req.CategoryCode).Return(category, nil)

	// Mock 创建字典项
	suite.dictRepo.On("CreateItem", suite.ctx, mock.AnythingOfType("*model.DictItem")).Return(nil)

	// Mock 清除缓存
	suite.cache.On("Del", suite.ctx, "dict_items:test_category").Return(nil)
	suite.cache.On("Del", suite.ctx, "dict_item:test_category:new_key").Return(nil)

	// Mock 日志
	suite.logger.On("Info", "Dict item created successfully", "category_code", req.CategoryCode, "item_key", req.ItemKey)
	suite.logger.On("Debug", "Dict cache cleared", "category_code", req.CategoryCode, "item_key", req.ItemKey)

	// 执行测试
	result, err := suite.service.CreateDictItem(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.CategoryCode, result.CategoryCode)
	assert.Equal(suite.T(), req.ItemKey, result.ItemKey)
	assert.Equal(suite.T(), req.ItemValue, result.ItemValue)

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
	suite.cache.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestCreateDictItem_CategoryNotExists 测试创建字典项时分类不存在
func (suite *DictServiceTestSuite) TestCreateDictItem_CategoryNotExists() {
	req := &service.CreateItemRequest{
		CategoryCode: "nonexistent_category",
		ItemKey:      "new_key",
		ItemValue:    "新值",
		Description:  "新字典项描述",
		SortOrder:    1,
		IsActive:     boolPtr(true),
	}

	// Mock 检查分类不存在
	suite.dictRepo.On("GetCategoryByCode", suite.ctx, req.CategoryCode).Return(nil, errors.New("not found"))

	// 执行测试
	result, err := suite.service.CreateDictItem(suite.ctx, req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "does not exist")

	// 验证mock调用
	suite.dictRepo.AssertExpectations(suite.T())
}

// TestDictServiceTestSuite 运行测试套件
func TestDictServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DictServiceTestSuite))
}
