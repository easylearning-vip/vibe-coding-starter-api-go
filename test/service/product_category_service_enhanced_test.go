package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

type ProductCategoryServiceEnhancedTestSuite struct {
	suite.Suite
	service    service.ProductCategoryService
	mockRepo   *mocks.MockProductCategoryRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *ProductCategoryServiceEnhancedTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockProductCategoryRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()
	suite.mockLogger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewProductCategoryService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestGetCategoryTree_Success() {
	// 准备测试数据
	expectedTree := []*model.ProductCategory{
		{
			BaseModel:   model.BaseModel{ID: 1},
			Name:        "Electronics",
			Description: "Electronic products",
			ParentId:    0,
			SortOrder:   1,
			IsActive:    true,
			Children: []*model.ProductCategory{
				{
					BaseModel:   model.BaseModel{ID: 2},
					Name:        "Phones",
					Description: "Mobile phones",
					ParentId:    1,
					SortOrder:   1,
					IsActive:    true,
				},
			},
		},
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetCategoryTree", suite.ctx).Return(expectedTree, nil)
	
	// 执行测试
	result, err := suite.service.GetCategoryTree(suite.ctx)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTree, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestGetCategoryPath_Success() {
	// 准备测试数据
	expectedPath := []*model.ProductCategory{
		{
			BaseModel:   model.BaseModel{ID: 1},
			Name:        "Electronics",
			Description: "Electronic products",
			ParentId:    0,
		},
		{
			BaseModel:   model.BaseModel{ID: 2},
			Name:        "Phones",
			Description: "Mobile phones",
			ParentId:    1,
		},
		{
			BaseModel:   model.BaseModel{ID: 3},
			Name:        "Smartphones",
			Description: "Smart phones",
			ParentId:    2,
		},
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetCategoryPath", suite.ctx, uint(3)).Return(expectedPath, nil)
	
	// 执行测试
	result, err := suite.service.GetCategoryPath(suite.ctx, 3)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPath, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestGetByParentID_Success() {
	// 准备测试数据
	expectedCategories := []*model.ProductCategory{
		{
			BaseModel:   model.BaseModel{ID: 2},
			Name:        "Phones",
			Description: "Mobile phones",
			ParentId:    1,
			SortOrder:   1,
			IsActive:    true,
		},
		{
			BaseModel:   model.BaseModel{ID: 3},
			Name:        "Laptops",
			Description: "Laptop computers",
			ParentId:    1,
			SortOrder:   2,
			IsActive:    true,
		},
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByParentID", suite.ctx, uint(1)).Return(expectedCategories, nil)
	
	// 执行测试
	result, err := suite.service.GetByParentID(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCategories, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestBatchUpdateSortOrder_Success() {
	// 准备测试数据
	updates := []*service.UpdateSortOrderRequest{
		{CategoryID: 1, SortOrder: 2},
		{CategoryID: 2, SortOrder: 1},
	}
	
	existingCategory1 := &model.ProductCategory{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Electronics",
	}
	
	existingCategory2 := &model.ProductCategory{
		BaseModel: model.BaseModel{ID: 2},
		Name:      "Phones",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingCategory1, nil)
	suite.mockRepo.On("GetByID", suite.ctx, uint(2)).Return(existingCategory2, nil)
	suite.mockRepo.On("Update", suite.ctx, existingCategory1).Return(nil)
	suite.mockRepo.On("Update", suite.ctx, existingCategory2).Return(nil)
	
	// 执行测试
	err := suite.service.BatchUpdateSortOrder(suite.ctx, updates)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, existingCategory1.SortOrder)
	assert.Equal(suite.T(), 1, existingCategory2.SortOrder)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestValidateParentChildRelationship_CircularReference() {
	// 准备测试数据 - 创建循环引用的路径
	path := []*model.ProductCategory{
		{BaseModel: model.BaseModel{ID: 3}, ParentId: 2},
		{BaseModel: model.BaseModel{ID: 2}, ParentId: 1},
		{BaseModel: model.BaseModel{ID: 1}, ParentId: 0},
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetCategoryPath", suite.ctx, uint(3)).Return(path, nil)
	
	// 执行测试 - 尝试让分类1成为分类3的父级，这会创建循环引用
	err := suite.service.ValidateParentChildRelationship(suite.ctx, 3, 1)
	
	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "circular reference detected")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestValidateParentChildRelationship_SameCategory() {
	// 执行测试 - 尝试让分类成为自己的父级
	err := suite.service.ValidateParentChildRelationship(suite.ctx, 1, 1)
	
	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be its own parent")
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestCanDeleteCategory_HasChildren() {
	// 设置 mock 期望
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(true, nil)
	
	// 执行测试
	canDelete, err := suite.service.CanDeleteCategory(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), canDelete)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestCanDeleteCategory_HasProducts() {
	// 设置 mock 期望
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(false, nil)
	suite.mockRepo.On("CountProductsByCategory", suite.ctx, uint(1)).Return(int64(5), nil)
	
	// 执行测试
	canDelete, err := suite.service.CanDeleteCategory(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), canDelete)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestCanDeleteCategory_CanDelete() {
	// 设置 mock 期望
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(false, nil)
	suite.mockRepo.On("CountProductsByCategory", suite.ctx, uint(1)).Return(int64(0), nil)
	
	// 执行测试
	canDelete, err := suite.service.CanDeleteCategory(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), canDelete)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestCreate_WithParentValidation() {
	// 准备测试数据
	req := &service.CreateProductCategoryRequest{
		Name:        "Smartphones",
		Description: "Smart phones",
		ParentId:    1,
		SortOrder:   1,
		IsActive:    true,
	}
	
	path := []*model.ProductCategory{
		{BaseModel: model.BaseModel{ID: 1}, ParentId: 0},
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetCategoryPath", suite.ctx, uint(1)).Return(path, nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.ProductCategory")).Return(nil).Run(func(args mock.Arguments) {
		productCategory := args.Get(1).(*model.ProductCategory)
		productCategory.ID = 4
	})
	
	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Name, result.Name)
	assert.Equal(suite.T(), req.ParentId, result.ParentId)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestDelete_WithValidation() {
	// 准备测试数据
	existingCategory := &model.ProductCategory{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Electronics",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingCategory, nil)
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(false, nil)
	suite.mockRepo.On("CountProductsByCategory", suite.ctx, uint(1)).Return(int64(0), nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)
	
	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceEnhancedTestSuite) TestDelete_CannotDelete() {
	// 准备测试数据
	existingCategory := &model.ProductCategory{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Electronics",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingCategory, nil)
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(true, nil)
	
	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)
	
	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be deleted")
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProductCategoryServiceEnhancedTestSuite(t *testing.T) {
	suite.Run(t, new(ProductCategoryServiceEnhancedTestSuite))
}