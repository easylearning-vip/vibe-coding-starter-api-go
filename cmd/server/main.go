package main

import (
	"context"
	"flag"
	"log"
	"os"

	"go.uber.org/fx"

	_ "vibe-coding-starter/docs" // 导入Swagger文档
	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/middleware"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/server"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// @title Vibe Coding Starter API
// @version 1.0
// @description A Go web application starter with AI-assisted development features
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// configProvider 创建配置提供者，支持命令行参数和环境变量
func configProvider() (*config.Config, error) {
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

	// 使用 LoadConfig 函数加载配置
	return config.LoadConfig(configPath)
}

func main() {
	app := fx.New(
		// 配置模块
		fx.Provide(configProvider),

		// 基础设施模块
		fx.Provide(
			logger.New,
			database.New,
			cache.New,
		),

		// 中间件模块
		fx.Provide(
			middleware.NewMiddleware,
		),

		// 仓储模块
		fx.Provide(
			repository.NewUserRepository,
			repository.NewArticleRepository,
			repository.NewFileRepository,
			repository.NewDictRepository,
		),

		// 服务模块
		fx.Provide(
			service.NewUserService,
			service.NewArticleService,
			service.NewFileService,
			service.NewDictService,
		),

		// 处理器模块
		fx.Provide(
			handler.NewUserHandler,
			handler.NewArticleHandler,
			handler.NewFileHandler,
			handler.NewHealthHandler,
			handler.NewDictHandler,
		),

		// 服务器模块
		fx.Provide(server.New),

		// 启动服务器
		fx.Invoke(func(srv *server.Server) {
			// 服务器启动在 OnStart hook 中处理
		}),

		// 优雅关闭
		fx.Invoke(func(lifecycle fx.Lifecycle, srv *server.Server) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := srv.Start(); err != nil {
							log.Printf("Server start error: %v", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Shutdown(ctx)
				},
			})
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	// 等待应用程序结束
	<-app.Done()
	log.Println("Application stopped")
}
