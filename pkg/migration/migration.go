package migration

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"vibe-coding-starter/pkg/logger"
	"gorm.io/gorm"
)

// Migrator 数据库迁移管理器
type Migrator struct {
	db     *gorm.DB
	logger logger.Logger
	config *Config
}

// Config 迁移配置
type Config struct {
	MigrationsPath string // 迁移文件路径
	DatabaseName   string // 数据库名称
	DatabaseDriver string // 数据库驱动类型
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB, logger logger.Logger, config *Config) *Migrator {
	if config.MigrationsPath == "" {
		// 获取项目根目录
		_, filename, _, _ := runtime.Caller(0)
		projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

		// 根据数据库驱动选择迁移路径
		var dbPath string
		switch config.DatabaseDriver {
		case "mysql":
			dbPath = "mysql"
		case "postgres":
			dbPath = "postgresql"
		case "sqlite":
			dbPath = "sqlite"
		default:
			dbPath = "mysql" // 默认使用mysql
		}

		config.MigrationsPath = filepath.Join(projectRoot, "migrations", dbPath)
	}

	return &Migrator{
		db:     db,
		logger: logger,
		config: config,
	}
}

// Up 执行向上迁移
func (m *Migrator) Up() error {
	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	m.logger.Info("Starting database migration...")

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Database migration completed successfully")
	return nil
}

// Down 执行向下迁移
func (m *Migrator) Down() error {
	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	m.logger.Info("Starting database rollback...")

	if err := migrator.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	m.logger.Info("Database rollback completed successfully")
	return nil
}

// Steps 执行指定步数的迁移
func (m *Migrator) Steps(n int) error {
	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	action := "forward"
	if n < 0 {
		action = "backward"
	}

	m.logger.Info("Starting database migration", "steps", n, "direction", action)

	if err := migrator.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	m.logger.Info("Database migration completed successfully", "steps", n)
	return nil
}

// Force 强制设置迁移版本
func (m *Migrator) Force(version int) error {
	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	m.logger.Warn("Forcing migration version", "version", version)

	if err := migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("Migration version forced successfully", "version", version)
	return nil
}

// Version 获取当前迁移版本
func (m *Migrator) Version() (uint, bool, error) {
	migrator, err := m.createMigrator()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Drop 删除所有表
func (m *Migrator) Drop() error {
	migrator, err := m.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	m.logger.Warn("Dropping all database tables...")

	if err := migrator.Drop(); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	m.logger.Info("All database tables dropped successfully")
	return nil
}

// Status 获取迁移状态
func (m *Migrator) Status() error {
	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	if version == 0 {
		m.logger.Info("Database migration status", "version", "No migrations applied", "dirty", dirty)
	} else {
		m.logger.Info("Database migration status", "version", version, "dirty", dirty)
	}

	return nil
}

// createMigrator 创建migrate实例
func (m *Migrator) createMigrator() (*migrate.Migrate, error) {
	// 获取原始数据库连接
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	var driver database.Driver
	var driverName string

	// 根据数据库驱动类型创建相应的驱动实例
	switch m.config.DatabaseDriver {
	case "mysql":
		driver, err = mysql.WithInstance(sqlDB, &mysql.Config{
			DatabaseName: m.config.DatabaseName,
		})
		driverName = "mysql"
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{
			DatabaseName: m.config.DatabaseName,
		})
		driverName = "postgres"
	case "sqlite":
		driver, err = sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
		driverName = "sqlite3"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", m.config.DatabaseDriver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s driver: %w", m.config.DatabaseDriver, err)
	}

	// 创建migrate实例
	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.config.MigrationsPath),
		driverName,
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return migrator, nil
}

// CreateMigrationFile 创建新的迁移文件
func (m *Migrator) CreateMigrationFile(name string) error {
	// 这个功能可以后续实现，用于生成新的迁移文件
	m.logger.Info("Creating migration file", "name", name)
	// TODO: 实现创建迁移文件的逻辑
	return nil
}
