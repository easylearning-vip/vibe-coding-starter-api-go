package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// ServiceGenerator 服务生成器
type ServiceGenerator struct {
	templateEngine *templates.Engine
}

// NewServiceGenerator 创建服务生成器
func NewServiceGenerator() *ServiceGenerator {
	return &ServiceGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成服务层
func (g *ServiceGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*ServiceConfig)
	if !ok {
		return fmt.Errorf("invalid config type for service generator")
	}

	// 如果没有指定模型，使用服务名称
	model := cfg.Model
	if model == "" {
		model = strings.TrimSuffix(cfg.Name, "Service")
	}

	// 获取字段信息（优先使用模型反射）
	reflector := NewModelReflector()
	fields, err := reflector.GetFieldsFromModelOrFallback(cfg.Model, cfg.Fields)
	if err != nil {
		return fmt.Errorf("failed to get model fields: %w", err)
	}

	parser := NewFieldParser()

	// 准备模板数据
	data := map[string]interface{}{
		"Name":                ToPascalCase(cfg.Name),
		"NameCamel":           ToCamelCase(cfg.Name),
		"NameSnake":           ToSnakeCase(cfg.Name),
		"Model":               ToPascalCase(model),
		"ModelLower":          strings.ToLower(model),
		"ModelCamel":          ToCamelCase(model),
		"ModelSnake":          ToSnakeCase(model),
		"ModelPlural":         Pluralize(strings.ToLower(model)),
		"WithCache":           cfg.WithCache,
		"Year":                GetCurrentYear(),
		"Fields":              fields,
		"CreateRequestFields": parser.GenerateCreateRequestFields(fields),
		"UpdateRequestFields": parser.GenerateUpdateRequestFields(fields),
		"ModelAssignment":     parser.GenerateModelAssignment(fields, "req"),
		"UpdateAssignment":    parser.GenerateUpdateAssignment(fields, "entity", "req"),
		"RequiredImports":     parser.GetRequiredImports(fields),
	}

	// 生成服务文件
	content, err := g.templateEngine.Render("service.go.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render service template: %w", err)
	}

	filename := fmt.Sprintf("%s.go", data["ModelSnake"])
	path := filepath.Join("internal", "service", filename)

	if err := g.writeFile(path, content); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// 生成服务测试
	if err := g.generateServiceTest(data); err != nil {
		return fmt.Errorf("failed to generate service test: %w", err)
	}

	// 生成服务 Mock
	if err := g.generateServiceMock(data); err != nil {
		return fmt.Errorf("failed to generate service mock: %w", err)
	}

	return nil
}

// generateServiceTest 生成服务测试
func (g *ServiceGenerator) generateServiceTest(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("service_test.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_test.go", data["NameSnake"])
	path := filepath.Join("test", "service", filename)

	return g.writeFile(path, content)
}

// writeFile 写入文件
func (g *ServiceGenerator) writeFile(path, content string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 检查文件是否已存在
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	// 写入文件
	return os.WriteFile(path, []byte(content), 0644)
}

// generateServiceMock 生成服务 Mock
func (g *ServiceGenerator) generateServiceMock(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("service_mock.go.tmpl", data)
	if err != nil {
		return err
	}

	// 追加到 service_mocks.go 文件
	mockPath := filepath.Join("test", "mocks", "service_mocks.go")
	return g.appendToFile(mockPath, content)
}

// appendToFile 追加到文件
func (g *ServiceGenerator) appendToFile(path, content string) error {
	// 如果文件不存在，创建它
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return g.writeFile(path, content)
	}

	// 读取现有内容
	existing, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 追加新内容
	newContent := string(existing) + "\n" + content

	return os.WriteFile(path, []byte(newContent), 0644)
}
