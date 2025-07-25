package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // 纯Go SQLite驱动

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	testConfig "vibe-coding-starter/test/config"
)

// TestDatabase 测试数据库包装器
type TestDatabase struct {
	DB     *gorm.DB
	config *testConfig.TestConfig
}

// NewTestDatabase 创建测试数据库连接
func NewTestDatabase(t *testing.T) *TestDatabase {
	config := testConfig.NewTestConfig()

	// 创建GORM连接 (SQLite内存数据库，无需等待)
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        ":memory:",
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 测试时静默日志
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// SQLite连接池配置
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	testDB := &TestDatabase{
		DB:     db,
		config: config,
	}

	// 自动迁移测试表
	testDB.Migrate(t)

	return testDB
}

// Migrate 执行数据库迁移
func (td *TestDatabase) Migrate(t *testing.T) {
	err := td.DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.Comment{},
		&model.File{},
		&model.DictCategory{},
		&model.DictItem{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
}

// Clean 清理测试数据
func (td *TestDatabase) Clean(t *testing.T) {
	// 按依赖关系顺序删除数据
	tables := []string{
		"article_tags",
		"comments",
		"files",
		"articles",
		"tags",
		"categories",
		"users",
		"dict_items",
		"dict_categories",
	}

	for _, table := range tables {
		if err := td.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			t.Logf("Warning: Failed to clean table %s: %v", table, err)
		}
	}
}

// Close 关闭数据库连接
func (td *TestDatabase) Close() error {
	sqlDB, err := td.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB 获取GORM数据库实例
func (td *TestDatabase) GetDB() *gorm.DB {
	return td.DB
}

// Health 健康检查
func (td *TestDatabase) Health() error {
	sqlDB, err := td.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// BeginTx 开始事务
func (td *TestDatabase) BeginTx(ctx context.Context) *gorm.DB {
	return td.DB.WithContext(ctx).Begin()
}

// CreateTestDatabase 实现database.Database接口
func (td *TestDatabase) CreateTestDatabase() database.Database {
	return &testDatabaseAdapter{td}
}

// testDatabaseAdapter 适配器，实现database.Database接口
type testDatabaseAdapter struct {
	*TestDatabase
}

// GetDB 获取 GORM 数据库实例
func (tda *testDatabaseAdapter) GetDB() *gorm.DB {
	return tda.DB
}

// Close 关闭数据库连接
func (tda *testDatabaseAdapter) Close() error {
	sqlDB, err := tda.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库表
func (tda *testDatabaseAdapter) AutoMigrate(dst ...interface{}) error {
	return tda.DB.AutoMigrate(dst...)
}

// Health 检查数据库健康状态
func (tda *testDatabaseAdapter) Health() error {
	sqlDB, err := tda.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// 注意：SQLite使用内存数据库，不需要DSN配置

// SeedTestData 插入测试数据
func (td *TestDatabase) SeedTestData(t *testing.T) {
	// 创建测试用户
	users := []*model.User{
		{
			Username: "testuser1",
			Email:    "test1@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjPeGvGzjYwSgjkMIUmEWqbp4YBARC", // secret
			Nickname: "Test User 1",
			Role:     model.UserRoleUser,
			Status:   model.UserStatusActive,
		},
		{
			Username: "testuser2",
			Email:    "test2@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjPeGvGzjYwSgjkMIUmEWqbp4YBARC", // secret
			Nickname: "Test User 2",
			Role:     model.UserRoleAdmin,
			Status:   model.UserStatusActive,
		},
	}

	for _, user := range users {
		if err := td.DB.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// 创建测试分类
	categories := []*model.Category{
		{
			Name:        "Technology",
			Slug:        "technology",
			Description: "Technology related articles",
		},
		{
			Name:        "Programming",
			Slug:        "programming",
			Description: "Programming tutorials and tips",
		},
	}

	for _, category := range categories {
		if err := td.DB.Create(category).Error; err != nil {
			t.Fatalf("Failed to create test category: %v", err)
		}
	}

	// 创建测试标签
	tags := []*model.Tag{
		{
			Name: "Go",
			Slug: "go",
		},
		{
			Name: "Web Development",
			Slug: "web-development",
		},
	}

	for _, tag := range tags {
		if err := td.DB.Create(tag).Error; err != nil {
			t.Fatalf("Failed to create test tag: %v", err)
		}
	}
}
