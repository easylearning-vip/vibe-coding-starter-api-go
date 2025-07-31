package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// DatabaseTableGenerator 数据库表生成器
type DatabaseTableGenerator struct {
	templateEngine *templates.Engine
}

// NewDatabaseTableGenerator 创建数据库表生成器
func NewDatabaseTableGenerator() *DatabaseTableGenerator {
	return &DatabaseTableGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 从数据库表生成模型代码
func (g *DatabaseTableGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*DatabaseTableConfig)
	if !ok {
		return fmt.Errorf("invalid config type for database table generator")
	}

	// 创建表结构读取器
	reader, err := NewTableReader(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	if err != nil {
		return fmt.Errorf("failed to create table reader: %w", err)
	}
	defer reader.Close()

	// 读取表结构
	fields, err := reader.ReadTableStructure(cfg.TableName)
	if err != nil {
		return fmt.Errorf("failed to read table structure: %w", err)
	}

	// 获取表注释
	tableComment, err := reader.GetTableComment(cfg.TableName)
	if err != nil {
		return fmt.Errorf("failed to get table comment: %w", err)
	}

	// 确定模型名称
	modelName := cfg.ModelName
	if modelName == "" {
		modelName = g.generateModelName(cfg.TableName)
	}

	// 转换字段格式
	modelFields := g.convertDatabaseFieldsToModelFields(fields)

	// 获取需要的导入包
	requiredImports := reader.GetRequiredImports(fields)

	// 准备模板数据
	data := map[string]interface{}{
		"Name":            ToPascalCase(modelName),
		"NameLower":       strings.ToLower(modelName),
		"NameCamel":       ToCamelCase(modelName),
		"NameSnake":       ToSnakeCase(modelName),
		"NamePlural":      Pluralize(ToSnakeCase(modelName)),
		"TableName":       cfg.TableName,
		"TableComment":    tableComment,
		"Fields":          modelFields,
		"RequiredImports": requiredImports,
		"WithTimestamps":  cfg.WithTimestamps,
		"WithSoftDelete":  cfg.WithSoftDelete,
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

	fmt.Printf("✅ Model '%s' generated from table '%s' successfully!\n", modelName, cfg.TableName)
	fmt.Printf("   📁 File: %s\n", path)
	fmt.Printf("   📊 Fields: %d\n", len(modelFields))
	if tableComment != "" {
		fmt.Printf("   💬 Comment: %s\n", tableComment)
	}

	return nil
}

// generateModelName 从表名生成模型名称
func (g *DatabaseTableGenerator) generateModelName(tableName string) string {
	// 移除常见的表前缀
	prefixes := []string{"tbl_", "tb_", "t_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(tableName, prefix) {
			tableName = strings.TrimPrefix(tableName, prefix)
			break
		}
	}

	// 移除复数后缀，转换为单数形式
	if strings.HasSuffix(tableName, "ies") {
		tableName = strings.TrimSuffix(tableName, "ies") + "y"
	} else if strings.HasSuffix(tableName, "es") {
		tableName = strings.TrimSuffix(tableName, "es")
	} else if strings.HasSuffix(tableName, "s") && !strings.HasSuffix(tableName, "ss") {
		tableName = strings.TrimSuffix(tableName, "s")
	}

	return ToPascalCase(tableName)
}

// convertDatabaseFieldsToModelFields 将数据库字段转换为模型字段格式
func (g *DatabaseTableGenerator) convertDatabaseFieldsToModelFields(dbFields []*DatabaseField) []*Field {
	var fields []*Field

	for _, dbField := range dbFields {
		field := &Field{
			Name:     dbField.Name,
			Type:     dbField.Type,
			JSONName: dbField.JSONName,
			GormTag:  dbField.GormTag,
			Comment:  dbField.Comment,
		}
		fields = append(fields, field)
	}

	return fields
}

// writeFile 写入文件
func (g *DatabaseTableGenerator) writeFile(path, content string) error {
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

// GenerateFromAllTables 从数据库中的所有表生成模型
func (g *DatabaseTableGenerator) GenerateFromAllTables(cfg *DatabaseTableConfig) error {
	// 创建表结构读取器
	reader, err := NewTableReader(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	if err != nil {
		return fmt.Errorf("failed to create table reader: %w", err)
	}
	defer reader.Close()

	// 获取所有表
	tables, err := reader.ListTables()
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	fmt.Printf("🚀 Found %d tables in database '%s'\n\n", len(tables), cfg.DatabaseName)

	// 为每个表生成模型
	for i, tableName := range tables {
		fmt.Printf("📦 [%d/%d] Generating model for table '%s'...\n", i+1, len(tables), tableName)

		// 创建单表配置
		tableConfig := &DatabaseTableConfig{
			DatabaseHost:     cfg.DatabaseHost,
			DatabasePort:     cfg.DatabasePort,
			DatabaseUser:     cfg.DatabaseUser,
			DatabasePassword: cfg.DatabasePassword,
			DatabaseName:     cfg.DatabaseName,
			TableName:        tableName,
			ModelName:        "", // 自动生成
			WithTimestamps:   cfg.WithTimestamps,
			WithSoftDelete:   cfg.WithSoftDelete,
		}

		if err := g.Generate(tableConfig); err != nil {
			fmt.Printf("❌ Failed to generate model for table '%s': %v\n", tableName, err)
			continue
		}
	}

	fmt.Printf("\n🎉 All models generated successfully!\n")
	return nil
}
