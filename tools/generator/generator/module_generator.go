package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// ModuleGenerator 模块生成器
type ModuleGenerator struct {
	templateEngine *templates.Engine
}

// NewModuleGenerator 创建模块生成器
func NewModuleGenerator() *ModuleGenerator {
	return &ModuleGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成完整的业务模块
func (g *ModuleGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*ModuleConfig)
	if !ok {
		return fmt.Errorf("invalid config type for module generator")
	}

	// 解析字段
	parser := NewFieldParser()
	fields, err := parser.ParseFields(cfg.Fields)
	if err != nil {
		return fmt.Errorf("failed to parse fields: %w", err)
	}

	// 准备模板数据
	data := g.prepareTemplateData(cfg, fields)

	// 生成各个组件
	if err := g.generateModel(data); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	if err := g.generateRepository(data); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	if err := g.generateService(data); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	if err := g.generateHandler(data); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	if err := g.generateTests(data); err != nil {
		return fmt.Errorf("failed to generate tests: %w", err)
	}

	if err := g.generateMigration(data); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	// 更新路由注册
	if err := g.updateRoutes(data); err != nil {
		return fmt.Errorf("failed to update routes: %w", err)
	}

	return nil
}

// prepareTemplateData 准备模板数据
func (g *ModuleGenerator) prepareTemplateData(cfg *ModuleConfig, fields []*Field) map[string]interface{} {
	name := ToPascalCase(cfg.Name)
	return map[string]interface{}{
		// 原有的 Name 系列变量
		"Name":            name,
		"NameLower":       strings.ToLower(cfg.Name),
		"NameCamel":       ToCamelCase(cfg.Name),
		"NameSnake":       ToSnakeCase(cfg.Name),
		"NameKebab":       ToKebabCase(cfg.Name),
		"NamePlural":      Pluralize(strings.ToLower(cfg.Name)),
		"NamePluralCamel": ToCamelCase(Pluralize(cfg.Name)),

		// 新增的 Model 系列变量（与 Name 系列保持一致）
		"Model":            name,
		"ModelLower":       strings.ToLower(cfg.Name),
		"ModelCamel":       ToCamelCase(cfg.Name),
		"ModelSnake":       ToSnakeCase(cfg.Name),
		"ModelKebab":       ToKebabCase(cfg.Name),
		"ModelPlural":      Pluralize(strings.ToLower(cfg.Name)),
		"ModelPluralCamel": ToCamelCase(Pluralize(cfg.Name)),

		// 表名
		"TableName": Pluralize(ToSnakeCase(cfg.Name)),

		"Fields":    fields,
		"WithAuth":  cfg.WithAuth,
		"WithCache": cfg.WithCache,
		"Timestamp": GenerateTimestamp(),
		"Year":      GetCurrentYear(),
	}
}

// generateModel 生成模型
func (g *ModuleGenerator) generateModel(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("model.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.go", data["NameSnake"])
	path := filepath.Join("internal", "model", filename)

	return g.writeFile(path, content)
}

// generateRepository 生成仓储
func (g *ModuleGenerator) generateRepository(data map[string]interface{}) error {
	// 生成接口
	interfaceContent, err := g.templateEngine.Render("repository_interface.go.tmpl", data)
	if err != nil {
		return err
	}

	// 更新interfaces.go文件
	interfacePath := filepath.Join("internal", "repository", "interfaces.go")
	if err := g.appendToFile(interfacePath, interfaceContent); err != nil {
		return err
	}

	// 生成实现
	implContent, err := g.templateEngine.Render("repository_impl.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_repository.go", data["NameSnake"])
	implPath := filepath.Join("internal", "repository", filename)

	return g.writeFile(implPath, implContent)
}

// generateService 生成服务
func (g *ModuleGenerator) generateService(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("service.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_service.go", data["NameSnake"])
	path := filepath.Join("internal", "service", filename)

	return g.writeFile(path, content)
}

// generateHandler 生成处理器
func (g *ModuleGenerator) generateHandler(data map[string]interface{}) error {
	// 使用专门的 HandlerGenerator
	handlerGen := NewHandlerGenerator()

	// 准备 HandlerConfig
	cfg := &HandlerConfig{
		Model:          data["Model"].(string),
		WithAuth:       data["WithAuth"].(bool),
		WithValidation: false, // 默认不启用验证
	}

	return handlerGen.Generate(cfg)
}

// generateTests 生成测试
func (g *ModuleGenerator) generateTests(data map[string]interface{}) error {
	// 生成Repository测试
	repoTestContent, err := g.templateEngine.Render("repository_test.go.tmpl", data)
	if err != nil {
		return err
	}

	repoTestFile := fmt.Sprintf("%s_repository_test.go", data["NameSnake"])
	repoTestPath := filepath.Join("test", "repository", repoTestFile)

	if err := g.writeFile(repoTestPath, repoTestContent); err != nil {
		return err
	}

	// 生成Service测试
	serviceTestContent, err := g.templateEngine.Render("service_test.go.tmpl", data)
	if err != nil {
		return err
	}

	serviceTestFile := fmt.Sprintf("%s_service_test.go", data["NameSnake"])
	serviceTestPath := filepath.Join("test", "service", serviceTestFile)

	if err := g.writeFile(serviceTestPath, serviceTestContent); err != nil {
		return err
	}

	// 生成Handler测试
	handlerTestContent, err := g.templateEngine.Render("handler_test.go.tmpl", data)
	if err != nil {
		return err
	}

	handlerTestFile := fmt.Sprintf("%s_handler_test.go", data["NameSnake"])
	handlerTestPath := filepath.Join("test", "handler", handlerTestFile)

	return g.writeFile(handlerTestPath, handlerTestContent)
}

// generateMigration 生成数据库迁移
func (g *ModuleGenerator) generateMigration(data map[string]interface{}) error {
	// 使用专门的 MigrationGenerator
	migrationGen := NewMigrationGenerator()

	// 准备 MigrationConfig
	cfg := &MigrationConfig{
		Name:   fmt.Sprintf("create_%s_table", data["NamePlural"]),
		Table:  data["NamePlural"].(string),
		Action: "create",
	}

	return migrationGen.Generate(cfg)
}

// updateRoutes 更新路由注册
func (g *ModuleGenerator) updateRoutes(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("routes.go.tmpl", data)
	if err != nil {
		return err
	}

	// 生成路由注册代码到单独的文件
	filename := fmt.Sprintf("%s_routes.go", data["NameSnake"])
	routesPath := filepath.Join("docs", "generated", filename)

	return g.writeFile(routesPath, content)
}

// writeFile 写入文件
func (g *ModuleGenerator) writeFile(path, content string) error {
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
func (g *ModuleGenerator) appendToFile(path, content string) error {
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
