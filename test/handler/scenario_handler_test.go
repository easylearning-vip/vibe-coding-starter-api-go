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

type ScenarioHandlerTestSuite struct {
	suite.Suite
	handler     *handler.ScenarioHandler
	mockService *mocks.MockScenarioService
	mockLogger  *mocks.MockLogger
	router      *gin.Engine
}

func (suite *ScenarioHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockService = &mocks.MockScenarioService{}
	suite.mockLogger = &mocks.MockLogger{}

	suite.handler = handler.NewScenarioHandler(
		suite.mockService,
		suite.mockLogger,
	)
	
	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(v1)
}

func (suite *ScenarioHandlerTestSuite) TestCreateScenario_Success() {
	// 准备测试数据
	req := service.CreateScenarioRequest{
		Name:        "Test Scenario",
		Description: "Test Description",
	}
	
	expectedScenario := &model.Scenario{
		BaseModel: model.BaseModel{ID: 1},
		Name:      req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockService.On("Create", mock.Anything, &req).Return(expectedScenario, nil)
	
	// 准备请求
	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/scenarios", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ScenarioHandlerTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedScenario := &model.Scenario{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Scenario",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedScenario, nil)
	
	// 准备请求
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/scenarios/1", nil)
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestScenarioHandlerSuite(t *testing.T) {
	suite.Run(t, new(ScenarioHandlerTestSuite))
}
