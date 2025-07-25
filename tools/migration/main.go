package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	migrationsDir = "migrations"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Migration file generator for vibe-coding-starter",
	Long:  `A tool to generate migration files with proper naming and structure.`,
}

// createCmd 创建迁移文件命令
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file with the specified name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := createMigration(name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

// createMigration 创建迁移文件
func createMigration(name string) error {
	// 获取下一个版本号
	version, err := getNextVersion()
	if err != nil {
		return fmt.Errorf("failed to get next version: %w", err)
	}

	// 格式化名称
	formattedName := formatName(name)
	
	// 生成文件名
	upFile := fmt.Sprintf("%06d_%s.up.sql", version, formattedName)
	downFile := fmt.Sprintf("%06d_%s.down.sql", version, formattedName)

	// 创建迁移目录（如果不存在）
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// 创建 up 文件
	upPath := filepath.Join(migrationsDir, upFile)
	upContent := generateUpTemplate(name)
	if err := ioutil.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up file: %w", err)
	}

	// 创建 down 文件
	downPath := filepath.Join(migrationsDir, downFile)
	downContent := generateDownTemplate(name)
	if err := ioutil.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down file: %w", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upPath)
	fmt.Printf("  %s\n", downPath)

	return nil
}

// getNextVersion 获取下一个版本号
func getNextVersion() (int, error) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil // 如果目录不存在，从版本1开始
		}
		return 0, err
	}

	maxVersion := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		// 提取版本号
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		if version > maxVersion {
			maxVersion = version
		}
	}

	return maxVersion + 1, nil
}

// formatName 格式化迁移名称
func formatName(name string) string {
	// 转换为小写
	name = strings.ToLower(name)
	
	// 替换空格和特殊字符为下划线
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	
	// 移除多余的下划线
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	
	// 移除开头和结尾的下划线
	name = strings.Trim(name, "_")
	
	return name
}

// generateUpTemplate 生成 up 迁移模板
func generateUpTemplate(name string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	template := fmt.Sprintf(`-- Migration: %s
-- Created: %s
-- Description: Add your migration description here

-- Add your SQL statements here
-- Example:
-- CREATE TABLE example (
--     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`, name, timestamp)

	return template
}

// generateDownTemplate 生成 down 迁移模板
func generateDownTemplate(name string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	template := fmt.Sprintf(`-- Rollback: %s
-- Created: %s
-- Description: Rollback migration

-- Add your rollback SQL statements here
-- Example:
-- DROP TABLE IF EXISTS example;
`, name, timestamp)

	return template
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
