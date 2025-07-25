package cmd

import (
	"fmt"
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

// Field 字段定义
type Field struct {
	Name     string
	Type     string
	Tag      string
	JsonTag  string
	GormTag  string
	Comment  string
	Required bool
}

// ParseFields 解析字段字符串
func ParseFields(fieldsStr string) ([]*Field, error) {
	if fieldsStr == "" {
		return []*Field{}, nil
	}

	var fields []*Field
	fieldPairs := strings.Split(fieldsStr, ",")

	for _, pair := range fieldPairs {
		parts := strings.Split(strings.TrimSpace(pair), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field format: %s (expected name:type)", pair)
		}

		name := strings.TrimSpace(parts[0])
		fieldType := strings.TrimSpace(parts[1])

		// 处理可选的required标记
		required := false
		if strings.HasSuffix(fieldType, "!") {
			required = true
			fieldType = strings.TrimSuffix(fieldType, "!")
		}

		field := &Field{
			Name:     ToPascalCase(name),
			Type:     mapGoType(fieldType),
			Required: required,
		}

		// 生成标签
		field.JsonTag = fmt.Sprintf(`json:"%s"`, ToSnakeCase(name))
		field.GormTag = generateGormTag(field)
		field.Tag = fmt.Sprintf(`%s %s`, field.JsonTag, field.GormTag)
		field.Comment = fmt.Sprintf("// %s %s", ToPascalCase(name), getTypeComment(fieldType))

		fields = append(fields, field)
	}

	return fields, nil
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

// generateGormTag 生成GORM标签
func generateGormTag(field *Field) string {
	var tags []string

	// 列名
	tags = append(tags, fmt.Sprintf("column:%s", ToSnakeCase(field.Name)))

	// 类型映射
	switch field.Type {
	case "string":
		if field.Name == "email" {
			tags = append(tags, "type:varchar(255)", "uniqueIndex")
		} else if strings.Contains(strings.ToLower(field.Name), "url") {
			tags = append(tags, "type:varchar(500)")
		} else if strings.Contains(strings.ToLower(field.Name), "description") || strings.Contains(strings.ToLower(field.Name), "content") {
			tags = append(tags, "type:text")
		} else {
			tags = append(tags, "type:varchar(255)")
		}
	case "time.Time":
		tags = append(tags, "type:datetime")
	case "bool":
		tags = append(tags, "type:boolean", "default:false")
	case "float64":
		tags = append(tags, "type:decimal(10,2)")
	}

	// 必填字段
	if field.Required {
		tags = append(tags, "not null")
	}

	return fmt.Sprintf(`gorm:"%s"`, strings.Join(tags, ";"))
}

// getTypeComment 获取类型注释
func getTypeComment(fieldType string) string {
	commentMap := map[string]string{
		"string":    "字符串",
		"int":       "整数",
		"int32":     "32位整数",
		"int64":     "64位整数",
		"uint":      "无符号整数",
		"uint32":    "32位无符号整数",
		"uint64":    "64位无符号整数",
		"float32":   "32位浮点数",
		"float64":   "64位浮点数",
		"bool":      "布尔值",
		"time":      "时间",
		"datetime":  "日期时间",
		"timestamp": "时间戳",
		"text":      "文本",
		"json":      "JSON数据",
		"decimal":   "小数",
	}

	if comment, exists := commentMap[fieldType]; exists {
		return comment
	}
	return "自定义类型"
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
