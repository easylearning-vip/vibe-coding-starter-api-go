package generator

import (
	"strings"
	"time"
)

// Generator 接口定义
type Generator interface {
	Generate(config interface{}) error
}

// ModuleConfig 模块生成配置
type ModuleConfig struct {
	Name      string
	Fields    string
	WithAuth  bool
	WithCache bool
}

// HandlerConfig 处理器生成配置
type HandlerConfig struct {
	Model          string
	WithAuth       bool
	WithValidation bool
}

// ServiceConfig 服务生成配置
type ServiceConfig struct {
	Name      string
	Model     string
	Fields    string // 字段定义字符串
	WithCache bool
}

// RepositoryConfig 仓储生成配置
type RepositoryConfig struct {
	Name  string
	Model string
}

// ModelConfig 模型生成配置
type ModelConfig struct {
	Name           string
	Fields         string
	WithTimestamps bool
}

// TestConfig 测试生成配置
type TestConfig struct {
	Service    string
	Handler    string
	Repository string
	Type       string // unit, integration, e2e
}

// MigrationConfig 迁移生成配置
type MigrationConfig struct {
	Name   string
	Table  string
	Action string // create, alter, drop
}

// DatabaseTableConfig 数据库表生成配置
type DatabaseTableConfig struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	TableName        string
	ModelName        string // 可选，如果为空则从表名生成
	WithTimestamps   bool
	WithSoftDelete   bool
}

// FrontendFramework 前端框架类型
type FrontendFramework string

const (
	FrameworkAntd FrontendFramework = "antd"
	FrameworkVue  FrontendFramework = "vue"
)

// ModuleType 模块类型
type ModuleType string

const (
	ModuleTypeAdmin  ModuleType = "admin"  // 管理后台模块
	ModuleTypePublic ModuleType = "public" // 普通用户模块
)

// FrontendConfig 前端代码生成配置
type FrontendConfig struct {
	Model      string            // 模型名称
	Framework  FrontendFramework // 前端框架类型
	OutputDir  string            // 输出目录
	ModuleType ModuleType        // 模块类型 (admin/public)
	WithAuth   bool              // 是否包含权限控制
	WithSearch bool              // 是否包含搜索功能
	WithExport bool              // 是否包含导出功能
	WithBatch  bool              // 是否包含批量操作
	ApiPrefix  string            // API 前缀，默认 /api/v1
	ModuleName string            // 模块名称，用于路由和菜单
}

// mapGoType 映射Go类型
func mapGoType(fieldType string) string {
	typeMap := map[string]string{
		"string":    "string",
		"int":       "int",
		"int32":     "int32",
		"int64":     "int64",
		"uint":      "uint",
		"uint32":    "uint32",
		"uint64":    "uint64",
		"float32":   "float32",
		"float64":   "float64",
		"bool":      "bool",
		"time":      "time.Time",
		"datetime":  "time.Time",
		"timestamp": "time.Time",
		"text":      "string",
		"json":      "string",
		"decimal":   "float64",
	}

	if goType, exists := typeMap[fieldType]; exists {
		return goType
	}
	return fieldType
}

// 字符串工具函数

// ToPascalCase 转换为PascalCase
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

// ToCamelCase 转换为camelCase
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// ToSnakeCase 转换为snake_case
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}
	return string(result)
}

// ToKebabCase 转换为kebab-case
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// Pluralize 复数形式
func Pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}
	return s + "s"
}

// GenerateTimestamp 生成时间戳
func GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

// GetCurrentYear 获取当前年份
func GetCurrentYear() int {
	return time.Now().Year()
}
