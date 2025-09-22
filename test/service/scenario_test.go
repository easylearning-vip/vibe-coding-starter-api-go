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

type ScenarioServiceTestSuite struct {
	suite.Suite
	service    service.ScenarioService
	mockRepo   *mocks.MockScenarioRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *ScenarioServiceTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockScenarioRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewScenarioService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *ScenarioServiceTestSuite) TestCreate_Success() {
	// 准备测试数据
	req := &service.CreateScenarioRequest{
		Name:        "Test Scenario",
		Description: "Test Description",
	}
	
	expectedScenario := &model.Scenario{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Scenario")).Return(nil).Run(func(args mock.Arguments) {
		scenario := args.Get(1).(*model.Scenario)
		scenario.ID = 1
	})
	
	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedScenario.Name, result.Name)
	assert.Equal(suite.T(), expectedScenario.Description, result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ScenarioServiceTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedScenario := &model.Scenario{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test Scenario",
		Description: "Test Description",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedScenario, nil)
	
	// 执行测试
	result, err := suite.service.GetByID(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedScenario, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ScenarioServiceTestSuite) TestUpdate_Success() {
	// 准备测试数据
	existingScenario := &model.Scenario{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Old Name",
		Description: "Old Description",
	}
	
	newName := "New Name"
	req := &service.UpdateScenarioRequest{
		Name: &newName,
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingScenario, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.Scenario")).Return(nil)
	
	// 执行测试
	result, err := suite.service.Update(suite.ctx, 1, req)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, result.Name)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ScenarioServiceTestSuite) TestDelete_Success() {
	// 准备测试数据
	existingScenario := &model.Scenario{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Scenario",
	}
	
	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingScenario, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)
	
	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ScenarioServiceTestSuite) TestList_Success() {
	// 准备测试数据
	expectedscenarios := []*model.Scenario{
		{BaseModel: model.BaseModel{ID: 1}, Name: "Scenario 1"},
		{BaseModel: model.BaseModel{ID: 2}, Name: "Scenario 2"},
	}
	
	opts := &service.ListScenarioOptions{
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
	suite.mockRepo.On("List", suite.ctx, repoOpts).Return(expectedscenarios, int64(2), nil)
	
	// 执行测试
	result, total, err := suite.service.List(suite.ctx, opts)
	
	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedscenarios, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestScenarioServiceSuite(t *testing.T) {
	suite.Run(t, new(ScenarioServiceTestSuite))
}
