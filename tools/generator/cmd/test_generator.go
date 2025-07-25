package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// TestGenerator 测试生成器
type TestGenerator struct {
	templateEngine *templates.Engine
}

// NewTestGenerator 创建测试生成器
func NewTestGenerator() *TestGenerator {
	return &TestGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成测试代码
func (g *TestGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*TestConfig)
	if !ok {
		return fmt.Errorf("invalid config type for test generator")
	}

	// 根据不同的测试目标生成测试
	if cfg.Service != "" {
		return g.generateServiceTest(cfg)
	}

	if cfg.Handler != "" {
		return g.generateHandlerTest(cfg)
	}

	if cfg.Repository != "" {
		return g.generateRepositoryTest(cfg)
	}

	return fmt.Errorf("no test target specified")
}

// generateServiceTest 生成服务测试
func (g *TestGenerator) generateServiceTest(cfg *TestConfig) error {
	serviceName := strings.TrimSuffix(cfg.Service, "Service")
	
	data := map[string]interface{}{
		"Service":     ToPascalCase(cfg.Service),
		"ServiceName": ToPascalCase(serviceName),
		"NameCamel":   ToCamelCase(serviceName),
		"NameSnake":   ToSnakeCase(serviceName),
		"Type":        cfg.Type,
		"Year":        GetCurrentYear(),
	}

	var templateName string
	var outputDir string

	switch cfg.Type {
	case "unit":
		templateName = "service_unit_test.go.tmpl"
		outputDir = "test/service"
	case "integration":
		templateName = "service_integration_test.go.tmpl"
		outputDir = "test/integration"
	default:
		templateName = "service_test.go.tmpl"
		outputDir = "test/service"
	}

	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render service test template: %w", err)
	}

	filename := fmt.Sprintf("%s_test.go", data["NameSnake"])
	path := filepath.Join(outputDir, filename)

	return g.writeFile(path, content)
}

// generateHandlerTest 生成处理器测试
func (g *TestGenerator) generateHandlerTest(cfg *TestConfig) error {
	handlerName := strings.TrimSuffix(cfg.Handler, "Handler")
	
	data := map[string]interface{}{
		"Handler":     ToPascalCase(cfg.Handler),
		"HandlerName": ToPascalCase(handlerName),
		"NameCamel":   ToCamelCase(handlerName),
		"NameSnake":   ToSnakeCase(handlerName),
		"Type":        cfg.Type,
		"Year":        GetCurrentYear(),
	}

	var templateName string
	var outputDir string

	switch cfg.Type {
	case "unit":
		templateName = "handler_unit_test.go.tmpl"
		outputDir = "test/handler"
	case "integration":
		templateName = "handler_integration_test.go.tmpl"
		outputDir = "test/integration"
	case "e2e":
		templateName = "handler_e2e_test.go.tmpl"
		outputDir = "test/e2e"
	default:
		templateName = "handler_test.go.tmpl"
		outputDir = "test/handler"
	}

	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render handler test template: %w", err)
	}

	filename := fmt.Sprintf("%s_test.go", data["NameSnake"])
	path := filepath.Join(outputDir, filename)

	return g.writeFile(path, content)
}

// generateRepositoryTest 生成仓储测试
func (g *TestGenerator) generateRepositoryTest(cfg *TestConfig) error {
	repoName := strings.TrimSuffix(cfg.Repository, "Repository")
	
	data := map[string]interface{}{
		"Repository":     ToPascalCase(cfg.Repository),
		"RepositoryName": ToPascalCase(repoName),
		"NameCamel":      ToCamelCase(repoName),
		"NameSnake":      ToSnakeCase(repoName),
		"Type":           cfg.Type,
		"Year":           GetCurrentYear(),
	}

	var templateName string
	var outputDir string

	switch cfg.Type {
	case "unit":
		templateName = "repository_unit_test.go.tmpl"
		outputDir = "test/repository"
	case "integration":
		templateName = "repository_integration_test.go.tmpl"
		outputDir = "test/integration"
	default:
		templateName = "repository_test.go.tmpl"
		outputDir = "test/repository"
	}

	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render repository test template: %w", err)
	}

	filename := fmt.Sprintf("%s_test.go", data["NameSnake"])
	path := filepath.Join(outputDir, filename)

	return g.writeFile(path, content)
}

// writeFile 写入文件
func (g *TestGenerator) writeFile(path, content string) error {
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
