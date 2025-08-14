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

// getModelInfo 获取模型信息 - 使用反射机制动态获取字段
func (g *FrontendGenerator) getModelInfo(modelName string) ([]*FrontendField, error) {
	// 尝试从现有模型文件中反射获取字段信息
	fields, err := g.getFieldsFromModel(modelName)
	if err != nil {
		// 如果反射失败，使用默认字段
		return g.getDefaultFields(), nil
	}

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

	// 转换反射获取的字段为前端字段
	for _, field := range fields {
		frontendField := g.convertToFrontendField(field)
		if frontendField != nil {
			frontendFields = append(frontendFields, frontendField)
		}
	}

	// 添加时间戳字段
	frontendFields = append(frontendFields, g.getTimestampFields()...)

	return frontendFields, nil
}

// getFieldsFromModel 从模型文件中反射获取字段信息
func (g *FrontendGenerator) getFieldsFromModel(modelName string) ([]*Field, error) {
	// 尝试使用模型反射器获取字段
	reflector := NewModelReflector()

	fields, err := reflector.ReflectModelFields(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to reflect model fields: %w", err)
	}

	return fields, nil
}

// convertToFrontendField 将模型字段转换为前端字段
func (g *FrontendGenerator) convertToFrontendField(field *Field) *FrontendField {
	// 跳过系统字段和时间戳字段
	if g.isSystemField(field.Name) {
		return nil
	}

	frontendField := &FrontendField{
		Field:      field,
		TSType:     g.getTypeScriptType(field.Type),
		FormType:   g.getFormType(field.Type, field.Name),
		TableShow:  g.shouldShowInTable(field.Name, field.Type),
		SearchShow: g.shouldShowInSearch(field.Name, field.Type),
		FormShow:   true,
		Required:   g.isRequiredField(field.Name, field.Type),
	}

	return frontendField
}

// isSystemField 判断是否为系统字段
func (g *FrontendGenerator) isSystemField(fieldName string) bool {
	systemFields := []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"}
	for _, sysField := range systemFields {
		if strings.EqualFold(fieldName, sysField) {
			return true
		}
	}
	return false
}

// getTypeScriptType 获取TypeScript类型
func (g *FrontendGenerator) getTypeScriptType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	default:
		return "string"
	}
}

// getFormType 获取表单控件类型
func (g *FrontendGenerator) getFormType(goType, fieldName string) string {
	fieldNameLower := strings.ToLower(fieldName)

	// 根据字段名推断控件类型
	if strings.Contains(fieldNameLower, "password") {
		return "password"
	}
	if strings.Contains(fieldNameLower, "email") {
		return "email"
	}
	if strings.Contains(fieldNameLower, "description") || strings.Contains(fieldNameLower, "content") {
		return "textarea"
	}
	if strings.Contains(fieldNameLower, "active") || strings.Contains(fieldNameLower, "enabled") {
		return "switch"
	}

	// 根据Go类型推断控件类型
	switch goType {
	case "bool":
		return "switch"
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "number"
	case "time.Time":
		return "datetime"
	default:
		return "input"
	}
}

// shouldShowInTable 判断是否在表格中显示
func (g *FrontendGenerator) shouldShowInTable(fieldName, goType string) bool {
	fieldNameLower := strings.ToLower(fieldName)

	// 不在表格中显示的字段
	hideInTable := []string{"description", "content", "password", "dimensions", "weight"}
	for _, hide := range hideInTable {
		if strings.Contains(fieldNameLower, hide) {
			return false
		}
	}

	// 外键字段通常不显示
	if strings.HasSuffix(fieldNameLower, "_id") && fieldNameLower != "id" {
		return false
	}

	return true
}

// shouldShowInSearch 判断是否在搜索中显示
func (g *FrontendGenerator) shouldShowInSearch(fieldName, goType string) bool {
	fieldNameLower := strings.ToLower(fieldName)

	// 可搜索的字段类型
	searchableFields := []string{"name", "title", "sku", "code", "email", "username"}
	for _, searchable := range searchableFields {
		if strings.Contains(fieldNameLower, searchable) {
			return true
		}
	}

	// 字符串类型的字段通常可搜索
	if goType == "string" && !strings.Contains(fieldNameLower, "password") {
		return true
	}

	return false
}

// isRequiredField 判断是否为必填字段
func (g *FrontendGenerator) isRequiredField(fieldName, goType string) bool {
	fieldNameLower := strings.ToLower(fieldName)

	// 通常必填的字段
	requiredFields := []string{"name", "title", "sku", "price", "stock_quantity"}
	for _, required := range requiredFields {
		if strings.Contains(fieldNameLower, required) {
			return true
		}
	}

	return false
}

// getTimestampFields 获取时间戳字段
func (g *FrontendGenerator) getTimestampFields() []*FrontendField {
	return []*FrontendField{
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
}

// getDefaultFields 获取默认字段（当反射失败时使用）
func (g *FrontendGenerator) getDefaultFields() []*FrontendField {
	return []*FrontendField{
		{
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
		},
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
}

// 保留原有的硬编码逻辑作为备用（已废弃，但保留以防需要）
func (g *FrontendGenerator) getModelInfoLegacy(modelName string) ([]*FrontendField, error) {
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
					Name:     "ParentID",
					Type:     "uint",
					JSONName: "parent_id",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  false,
				SearchShow: false,
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
	case "product":
		// 产品字段
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
					Name:     "CategoryID",
					Type:     "uint",
					JSONName: "category_id",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  false,
				SearchShow: false,
				FormShow:   true,
				Required:   false,
			},
			{
				Field: &Field{
					Name:     "SKU",
					Type:     "string",
					JSONName: "sku",
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
					Name:     "Price",
					Type:     "float64",
					JSONName: "price",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  true,
				SearchShow: false,
				FormShow:   true,
				Required:   true,
			},
			{
				Field: &Field{
					Name:     "CostPrice",
					Type:     "float64",
					JSONName: "cost_price",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  false,
				SearchShow: false,
				FormShow:   true,
				Required:   false,
			},
			{
				Field: &Field{
					Name:     "StockQuantity",
					Type:     "int",
					JSONName: "stock_quantity",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  true,
				SearchShow: false,
				FormShow:   true,
				Required:   true,
			},
			{
				Field: &Field{
					Name:     "MinStock",
					Type:     "int",
					JSONName: "min_stock",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  false,
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
			{
				Field: &Field{
					Name:     "Weight",
					Type:     "float64",
					JSONName: "weight",
				},
				TSType:     "number",
				FormType:   "number",
				TableShow:  false,
				SearchShow: false,
				FormShow:   true,
				Required:   false,
			},
			{
				Field: &Field{
					Name:     "Dimensions",
					Type:     "string",
					JSONName: "dimensions",
				},
				TSType:     "string",
				FormType:   "input",
				TableShow:  false,
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
		"DisplayName": pascalName,                  // 用于显示的名称
		"NameLower":   strings.ToLower(pascalName), // 小写名称，用于国际化key
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

	// 生成路由配置
	if err := g.generateAntdRoute(data, cfg); err != nil {
		return fmt.Errorf("failed to generate Antd route: %w", err)
	}

	// 生成国际化配置
	if err := g.generateAntdLocales(data, cfg); err != nil {
		return fmt.Errorf("failed to generate Antd locales: %w", err)
	}

	// 自动更新路由配置
	if err := g.updateAntdRoutes(data, cfg); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update routes automatically: %v\n", err)
		fmt.Printf("   Please manually add the route to config/routes.ts\n")
	} else {
		fmt.Printf("✅ Routes updated automatically\n")
	}

	// 自动更新国际化配置
	if err := g.updateAntdLocales(data, cfg); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update locales automatically: %v\n", err)
		fmt.Printf("   Please manually import the locale files\n")
	} else {
		fmt.Printf("✅ Locales updated automatically\n")
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

// generateAntdRoute 生成 Antd 路由配置
func (g *FrontendGenerator) generateAntdRoute(data map[string]interface{}, cfg *FrontendConfig) error {
	content, err := g.templateEngine.Render("antd_route.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render Antd route template: %w", err)
	}

	// 确定文件路径
	moduleName := data["ModuleName"].(string)
	routePath := filepath.Join(cfg.OutputDir, "docs", "generated", fmt.Sprintf("%s_route.ts", moduleName))

	if err := g.writeFile(routePath, content); err != nil {
		return fmt.Errorf("failed to write Antd route file: %w", err)
	}

	fmt.Printf("✅ Generated Antd route config: %s\n", routePath)
	fmt.Printf("📝 Please manually add the route configuration to config/routes.ts\n")
	return nil
}

// generateAntdLocales 生成 Antd 国际化配置
func (g *FrontendGenerator) generateAntdLocales(data map[string]interface{}, cfg *FrontendConfig) error {
	moduleName := data["ModuleName"].(string)

	// 生成中文国际化配置
	zhContent, err := g.templateEngine.Render("locale.zh-CN.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render zh-CN locale template: %w", err)
	}

	zhPath := filepath.Join(cfg.OutputDir, "src", "locales", "zh-CN", fmt.Sprintf("%s.ts", moduleName))
	if err := g.writeFile(zhPath, zhContent); err != nil {
		return fmt.Errorf("failed to write zh-CN locale file: %w", err)
	}

	// 生成英文国际化配置
	enContent, err := g.templateEngine.Render("locale.en-US.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render en-US locale template: %w", err)
	}

	enPath := filepath.Join(cfg.OutputDir, "src", "locales", "en-US", fmt.Sprintf("%s.ts", moduleName))
	if err := g.writeFile(enPath, enContent); err != nil {
		return fmt.Errorf("failed to write en-US locale file: %w", err)
	}

	fmt.Printf("✅ Generated locale files: %s, %s\n", zhPath, enPath)
	return nil
}

// updateAntdRoutes 自动更新路由配置
func (g *FrontendGenerator) updateAntdRoutes(data map[string]interface{}, cfg *FrontendConfig) error {
	routesPath := filepath.Join(cfg.OutputDir, "config", "routes.ts")

	// 读取现有路由配置
	content, err := os.ReadFile(routesPath)
	if err != nil {
		return fmt.Errorf("failed to read routes.ts: %w", err)
	}

	routesContent := string(content)
	moduleName := data["ModuleName"].(string)
	displayName := data["DisplayName"].(string)

	// 检查路由是否已存在
	routePattern := fmt.Sprintf("/admin/%s", moduleName)
	if strings.Contains(routesContent, routePattern) {
		return fmt.Errorf("route already exists: %s", routePattern)
	}

	// 找到插入位置（在 dict 路由后面）
	dictPattern := "{ path: '/admin/dict', name: '数据字典', component: './admin/dict' },"
	dictIndex := strings.Index(routesContent, dictPattern)
	if dictIndex == -1 {
		return fmt.Errorf("could not find dict route pattern in routes.ts")
	}

	// 在 dict 路由后添加新路由
	insertPos := dictIndex + len(dictPattern)
	newRoute := fmt.Sprintf("\n      { path: '/admin/%s', name: '%s管理', component: './admin/%s' },",
		moduleName, displayName, moduleName)

	updatedContent := routesContent[:insertPos] + newRoute + routesContent[insertPos:]

	// 写回文件
	if err := os.WriteFile(routesPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write routes.ts: %w", err)
	}

	return nil
}

// updateAntdLocales 自动更新国际化配置
func (g *FrontendGenerator) updateAntdLocales(data map[string]interface{}, cfg *FrontendConfig) error {
	moduleName := data["ModuleName"].(string)

	// 更新中文国际化配置
	if err := g.updateLocaleFile(cfg.OutputDir, "zh-CN", moduleName); err != nil {
		return fmt.Errorf("failed to update zh-CN locale: %w", err)
	}

	// 更新英文国际化配置
	if err := g.updateLocaleFile(cfg.OutputDir, "en-US", moduleName); err != nil {
		return fmt.Errorf("failed to update en-US locale: %w", err)
	}

	return nil
}

// updateLocaleFile 更新单个国际化文件
func (g *FrontendGenerator) updateLocaleFile(outputDir, locale, moduleName string) error {
	// 主国际化文件路径
	mainLocalePath := filepath.Join(outputDir, "src", "locales", fmt.Sprintf("%s.ts", locale))

	// 读取现有内容
	content, err := os.ReadFile(mainLocalePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", mainLocalePath, err)
	}

	localeContent := string(content)

	// 检查是否已经导入
	importStatement := fmt.Sprintf("import %s from './%s/%s';", moduleName, locale, moduleName)
	if strings.Contains(localeContent, importStatement) {
		return fmt.Errorf("locale import already exists for %s", moduleName)
	}

	// 找到导入部分的结束位置
	importEndPattern := "import settings from"
	importEndIndex := strings.Index(localeContent, importEndPattern)
	if importEndIndex == -1 {
		return fmt.Errorf("could not find import section end in %s", mainLocalePath)
	}

	// 在导入部分末尾添加新的导入
	// 找到这一行的结尾
	lineEnd := strings.Index(localeContent[importEndIndex:], "\n")
	if lineEnd == -1 {
		return fmt.Errorf("could not find line end after import section")
	}
	insertPos := importEndIndex + lineEnd

	newImport := fmt.Sprintf("\nimport %s from './%s/%s';", moduleName, locale, moduleName)
	updatedContent := localeContent[:insertPos] + newImport + localeContent[insertPos:]

	// 找到导出部分，添加新的模块
	exportPattern := "...component,"
	exportIndex := strings.Index(updatedContent, exportPattern)
	if exportIndex == -1 {
		return fmt.Errorf("could not find export section in %s", mainLocalePath)
	}

	// 在 component 后添加新模块
	insertPos = exportIndex + len(exportPattern)
	newExport := fmt.Sprintf("\n  ...%s,", moduleName)
	finalContent := updatedContent[:insertPos] + newExport + updatedContent[insertPos:]

	// 写回文件
	if err := os.WriteFile(mainLocalePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", mainLocalePath, err)
	}

	return nil
}

// writeFile 写入文件
func (g *FrontendGenerator) writeFile(path, content string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 检查文件是否已存在，如果存在则覆盖
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  File already exists, overwriting: %s\n", path)
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
