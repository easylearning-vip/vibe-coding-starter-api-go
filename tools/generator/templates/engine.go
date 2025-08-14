package templates

import (
	"embed"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templateFS embed.FS

// Engine 模板引擎
type Engine struct {
	templates map[string]*template.Template
}

// NewEngine 创建模板引擎
func NewEngine() *Engine {
	engine := &Engine{
		templates: make(map[string]*template.Template),
	}

	// 加载所有模板
	engine.loadTemplates()

	return engine
}

// loadTemplates 加载所有模板文件
func (e *Engine) loadTemplates() {
	// 获取所有模板文件
	entries, err := templateFS.ReadDir(".")
	if err != nil {
		panic(fmt.Sprintf("failed to read template directory: %v", err))
	}

	// 创建函数映射
	funcMap := template.FuncMap{
		"toPascalCase":      toPascalCase,
		"toCamelCase":       toCamelCase,
		"toSnakeCase":       toSnakeCase,
		"toKebabCase":       toKebabCase,
		"pluralize":         pluralize,
		"lower":             strings.ToLower,
		"upper":             strings.ToUpper,
		"title":             strings.Title,
		"join":              strings.Join,
		"contains":          strings.Contains,
		"hasPrefix":         strings.HasPrefix,
		"hasSuffix":         strings.HasSuffix,
		"trimPrefix":        strings.TrimPrefix,
		"trimSuffix":        strings.TrimSuffix,
		"replace":           strings.ReplaceAll,
		"GenerateSQLColumn": generateSQLColumn,
	}

	// 加载每个模板文件
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			content, err := templateFS.ReadFile(entry.Name())
			if err != nil {
				panic(fmt.Sprintf("failed to read template %s: %v", entry.Name(), err))
			}

			tmpl, err := template.New(entry.Name()).Funcs(funcMap).Parse(string(content))
			if err != nil {
				panic(fmt.Sprintf("failed to parse template %s: %v", entry.Name(), err))
			}

			e.templates[entry.Name()] = tmpl
		}
	}
}

// Render 渲染模板
func (e *Engine) Render(templateName string, data interface{}) (string, error) {
	tmpl, exists := e.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// ListTemplates 列出所有可用的模板
func (e *Engine) ListTemplates() []string {
	var names []string
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}

// HasTemplate 检查模板是否存在
func (e *Engine) HasTemplate(name string) bool {
	_, exists := e.templates[name]
	return exists
}

// 模板函数

// toPascalCase 转换为PascalCase
func toPascalCase(s string) string {
	// 如果字符串已经是PascalCase格式（首字母大写，没有分隔符），直接返回
	if isPascalCaseTemplate(s) {
		return s
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}

	return strings.Join(words, "")
}

// isPascalCaseTemplate 检查字符串是否已经是PascalCase格式
func isPascalCaseTemplate(s string) bool {
	if len(s) == 0 {
		return false
	}

	// 首字母必须是大写
	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}

	// 不能包含分隔符
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			return false
		}
	}

	return true
}

// toCamelCase 转换为camelCase
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// toSnakeCase 转换为snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}
	return string(result)
}

// toKebabCase 转换为kebab-case
func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

// pluralize 复数形式
func pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}
	return s + "s"
}

// Field 表示模型字段（简化版本，避免循环导入）
type Field struct {
	Name string
	Type string
}

// generateSQLColumn 生成SQL列定义
func generateSQLColumn(field interface{}, databaseType string) string {
	// 类型断言，支持不同的字段结构
	var name, fieldType string

	// Use reflection to extract Name and Type fields from any struct
	fieldValue := reflect.ValueOf(field)
	if fieldValue.Kind() == reflect.Ptr {
		fieldValue = fieldValue.Elem()
	}

	if fieldValue.Kind() == reflect.Struct {
		nameField := fieldValue.FieldByName("Name")
		typeField := fieldValue.FieldByName("Type")

		if nameField.IsValid() && typeField.IsValid() {
			name = nameField.String()
			fieldType = typeField.String()
		} else {
			return "    -- Invalid field structure"
		}
	} else {
		switch f := field.(type) {
		case map[string]interface{}:
			name = fmt.Sprintf("%v", f["Name"])
			fieldType = fmt.Sprintf("%v", f["Type"])
		default:
			return "    -- Invalid field type"
		}
	}

	columnName := toSnakeCase(name)

	var sqlType string
	var defaultValue string

	switch fieldType {
	case "string":
		sqlType = "VARCHAR(255)"
		defaultValue = ""
	case "int", "int32":
		if databaseType == "postgres" {
			sqlType = "INTEGER"
		} else {
			sqlType = "INT"
		}
		defaultValue = ""
	case "int64":
		if databaseType == "postgres" {
			sqlType = "BIGINT"
		} else {
			sqlType = "BIGINT"
		}
		defaultValue = ""
	case "uint", "uint32":
		if databaseType == "postgres" {
			sqlType = "INTEGER"
		} else {
			sqlType = "INT UNSIGNED"
		}
		defaultValue = ""
	case "uint64":
		if databaseType == "postgres" {
			sqlType = "BIGINT"
		} else {
			sqlType = "BIGINT UNSIGNED"
		}
		defaultValue = ""
	case "float32":
		sqlType = "FLOAT"
		defaultValue = ""
	case "float64":
		sqlType = "DECIMAL(10,2)"
		defaultValue = ""
	case "bool":
		if databaseType == "postgres" {
			sqlType = "BOOLEAN"
		} else {
			sqlType = "BOOLEAN"
		}
		defaultValue = " DEFAULT FALSE"
	case "time.Time":
		if databaseType == "postgres" {
			sqlType = "TIMESTAMP WITH TIME ZONE"
		} else {
			sqlType = "DATETIME"
		}
		defaultValue = ""
	default:
		sqlType = "TEXT"
		defaultValue = ""
	}

	// 构建完整的列定义
	columnDef := fmt.Sprintf("    %s %s", columnName, sqlType)

	// 添加NOT NULL约束（除了某些特殊情况）
	if fieldType != "bool" && !strings.Contains(strings.ToLower(name), "optional") {
		columnDef += " NOT NULL"
	}

	// 添加默认值
	if defaultValue != "" {
		columnDef += defaultValue
	}

	// 注意：Name字段的唯一约束现在在Indexes部分显式声明，不在列定义中

	return columnDef
}
