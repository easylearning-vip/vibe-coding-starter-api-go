package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// FileHandlerTestSuite 文件处理器测试套件
type FileHandlerTestSuite struct {
	suite.Suite
	fileService *mocks.MockFileService
	logger      *mocks.MockLogger
	handler     *handler.FileHandler
	router      *gin.Engine
}

// SetupSuite 设置测试套件
func (suite *FileHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	suite.fileService = new(mocks.MockFileService)
	suite.logger = new(mocks.MockLogger)

	// 创建文件处理器
	suite.handler = handler.NewFileHandler(
		suite.fileService,
		suite.logger,
	)

	// 设置路由
	suite.router = gin.New()
	api := suite.router.Group("/api/v1")
	suite.handler.RegisterRoutes(api)
}

// SetupTest 每个测试前的设置
func (suite *FileHandlerTestSuite) SetupTest() {
	suite.fileService.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestUpload 测试上传文件
func (suite *FileHandlerTestSuite) TestUpload() {
	userID := uint(1)
	file := &model.File{
		BaseModel:    model.BaseModel{ID: 1},
		Name:         "test-file.jpg",
		OriginalName: "test.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		IsPublic:     true,
		OwnerID:      userID,
	}

	// Mock 文件服务
	suite.fileService.On("Upload", mock.Anything, mock.AnythingOfType("*service.UploadRequest")).Return(file, nil)

	// 创建multipart请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.WriteField("is_public", "true")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)

	// 直接调用处理器方法
	suite.handler.Upload(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response model.File
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), file.ID, response.ID)
	assert.Equal(suite.T(), file.Name, response.Name)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
}

// TestUploadNoFile 测试上传时没有文件
func (suite *FileHandlerTestSuite) TestUploadNoFile() {
	userID := uint(1)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	// 创建空的multipart请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)

	// 直接调用处理器方法
	suite.handler.Upload(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "file_required", response.Error)

	// 验证mock调用
	suite.logger.AssertExpectations(suite.T())
}

// TestGetByID 测试根据ID获取文件
func (suite *FileHandlerTestSuite) TestGetByID() {
	fileID := uint(1)
	file := &model.File{
		BaseModel:    model.BaseModel{ID: fileID},
		Name:         "test-file.jpg",
		OriginalName: "test.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		IsPublic:     true,
	}

	// Mock 文件服务
	suite.fileService.On("GetByID", mock.Anything, fileID).Return(file, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files/"+strconv.Itoa(int(fileID)), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response model.File
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), file.ID, response.ID)
	assert.Equal(suite.T(), file.Name, response.Name)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
}

// TestGetByIDNotFound 测试获取不存在的文件
func (suite *FileHandlerTestSuite) TestGetByIDNotFound() {
	fileID := uint(999)

	// Mock 文件服务返回错误
	suite.fileService.On("GetByID", mock.Anything, fileID).Return(nil, assert.AnError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files/"+strconv.Itoa(int(fileID)), nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "file_not_found", response.Error)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestDownload 测试下载文件
func (suite *FileHandlerTestSuite) TestDownload() {
	fileID := uint(1)
	file := &model.File{
		BaseModel:    model.BaseModel{ID: fileID},
		Name:         "test-file.jpg",
		OriginalName: "test.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		IsPublic:     true,
	}

	downloadResponse := &service.DownloadResponse{
		File:     file,
		FileData: []byte("fake image data"),
	}

	// Mock 文件服务
	suite.fileService.On("Download", mock.Anything, fileID).Return(downloadResponse, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files/"+strconv.Itoa(int(fileID))+"/download", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), file.MimeType, w.Header().Get("Content-Type"))
	assert.Equal(suite.T(), "attachment; filename="+file.OriginalName, w.Header().Get("Content-Disposition"))
	assert.Equal(suite.T(), downloadResponse.FileData, w.Body.Bytes())

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
}

// TestList 测试获取文件列表
func (suite *FileHandlerTestSuite) TestList() {
	files := []*model.File{
		{
			BaseModel:    model.BaseModel{ID: 1},
			Name:         "file1.jpg",
			OriginalName: "file1.jpg",
			Size:         1024,
			IsPublic:     true,
		},
		{
			BaseModel:    model.BaseModel{ID: 2},
			Name:         "file2.png",
			OriginalName: "file2.png",
			Size:         2048,
			IsPublic:     false,
		},
	}

	// Mock 文件服务
	suite.fileService.On("List", mock.Anything, mock.AnythingOfType("repository.ListOptions")).Return(files, int64(2), nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/files?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.ListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), response.Total)
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 10, response.Size)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
}

// TestDelete 测试删除文件
func (suite *FileHandlerTestSuite) TestDelete() {
	fileID := uint(1)
	userID := uint(1)

	// Mock 文件服务
	suite.fileService.On("Delete", mock.Anything, fileID).Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/files/"+strconv.Itoa(int(fileID)), nil)
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID和路由参数
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)
	c.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(int(fileID))}}

	// 直接调用处理器方法
	suite.handler.Delete(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response handler.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "File deleted successfully", response.Message)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
}

// TestDeleteServiceError 测试删除文件时服务错误
func (suite *FileHandlerTestSuite) TestDeleteServiceError() {
	fileID := uint(1)
	userID := uint(1)

	// Mock 文件服务返回错误
	suite.fileService.On("Delete", mock.Anything, fileID).Return(assert.AnError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("uint64"), mock.AnythingOfType("string"), mock.Anything).Return()

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/files/"+strconv.Itoa(int(fileID)), nil)
	w := httptest.NewRecorder()

	// 创建Gin上下文并设置用户ID和路由参数
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", userID)
	c.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(int(fileID))}}

	// 直接调用处理器方法
	suite.handler.Delete(c)

	// 验证响应
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response handler.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "delete_failed", response.Error)

	// 验证mock调用
	suite.fileService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestFileHandlerTestSuite 运行测试套件
func TestFileHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FileHandlerTestSuite))
}
