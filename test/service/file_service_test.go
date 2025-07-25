package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

// FileServiceTestSuite 文件服务测试套件
type FileServiceTestSuite struct {
	suite.Suite
	fileRepo *mocks.MockFileRepository
	cache    *mocks.MockCache
	logger   *mocks.MockLogger
	service  service.FileService
	ctx      context.Context
}

// SetupSuite 设置测试套件
func (suite *FileServiceTestSuite) SetupSuite() {
	suite.fileRepo = new(mocks.MockFileRepository)
	suite.cache = new(mocks.MockCache)
	suite.logger = new(mocks.MockLogger)
	suite.ctx = context.Background()

	// 创建文件服务
	suite.service = service.NewFileService(
		suite.fileRepo,
		suite.cache,
		suite.logger,
	)
}

// SetupTest 每个测试前的设置
func (suite *FileServiceTestSuite) SetupTest() {
	// 重置所有mock
	suite.fileRepo.ExpectedCalls = nil
	suite.cache.ExpectedCalls = nil
	suite.logger.ExpectedCalls = nil
}

// TestUpload 测试上传文件
func (suite *FileServiceTestSuite) TestUpload() {
	req := &service.UploadRequest{
		FileName:    "test.jpg",
		FileSize:    1024,
		MimeType:    "image/jpeg",
		FileData:    []byte("fake image data"),
		IsPublic:    true,
		StorageType: model.StorageTypeLocal,
	}

	// Mock 检查文件是否已存在（不存在）
	suite.fileRepo.On("GetByHash", suite.ctx, mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)

	// Mock 创建文件记录
	suite.fileRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.File")).Return(nil).Run(func(args mock.Arguments) {
		file := args.Get(1).(*model.File)
		file.ID = 1
	})

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行上传
	file, err := suite.service.Upload(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), file)
	assert.Equal(suite.T(), req.FileName, file.OriginalName)
	assert.Equal(suite.T(), req.FileSize, file.Size)
	assert.Equal(suite.T(), req.MimeType, file.MimeType)
	assert.Equal(suite.T(), req.IsPublic, file.IsPublic)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestUploadExistingFile 测试上传已存在的文件
func (suite *FileServiceTestSuite) TestUploadExistingFile() {
	req := &service.UploadRequest{
		FileName: "existing.jpg",
		FileSize: 1024,
		MimeType: "image/jpeg",
		FileData: []byte("fake image data"),
	}

	existingFile := &model.File{
		BaseModel:    model.BaseModel{ID: 1},
		Name:         "existing-file.jpg",
		OriginalName: "existing.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		Hash:         "existing-hash",
	}

	// Mock 检查文件已存在
	suite.fileRepo.On("GetByHash", suite.ctx, mock.AnythingOfType("string")).Return(existingFile, nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行上传
	file, err := suite.service.Upload(suite.ctx, req)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), file)
	assert.Equal(suite.T(), existingFile.ID, file.ID)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestGetByID 测试根据ID获取文件
func (suite *FileServiceTestSuite) TestGetByID() {
	fileID := uint(1)
	file := &model.File{
		BaseModel:    model.BaseModel{ID: fileID},
		Name:         "test-file.jpg",
		OriginalName: "test.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		IsPublic:     true,
	}

	// Mock 获取文件
	suite.fileRepo.On("GetByID", suite.ctx, fileID).Return(file, nil)

	// 执行获取
	result, err := suite.service.GetByID(suite.ctx, fileID)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), file.ID, result.ID)
	assert.Equal(suite.T(), file.Name, result.Name)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
}

// TestGetByIDNotFound 测试获取不存在的文件
func (suite *FileServiceTestSuite) TestGetByIDNotFound() {
	fileID := uint(999)

	// Mock 文件不存在
	suite.fileRepo.On("GetByID", suite.ctx, fileID).Return(nil, gorm.ErrRecordNotFound)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 执行获取
	result, err := suite.service.GetByID(suite.ctx, fileID)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestDownload 测试下载文件
func (suite *FileServiceTestSuite) TestDownload() {
	fileID := uint(1)
	file := &model.File{
		BaseModel:    model.BaseModel{ID: fileID},
		Name:         "test-file.jpg",
		OriginalName: "test.jpg",
		Path:         "/uploads/test-file.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		IsPublic:     true,
	}

	// Mock 获取文件
	suite.fileRepo.On("GetByID", suite.ctx, fileID).Return(file, nil)

	// Mock 日志（当文件读取失败时）
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	// 执行下载
	response, err := suite.service.Download(suite.ctx, fileID)

	// 验证结果（文件不存在，应该失败）
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Contains(suite.T(), err.Error(), "failed to read file")

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestDelete 测试删除文件
func (suite *FileServiceTestSuite) TestDelete() {
	fileID := uint(1)
	file := &model.File{
		BaseModel:   model.BaseModel{ID: fileID},
		Name:        "test-file.jpg",
		Path:        "/uploads/test-file.jpg",
		StorageType: model.StorageTypeLocal,
	}

	// Mock 获取文件
	suite.fileRepo.On("GetByID", suite.ctx, fileID).Return(file, nil)

	// Mock 删除文件记录
	suite.fileRepo.On("Delete", suite.ctx, fileID).Return(nil)

	// Mock 日志
	suite.logger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()
	suite.logger.On("Warn", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	// 执行删除
	err := suite.service.Delete(suite.ctx, fileID)

	// 验证结果
	require.NoError(suite.T(), err)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestList 测试获取文件列表
func (suite *FileServiceTestSuite) TestList() {
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

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	// Mock 获取文件列表
	suite.fileRepo.On("List", suite.ctx, opts).Return(files, int64(2), nil)

	// 执行获取列表
	result, total, err := suite.service.List(suite.ctx, opts)

	// 验证结果
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int64(2), total)
	assert.Equal(suite.T(), files[0].ID, result[0].ID)
	assert.Equal(suite.T(), files[1].ID, result[1].ID)

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
}

// TestListError 测试获取文件列表失败
func (suite *FileServiceTestSuite) TestListError() {
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	expectedError := errors.New("database error")

	// Mock 获取文件列表失败
	suite.fileRepo.On("List", suite.ctx, opts).Return(nil, int64(0), expectedError)

	// Mock 日志
	suite.logger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	// 执行获取列表
	result, total, err := suite.service.List(suite.ctx, opts)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), int64(0), total)
	assert.Contains(suite.T(), err.Error(), "failed to get files")

	// 验证mock调用
	suite.fileRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

// TestFileServiceTestSuite 运行测试套件
func TestFileServiceTestSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}
