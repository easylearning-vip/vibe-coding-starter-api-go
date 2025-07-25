package templates

import (
	"embed"
	"fmt"
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
		"toPascalCase": toPascalCase,
		"toCamelCase":  toCamelCase,
		"toSnakeCase":  toSnakeCase,
		"toKebabCase":  toKebabCase,
		"pluralize":    pluralize,
		"lower":        strings.ToLower,
		"upper":        strings.ToUpper,
		"title":        strings.Title,
		"join":         strings.Join,
		"contains":     strings.Contains,
		"hasPrefix":    strings.HasPrefix,
		"hasSuffix":    strings.HasSuffix,
		"trimPrefix":   strings.TrimPrefix,
		"trimSuffix":   strings.TrimSuffix,
		"replace":      strings.ReplaceAll,
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
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}

	return strings.Join(words, "")
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
