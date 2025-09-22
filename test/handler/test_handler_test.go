package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

type TestHandlerTestSuite struct {
	suite.Suite
	handler     *handler.TestHandler
	mockService *mocks.MockTestService
	mockLogger  *mocks.MockLogger
	router      *gin.Engine
}

func (suite *TestHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockService = &mocks.MockTestService{}
	suite.mockLogger = &mocks.MockLogger{}

	suite.handler = handler.NewTestHandler(
		suite.mockService,
		suite.mockLogger,
	)

	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(v1)
}

func (suite *TestHandlerTestSuite) TestCreateTest_Success() {
	// 准备测试数据
	req := service.CreateTestRequest{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}

	expectedTest := &model.Test{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}

	// 设置 mock 期望
	suite.mockService.On("Create", mock.Anything, &req).Return(expectedTest, nil)

	// 准备请求
	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/tests", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	// 执行请求
	suite.router.ServeHTTP(w, request)

	// 验证结果
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *TestHandlerTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedTest := &model.Test{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}

	// 设置 mock 期望
	suite.mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedTest, nil)

	// 准备请求
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/tests/1", nil)

	// 执行请求
	suite.router.ServeHTTP(w, request)

	// 验证结果
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestTestHandlerSuite(t *testing.T) {
	suite.Run(t, new(TestHandlerTestSuite))
}
