package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// RepositoryGenerator 仓储生成器
type RepositoryGenerator struct {
	templateEngine *templates.Engine
}

// NewRepositoryGenerator 创建仓储生成器
func NewRepositoryGenerator() *RepositoryGenerator {
	return &RepositoryGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成仓储层
func (g *RepositoryGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*RepositoryConfig)
	if !ok {
		return fmt.Errorf("invalid config type for repository generator")
	}

	// 如果没有指定模型，使用仓储名称
	model := cfg.Model
	if model == "" {
		model = strings.TrimSuffix(cfg.Name, "Repository")
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Name":        ToPascalCase(cfg.Name),
		"NameCamel":   ToCamelCase(cfg.Name),
		"NameSnake":   ToSnakeCase(cfg.Name),
		"Model":       ToPascalCase(model),
		"ModelLower":  strings.ToLower(model),
		"ModelCamel":  ToCamelCase(model),
		"ModelSnake":  ToSnakeCase(model),
		"ModelPlural": Pluralize(strings.ToLower(model)),
		"Year":        GetCurrentYear(),
	}

	// 生成仓储接口
	if err := g.generateRepositoryInterface(data); err != nil {
		return fmt.Errorf("failed to generate repository interface: %w", err)
	}

	// 生成仓储实现
	if err := g.generateRepositoryImpl(data); err != nil {
		return fmt.Errorf("failed to generate repository implementation: %w", err)
	}

	// 生成仓储测试
	if err := g.generateRepositoryTest(data); err != nil {
		return fmt.Errorf("failed to generate repository test: %w", err)
	}

	// 生成仓储Mock
	if err := g.generateRepositoryMock(data); err != nil {
		return fmt.Errorf("failed to generate repository mock: %w", err)
	}

	return nil
}

// generateRepositoryInterface 生成仓储接口
func (g *RepositoryGenerator) generateRepositoryInterface(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("repository_interface.go.tmpl", data)
	if err != nil {
		return err
	}

	// 更新interfaces.go文件
	interfacePath := filepath.Join("internal", "repository", "interfaces.go")
	return g.appendToFile(interfacePath, content)
}

// generateRepositoryImpl 生成仓储实现
func (g *RepositoryGenerator) generateRepositoryImpl(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("repository_impl.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.go", data["ModelSnake"])
	path := filepath.Join("internal", "repository", filename)

	return g.writeFile(path, content)
}

// generateRepositoryTest 生成仓储测试
func (g *RepositoryGenerator) generateRepositoryTest(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("repository_test.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_test.go", data["NameSnake"])
	path := filepath.Join("test", "repository", filename)

	return g.writeFile(path, content)
}

// writeFile 写入文件
func (g *RepositoryGenerator) writeFile(path, content string) error {
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

// appendToFile 追加到文件
func (g *RepositoryGenerator) appendToFile(path, content string) error {
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

// generateRepositoryMock 生成仓储Mock
func (g *RepositoryGenerator) generateRepositoryMock(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("repository_mock.go.tmpl", data)
	if err != nil {
		return err
	}

	// 更新repository_mocks.go文件
	mockPath := filepath.Join("test", "mocks", "repository_mocks.go")
	return g.appendToFile(mockPath, content)
}
