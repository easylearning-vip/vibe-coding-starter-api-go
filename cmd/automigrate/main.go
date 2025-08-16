package main

import (
	"flag"
	"log"
	"os"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

func main() {
	// 解析命令行参数
	configFile := flag.String("c", "", "Configuration file path")
	flag.Parse()

	// 如果命令行参数为空，检查环境变量 CONFIG_FILE
	configPath := *configFile
	if configPath == "" {
		if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
			configPath = envConfigFile
		}
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建日志器
	appLogger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// 创建数据库连接
	db, err := database.New(cfg, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 执行自动迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.File{},
		&model.DictCategory{},
		&model.DictItem{},
		&model.Department{},
		&model.ProductCategory{},
		&model.Product{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("Auto migration completed successfully!")
}
