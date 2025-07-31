package generator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

// AppConfig 应用配置
type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

// ConfigReader 配置读取器
type ConfigReader struct {
	configPath string
}

// NewConfigReader 创建配置读取器
func NewConfigReader(configPath string) *ConfigReader {
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return &ConfigReader{
		configPath: configPath,
	}
}

// ReadDatabaseConfig 读取数据库配置
func (r *ConfigReader) ReadDatabaseConfig() (*DatabaseConfig, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// 如果默认配置文件不存在，尝试其他配置文件
		alternatives := []string{
			"configs/config-k3d.yaml",
			"configs/config-docker.yaml",
		}

		found := false
		for _, alt := range alternatives {
			if _, err := os.Stat(alt); err == nil {
				r.configPath = alt
				found = true
				break
			}
		}

		if !found {
			// 返回默认配置
			return &DatabaseConfig{
				Driver: "mysql",
			}, nil
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config.Database, nil
}

// GetMigrationDir 获取迁移文件目录
func (r *ConfigReader) GetMigrationDir() (string, error) {
	dbConfig, err := r.ReadDatabaseConfig()
	if err != nil {
		return "", err
	}

	switch dbConfig.Driver {
	case "postgres", "postgresql":
		return "migrations/postgres", nil
	case "mysql":
		return "migrations/mysql", nil
	case "sqlite", "sqlite3":
		return "migrations/sqlite", nil
	default:
		return "migrations", nil
	}
}

// GetDatabaseType 获取数据库类型
func (r *ConfigReader) GetDatabaseType() (string, error) {
	dbConfig, err := r.ReadDatabaseConfig()
	if err != nil {
		return "", err
	}

	switch dbConfig.Driver {
	case "postgres", "postgresql":
		return "postgres", nil
	case "mysql":
		return "mysql", nil
	case "sqlite", "sqlite3":
		return "sqlite", nil
	default:
		return "mysql", nil // 默认为 MySQL
	}
}
