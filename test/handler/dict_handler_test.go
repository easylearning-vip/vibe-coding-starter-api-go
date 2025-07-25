package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// boolPtr 创建bool指针的辅助函数
func boolPtr(b bool) *bool {
	return &b
}

// DictHandlerTestSuite 数据字典处理器测试套件
type DictHandlerTestSuite struct {
	suite.Suite
	dictService *mocks.MockDictService
	logger      *mocks.MockLogger
	handler     *handler.DictHandler
	router      *gin.Engine
}

// SetupTest 设置每个测试
func (suite *DictHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.dictService = &mocks.MockDictService{}
	suite.logger = &mocks.MockLogger{}

	suite.handler = handler.NewDictHandler(
		suite.dictService,
		suite.logger,
	)

	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(v1)
}

// TestGetItemsByCategory_Success 测试成功获取字典项
func (suite *DictHandlerTestSuite) TestGetItemsByCategory_Success() {
	category := "test_category"
	expectedItems := []*model.DictItem{
		{
			BaseModel:    model.BaseModel{ID: 1},
			CategoryCode: category,
			ItemKey:      "key1",
			ItemValue:    "值1",
			IsActive:     boolPtr(true),
		},
		{
			BaseModel:    model.BaseModel{ID: 2},
			CategoryCode: category,
			ItemKey:      "key2",
			ItemValue:    "值2",
			IsActive:     boolPtr(true),
		},
	}

	// Mock 服务调用
	suite.dictService.On("GetDictItems", mock.Anything, category).Return(expectedItems, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dict/items/"+category, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.Equal(suite.T(), "success", response.Message)

	// 验证数据
	dataBytes, _ := json.Marshal(response.Data)
	var items []*model.DictItem
	err = json.Unmarshal(dataBytes, &items)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), items, 2)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestGetItemsByCategory_EmptyCategory 测试空分类参数
func (suite *DictHandlerTestSuite) TestGetItemsByCategory_EmptyCategory() {
	// Mock 日志
	suite.logger.On("Error", "Category parameter is required")

	// 创建请求（空分类参数）
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dict/items/", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - Gin会返回404对于不匹配的路由
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// 不验证mock调用，因为Handler没有被调用
}

// TestGetItemByKey_Success 测试成功获取特定字典项
func (suite *DictHandlerTestSuite) TestGetItemByKey_Success() {
	category := "test_category"
	key := "test_key"
	expectedItem := &model.DictItem{
		BaseModel:    model.BaseModel{ID: 1},
		CategoryCode: category,
		ItemKey:      key,
		ItemValue:    "测试值",
		IsActive:     boolPtr(true),
	}

	// Mock 服务调用
	suite.dictService.On("GetDictItemByKey", mock.Anything, category, key).Return(expectedItem, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dict/item/"+category+"/"+key, nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Code)

	// 验证数据
	dataBytes, _ := json.Marshal(response.Data)
	var item model.DictItem
	err = json.Unmarshal(dataBytes, &item)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedItem.ItemKey, item.ItemKey)
	assert.Equal(suite.T(), expectedItem.ItemValue, item.ItemValue)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestCreateCategory_Success 测试成功创建字典分类
func (suite *DictHandlerTestSuite) TestCreateCategory_Success() {
	reqBody := service.CreateCategoryRequest{
		Code:        "new_category",
		Name:        "新分类",
		Description: "新分类描述",
		SortOrder:   1,
	}

	expectedCategory := &model.DictCategory{
		BaseModel:   model.BaseModel{ID: 1},
		Code:        reqBody.Code,
		Name:        reqBody.Name,
		Description: reqBody.Description,
		SortOrder:   reqBody.SortOrder,
	}

	// Mock 服务调用
	suite.dictService.On("CreateDictCategory", mock.Anything, &reqBody).Return(expectedCategory, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dict/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, response.Code)
	assert.Equal(suite.T(), "Category created successfully", response.Message)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestCreateCategory_InvalidRequest 测试无效请求创建分类
func (suite *DictHandlerTestSuite) TestCreateCategory_InvalidRequest() {
	// Mock 日志
	suite.logger.On("Error", "Failed to bind request", "error", mock.AnythingOfType("*json.SyntaxError"))

	// 创建无效JSON请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dict/categories", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// 验证mock调用
	suite.logger.AssertExpectations(suite.T())
}

// TestCreateItem_Success 测试成功创建字典项
func (suite *DictHandlerTestSuite) TestCreateItem_Success() {
	reqBody := service.CreateItemRequest{
		CategoryCode: "test_category",
		ItemKey:      "new_key",
		ItemValue:    "新值",
		Description:  "新字典项描述",
		SortOrder:    1,
		IsActive:     boolPtr(true),
	}

	expectedItem := &model.DictItem{
		BaseModel:    model.BaseModel{ID: 1},
		CategoryCode: reqBody.CategoryCode,
		ItemKey:      reqBody.ItemKey,
		ItemValue:    reqBody.ItemValue,
		Description:  reqBody.Description,
		SortOrder:    reqBody.SortOrder,
		IsActive:     reqBody.IsActive,
	}

	// Mock 服务调用
	suite.dictService.On("CreateDictItem", mock.Anything, &reqBody).Return(expectedItem, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dict/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, response.Code)
	assert.Equal(suite.T(), "Dict item created successfully", response.Message)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestUpdateItem_Success 测试成功更新字典项
func (suite *DictHandlerTestSuite) TestUpdateItem_Success() {
	itemID := uint(1)
	reqBody := service.UpdateItemRequest{
		ItemValue:   "更新的值",
		Description: "更新的描述",
		SortOrder:   2,
		IsActive:    boolPtr(false),
	}

	expectedItem := &model.DictItem{
		BaseModel:    model.BaseModel{ID: itemID},
		CategoryCode: "test_category",
		ItemKey:      "test_key",
		ItemValue:    reqBody.ItemValue,
		Description:  reqBody.Description,
		SortOrder:    reqBody.SortOrder,
		IsActive:     reqBody.IsActive,
	}

	// Mock 服务调用
	suite.dictService.On("UpdateDictItem", mock.Anything, itemID, &reqBody).Return(expectedItem, nil)

	// 创建请求
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/dict/items/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.Equal(suite.T(), "Dict item updated successfully", response.Message)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestDeleteItem_Success 测试成功删除字典项
func (suite *DictHandlerTestSuite) TestDeleteItem_Success() {
	itemID := uint(1)

	// Mock 服务调用
	suite.dictService.On("DeleteDictItem", mock.Anything, itemID).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/dict/items/1", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestInitDefaultData_Success 测试成功初始化默认数据
func (suite *DictHandlerTestSuite) TestInitDefaultData_Success() {
	// Mock 服务调用
	suite.dictService.On("InitDefaultDictData", mock.Anything).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/dict/init", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.Equal(suite.T(), "Default dictionary data initialized successfully", response.Message)

	// 验证mock调用
	suite.dictService.AssertExpectations(suite.T())
}

// TestDictHandlerTestSuite 运行测试套件
func TestDictHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DictHandlerTestSuite))
}
