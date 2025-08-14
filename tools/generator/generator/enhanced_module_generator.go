package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// EnhancedModuleGenerator 增强的模块生成器
type EnhancedModuleGenerator struct {
	templateEngine *templates.Engine
}

// NewEnhancedModuleGenerator 创建增强的模块生成器
func NewEnhancedModuleGenerator() *EnhancedModuleGenerator {
	return &EnhancedModuleGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成完整的业务模块（增强版）
func (g *EnhancedModuleGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*EnhancedModuleConfig)
	if !ok {
		return fmt.Errorf("invalid config type for enhanced module generator")
	}

	// 解析字段
	parser := NewFieldParser()
	fields, err := parser.ParseFields(cfg.Fields)
	if err != nil {
		return fmt.Errorf("failed to parse fields: %w", err)
	}

	// 准备模板数据
	data := g.prepareTemplateData(cfg, fields)

	fmt.Printf("🚀 Starting enhanced module generation for: %s\n", cfg.Name)

	// 生成后端组件
	if err := g.generateBackendComponents(data, cfg); err != nil {
		return fmt.Errorf("failed to generate backend components: %w", err)
	}

	// 生成前端组件（如果配置了前端输出目录）
	if cfg.FrontendOutputDir != "" {
		if err := g.generateFrontendComponents(data, cfg); err != nil {
			return fmt.Errorf("failed to generate frontend components: %w", err)
		}
	}

	// 自动路由注册
	if cfg.AutoRouteRegister {
		if err := g.autoRegisterRoutes(data); err != nil {
			fmt.Printf("⚠️  Warning: Failed to auto-register routes: %v\n", err)
			g.printManualRouteInstructions(data)
		} else {
			fmt.Printf("✅ Routes registered automatically\n")
		}
	}

	// 自动数据库迁移
	if cfg.AutoMigration {
		if err := g.autoMigration(data); err != nil {
			fmt.Printf("⚠️  Warning: Failed to auto-migrate: %v\n", err)
		} else {
			fmt.Printf("✅ Database migration completed automatically\n")
		}
	}

	fmt.Printf("🎉 Enhanced module generation completed successfully!\n")
	return nil
}

// prepareTemplateData 准备模板数据（增强版）
func (g *EnhancedModuleGenerator) prepareTemplateData(cfg *EnhancedModuleConfig, fields []*Field) map[string]interface{} {
	name := ToPascalCase(cfg.Name)

	// 转换为前端字段
	frontendFields := g.convertToFrontendFields(fields)

	// 智能搜索字段配置
	searchFields := g.getSmartSearchFields(fields, cfg.SmartSearchFields)
	frontendSearchFields := g.convertToFrontendFields(searchFields)

	// 字段标签配置
	fieldLabels := g.prepareFieldLabels(fields, cfg.FieldLabels, cfg.FieldLabelsEn)

	// 生成请求字段
	createRequestFields := g.generateRequestFields(fields, false)
	updateRequestFields := g.generateRequestFields(fields, true)

	return map[string]interface{}{
		// 基础名称变量
		"Name":            name,
		"NameLower":       strings.ToLower(cfg.Name),
		"NameCamel":       ToCamelCase(cfg.Name),
		"NameSnake":       ToSnakeCase(cfg.Name),
		"NameKebab":       ToKebabCase(cfg.Name),
		"NamePlural":      Pluralize(strings.ToLower(cfg.Name)),
		"NamePluralCamel": ToCamelCase(Pluralize(cfg.Name)),

		// 模型变量
		"Model":            name,
		"ModelLower":       strings.ToLower(cfg.Name),
		"ModelCamel":       ToCamelCase(cfg.Name),
		"ModelSnake":       ToSnakeCase(cfg.Name),
		"ModelKebab":       ToKebabCase(cfg.Name),
		"ModelPlural":      Pluralize(strings.ToLower(cfg.Name)),
		"ModelPluralCamel": ToCamelCase(Pluralize(cfg.Name)),

		// 表名
		"TableName": Pluralize(ToSnakeCase(cfg.Name)),

		// 字段和配置（后端使用原始字段，前端使用转换后的字段）
		"Fields":              fields,               // 后端模板使用
		"FrontendFields":      frontendFields,       // 前端模板使用
		"SearchFields":        frontendSearchFields, // 前端搜索字段
		"FieldLabels":         fieldLabels,
		"CreateRequestFields": createRequestFields,
		"UpdateRequestFields": updateRequestFields,
		"WithAuth":            cfg.WithAuth,
		"WithCache":           cfg.WithCache,
		"Timestamp":           GenerateTimestamp(),
		"Year":                GetCurrentYear(),
		"DisplayName":         name, // 用于显示的名称

		// 增强功能标志
		"AutoRouteRegister":  cfg.AutoRouteRegister,
		"AutoMigration":      cfg.AutoMigration,
		"AutoI18n":           cfg.AutoI18n,
		"SmartSearchFields":  cfg.SmartSearchFields,
		"FrontendOutputDir":  cfg.FrontendOutputDir,
		"FrontendFramework":  cfg.FrontendFramework,
		"FrontendModuleType": cfg.FrontendModuleType,
	}
}

// convertToFrontendFields 将普通字段转换为前端字段
func (g *EnhancedModuleGenerator) convertToFrontendFields(fields []*Field) []*FrontendField {
	var frontendFields []*FrontendField

	for _, field := range fields {
		frontendField := &FrontendField{
			Field:      field,
			TSType:     g.getTypeScriptType(field.Type),
			FormType:   g.getFormType(field.Type, field.Name),
			TableShow:  g.shouldShowInTable(field.Name, field.Type),
			SearchShow: g.shouldShowInSearch(field.Name, field.Type),
			FormShow:   true,
			Required:   g.isRequiredField(field.Name, field.Type),
		}
		frontendFields = append(frontendFields, frontendField)
	}

	return frontendFields
}

// getTypeScriptType 获取TypeScript类型
func (g *EnhancedModuleGenerator) getTypeScriptType(goType string) string {
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
func (g *EnhancedModuleGenerator) getFormType(goType, fieldName string) string {
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
func (g *EnhancedModuleGenerator) shouldShowInTable(fieldName, goType string) bool {
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
func (g *EnhancedModuleGenerator) shouldShowInSearch(fieldName, goType string) bool {
	fieldNameLower := strings.ToLower(fieldName)

	// 可搜索的字段名模式
	searchablePatterns := []string{
		"name", "title", "sku", "code", "email", "username",
		"description", "reason", "reference_id", "change_type",
	}

	for _, pattern := range searchablePatterns {
		if strings.Contains(fieldNameLower, pattern) {
			return true
		}
	}

	// 字符串类型的字段通常可搜索（排除密码等敏感字段）
	if goType == "string" && !strings.Contains(fieldNameLower, "password") {
		return true
	}

	return false
}

// isRequiredField 判断是否为必填字段
func (g *EnhancedModuleGenerator) isRequiredField(fieldName, goType string) bool {
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

// getSmartSearchFields 智能获取搜索字段
func (g *EnhancedModuleGenerator) getSmartSearchFields(fields []*Field, smartSearch bool) []*Field {
	if !smartSearch {
		// 如果不启用智能搜索，返回默认的Name字段（如果存在）
		for _, field := range fields {
			if strings.EqualFold(field.Name, "name") {
				return []*Field{field}
			}
		}
		return []*Field{}
	}

	var searchFields []*Field

	// 智能选择搜索字段
	for _, field := range fields {
		if g.isSearchableField(field) {
			searchFields = append(searchFields, field)
		}
	}

	return searchFields
}

// isSearchableField 判断字段是否适合搜索
func (g *EnhancedModuleGenerator) isSearchableField(field *Field) bool {
	fieldNameLower := strings.ToLower(field.Name)

	// 可搜索的字段名模式
	searchablePatterns := []string{
		"name", "title", "sku", "code", "email", "username",
		"description", "reason", "reference_id", "change_type",
	}

	for _, pattern := range searchablePatterns {
		if strings.Contains(fieldNameLower, pattern) {
			return true
		}
	}

	// 字符串类型的字段通常可搜索（排除密码等敏感字段）
	if field.Type == "string" && !strings.Contains(fieldNameLower, "password") {
		return true
	}

	return false
}

// prepareFieldLabels 准备字段标签
func (g *EnhancedModuleGenerator) prepareFieldLabels(fields []*Field, zhLabels, enLabels map[string]string) map[string]interface{} {
	fieldLabels := make(map[string]interface{})

	for _, field := range fields {
		fieldName := field.Name

		// 中文标签
		zhLabel := zhLabels[fieldName]
		if zhLabel == "" {
			zhLabel = g.generateDefaultZhLabel(fieldName)
		}

		// 英文标签
		enLabel := enLabels[fieldName]
		if enLabel == "" {
			enLabel = g.generateDefaultEnLabel(fieldName)
		}

		fieldLabels[fieldName] = map[string]string{
			"zh": zhLabel,
			"en": enLabel,
		}
	}

	return fieldLabels
}

// generateDefaultZhLabel 生成默认中文标签
func (g *EnhancedModuleGenerator) generateDefaultZhLabel(fieldName string) string {
	// 常见字段的中文标签映射
	labelMap := map[string]string{
		"ProductId":      "产品ID",
		"ChangeType":     "变更类型",
		"QuantityChange": "变更数量",
		"QuantityBefore": "变更前数量",
		"QuantityAfter":  "变更后数量",
		"Reason":         "变更原因",
		"OperatorId":     "操作员ID",
		"ReferenceId":    "关联单据ID",
		"ReferenceType":  "关联单据类型",
		"Name":           "名称",
		"Description":    "描述",
		"Price":          "价格",
		"Stock":          "库存",
		"IsActive":       "是否启用",
		"CreatedAt":      "创建时间",
		"UpdatedAt":      "更新时间",
	}

	if label, exists := labelMap[fieldName]; exists {
		return label
	}

	// 如果没有预定义标签，返回字段名
	return fieldName
}

// generateDefaultEnLabel 生成默认英文标签
func (g *EnhancedModuleGenerator) generateDefaultEnLabel(fieldName string) string {
	// 常见字段的英文标签映射
	labelMap := map[string]string{
		"ProductId":      "Product ID",
		"ChangeType":     "Change Type",
		"QuantityChange": "Quantity Change",
		"QuantityBefore": "Quantity Before",
		"QuantityAfter":  "Quantity After",
		"Reason":         "Reason",
		"OperatorId":     "Operator ID",
		"ReferenceId":    "Reference ID",
		"ReferenceType":  "Reference Type",
		"Name":           "Name",
		"Description":    "Description",
		"Price":          "Price",
		"Stock":          "Stock",
		"IsActive":       "Active",
		"CreatedAt":      "Created At",
		"UpdatedAt":      "Updated At",
	}

	if label, exists := labelMap[fieldName]; exists {
		return label
	}

	// 如果没有预定义标签，使用字段名并添加空格
	return g.addSpacesToCamelCase(fieldName)
}

// addSpacesToCamelCase 在驼峰命名中添加空格
func (g *EnhancedModuleGenerator) addSpacesToCamelCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return string(result)
}

// generateBackendComponents 生成后端组件
func (g *EnhancedModuleGenerator) generateBackendComponents(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	fmt.Printf("📦 Generating backend components...\n")

	// 生成模型
	if err := g.generateModel(data); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}
	fmt.Printf("✅ Generated model\n")

	// 生成仓储
	if err := g.generateRepository(data); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}
	fmt.Printf("✅ Generated repository\n")

	// 生成服务
	if err := g.generateService(data); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}
	fmt.Printf("✅ Generated service\n")

	// 生成处理器
	if err := g.generateHandler(data); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}
	fmt.Printf("✅ Generated handler\n")

	// 生成测试
	if err := g.generateTests(data); err != nil {
		return fmt.Errorf("failed to generate tests: %w", err)
	}
	fmt.Printf("✅ Generated tests\n")

	// 生成迁移
	if err := g.generateMigration(data); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}
	fmt.Printf("✅ Generated migration\n")

	return nil
}

// generateFrontendComponents 生成前端组件
func (g *EnhancedModuleGenerator) generateFrontendComponents(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	fmt.Printf("🎨 Generating frontend components...\n")

	// 使用增强的前端生成器
	if err := g.generateEnhancedFrontend(data, cfg); err != nil {
		return fmt.Errorf("failed to generate enhanced frontend: %w", err)
	}

	fmt.Printf("✅ Generated frontend components\n")
	return nil
}

// generateEnhancedFrontend 生成增强的前端代码
func (g *EnhancedModuleGenerator) generateEnhancedFrontend(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	// 生成页面组件
	if err := g.generateEnhancedPage(data, cfg); err != nil {
		return fmt.Errorf("failed to generate enhanced page: %w", err)
	}

	// 生成 API 服务
	if err := g.generateEnhancedService(data, cfg); err != nil {
		return fmt.Errorf("failed to generate enhanced service: %w", err)
	}

	// 生成类型定义
	if err := g.generateEnhancedTypes(data, cfg); err != nil {
		return fmt.Errorf("failed to generate enhanced types: %w", err)
	}

	// 生成路由配置
	if err := g.generateEnhancedRoute(data, cfg); err != nil {
		return fmt.Errorf("failed to generate enhanced route: %w", err)
	}

	// 生成国际化配置
	if cfg.AutoI18n {
		if err := g.generateEnhancedLocales(data, cfg); err != nil {
			return fmt.Errorf("failed to generate enhanced locales: %w", err)
		}

		// 自动更新国际化配置
		if err := g.updateEnhancedLocales(data, cfg); err != nil {
			fmt.Printf("⚠️  Warning: Failed to update locales automatically: %v\n", err)
			fmt.Printf("   Please manually import the locale files\n")
		} else {
			fmt.Printf("✅ Locales updated automatically\n")
		}
	}

	// 自动更新路由配置
	if err := g.updateEnhancedRoutes(data, cfg); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update routes automatically: %v\n", err)
		fmt.Printf("   Please manually add the route to config/routes.ts\n")
	} else {
		fmt.Printf("✅ Routes updated automatically\n")
	}

	return nil
}

// generateEnhancedPage 生成增强的页面组件
func (g *EnhancedModuleGenerator) generateEnhancedPage(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	content, err := g.templateEngine.Render("enhanced_antd_page.tsx.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render enhanced page template: %w", err)
	}

	// 确定文件路径
	moduleName := strings.ToLower(cfg.Name)

	// 根据模块类型确定页面路径
	var pagePath string
	switch cfg.FrontendModuleType {
	case ModuleTypeAdmin:
		// 管理后台模块：src/pages/admin/{module}/index.tsx
		pagePath = filepath.Join(cfg.FrontendOutputDir, "src", "pages", "admin", moduleName, "index.tsx")
	case ModuleTypePublic:
		// 普通用户模块：src/pages/{module}/index.tsx
		pagePath = filepath.Join(cfg.FrontendOutputDir, "src", "pages", moduleName, "index.tsx")
	default:
		// 默认为管理后台模块
		pagePath = filepath.Join(cfg.FrontendOutputDir, "src", "pages", "admin", moduleName, "index.tsx")
	}

	if err := g.writeFile(pagePath, content); err != nil {
		return fmt.Errorf("failed to write enhanced page file: %w", err)
	}

	fmt.Printf("✅ Generated enhanced page: %s\n", pagePath)
	return nil
}

// generateEnhancedService 生成增强的 API 服务
func (g *EnhancedModuleGenerator) generateEnhancedService(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	content, err := g.templateEngine.Render("antd_service.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render service template: %w", err)
	}

	// 服务文件路径：src/services/{module}/api.ts
	moduleName := strings.ToLower(cfg.Name)
	servicePath := filepath.Join(cfg.FrontendOutputDir, "src", "services", moduleName, "api.ts")

	if err := g.writeFile(servicePath, content); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	fmt.Printf("✅ Generated enhanced service: %s\n", servicePath)
	return nil
}

// generateEnhancedTypes 生成增强的类型定义
func (g *EnhancedModuleGenerator) generateEnhancedTypes(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	content, err := g.templateEngine.Render("antd_types.d.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render types template: %w", err)
	}

	// 类型文件路径：src/services/{module}/typings.d.ts
	moduleName := strings.ToLower(cfg.Name)
	typesPath := filepath.Join(cfg.FrontendOutputDir, "src", "services", moduleName, "typings.d.ts")

	if err := g.writeFile(typesPath, content); err != nil {
		return fmt.Errorf("failed to write types file: %w", err)
	}

	fmt.Printf("✅ Generated enhanced types: %s\n", typesPath)
	return nil
}

// generateEnhancedRoute 生成增强的路由配置
func (g *EnhancedModuleGenerator) generateEnhancedRoute(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	content, err := g.templateEngine.Render("antd_route.ts.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render route template: %w", err)
	}

	// 确定文件路径
	moduleName := strings.ToLower(cfg.Name)
	routePath := filepath.Join(cfg.FrontendOutputDir, "docs", "generated", fmt.Sprintf("%s_route.ts", moduleName))

	if err := g.writeFile(routePath, content); err != nil {
		return fmt.Errorf("failed to write route file: %w", err)
	}

	fmt.Printf("✅ Generated enhanced route config: %s\n", routePath)
	return nil
}

// generateEnhancedLocales 生成增强的国际化配置
func (g *EnhancedModuleGenerator) generateEnhancedLocales(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	moduleName := strings.ToLower(cfg.Name)

	// 生成中文国际化配置
	zhContent, err := g.templateEngine.Render("enhanced_locale.zh-CN.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render zh-CN locale template: %w", err)
	}

	zhPath := filepath.Join(cfg.FrontendOutputDir, "src", "locales", "zh-CN", fmt.Sprintf("%s.ts", moduleName))
	if err := g.writeFile(zhPath, zhContent); err != nil {
		return fmt.Errorf("failed to write zh-CN locale file: %w", err)
	}

	// 生成英文国际化配置
	enContent, err := g.templateEngine.Render("enhanced_locale.en-US.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render en-US locale template: %w", err)
	}

	enPath := filepath.Join(cfg.FrontendOutputDir, "src", "locales", "en-US", fmt.Sprintf("%s.ts", moduleName))
	if err := g.writeFile(enPath, enContent); err != nil {
		return fmt.Errorf("failed to write en-US locale file: %w", err)
	}

	fmt.Printf("✅ Generated enhanced locale files: %s, %s\n", zhPath, enPath)
	return nil
}

// updateEnhancedRoutes 自动更新路由配置
func (g *EnhancedModuleGenerator) updateEnhancedRoutes(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	routesPath := filepath.Join(cfg.FrontendOutputDir, "config", "routes.ts")

	// 读取现有路由配置
	content, err := os.ReadFile(routesPath)
	if err != nil {
		return fmt.Errorf("failed to read routes.ts: %w", err)
	}

	routesContent := string(content)
	moduleName := strings.ToLower(cfg.Name)
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

// updateEnhancedLocales 自动更新国际化配置
func (g *EnhancedModuleGenerator) updateEnhancedLocales(data map[string]interface{}, cfg *EnhancedModuleConfig) error {
	moduleName := strings.ToLower(cfg.Name)

	// 更新中文国际化配置
	if err := g.updateEnhancedLocaleFile(cfg.FrontendOutputDir, "zh-CN", moduleName); err != nil {
		return fmt.Errorf("failed to update zh-CN locale: %w", err)
	}

	// 更新英文国际化配置
	if err := g.updateEnhancedLocaleFile(cfg.FrontendOutputDir, "en-US", moduleName); err != nil {
		return fmt.Errorf("failed to update en-US locale: %w", err)
	}

	return nil
}

// updateEnhancedLocaleFile 更新单个国际化文件
func (g *EnhancedModuleGenerator) updateEnhancedLocaleFile(outputDir, locale, moduleName string) error {
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

// autoRegisterRoutes 自动注册路由
func (g *EnhancedModuleGenerator) autoRegisterRoutes(data map[string]interface{}) error {
	fmt.Printf("🔗 Auto-registering routes...\n")

	// 更新 server.go
	if err := g.updateServerFile(data); err != nil {
		return fmt.Errorf("failed to update server.go: %w", err)
	}

	// 更新 main.go
	if err := g.updateMainFile(data); err != nil {
		return fmt.Errorf("failed to update main.go: %w", err)
	}

	return nil
}

// autoMigration 自动执行数据库迁移
func (g *EnhancedModuleGenerator) autoMigration(data map[string]interface{}) error {
	fmt.Printf("🗄️  Auto-migrating database...\n")

	// 使用 GORM AutoMigrate
	modelName := data["Model"].(string)

	// 创建临时的 AutoMigrate 脚本
	scriptContent := fmt.Sprintf(`package main

import (
	"log"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	logger, err := logger.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// 连接数据库
	db, err := database.New(cfg, logger)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移%s表
	err = db.AutoMigrate(&model.%s{})
	if err != nil {
		log.Fatal("Failed to migrate %s table:", err)
	}

	log.Println("%s table migrated successfully!")
}`, modelName, modelName, modelName, modelName)

	// 写入临时脚本
	scriptPath := fmt.Sprintf("cmd/automigrate_%s/main.go", strings.ToLower(modelName))
	if err := g.writeFile(scriptPath, scriptContent); err != nil {
		return fmt.Errorf("failed to write migration script: %w", err)
	}

	// 执行迁移脚本
	// 注意：这里应该在实际环境中执行，但为了安全起见，我们只生成脚本
	fmt.Printf("📝 Migration script generated: %s\n", scriptPath)
	fmt.Printf("💡 Run: go run %s\n", scriptPath)

	return nil
}

// generateModel 生成模型
func (g *EnhancedModuleGenerator) generateModel(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("model.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.go", data["NameSnake"])
	path := filepath.Join("internal", "model", filename)

	return g.writeFile(path, content)
}

// generateRepository 生成仓储
func (g *EnhancedModuleGenerator) generateRepository(data map[string]interface{}) error {
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
func (g *EnhancedModuleGenerator) generateService(data map[string]interface{}) error {
	content, err := g.templateEngine.Render("service.go.tmpl", data)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_service.go", data["NameSnake"])
	path := filepath.Join("internal", "service", filename)

	return g.writeFile(path, content)
}

// generateHandler 生成处理器
func (g *EnhancedModuleGenerator) generateHandler(data map[string]interface{}) error {
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
func (g *EnhancedModuleGenerator) generateTests(data map[string]interface{}) error {
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
func (g *EnhancedModuleGenerator) generateMigration(data map[string]interface{}) error {
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

// updateServerFile 更新 server.go 文件
func (g *EnhancedModuleGenerator) updateServerFile(data map[string]interface{}) error {
	serverPath := "internal/server/server.go"

	// 读取现有内容
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	serverContent := string(content)
	modelName := data["Model"].(string)
	modelCamel := data["ModelCamel"].(string)

	// 检查是否已经存在
	if strings.Contains(serverContent, fmt.Sprintf("%sHandler", modelCamel)) {
		return fmt.Errorf("%s handler already exists in server.go", modelName)
	}

	// 1. 添加字段到 Server 结构体
	structPattern := "dictHandler            *handler.DictHandler"
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
	initPattern := "dictHandler:                dictHandler,"
	initIndex := strings.Index(serverContent, initPattern)
	if initIndex == -1 {
		return fmt.Errorf("could not find field initialization section")
	}

	initEnd := initIndex + len(initPattern)
	newInit := fmt.Sprintf("\n\t\t%sHandler: %sHandler,", modelCamel, modelCamel)
	serverContent = serverContent[:initEnd] + newInit + serverContent[initEnd:]

	// 4. 添加路由注册
	routePattern := "// ProductStockHistory管理路由"
	routeIndex := strings.Index(serverContent, routePattern)
	if routeIndex != -1 {
		// 如果找到了ProductStockHistory的注释，在其后添加
		routeEnd := strings.Index(serverContent[routeIndex:], "\n")
		if routeEnd != -1 {
			insertPos := routeIndex + routeEnd
			newRoute := fmt.Sprintf("\n\t\t\t\t\t// %s管理路由\n\t\t\t\t\ts.%sHandler.RegisterRoutes(admin)", modelName, modelCamel)
			serverContent = serverContent[:insertPos] + newRoute + serverContent[insertPos:]
		}
	} else {
		// 如果没有找到，在admin路由组的末尾添加
		adminEndPattern := "}\n\n\t\t\t}"
		adminEndIndex := strings.Index(serverContent, adminEndPattern)
		if adminEndIndex == -1 {
			return fmt.Errorf("could not find admin routes section end")
		}

		newRoute := fmt.Sprintf("\n\n\t\t\t\t\t// %s管理路由\n\t\t\t\t\ts.%sHandler.RegisterRoutes(admin)", modelName, modelCamel)
		serverContent = serverContent[:adminEndIndex] + newRoute + serverContent[adminEndIndex:]
	}

	// 写回文件
	if err := os.WriteFile(serverPath, []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to write server.go: %w", err)
	}

	return nil
}

// updateMainFile 更新 main.go 文件
func (g *EnhancedModuleGenerator) updateMainFile(data map[string]interface{}) error {
	mainPath := "cmd/server/main.go"

	// 读取现有内容
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	mainContent := string(content)
	modelName := data["Model"].(string)

	// 检查是否已经存在
	if strings.Contains(mainContent, fmt.Sprintf("repository.New%sRepository", modelName)) {
		return fmt.Errorf("%s repository already exists in main.go", modelName)
	}

	// 1. 添加仓储提供者
	repoPattern := "repository.NewDictRepository,"
	repoIndex := strings.Index(mainContent, repoPattern)
	if repoIndex == -1 {
		return fmt.Errorf("could not find repository providers section in main.go")
	}

	repoEnd := repoIndex + len(repoPattern)
	newRepo := fmt.Sprintf("\n\t\t\trepository.New%sRepository,", modelName)
	mainContent = mainContent[:repoEnd] + newRepo + mainContent[repoEnd:]

	// 2. 添加服务提供者
	servicePattern := "service.NewDictService,"
	serviceIndex := strings.Index(mainContent, servicePattern)
	if serviceIndex == -1 {
		return fmt.Errorf("could not find service providers section in main.go")
	}

	serviceEnd := serviceIndex + len(servicePattern)
	newService := fmt.Sprintf("\n\t\t\tservice.New%sService,", modelName)
	mainContent = mainContent[:serviceEnd] + newService + mainContent[serviceEnd:]

	// 3. 添加处理器提供者
	handlerPattern := "handler.NewDictHandler,"
	handlerIndex := strings.Index(mainContent, handlerPattern)
	if handlerIndex == -1 {
		return fmt.Errorf("could not find handler providers section in main.go")
	}

	handlerEnd := handlerIndex + len(handlerPattern)
	newHandler := fmt.Sprintf("\n\t\t\thandler.New%sHandler,", modelName)
	mainContent = mainContent[:handlerEnd] + newHandler + mainContent[handlerEnd:]

	// 写回文件
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}

// printManualRouteInstructions 打印手动路由注册说明
func (g *EnhancedModuleGenerator) printManualRouteInstructions(data map[string]interface{}) {
	modelName := data["Model"].(string)
	modelCamel := data["ModelCamel"].(string)

	fmt.Printf("⚠️  Please manually add the handler to internal/server/server.go:\n")
	fmt.Printf("   1. Add field to Server struct: %sHandler *handler.%sHandler\n", modelCamel, modelName)
	fmt.Printf("   2. Add parameter to New function: %sHandler *handler.%sHandler\n", modelCamel, modelName)
	fmt.Printf("   3. Initialize field in New function: %sHandler: %sHandler,\n", modelCamel, modelCamel)
	fmt.Printf("   4. Register routes in setupRoutes: s.%sHandler.RegisterRoutes(admin)\n", modelCamel)
	fmt.Printf("   5. Add providers to cmd/server/main.go:\n")
	fmt.Printf("      - repository.New%sRepository,\n", modelName)
	fmt.Printf("      - service.New%sService,\n", modelName)
	fmt.Printf("      - handler.New%sHandler,\n", modelName)
}

// writeFile 写入文件
func (g *EnhancedModuleGenerator) writeFile(path, content string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 检查文件是否已存在
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  File already exists, overwriting: %s\n", path)
	}

	// 写入文件
	return os.WriteFile(path, []byte(content), 0644)
}

// appendToFile 追加到文件
func (g *EnhancedModuleGenerator) appendToFile(path, content string) error {
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

// generateRequestFields 生成请求字段
func (g *EnhancedModuleGenerator) generateRequestFields(fields []*Field, isUpdate bool) string {
	var fieldLines []string

	for _, field := range fields {
		// 跳过系统字段
		if field.JSONName == "id" || field.JSONName == "created_at" || field.JSONName == "updated_at" {
			continue
		}

		// 构建字段行
		fieldLine := fmt.Sprintf("\t%s %s `json:\"%s\"`", field.Name, field.Type, field.JSONName)
		fieldLines = append(fieldLines, fieldLine)
	}

	if len(fieldLines) == 0 {
		return "\t// No fields"
	}

	return strings.Join(fieldLines, "\n")
}
