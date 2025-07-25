package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用程序配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	AI       AIConfig       `mapstructure:"ai"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Security SecurityConfig `mapstructure:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Issuer     string `mapstructure:"issuer"`
	Expiration int    `mapstructure:"expiration"`
}

// AIConfig AI 辅助开发配置
type AIConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	Temperature float32 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableHTTPS           bool     `mapstructure:"enable_https"`
	TLSCertFile           string   `mapstructure:"tls_cert_file"`
	TLSKeyFile            string   `mapstructure:"tls_key_file"`
	HSTSMaxAge            int      `mapstructure:"hsts_max_age"`
	HSTSIncludeSubDomains bool     `mapstructure:"hsts_include_subdomains"`
	HSTSPreload           bool     `mapstructure:"hsts_preload"`
	CSPPolicy             string   `mapstructure:"csp_policy"`
	IPWhitelist           []string `mapstructure:"ip_whitelist"`
	IPBlacklist           []string `mapstructure:"ip_blacklist"`
	BlockedUserAgents     []string `mapstructure:"blocked_user_agents"`
	MaxRequestSize        int64    `mapstructure:"max_request_size"`
	RequestTimeout        int      `mapstructure:"request_timeout"`
}

// New 创建新的配置实例
func New() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置环境变量前缀
	viper.SetEnvPrefix("VIBE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在时使用默认值
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 60)

	// 数据库默认配置
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "vibe_coding_starter")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	// 缓存默认配置
	viper.SetDefault("cache.driver", "redis")
	viper.SetDefault("cache.host", "localhost")
	viper.SetDefault("cache.port", 6379)
	viper.SetDefault("cache.password", "")
	viper.SetDefault("cache.database", 0)
	viper.SetDefault("cache.pool_size", 10)

	// 日志默认配置
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("logger.output", "stdout")
	viper.SetDefault("logger.filename", "logs/app.log")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_age", 30)
	viper.SetDefault("logger.max_backups", 10)
	viper.SetDefault("logger.compress", true)

	// JWT 默认配置
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.issuer", "vibe-coding-starter")
	viper.SetDefault("jwt.expiration", 86400) // 24 hours

	// AI 默认配置
	viper.SetDefault("ai.enabled", false)
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-3.5-turbo")
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 1000)

	// CORS 默认配置
	viper.SetDefault("cors.allow_origins", []string{"http://localhost:3000", "http://localhost:3001"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Request-ID"})
	viper.SetDefault("cors.expose_headers", []string{"Content-Length", "X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 43200) // 12 hours

	// Security 默认配置
	viper.SetDefault("security.enable_https", false)
	viper.SetDefault("security.hsts_max_age", 31536000) // 1 year
	viper.SetDefault("security.hsts_include_subdomains", true)
	viper.SetDefault("security.hsts_preload", false)
	viper.SetDefault("security.csp_policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
	viper.SetDefault("security.max_request_size", 10485760) // 10MB
	viper.SetDefault("security.request_timeout", 30)        // 30 seconds
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&multiStatements=true",
			c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.Username, c.Password, c.Database)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}

// GetAddress 获取服务器地址
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetRedisAddress 获取 Redis 地址
func (c *CacheConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 设置默认值
	setDefaults()

	// 设置配置文件路径
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		viper.AddConfigPath("../../configs")
	}

	// 读取环境变量
	viper.SetEnvPrefix("VIBE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置
			fmt.Println("Config file not found, using default configuration")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
