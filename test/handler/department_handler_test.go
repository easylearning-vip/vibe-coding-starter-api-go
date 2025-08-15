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

type DepartmentHandlerTestSuite struct {
	suite.Suite
	handler     *handler.DepartmentHandler
	mockService *mocks.MockDepartmentService
	mockLogger  *mocks.MockLogger
	router      *gin.Engine
}

func (suite *DepartmentHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.mockService = &mocks.MockDepartmentService{}
	suite.mockLogger = &mocks.MockLogger{}

	suite.handler = handler.NewDepartmentHandler(
		suite.mockService,
		suite.mockLogger,
	)
	
	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(v1)
}

func (suite *DepartmentHandlerTestSuite) TestCreateDepartment_Success() {
	// 准备测试数据
	req := service.CreateDepartmentRequest{
		Name:        "Test Department",
		Description: "Test Description",
	}
	
	expectedDepartment := &model.Department{
		BaseModel: model.BaseModel{ID: 1},
		Name:      req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockService.On("Create", mock.Anything, &req).Return(expectedDepartment, nil)
	
	// 准备请求
	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/departments", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *DepartmentHandlerTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedDepartment := &model.Department{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Department",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedDepartment, nil)
	
	// 准备请求
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/departments/1", nil)
	
	// 执行请求
	suite.router.ServeHTTP(w, request)
	
	// 验证结果
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestDepartmentHandlerSuite(t *testing.T) {
	suite.Run(t, new(DepartmentHandlerTestSuite))
}
