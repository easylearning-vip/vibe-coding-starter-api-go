package cmd

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

	// 如果没有指定表名，从迁移名称推断
	table := cfg.Table
	if table == "" {
		table = g.extractTableFromName(cfg.Name)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Name":      cfg.Name,
		"Table":     table,
		"TableName": ToSnakeCase(table),
		"Action":    cfg.Action,
		"Timestamp": GenerateTimestamp(),
		"Year":      GetCurrentYear(),
	}

	// 根据操作类型选择模板
	var templateName string
	switch cfg.Action {
	case "create":
		templateName = "migration_create.sql.tmpl"
	case "alter":
		templateName = "migration_alter.sql.tmpl"
	case "drop":
		templateName = "migration_drop.sql.tmpl"
	default:
		templateName = "migration_create.sql.tmpl"
	}

	// 生成迁移文件
	content, err := g.templateEngine.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render migration template: %w", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("%s_%s.sql", data["Timestamp"], ToSnakeCase(cfg.Name))
	path := filepath.Join("migrations", filename)

	if err := g.writeFile(path, content); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	// 生成回滚文件
	if err := g.generateRollback(data, cfg); err != nil {
		return fmt.Errorf("failed to generate rollback file: %w", err)
	}

	return nil
}

// generateRollback 生成回滚文件
func (g *MigrationGenerator) generateRollback(data map[string]interface{}, cfg *MigrationConfig) error {
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

	filename := fmt.Sprintf("%s_%s_rollback.sql", data["Timestamp"], ToSnakeCase(cfg.Name))
	path := filepath.Join("migrations", "rollbacks", filename)

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
