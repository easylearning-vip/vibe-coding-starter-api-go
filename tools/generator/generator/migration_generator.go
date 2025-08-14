package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vibe-coding-starter/tools/generator/templates"
)

// MigrationGenerator 迁移生成器
type MigrationGenerator struct {
	templateEngine *templates.Engine
}

// NewMigrationGenerator 创建迁移生成器
func NewMigrationGenerator() *MigrationGenerator {
	return &MigrationGenerator{
		templateEngine: templates.NewEngine(),
	}
}

// Generate 生成数据库迁移
func (g *MigrationGenerator) Generate(config interface{}) error {
	cfg, ok := config.(*MigrationConfig)
	if !ok {
		return fmt.Errorf("invalid config type for migration generator")
	}

	// 读取配置文件获取数据库类型
	configReader := NewConfigReader("")
	databaseType, err := configReader.GetDatabaseType()
	if err != nil {
		return fmt.Errorf("failed to read database config: %w", err)
	}

	migrationDir, err := configReader.GetMigrationDir()
	if err != nil {
		return fmt.Errorf("failed to get migration directory: %w", err)
	}

	// 如果没有指定表名，从迁移名称推断
	table := cfg.Table
	if table == "" {
		table = g.extractTableFromName(cfg.Name)
		// 只有从迁移名称推断的表名才需要应用命名规则
		tableName := ToSnakeCase(Pluralize(ToPascalCase(table)))
		table = tableName
	}

	// 使用提供的表名（已经是正确格式）
	tableName := table

	// 解析字段信息
	var fields []*Field
	if cfg.Fields != "" {
		parser := NewFieldParser()
		var err error
		fields, err = parser.ParseFields(cfg.Fields)
		if err != nil {
			return fmt.Errorf("failed to parse fields: %w", err)
		}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Name":         cfg.Name,
		"Table":        table,
		"TableName":    tableName, // 使用正确的表名
		"Action":       cfg.Action,
		"DatabaseType": databaseType,
		"Timestamp":    GenerateTimestamp(),
		"Year":         GetCurrentYear(),
		"Fields":       fields,
	}

	// 根据操作类型选择模板
	var templateName string
	switch cfg.Action {
	case "create":
		templateName = "migration.sql.tmpl"
	case "alter":
		templateName = "migration_alter.sql.tmpl"
	case "drop":
		templateName = "migration_drop.sql.tmpl"
	default:
		templateName = "migration.sql.tmpl"
	}

	// 生成迁移文件
	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render migration template: %w", err)
	}

	// 生成文件名和路径
	filename := fmt.Sprintf("%s_%s.up.sql", data["Timestamp"], ToSnakeCase(cfg.Name))
	path := filepath.Join(migrationDir, filename)

	if err := g.writeFile(path, content); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	// 生成回滚文件
	if err := g.generateRollback(data, cfg, migrationDir); err != nil {
		return fmt.Errorf("failed to generate rollback file: %w", err)
	}

	return nil
}

// generateRollback 生成回滚文件
func (g *MigrationGenerator) generateRollback(data map[string]interface{}, cfg *MigrationConfig, migrationDir string) error {
	var templateName string
	switch cfg.Action {
	case "create":
		templateName = "rollback_drop.sql.tmpl"
	case "alter":
		templateName = "rollback_alter.sql.tmpl"
	case "drop":
		templateName = "rollback_create.sql.tmpl"
	default:
		return nil // 不生成回滚文件
	}

	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return err
	}

	// 使用down后缀而不是rollback，这是更标准的命名约定
	filename := fmt.Sprintf("%s_%s.down.sql", data["Timestamp"], ToSnakeCase(cfg.Name))
	path := filepath.Join(migrationDir, filename)

	return g.writeFile(path, content)
}

// extractTableFromName 从迁移名称提取表名
func (g *MigrationGenerator) extractTableFromName(name string) string {
	name = strings.ToLower(name)

	// 移除常见的前缀
	prefixes := []string{"create_", "add_", "drop_", "alter_", "modify_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}

	// 移除常见的后缀
	suffixes := []string{"_table", "_column", "_index", "_constraint"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			name = strings.TrimSuffix(name, suffix)
			break
		}
	}

	return name
}

// writeFile 写入文件
func (g *MigrationGenerator) writeFile(path, content string) error {
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
