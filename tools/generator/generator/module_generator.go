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

	// 自动更新 server.go
	if err := g.UpdateServerFile(data); err != nil {
		return fmt.Errorf("failed to update server.go: %w", err)
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

// UpdateServerFile 自动更新 server.go 文件
func (g *ModuleGenerator) UpdateServerFile(data map[string]interface{}) error {
	serverPath := "internal/server/server.go"

	// 读取现有内容
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	serverContent := string(content)
	modelName := data["Model"].(string)
	modelCamel := ToCamelCase(strings.ToLower(modelName))

	// 检查是否已经存在
	if strings.Contains(serverContent, fmt.Sprintf("%sHandler", modelCamel)) {
		fmt.Printf("⚠️  %s handler already exists in server.go\n", modelName)
		return nil
	}

	// 1. 添加字段到 Server 结构体
	structPattern := "dictHandler    *handler.DictHandler"
	structIndex := strings.Index(serverContent, structPattern)
	if structIndex == -1 {
		return fmt.Errorf("could not find Server struct fields section")
	}

	structEnd := structIndex + len(structPattern)
	newField := fmt.Sprintf("\n\t%sHandler *handler.%sHandler", modelCamel, modelName)
	serverContent = serverContent[:structEnd] + newField + serverContent[structEnd:]

	// 2. 添加参数到构造函数
	paramPattern := "dictHandler *handler.DictHandler,"
	paramIndex := strings.Index(serverContent, paramPattern)
	if paramIndex == -1 {
		return fmt.Errorf("could not find constructor parameters section")
	}

	paramEnd := paramIndex + len(paramPattern)
	newParam := fmt.Sprintf("\n\t%sHandler *handler.%sHandler,", modelCamel, modelName)
	serverContent = serverContent[:paramEnd] + newParam + serverContent[paramEnd:]

	// 3. 添加字段初始化
	initPattern := "dictHandler:    dictHandler,"
	initIndex := strings.Index(serverContent, initPattern)
	if initIndex == -1 {
		return fmt.Errorf("could not find constructor initialization section")
	}

	initEnd := initIndex + len(initPattern)
	newInit := fmt.Sprintf("\n\t\t%sHandler: %sHandler,", modelCamel, modelCamel)
	serverContent = serverContent[:initEnd] + newInit + serverContent[initEnd:]

	// 4. 添加路由注册
	// 寻找admin路由组的结束位置，在管理员文章路由之后添加
	adminArticlesPattern := "adminArticles.DELETE(\"/:id\", s.articleHandler.Delete)\n\t\t\t\t}"
	adminArticlesIndex := strings.Index(serverContent, adminArticlesPattern)
	if adminArticlesIndex != -1 {
		// 在管理员文章路由组之后添加新的路由
		insertPos := adminArticlesIndex + len(adminArticlesPattern)
		newRoute := fmt.Sprintf("\n\n\t\t\t\t// %s管理路由\n\t\t\t\ts.%sHandler.RegisterRoutes(admin)", modelName, modelCamel)
		serverContent = serverContent[:insertPos] + newRoute + serverContent[insertPos:]
	} else {
		// 如果找不到管理员文章路由，尝试在admin路由组的末尾添加
		// 寻找admin路由组的结束位置
		adminPattern := "// 管理员专用路由\n\t\t\t\ts.userHandler.RegisterRoutes(admin)"
		adminIndex := strings.Index(serverContent, adminPattern)
		if adminIndex != -1 {
			insertPos := adminIndex + len(adminPattern)
			newRoute := fmt.Sprintf("\n\n\t\t\t\t// %s管理路由\n\t\t\t\ts.%sHandler.RegisterRoutes(admin)", modelName, modelCamel)
			serverContent = serverContent[:insertPos] + newRoute + serverContent[insertPos:]
		} else {
			return fmt.Errorf("could not find admin routes section")
		}
	}

	// 写回文件
	if err := os.WriteFile(serverPath, []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to write server.go: %w", err)
	}

	return nil
}

// UpdateMainFile 自动更新 main.go 文件
func (g *ModuleGenerator) UpdateMainFile(data map[string]interface{}) error {
	modelName := data["Model"].(string)

	// 读取 main.go 文件
	mainPath := "cmd/server/main.go"
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	content := string(mainContent)

	// 检查是否已经存在
	if strings.Contains(content, fmt.Sprintf("repository.New%sRepository", modelName)) {
		return fmt.Errorf("%s repository already exists in main.go", modelName)
	}

	// 1. 添加仓储提供者
	repoPattern := "repository.NewDictRepository,"
	repoIndex := strings.Index(content, repoPattern)
	if repoIndex == -1 {
		return fmt.Errorf("could not find repository providers section in main.go")
	}

	repoEnd := repoIndex + len(repoPattern)
	newRepo := fmt.Sprintf("\n\t\t\trepository.New%sRepository,", modelName)
	content = content[:repoEnd] + newRepo + content[repoEnd:]

	// 2. 添加服务提供者
	servicePattern := "service.NewDictService,"
	serviceIndex := strings.Index(content, servicePattern)
	if serviceIndex == -1 {
		return fmt.Errorf("could not find service providers section in main.go")
	}

	serviceEnd := serviceIndex + len(servicePattern)
	newService := fmt.Sprintf("\n\t\t\tservice.New%sService,", modelName)
	content = content[:serviceEnd] + newService + content[serviceEnd:]

	// 3. 添加处理器提供者
	handlerPattern := "handler.NewDictHandler,"
	handlerIndex := strings.Index(content, handlerPattern)
	if handlerIndex == -1 {
		return fmt.Errorf("could not find handler providers section in main.go")
	}

	handlerEnd := handlerIndex + len(handlerPattern)
	newHandler := fmt.Sprintf("\n\t\t\thandler.New%sHandler,", modelName)
	content = content[:handlerEnd] + newHandler + content[handlerEnd:]

	// 写回文件
	if err := os.WriteFile(mainPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}
