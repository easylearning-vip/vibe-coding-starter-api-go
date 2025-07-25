package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
	"vibe-coding-starter/pkg/migration"
)

var (
	configFile string
	cfg        *config.Config
	migrator   *migration.Migrator
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration tool for vibe-coding-starter",
	Long:  `A comprehensive database migration tool for managing database schema changes.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		initMigrator()
	},
}

// upCmd 向上迁移命令
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Long:  `Run all pending database migrations to bring the database up to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	},
}

// downCmd 向下迁移命令
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Long:  `Rollback the last applied migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrator.Down(); err != nil {
			log.Fatalf("Migration rollback failed: %v", err)
		}
	},
}

// stepsCmd 步进迁移命令
var stepsCmd = &cobra.Command{
	Use:   "steps [n]",
	Short: "Run n migrations forward or backward",
	Long:  `Run n migrations forward (positive number) or backward (negative number).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid step number: %v", err)
		}

		if err := migrator.Steps(n); err != nil {
			log.Fatalf("Migration steps failed: %v", err)
		}
	},
}

// forceCmd 强制版本命令
var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force set migration version",
	Long:  `Force set the migration version without running migrations. Use with caution!`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}

		if err := migrator.Force(version); err != nil {
			log.Fatalf("Force migration version failed: %v", err)
		}
	},
}

// versionCmd 版本查询命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Show the current migration version and dirty state.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	},
}

// dropCmd 删除所有表命令
var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all database tables",
	Long:  `Drop all database tables. This is a destructive operation!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Are you sure you want to drop all database tables? (yes/no): ")
		var response string
		fmt.Scanln(&response)

		if response != "yes" {
			fmt.Println("Operation cancelled.")
			return
		}

		if err := migrator.Drop(); err != nil {
			log.Fatalf("Drop database failed: %v", err)
		}
	},
}

// freshCmd 重新创建数据库命令
var freshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run all migrations",
	Long:  `Drop all database tables and re-run all migrations from the beginning.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Are you sure you want to drop all tables and re-run migrations? (yes/no): ")
		var response string
		fmt.Scanln(&response)

		if response != "yes" {
			fmt.Println("Operation cancelled.")
			return
		}

		// 先删除所有表
		if err := migrator.Drop(); err != nil {
			log.Fatalf("Drop database failed: %v", err)
		}

		// 然后运行所有迁移
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("Database refreshed successfully!")
	},
}

func init() {
	// 添加持久化标志
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")

	// 添加子命令
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(stepsCmd)
	rootCmd.AddCommand(forceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(dropCmd)
	rootCmd.AddCommand(freshCmd)
}

// initConfig 初始化配置
func initConfig() {
	var err error
	if configFile != "" {
		cfg, err = config.LoadConfig(configFile)
	} else {
		cfg, err = config.LoadConfig("")
	}

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

// initMigrator 初始化迁移器
func initMigrator() {
	// 初始化日志
	log, err := logger.New(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	// 初始化数据库
	dbInstance, err := database.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	// 创建迁移器
	migrator = migration.NewMigrator(dbInstance.GetDB(), log, &migration.Config{
		MigrationsPath: "", // 将根据数据库驱动自动选择路径
		DatabaseName:   cfg.Database.Database,
		DatabaseDriver: cfg.Database.Driver,
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
