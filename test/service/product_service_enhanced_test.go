package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

type ProductServiceEnhancedTestSuite struct {
	suite.Suite
	service    service.ProductService
	mockRepo   *mocks.MockProductRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *ProductServiceEnhancedTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockProductRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用 - 使用Maybe()让mock更灵活
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	suite.service = service.NewProductService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *ProductServiceEnhancedTestSuite) TestGetBySKU_Success() {
	// 准备测试数据
	expectedProduct := &model.Product{
		BaseModel:     model.BaseModel{ID: 1},
		Name:         "Test Product",
		Sku:          "TEST-001",
		Price:        99.99,
		StockQuantity: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetBySKU", suite.ctx, "TEST-001").Return(expectedProduct, nil)

	// 执行测试
	result, err := suite.service.GetBySKU(suite.ctx, "TEST-001")

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProduct, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetByCategoryID_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Product 1",
			CategoryId:    1,
			Price:        99.99,
			StockQuantity: 10,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Product 2",
			CategoryId:    1,
			Price:        199.99,
			StockQuantity: 5,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByCategoryID", suite.ctx, uint(1), mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.GetByCategoryID(suite.ctx, 1, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetByCategoryWithSubcategories_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Product 1",
			CategoryId:    1,
			Price:        99.99,
			StockQuantity: 10,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Product 2",
			CategoryId:    2, // 子分类
			Price:        199.99,
			StockQuantity: 5,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByCategoryWithSubcategories", suite.ctx, uint(1), mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.GetByCategoryWithSubcategories(suite.ctx, 1, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetHotSellingProducts_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Hot Product 1",
			Price:        99.99,
			StockQuantity: 2, // 低库存
			IsActive:      true,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Hot Product 2",
			Price:        199.99,
			StockQuantity: 1, // 低库存
			IsActive:      true,
		},
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetHotSellingProducts", suite.ctx, 10).Return(expectedProducts, nil)

	// 执行测试
	result, err := suite.service.GetHotSellingProducts(suite.ctx, 10)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetHotSellingProducts_InvalidLimit() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Hot Product",
			Price:        99.99,
			StockQuantity: 2,
			IsActive:      true,
		},
	}

	// 设置 mock 期望 - 应该使用默认限制 10
	suite.mockRepo.On("GetHotSellingProducts", suite.ctx, 10).Return(expectedProducts, nil)

	// 执行测试 - 传入无效限制
	result, err := suite.service.GetHotSellingProducts(suite.ctx, 0)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetByPriceRange_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Product 1",
			Price:        99.99,
			StockQuantity: 10,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Product 2",
			Price:        149.99,
			StockQuantity: 5,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByPriceRange", suite.ctx, 50.0, 200.0, mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.GetByPriceRange(suite.ctx, 50.0, 200.0, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetByPriceRange_InvalidRange() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Product 1",
			Price:        99.99,
			StockQuantity: 10,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望 - 应该自动交换最小值和最大值
	suite.mockRepo.On("GetByPriceRange", suite.ctx, 50.0, 200.0, mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(1), nil)

	// 执行测试 - 传入反转的范围
	result, total, err := suite.service.GetByPriceRange(suite.ctx, 200.0, 50.0, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(1), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestSearchProducts_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Search Product 1",
			Description:   "A great product for searching",
			Sku:           "SEARCH-001",
			Price:        99.99,
			StockQuantity: 10,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Search Product 2",
			Description:   "Another search product",
			Sku:           "SEARCH-002",
			Price:        199.99,
			StockQuantity: 5,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("SearchProducts", suite.ctx, "search", mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.SearchProducts(suite.ctx, "search", opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestSearchProducts_EmptyQuery() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Product 1",
			Price:        99.99,
			StockQuantity: 10,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望 - 空查询应该调用 List 方法
	suite.mockRepo.On("List", suite.ctx, mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(1), nil)

	// 执行测试 - 传入空查询
	result, total, err := suite.service.SearchProducts(suite.ctx, "", opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(1), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetLowStockProducts_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Low Stock Product 1",
			Price:        99.99,
			StockQuantity: 2,
			MinStock:      5,
			IsActive:      true,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Low Stock Product 2",
			Price:        199.99,
			StockQuantity: 1,
			MinStock:      3,
			IsActive:      true,
		},
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetLowStockProducts", suite.ctx).Return(expectedProducts, nil)

	// 执行测试
	result, err := suite.service.GetLowStockProducts(suite.ctx)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestGetProductsByStatus_Success() {
	// 准备测试数据
	expectedProducts := []*model.Product{
		{
			BaseModel:     model.BaseModel{ID: 1},
			Name:         "Active Product 1",
			Price:        99.99,
			StockQuantity: 10,
			IsActive:      true,
		},
		{
			BaseModel:     model.BaseModel{ID: 2},
			Name:         "Active Product 2",
			Price:        199.99,
			StockQuantity: 5,
			IsActive:      true,
		},
	}

	opts := &service.ListProductOptions{
		Page:     1,
		PageSize: 10,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetProductsByStatus", suite.ctx, true, mock.AnythingOfType("repository.ListOptions")).Return(expectedProducts, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.GetProductsByStatus(suite.ctx, true, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestBatchUpdatePrices_Success() {
	// 准备测试数据
	updates := []*service.UpdatePriceRequest{
		{ProductID: 1, Price: 149.99, CostPrice: &[]float64{99.99}[0]},
		{ProductID: 2, Price: 199.99},
	}

	existingProduct1 := &model.Product{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Product 1",
		Price:     99.99,
		CostPrice: 79.99,
	}

	existingProduct2 := &model.Product{
		BaseModel: model.BaseModel{ID: 2},
		Name:      "Product 2",
		Price:     179.99,
		CostPrice: 149.99,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct1, nil)
	suite.mockRepo.On("GetByID", suite.ctx, uint(2)).Return(existingProduct2, nil)
	suite.mockRepo.On("Update", suite.ctx, existingProduct1).Return(nil)
	suite.mockRepo.On("Update", suite.ctx, existingProduct2).Return(nil)

	// 执行测试
	err := suite.service.BatchUpdatePrices(suite.ctx, updates)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 149.99, existingProduct1.Price)
	assert.Equal(suite.T(), 99.99, existingProduct1.CostPrice)
	assert.Equal(suite.T(), 199.99, existingProduct2.Price)
	assert.Equal(suite.T(), 149.99, existingProduct2.CostPrice) // 成本价应该保持不变
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestBatchUpdateStatus_Success() {
	// 准备测试数据
	updates := []*service.UpdateStatusRequest{
		{ProductID: 1, IsActive: true},
		{ProductID: 2, IsActive: false},
	}

	existingProduct1 := &model.Product{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Product 1",
		IsActive:  false,
	}

	existingProduct2 := &model.Product{
		BaseModel: model.BaseModel{ID: 2},
		Name:      "Product 2",
		IsActive:  true,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct1, nil)
	suite.mockRepo.On("GetByID", suite.ctx, uint(2)).Return(existingProduct2, nil)
	suite.mockRepo.On("Update", suite.ctx, existingProduct1).Return(nil)
	suite.mockRepo.On("Update", suite.ctx, existingProduct2).Return(nil)

	// 执行测试
	err := suite.service.BatchUpdateStatus(suite.ctx, updates)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), existingProduct1.IsActive)
	assert.False(suite.T(), existingProduct2.IsActive)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceEnhancedTestSuite) TestBatchUpdatePrices_ProductNotFound() {
	// 准备测试数据
	updates := []*service.UpdatePriceRequest{
		{ProductID: 999, Price: 149.99}, // 不存在的产品ID
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(999)).Return(nil, fmt.Errorf("product not found"))

	// 执行测试
	err := suite.service.BatchUpdatePrices(suite.ctx, updates)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get product 999")
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProductServiceEnhancedTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceEnhancedTestSuite))
}