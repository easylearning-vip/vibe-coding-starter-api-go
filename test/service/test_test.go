package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/test/mocks"
)

type TestServiceTestSuite struct {
	suite.Suite
	service    service.TestService
	mockRepo   *mocks.MockTestRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *TestServiceTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockTestRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewTestService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *TestServiceTestSuite) TestCreate_Success() {
	// 准备测试数据
	req := &service.CreateTestRequest{
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}

	expectedTest := &model.Test{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}

	// 设置 mock 期望
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Test")).Return(nil).Run(func(args mock.Arguments) {
		test := args.Get(1).(*model.Test)
		test.ID = 1
	})

	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedTest.Name, result.Name)
	assert.Equal(suite.T(), expectedTest.Description, result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TestServiceTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedTest := &model.Test{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test Test",
		Description: sql.NullString{String: "Test Description", Valid: true},
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedTest, nil)

	// 执行测试
	result, err := suite.service.GetByID(suite.ctx, 1)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTest, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TestServiceTestSuite) TestUpdate_Success() {
	// 准备测试数据
	existingTest := &model.Test{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Old Name",
		Description: sql.NullString{String: "Old Description", Valid: true},
	}

	newName := "New Name"
	req := &service.UpdateTestRequest{
		Name: &newName,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingTest, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.Test")).Return(nil)

	// 执行测试
	result, err := suite.service.Update(suite.ctx, 1, req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, result.Name)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TestServiceTestSuite) TestDelete_Success() {
	// 准备测试数据
	existingTest := &model.Test{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Test",
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingTest, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)

	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TestServiceTestSuite) TestList_Success() {
	// 准备测试数据
	expectedtests := []*model.Test{
		{BaseModel: model.BaseModel{ID: 1}, Name: "Test 1"},
		{BaseModel: model.BaseModel{ID: 2}, Name: "Test 2"},
	}

	opts := &service.ListTestOptions{
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
	suite.mockRepo.On("List", suite.ctx, repoOpts).Return(expectedtests, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.List(suite.ctx, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedtests, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestTestServiceSuite(t *testing.T) {
	suite.Run(t, new(TestServiceTestSuite))
}
