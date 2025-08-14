package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	modelName := ToPascalCase(cfg.Name)
	tableName := ToSnakeCase(Pluralize(modelName))

	data := map[string]interface{}{
		"Name":            modelName,
		"NameLower":       strings.ToLower(cfg.Name),
		"NameCamel":       ToCamelCase(cfg.Name),
		"NameSnake":       ToSnakeCase(cfg.Name),
		"NamePlural":      tableName, // 使用snake_case的复数形式作为表名
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

	// 更新测试数据库配置
	if err := g.updateTestDatabase(modelName, tableName); err != nil {
		return fmt.Errorf("failed to update test database: %w", err)
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

// updateTestDatabase 更新测试数据库配置，添加新模型到迁移和清理列表
func (g *ModelGenerator) updateTestDatabase(modelName, tableName string) error {
	databasePath := filepath.Join("test", "testutil", "database.go")

	// 读取现有文件
	content, err := os.ReadFile(databasePath)
	if err != nil {
		return fmt.Errorf("failed to read database.go: %w", err)
	}

	contentStr := string(content)

	// 检查模型是否已经存在于AutoMigrate中
	modelImport := fmt.Sprintf("&model.%s{}", modelName)
	if strings.Contains(contentStr, modelImport) {
		// 模型已存在，不需要更新
		return nil
	}

	// 添加模型到AutoMigrate
	migratePattern := `(&model\.DictItem{},\s*\n)(\s*&model\.\w+{},\s*\n)*(\s*\))`
	migrateReplacement := fmt.Sprintf("$1\t\t&model.%s{},\n$3", modelName)

	re := regexp.MustCompile(migratePattern)
	if re.MatchString(contentStr) {
		contentStr = re.ReplaceAllString(contentStr, migrateReplacement)
	} else {
		return fmt.Errorf("could not find AutoMigrate section to update")
	}

	// 添加表到Clean方法的tables列表
	cleanPattern := `("dict_categories",\s*\n)(\s*"[^"]+",\s*\n)*(\s*\})`
	cleanReplacement := fmt.Sprintf("$1\t\t\"%s\",\n$3", tableName)

	re = regexp.MustCompile(cleanPattern)
	if re.MatchString(contentStr) {
		contentStr = re.ReplaceAllString(contentStr, cleanReplacement)
	} else {
		return fmt.Errorf("could not find Clean tables section to update")
	}

	// 写回文件
	return os.WriteFile(databasePath, []byte(contentStr), 0644)
}
