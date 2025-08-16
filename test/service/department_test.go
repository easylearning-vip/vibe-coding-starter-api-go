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

type DepartmentServiceTestSuite struct {
	suite.Suite
	service    service.DepartmentService
	mockRepo   *mocks.MockDepartmentRepository
	mockLogger *mocks.MockLogger
	ctx        context.Context
}

func (suite *DepartmentServiceTestSuite) SetupTest() {
	suite.mockRepo = &mocks.MockDepartmentRepository{}
	suite.mockLogger = &mocks.MockLogger{}
	suite.ctx = context.Background()

	// 设置logger mock的期望调用
	suite.mockLogger.On("Info", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return()

	suite.service = service.NewDepartmentService(
		suite.mockRepo,
		suite.mockLogger,
	)
}

func (suite *DepartmentServiceTestSuite) TestCreate_Success() {
	// 准备测试数据
	req := &service.CreateDepartmentRequest{
		Name:        "Test Department",
		Description: "Test Description",
	}

	expectedDepartment := &model.Department{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        req.Name,
		Description: req.Description,
	}

	// 设置 mock 期望
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*model.Department")).Return(nil).Run(func(args mock.Arguments) {
		department := args.Get(1).(*model.Department)
		department.ID = 1
	})

	// 执行测试
	result, err := suite.service.Create(suite.ctx, req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedDepartment.Name, result.Name)
	assert.Equal(suite.T(), expectedDepartment.Description, result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestGetByID_Success() {
	// 准备测试数据
	expectedDepartment := &model.Department{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Test Department",
		Description: "Test Description",
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedDepartment, nil)

	// 执行测试
	result, err := suite.service.GetByID(suite.ctx, 1)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDepartment, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestUpdate_Success() {
	// 准备测试数据
	existingDepartment := &model.Department{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Old Name",
		Description: "Old Description",
	}

	newName := "New Name"
	req := &service.UpdateDepartmentRequest{
		Name: &newName,
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingDepartment, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*model.Department")).Return(nil)

	// 执行测试
	result, err := suite.service.Update(suite.ctx, 1, req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, result.Name)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestDelete_Success() {
	// 准备测试数据
	existingDepartment := &model.Department{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Department",
	}

	// 设置 mock 期望
	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingDepartment, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// 执行测试
	err := suite.service.Delete(suite.ctx, 1)

	// 验证结果
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestList_Success() {
	// 准备测试数据
	expecteddepartments := []*model.Department{
		{BaseModel: model.BaseModel{ID: 1}, Name: "Department 1"},
		{BaseModel: model.BaseModel{ID: 2}, Name: "Department 2"},
	}

	opts := &service.ListDepartmentOptions{
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
	suite.mockRepo.On("List", suite.ctx, repoOpts).Return(expecteddepartments, int64(2), nil)

	// 执行测试
	result, total, err := suite.service.List(suite.ctx, opts)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expecteddepartments, result)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestDepartmentServiceSuite(t *testing.T) {
	suite.Run(t, new(DepartmentServiceTestSuite))
}
