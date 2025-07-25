package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/test/testutil"
)

// FileRepositoryTestSuite 文件仓储测试套件
type FileRepositoryTestSuite struct {
	suite.Suite
	db     *testutil.TestDatabase
	cache  *testutil.TestCache
	logger *testutil.TestLogger
	repo   repository.FileRepository
	ctx    context.Context
	owner  *model.User
}

// SetupSuite 设置测试套件
func (suite *FileRepositoryTestSuite) SetupSuite() {
	suite.db = testutil.NewTestDatabase(suite.T())
	suite.cache = testutil.NewTestCache(suite.T())
	suite.logger = testutil.NewTestLogger(suite.T())
	suite.ctx = context.Background()

	// 创建文件仓储
	suite.repo = repository.NewFileRepository(
		suite.db.CreateTestDatabase(),
		suite.logger.CreateTestLogger(),
	)
}

// TearDownSuite 清理测试套件
func (suite *FileRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
	suite.cache.Close()
	suite.logger.Close()
}

// SetupTest 每个测试前的设置
func (suite *FileRepositoryTestSuite) SetupTest() {
	suite.db.Clean(suite.T())
	suite.cache.Clean(suite.T())
	suite.createTestData()
}

// createTestData 创建测试数据
func (suite *FileRepositoryTestSuite) createTestData() {
	// 创建文件所有者
	suite.owner = &model.User{
		Username: "fileowner",
		Email:    "owner@example.com",
		Password: "password123",
		Nickname: "File Owner",
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}
	err := suite.db.GetDB().Create(suite.owner).Error
	require.NoError(suite.T(), err)
}

// TestCreate 测试创建文件
func (suite *FileRepositoryTestSuite) TestCreate() {
	file := &model.File{
		Name:         "test-file.jpg",
		OriginalName: "original-test.jpg",
		Path:         "/uploads/test-file.jpg",
		URL:          "/files/test-file.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		Extension:    "jpg",
		Hash:         "abcd1234567890",
		StorageType:  model.StorageTypeLocal,
		OwnerID:      suite.owner.ID,
		IsPublic:     true,
	}

	err := suite.repo.Create(suite.ctx, file)
	require.NoError(suite.T(), err)
	assert.NotZero(suite.T(), file.ID)
	assert.NotZero(suite.T(), file.CreatedAt)
	assert.NotZero(suite.T(), file.UpdatedAt)
}

// TestCreateDuplicateHash 测试创建重复hash的文件
func (suite *FileRepositoryTestSuite) TestCreateDuplicateHash() {
	hash := "duplicate-hash-123"

	file1 := &model.File{
		Name:         "file1.jpg",
		OriginalName: "original1.jpg",
		Path:         "/uploads/file1.jpg",
		Size:         1024,
		MimeType:     "image/jpeg",
		Extension:    "jpg",
		Hash:         hash,
		StorageType:  model.StorageTypeLocal,
		OwnerID:      suite.owner.ID,
	}

	file2 := &model.File{
		Name:         "file2.jpg",
		OriginalName: "original2.jpg",
		Path:         "/uploads/file2.jpg",
		Size:         2048,
		MimeType:     "image/jpeg",
		Extension:    "jpg",
		Hash:         hash, // 相同hash
		StorageType:  model.StorageTypeLocal,
		OwnerID:      suite.owner.ID,
	}

	err := suite.repo.Create(suite.ctx, file1)
	require.NoError(suite.T(), err)

	err = suite.repo.Create(suite.ctx, file2)
	assert.Error(suite.T(), err)
}

// TestGetByID 测试根据ID获取文件
func (suite *FileRepositoryTestSuite) TestGetByID() {
	// 创建测试文件
	file := suite.createTestFile("test-file.jpg", "test-hash-123")

	// 获取文件
	foundFile, err := suite.repo.GetByID(suite.ctx, file.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), file.ID, foundFile.ID)
	assert.Equal(suite.T(), file.Name, foundFile.Name)
	assert.Equal(suite.T(), file.Hash, foundFile.Hash)
	assert.NotNil(suite.T(), foundFile.Owner)
	assert.Equal(suite.T(), suite.owner.ID, foundFile.Owner.ID)
}

// TestGetByIDNotFound 测试获取不存在的文件
func (suite *FileRepositoryTestSuite) TestGetByIDNotFound() {
	_, err := suite.repo.GetByID(suite.ctx, 999)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "file not found")
}

// TestGetByHash 测试根据hash获取文件
func (suite *FileRepositoryTestSuite) TestGetByHash() {
	// 创建测试文件
	file := suite.createTestFile("test-file.jpg", "unique-hash-456")

	// 根据hash获取文件
	foundFile, err := suite.repo.GetByHash(suite.ctx, file.Hash)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), file.ID, foundFile.ID)
	assert.Equal(suite.T(), file.Hash, foundFile.Hash)
}

// TestGetByHashNotFound 测试获取不存在hash的文件
func (suite *FileRepositoryTestSuite) TestGetByHashNotFound() {
	_, err := suite.repo.GetByHash(suite.ctx, "nonexistent-hash")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "file not found")
}

// TestUpdate 测试更新文件
func (suite *FileRepositoryTestSuite) TestUpdate() {
	// 创建测试文件
	file := suite.createTestFile("test-file.jpg", "update-hash-789")

	// 更新文件信息
	file.IsPublic = false
	file.DownloadCount = 5

	err := suite.repo.Update(suite.ctx, file)
	require.NoError(suite.T(), err)

	// 验证更新
	updatedFile, err := suite.repo.GetByID(suite.ctx, file.ID)
	require.NoError(suite.T(), err)
	assert.False(suite.T(), updatedFile.IsPublic)
	assert.Equal(suite.T(), 5, updatedFile.DownloadCount)
}

// TestDelete 测试删除文件
func (suite *FileRepositoryTestSuite) TestDelete() {
	// 创建测试文件
	file := suite.createTestFile("test-file.jpg", "delete-hash-101")

	// 删除文件
	err := suite.repo.Delete(suite.ctx, file.ID)
	require.NoError(suite.T(), err)

	// 验证删除
	_, err = suite.repo.GetByID(suite.ctx, file.ID)
	assert.Error(suite.T(), err)
}

// TestList 测试获取文件列表
func (suite *FileRepositoryTestSuite) TestList() {
	// 创建多个测试文件
	files := []struct {
		name    string
		hash    string
		public  bool
		storage string
	}{
		{"file1.jpg", "hash1", true, model.StorageTypeLocal},
		{"file2.png", "hash2", false, model.StorageTypeLocal},
		{"file3.pdf", "hash3", true, model.StorageTypeS3},
	}

	for _, f := range files {
		file := suite.createTestFile(f.name, f.hash)
		file.IsPublic = f.public
		file.StorageType = f.storage
		err := suite.repo.Update(suite.ctx, file)
		require.NoError(suite.T(), err)
	}

	// 测试基本列表查询
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), total)
	assert.Len(suite.T(), result, 3)
}

// TestGetByOwner 测试根据所有者获取文件列表
func (suite *FileRepositoryTestSuite) TestGetByOwner() {
	// 创建测试文件
	suite.createTestFile("file1.jpg", "owner-hash-1")
	suite.createTestFile("file2.png", "owner-hash-2")

	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := suite.repo.GetByOwner(suite.ctx, suite.owner.ID, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(2), total)
	assert.Len(suite.T(), result, 2)

	// 验证所有文件都属于指定所有者
	for _, file := range result {
		assert.Equal(suite.T(), suite.owner.ID, file.OwnerID)
	}
}

// TestListWithFilters 测试带过滤器的文件列表
func (suite *FileRepositoryTestSuite) TestListWithFilters() {
	// 创建不同类型的文件
	suite.createTestFileWithType("public.jpg", "public-hash", true, model.StorageTypeLocal)
	suite.createTestFileWithType("private.png", "private-hash", false, model.StorageTypeLocal)

	// 测试公开文件过滤
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"is_public": true,
		},
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), result, 1)
	assert.True(suite.T(), result[0].IsPublic)
}

// TestListWithSearch 测试带搜索的文件列表
func (suite *FileRepositoryTestSuite) TestListWithSearch() {
	// 创建测试文件
	suite.createTestFile("document.pdf", "doc-hash")
	suite.createTestFile("image.jpg", "img-hash")

	// 测试搜索
	opts := repository.ListOptions{
		Page:     1,
		PageSize: 10,
		Search:   "document",
	}

	result, total, err := suite.repo.List(suite.ctx, opts)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), total, int64(0))
	assert.Greater(suite.T(), len(result), 0)
}

// createTestFile 创建测试文件
func (suite *FileRepositoryTestSuite) createTestFile(name, hash string) *model.File {
	file := &model.File{
		Name:         name,
		OriginalName: fmt.Sprintf("original-%s", name),
		Path:         fmt.Sprintf("/uploads/%s", name),
		URL:          fmt.Sprintf("/files/%s", name),
		Size:         1024,
		MimeType:     "application/octet-stream",
		Extension:    "bin",
		Hash:         hash,
		StorageType:  model.StorageTypeLocal,
		OwnerID:      suite.owner.ID,
		IsPublic:     false,
	}

	err := suite.repo.Create(suite.ctx, file)
	require.NoError(suite.T(), err)
	return file
}

// createTestFileWithType 创建指定类型的测试文件
func (suite *FileRepositoryTestSuite) createTestFileWithType(name, hash string, isPublic bool, storageType string) *model.File {
	file := suite.createTestFile(name, hash)
	file.IsPublic = isPublic
	file.StorageType = storageType
	err := suite.repo.Update(suite.ctx, file)
	require.NoError(suite.T(), err)
	return file
}

// TestFileRepositoryTestSuite 运行测试套件
func TestFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileRepositoryTestSuite))
}
