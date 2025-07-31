package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// HandlerGenerator 处理器生成器
type HandlerGenerator struct {
	templateEngine *templates.Engine
}

// NewHandlerGenerator 创建处理器生成器
func NewHandlerGenerator() *HandlerGenerator {
	return &HandlerGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成API处理器
func (g *HandlerGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*HandlerConfig)
	if !ok {
		return fmt.Errorf("invalid config type for handler generator")
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Model":          ToPascalCase(cfg.Model),
		"ModelLower":     strings.ToLower(cfg.Model),
		"ModelCamel":     ToCamelCase(cfg.Model),
		"ModelSnake":     ToSnakeCase(cfg.Model),
		"ModelPlural":    Pluralize(strings.ToLower(cfg.Model)),
		"WithAuth":       cfg.WithAuth,
		"WithValidation": cfg.WithValidation,
		"Year":           GetCurrentYear(),
	}

	// 生成处理器文件
	content, err := g.templateEngine.Render("handler.go.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render handler template: %w", err)
	}

	filename := fmt.Sprintf("%s.go", data["ModelSnake"])
	path := filepath.Join("internal", "handler", filename)

	if err := g.writeFile(path, content); err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	// 请求结构体现在包含在主 handler 文件中，不需要单独生成

	// 生成处理器测试
	if err := g.generateHandlerTest(data); err != nil {
		return fmt.Errorf("failed to generate handler test: %w", err)
	}

	return nil
}

// generateRequestStructs 生成请求结构体
func (g *HandlerGenerator) generateRequestStructs(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("handler_requests.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_requests.go", data["ModelSnake"])
	path := filepath.Join("internal", "handler", filename)

	return g.writeFile(path, content)
}

// generateHandlerTest 生成处理器测试
func (g *HandlerGenerator) generateHandlerTest(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("handler_test.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_handler_test.go", data["ModelSnake"])
	path := filepath.Join("test", "handler", filename)

	return g.writeFile(path, content)
}

// writeFile 写入文件
func (g *HandlerGenerator) writeFile(path, content string) error {
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
