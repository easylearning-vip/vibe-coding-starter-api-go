package handler_test

import (
	"bytes"
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

type DifficultyLevelHandlerTestSuite struct {
	suite.Suite
	handler     *handler.DifficultyLevelHandler
	mockService *mocks.MockDifficultyLevelService
	mockLogger  *mocks.MockLogger
	router      *gin.Engine
}

func (suite *DifficultyLevelHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockService = &mocks.MockDifficultyLevelService{}
	suite.mockLogger = &mocks.MockLogger{}

	suite.handler = handler.NewDifficultyLevelHandler(
		suite.mockService,
		suite.mockLogger,
	)
	
	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(v1)
}

func (suite *DifficultyLevelHandlerTestSuite) TestCreateDifficultyLevel_Success() {
	// 准备测试数据
	req := service.CreateDifficultyLevelRequest{
		Name:        "Test DifficultyLevel",
		Description: "Test Description",
	}
	
	expectedDifficultyLevel := &model.DifficultyLevel{
		BaseModel: model.BaseModel{ID: 1},
		Name:      req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockService.On("Create", mock.Anything, &req).Return(expectedDifficultyLevel, nil)
	
	// 准备请求
	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/difficultylevels", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *DifficultyLevelHandlerTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedDifficultyLevel := &model.DifficultyLevel{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test DifficultyLevel",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedDifficultyLevel, nil)
	
	// 准备请求
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/difficultylevels/1", nil)
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestDifficultyLevelHandlerSuite(t *testing.T) {
	suite.Run(t, new(DifficultyLevelHandlerTestSuite))
}
