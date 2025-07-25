package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"vibe-coding-starter/internal/config"
	appLogger "vibe-coding-starter/pkg/logger"
)

// Database 数据库接口
type Database interface {
	GetDB() *gorm.DB
	Close() error
	AutoMigrate(dst ...interface{}) error
	Health() error
}

// database 数据库实现
type database struct {
	db *gorm.DB
}

// New 创建新的数据库连接
func New(cfg *config.Config, log appLogger.Logger) (Database, error) {
	var dialector gorm.Dialector

	dsn := cfg.Database.GetDSN()
	if dsn == "" {
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// 根据驱动类型选择方言
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected successfully",
		"driver", cfg.Database.Driver,
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.Database,
	)

	return &database{db: db}, nil
}

// GetDB 获取 GORM 数据库实例
func (d *database) GetDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库表
func (d *database) AutoMigrate(dst ...interface{}) error {
	return d.db.AutoMigrate(dst...)
}

// Health 检查数据库健康状态
func (d *database) Health() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
