package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// ModelGenerator 模型生成器
type ModelGenerator struct {
	templateEngine *templates.Engine
}

// NewModelGenerator 创建模型生成器
func NewModelGenerator() *ModelGenerator {
	return &ModelGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成数据模型
func (g *ModelGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*ModelConfig)
	if !ok {
		return fmt.Errorf("invalid config type for model generator")
	}

	// 解析字段
	parser := NewFieldParser()
	fields, err := parser.ParseFields(cfg.Fields)
	if err != nil {
		return fmt.Errorf("failed to parse fields: %w", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Name":            ToPascalCase(cfg.Name),
		"NameLower":       strings.ToLower(cfg.Name),
		"NameCamel":       ToCamelCase(cfg.Name),
		"NameSnake":       ToSnakeCase(cfg.Name),
		"NamePlural":      Pluralize(strings.ToLower(cfg.Name)),
		"Fields":          fields,
		"RequiredImports": parser.GetRequiredImports(fields),
		"WithTimestamps":  cfg.WithTimestamps,
		"Year":            GetCurrentYear(),
	}

	// 生成模型文件
	content, err := g.templateEngine.Render("model.go.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render model template: %w", err)
	}

	filename := fmt.Sprintf("%s.go", data["NameSnake"])
	path := filepath.Join("internal", "model", filename)

	if err := g.writeFile(path, content); err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	return nil
}

// writeFile 写入文件
func (g *ModelGenerator) writeFile(path, content string) error {
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
