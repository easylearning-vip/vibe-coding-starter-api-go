package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

type ProductCategoryServiceTestSuite struct {
	suite.Suite
	service    service.ProductCategoryService
	mockRepo   *mocks.MockProductCategoryRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *ProductCategoryServiceTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockProductCategoryRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewProductCategoryService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *ProductCategoryServiceTestSuite) TestCreate_Success() {
	// 准备测试数据
	req := &service.CreateProductCategoryRequest{
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}
	
	expectedProductCategory := &model.ProductCategory{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.ProductCategory")).Return(nil).Run(func(args mock.Arguments) {
		productCategory := args.Get(1).(*model.ProductCategory)
		productCategory.ID = 1
	})
	
	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedProductCategory.Name, result.Name)
	assert.Equal(suite.T(), expectedProductCategory.Description, result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedProductCategory := &model.ProductCategory{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test ProductCategory",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedProductCategory, nil)
	
	// 执行测试
	result, err := suite.service.GetByID(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProductCategory, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceTestSuite) TestUpdate_Success() {
	// 准备测试数据
	existingProductCategory := &model.ProductCategory{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Old Name",
		Description: "Old Description",
	}
	
	newName := "New Name"
	req := &service.UpdateProductCategoryRequest{
		Name: &newName,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProductCategory, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.ProductCategory")).Return(nil)
	
	// 执行测试
	result, err := suite.service.Update(suite.ctx, 1, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, result.Name)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceTestSuite) TestDelete_Success() {
	// 准备测试数据
	existingProductCategory := &model.ProductCategory{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test ProductCategory",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProductCategory, nil)
	suite.mockRepo.On("HasChildren", suite.ctx, uint(1)).Return(false, nil)
	suite.mockRepo.On("CountProductsByCategory", suite.ctx, uint(1)).Return(int64(0), nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)
	
	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductCategoryServiceTestSuite) TestList_Success() {
	// 准备测试数据
	expectedproductcategories := []*model.ProductCategory{
		{BaseModel: model.BaseModel{ID: 1}, Name: "ProductCategory 1"},
		{BaseModel: model.BaseModel{ID: 2}, Name: "ProductCategory 2"},
	}
	
	opts := &service.ListProductCategoryOptions{
		Page:     1,
		PageSize: 10,
	}
	
	repoOpts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Sort:     "",
		Order:    "",
		Filters:  nil, // 使用nil而不是空map，与service实际传递的值一致
		Search:   "",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("List", suite.ctx, repoOpts).Return(expectedproductcategories, int64(2), nil)
	
	// 执行测试
	result, total, err := suite.service.List(suite.ctx, opts)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedproductcategories, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProductCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductCategoryServiceTestSuite))
}
