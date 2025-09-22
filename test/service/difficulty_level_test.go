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

type DifficultyLevelServiceTestSuite struct {
	suite.Suite
	service    service.DifficultyLevelService
	mockRepo   *mocks.MockDifficultyLevelRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *DifficultyLevelServiceTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockDifficultyLevelRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewDifficultyLevelService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *DifficultyLevelServiceTestSuite) TestCreate_Success() {
	// 准备测试数据
	req := &service.CreateDifficultyLevelRequest{
		Name:        "Test DifficultyLevel",
		Description: "Test Description",
	}
	
	expectedDifficultyLevel := &model.DifficultyLevel{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.DifficultyLevel")).Return(nil).Run(func(args mock.Arguments) {
		difficultyLevel := args.Get(1).(*model.DifficultyLevel)
		difficultyLevel.ID = 1
	})
	
	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedDifficultyLevel.Name, result.Name)
	assert.Equal(suite.T(), expectedDifficultyLevel.Description, result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DifficultyLevelServiceTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedDifficultyLevel := &model.DifficultyLevel{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test DifficultyLevel",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedDifficultyLevel, nil)
	
	// 执行测试
	result, err := suite.service.GetByID(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDifficultyLevel, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DifficultyLevelServiceTestSuite) TestUpdate_Success() {
	// 准备测试数据
	existingDifficultyLevel := &model.DifficultyLevel{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Old Name",
		Description: "Old Description",
	}
	
	newName := "New Name"
	req := &service.UpdateDifficultyLevelRequest{
		Name: &newName,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingDifficultyLevel, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.DifficultyLevel")).Return(nil)
	
	// 执行测试
	result, err := suite.service.Update(suite.ctx, 1, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, result.Name)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DifficultyLevelServiceTestSuite) TestDelete_Success() {
	// 准备测试数据
	existingDifficultyLevel := &model.DifficultyLevel{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test DifficultyLevel",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingDifficultyLevel, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)
	
	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DifficultyLevelServiceTestSuite) TestList_Success() {
	// 准备测试数据
	expecteddifficultylevels := []*model.DifficultyLevel{
		{BaseModel: model.BaseModel{ID: 1}, Name: "DifficultyLevel 1"},
		{BaseModel: model.BaseModel{ID: 2}, Name: "DifficultyLevel 2"},
	}
	
	opts := &service.ListDifficultyLevelOptions{
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
	suite.mockRepo.On("List", suite.ctx, repoOpts).Return(expecteddifficultylevels, int64(2), nil)
	
	// 执行测试
	result, total, err := suite.service.List(suite.ctx, opts)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expecteddifficultylevels, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestDifficultyLevelServiceSuite(t *testing.T) {
	suite.Run(t, new(DifficultyLevelServiceTestSuite))
}
