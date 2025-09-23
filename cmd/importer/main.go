package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/fx"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// ImporterApp 导入器应用结构
type ImporterApp struct {
	importerService service.ToeicImporterService
	logger          logger.Logger
}

// NewImporterApp 创建导入器应用
func NewImporterApp(importerService service.ToeicImporterService, logger logger.Logger) *ImporterApp {
	return &ImporterApp{
		importerService: importerService,
		logger:          logger,
	}
}

// Run 运行导入器
func (app *ImporterApp) Run(filename string, cleanup bool) error {
	ctx := context.Background()

	// 验证文件
	if err := app.importerService.ValidateFile(filename); err != nil {
		return fmt.Errorf("file validation failed: %w", err)
	}

	// 提取测试名称
	testName := extractTestName(filename)

	// 清理现有数据（如果需要）
	if cleanup {
		app.logger.Info("Cleaning up existing test data", "testName", testName)
		if err := app.importerService.CleanupTestData(ctx, testName); err != nil {
			app.logger.Warn("Failed to cleanup test data", "testName", testName, "error", err)
		}
	}

	// 导入文件
	startTime := time.Now()
	result, err := app.importerService.ImportFile(ctx, filename)
	duration := time.Since(startTime)

	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	// 打印结果
	fmt.Printf("\n=== Import Results ===\n")
	fmt.Printf("File: %s\n", filepath.Base(filename))
	fmt.Printf("Test Name: %s\n", result.TestName)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("\nImported:\n")
	fmt.Printf("  Part 2 Questions: %d\n", result.Part2Questions)
	fmt.Printf("  Part 3 Conversations: %d\n", result.Part3Conversations)
	fmt.Printf("  Part 4 Talks: %d\n", result.Part4Talks)

	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, error := range result.Errors {
			fmt.Printf("  - %s\n", error)
		}
		return fmt.Errorf("import completed with %d errors", len(result.Errors))
	}

	fmt.Printf("\n✅ Import completed successfully!\n")
	return nil
}

// configProvider 创建配置提供者
func configProvider() (*config.Config, error) {
	// 使用默认配置文件或环境变量
	configPath := ""
	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configPath = envConfigFile
	}

	// 使用 LoadConfig 函数加载配置
	return config.LoadConfig(configPath)
}

// extractTestName 从文件名提取测试名称
func extractTestName(filename string) string {
	parts := filepath.Base(filename)
	if filepath.Ext(parts) == ".md" {
		return parts[:len(parts)-3]
	}
	return parts
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("Usage: toeic-importer [options] -file <path-to-test-file>")
	fmt.Println("\nOptions:")
	fmt.Println("  -file string    Path to the TOEIC test file to import (required)")
	fmt.Println("  -cleanup        Clean up existing test data before importing (default: false)")
	fmt.Println("  -c string       Configuration file path")
	fmt.Println("\nExample:")
	fmt.Println("  toeic-importer -file /home/ubuntu/workspace/vocabulary/toeic/test01.md")
	fmt.Println("  toeic-importer -file /home/ubuntu/workspace/vocabulary/toeic/test02.md -cleanup")
	fmt.Println("\nSupported files:")
	fmt.Println("  - Any .md (Markdown) file containing TOEIC test content")
}

func main() {
	// 解析命令行参数
	var filename string
	var cleanup bool
	flag.StringVar(&filename, "file", "", "Path to the TOEIC test file to import")
	flag.BoolVar(&cleanup, "cleanup", false, "Clean up existing test data before importing")
	flag.Parse()

	if filename == "" {
		printUsage()
		os.Exit(1)
	}

	// 验证文件存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", filename)
	}

	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("TOEIC Question Bank Importer starting...")
	log.Printf("Target file: %s", filename)

	// 创建FX应用
	app := fx.New(
		// 配置模块
		fx.Provide(configProvider),

		// 基础设施模块
		fx.Provide(
			logger.New,
			database.New,
		),

		// 仓储模块
		fx.Provide(
			repository.NewTestRepository,
			repository.NewScenarioRepository,
			repository.NewDifficultyLevelRepository,
			repository.NewPart2QuestionRepository,
			repository.NewPart3ConversationRepository,
			repository.NewPart3AnswerOptionRepository,
			repository.NewPart4TalkRepository,
			repository.NewPart4AnswerOptionRepository,
		),

		// 服务模块
		fx.Provide(
			service.NewToeicImporterService,
		),

		// 应用模块
		fx.Provide(NewImporterApp),

		// 运行导入器
		fx.Invoke(func(importerApp *ImporterApp) {
			if err := importerApp.Run(filename, cleanup); err != nil {
				log.Fatalf("Import failed: %v", err)
			}
		}),

		// 禁用FX日志
		fx.NopLogger,
	)

	// 启动应用
	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	// 停止应用
	if err := app.Stop(context.Background()); err != nil {
		log.Fatal("Failed to stop application:", err)
	}

	log.Println("Import completed")
}
