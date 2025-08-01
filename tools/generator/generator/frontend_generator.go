package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// FrontendGenerator 前端代码生成器
type FrontendGenerator struct {
	templateEngine *templates.Engine
}

// NewFrontendGenerator 创建前端代码生成器
func NewFrontendGenerator() *FrontendGenerator {
	return &FrontendGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成前端代码
func (g *FrontendGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*FrontendConfig)
	if !ok {
		return fmt.Errorf("invalid config type for frontend generator")
	}

	// 验证输出目录
	if err := g.validateOutputDir(cfg.OutputDir); err != nil {
		return fmt.Errorf("output directory validation failed: %w", err)
	}

	// 获取模型信息
	modelInfo, err := g.getModelInfo(cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to get model info: %w", err)
	}

	// 准备模板数据
	data := g.prepareTemplateData(cfg, modelInfo)

	// 根据框架类型生成代码
	switch cfg.Framework {
	case FrameworkAntd:
		return g.generateAntdCode(data, cfg)
	case FrameworkVue:
		return g.generateVueCode(data, cfg)
	default:
		return fmt.Errorf("unsupported framework: %s", cfg.Framework)
	}
}

// validateOutputDir 验证输出目录
func (g *FrontendGenerator) validateOutputDir(outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	// 检查目录是否存在
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}

	// 检查是否是有效的前端项目目录
	// 对于 Antd 项目，检查是否存在 package.json 和 src 目录
	packageJsonPath := filepath.Join(outputDir, "package.json")
	srcPath := filepath.Join(outputDir, "src")

	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("not a valid frontend project directory (missing package.json): %s", outputDir)
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("not a valid frontend project directory (missing src directory): %s", outputDir)
	}

	return nil
}

// FrontendField 前端字段信息
type FrontendField struct {
	*Field            // 嵌入基础字段信息
	TSType     string // TypeScript 类型
	FormType   string // 表单控件类型
	TableShow  bool   // 是否在表格中显示
	SearchShow bool   // 是否在搜索中显示
	FormShow   bool   // 是否在表单中显示
	Required   bool   // 是否必填
}

// getModelInfo 获取模型信息
func (g *FrontendGenerator) getModelInfo(modelName string) ([]*FrontendField, error) {
	// 创建示例字段，基于常见的产品分类模型
	var frontendFields []*FrontendField

	// 添加系统字段
	frontendFields = append(frontendFields, &FrontendField{
		Field: &Field{
			Name:     "ID",
			Type:     "uint",
			JSONName: "id",
		},
		TSType:     "number",
		FormType:   "hidden",
		TableShow:  true,
		SearchShow: false,
		FormShow:   false,
	})

	// 根据模型名称生成对应的业务字段
	switch strings.ToLower(modelName) {
	case "productcategory", "product_category":
		// 产品分类字段
		businessFields := []*FrontendField{
			{
				Field: &Field{
					Name:     "Name",
					Type:     "string",
					JSONName: "name",
				},
				TSType:     "string",
				FormType:   "input",
				TableShow:  true,
				SearchShow: true,
				FormShow:   true,
				Required:   true,
			},
			{
				Field: &Field{
					Name:     "Description",
					Type:     "string",
					JSONName: "description",
				},
				TSType:     "string",
				FormType:   "textarea",
				TableShow:  true,
				SearchShow: true,
				FormShow:   true,
				Required:   false,
			},
			{
				Field: &Field{
					Name:     "SortOrder",
					Type:     "int",
					JSONName: "sort_order",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  true,
				SearchShow: false,
				FormShow:   true,
				Required:   false,
			},
			{
				Field: &Field{
					Name:     "IsActive",
					Type:     "bool",
					JSONName: "is_active",
				},
				TSType:     "boolean",
				FormType:   "switch",
				TableShow:  true,
				SearchShow: false,
				FormShow:   true,
				Required:   false,
			},
		}
		frontendFields = append(frontendFields, businessFields...)
	default:
		// 默认字段
		defaultFields := []*FrontendField{
			{
				Field: &Field{
					Name:     "Name",
					Type:     "string",
					JSONName: "name",
				},
				TSType:     "string",
				FormType:   "input",
				TableShow:  true,
				SearchShow: true,
				FormShow:   true,
				Required:   true,
			},
		}
		frontendFields = append(frontendFields, defaultFields...)
	}

	// 添加时间戳字段
	timestampFields := []*FrontendField{
		{
			Field: &Field{
				Name:     "CreatedAt",
				Type:     "time.Time",
				JSONName: "created_at",
			},
			TSType:     "string",
			FormType:   "datetime",
			TableShow:  true,
			SearchShow: false,
			FormShow:   false,
		},
		{
			Field: &Field{
				Name:     "UpdatedAt",
				Type:     "time.Time",
				JSONName: "updated_at",
			},
			TSType:     "string",
			FormType:   "datetime",
			TableShow:  false,
			SearchShow: false,
			FormShow:   false,
		},
	}

	frontendFields = append(frontendFields, timestampFields...)

	return frontendFields, nil
}

// prepareTemplateData 准备模板数据
func (g *FrontendGenerator) prepareTemplateData(cfg *FrontendConfig, fields []*FrontendField) map[string]interface{} {
	modelName := cfg.Model

	// 确保正确的 PascalCase 转换
	pascalName := ToPascalCase(modelName)

	data := map[string]interface{}{
		"Model":       pascalName,
		"ModelLower":  strings.ToLower(pascalName),
		"ModelCamel":  ToCamelCase(pascalName),
		"ModelSnake":  ToSnakeCase(pascalName),
		"ModelKebab":  ToKebabCase(pascalName),
		"ModelPlural": Pluralize(strings.ToLower(pascalName)),
		"ModuleName":  cfg.ModuleName,
		"ModuleType":  string(cfg.ModuleType),
		"IsAdmin":     cfg.ModuleType == ModuleTypeAdmin,
		"IsPublic":    cfg.ModuleType == ModuleTypePublic,
		"ApiPrefix":   cfg.ApiPrefix,
		"Fields":      fields,
		"WithAuth":    cfg.WithAuth,
		"WithSearch":  cfg.WithSearch,
		"WithExport":  cfg.WithExport,
		"WithBatch":   cfg.WithBatch,
		"Year":        GetCurrentYear(),
	}

	// 设置默认值
	if data["ApiPrefix"] == "" || data["ApiPrefix"] == "/api/v1" {
		if cfg.ModuleType == ModuleTypeAdmin {
			data["ApiPrefix"] = "/api/v1/admin"
		} else {
			data["ApiPrefix"] = "/api/v1"
		}
	}
	if data["ModuleName"] == "" {
		data["ModuleName"] = strings.ToLower(modelName)
	}

	return data
}

// generateAntdCode 生成 Antd 前端代码
func (g *FrontendGenerator) generateAntdCode(data map[string]interface{}, cfg *FrontendConfig) error {
	// 生成页面组件
	if err := g.generateAntdPage(data, cfg); err != nil {
		return fmt.Errorf("failed to generate Antd page: %w", err)
	}

	// 生成 API 服务
	if err := g.generateAntdService(data, cfg); err != nil {
		return fmt.Errorf("failed to generate Antd service: %w", err)
	}

	// 生成类型定义
	if err := g.generateAntdTypes(data, cfg); err != nil {
		return fmt.Errorf("failed to generate Antd types: %w", err)
	}

	return nil
}

// generateVueCode 生成 Vue 前端代码
func (g *FrontendGenerator) generateVueCode(data map[string]interface{}, cfg *FrontendConfig) error {
	// TODO: 实现 Vue 代码生成
	return fmt.Errorf("Vue framework is not implemented yet")
}

// generateAntdPage 生成 Antd 页面组件
func (g *FrontendGenerator) generateAntdPage(data map[string]interface{}, cfg *FrontendConfig) error {
	content, err := g.templateEngine.Render("antd_page.tsx.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render Antd page template: %w", err)
	}

	// 确定文件路径
	moduleName := data["ModuleName"].(string)

	// 根据模块类型确定页面路径
	var pagePath string
	switch cfg.ModuleType {
	case ModuleTypeAdmin:
		// 管理后台模块：src/pages/admin/{module}/index.tsx
		pagePath = filepath.Join(cfg.OutputDir, "src", "pages", "admin", moduleName, "index.tsx")
	case ModuleTypePublic:
		// 普通用户模块：src/pages/{module}/index.tsx
		pagePath = filepath.Join(cfg.OutputDir, "src", "pages", moduleName, "index.tsx")
	default:
		// 默认为管理后台模块
		pagePath = filepath.Join(cfg.OutputDir, "src", "pages", "admin", moduleName, "index.tsx")
	}

	if err := g.writeFile(pagePath, content); err != nil {
		return fmt.Errorf("failed to write Antd page file: %w", err)
	}

	fmt.Printf("✅ Generated Antd page: %s\n", pagePath)
	return nil
}

// generateAntdService 生成 Antd API 服务
func (g *FrontendGenerator) generateAntdService(data map[string]interface{}, cfg *FrontendConfig) error {
	content, err := g.templateEngine.Render("antd_service.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render Antd service template: %w", err)
	}

	// 服务文件路径：src/services/{module}/api.ts
	moduleName := data["ModuleName"].(string)
	servicePath := filepath.Join(cfg.OutputDir, "src", "services", moduleName, "api.ts")

	if err := g.writeFile(servicePath, content); err != nil {
		return fmt.Errorf("failed to write Antd service file: %w", err)
	}

	fmt.Printf("✅ Generated Antd service: %s\n", servicePath)
	return nil
}

// generateAntdTypes 生成 Antd 类型定义
func (g *FrontendGenerator) generateAntdTypes(data map[string]interface{}, cfg *FrontendConfig) error {
	content, err := g.templateEngine.Render("antd_types.d.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render Antd types template: %w", err)
	}

	// 类型文件路径：src/services/{module}/typings.d.ts
	moduleName := data["ModuleName"].(string)
	typesPath := filepath.Join(cfg.OutputDir, "src", "services", moduleName, "typings.d.ts")

	if err := g.writeFile(typesPath, content); err != nil {
		return fmt.Errorf("failed to write Antd types file: %w", err)
	}

	fmt.Printf("✅ Generated Antd types: %s\n", typesPath)
	return nil
}

// writeFile 写入文件
func (g *FrontendGenerator) writeFile(path, content string) error {
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

// mapGoTypeToTS 将 Go 类型映射为 TypeScript 类型
func (g *FrontendGenerator) mapGoTypeToTS(goType string) string {
	typeMap := map[string]string{
		"string":    "string",
		"int":       "number",
		"int32":     "number",
		"int64":     "number",
		"uint":      "number",
		"uint32":    "number",
		"uint64":    "number",
		"float32":   "number",
		"float64":   "number",
		"bool":      "boolean",
		"time.Time": "string",
	}

	if tsType, exists := typeMap[goType]; exists {
		return tsType
	}
	return "string" // 默认为 string
}

// mapGoTypeToFormType 将 Go 类型映射为表单控件类型
func (g *FrontendGenerator) mapGoTypeToFormType(goType string) string {
	typeMap := map[string]string{
		"string":    "input",
		"int":       "number",
		"int32":     "number",
		"int64":     "number",
		"uint":      "number",
		"uint32":    "number",
		"uint64":    "number",
		"float32":   "number",
		"float64":   "number",
		"bool":      "switch",
		"time.Time": "datetime",
	}

	if formType, exists := typeMap[goType]; exists {
		return formType
	}
	return "input" // 默认为 input
}

// shouldShowInTable 判断字段是否应该在表格中显示
func (g *FrontendGenerator) shouldShowInTable(fieldName, fieldType string) bool {
	// 系统字段显示规则
	systemFields := map[string]bool{
		"ID":        true,
		"CreatedAt": true,
		"UpdatedAt": false,
		"DeletedAt": false,
	}

	if show, exists := systemFields[fieldName]; exists {
		return show
	}

	// 业务字段默认显示
	return true
}

// shouldShowInSearch 判断字段是否应该在搜索中显示
func (g *FrontendGenerator) shouldShowInSearch(fieldName, fieldType string) bool {
	// 系统字段不参与搜索
	systemFields := map[string]bool{
		"ID":        false,
		"CreatedAt": false,
		"UpdatedAt": false,
		"DeletedAt": false,
	}

	if show, exists := systemFields[fieldName]; exists {
		return show
	}

	// 字符串类型字段适合搜索
	return fieldType == "string"
}

// shouldShowInForm 判断字段是否应该在表单中显示
func (g *FrontendGenerator) shouldShowInForm(fieldName, fieldType string) bool {
	// 系统字段不在表单中编辑
	systemFields := map[string]bool{
		"ID":        false,
		"CreatedAt": false,
		"UpdatedAt": false,
		"DeletedAt": false,
	}

	if show, exists := systemFields[fieldName]; exists {
		return show
	}

	// 业务字段默认可编辑
	return true
}

// isFieldRequired 判断字段是否必填
func (g *FrontendGenerator) isFieldRequired(fieldName, fieldType string) bool {
	// 根据字段名称判断
	requiredFields := map[string]bool{
		"Name":  true,
		"Title": true,
		"Email": true,
	}

	return requiredFields[fieldName]
}
