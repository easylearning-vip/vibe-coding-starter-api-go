package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/fx"

	"vibe-coding-starter/internal/config"
	mcpHandler "vibe-coding-starter/internal/handler/mcp"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/database"
	appLogger "vibe-coding-starter/pkg/logger"
)

var (
	hostFlag *string
	portFlag *int
)

// configProvider 创建配置提供者，支持命令行参数和环境变量
func configProvider() (*config.Config, error) {
	// 解析命令行参数
	configFile := flag.String("c", "", "Configuration file path")
	hostFlag = flag.String("host", "0.0.0.0", "HTTP host to listen on")
	portFlag = flag.Int("port", 6275, "HTTP port to listen on")
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

// runHTTPServer 启动HTTP服务器，使用Fx生命周期钩子
func runHTTPServer(lc fx.Lifecycle, logr appLogger.Logger, handler http.Handler) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *hostFlag, *portFlag),
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logr.Info("Starting TOEIC MCP server (streamable HTTP)",
					"addr", srv.Addr,
					"auth_header", "Authorization: Token <token> or X-User-Token: <token>",
				)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("MCP server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		// 配置与基础设施
		fx.Provide(
			configProvider,
			appLogger.New,
			database.New,
		),

		// 仓储层
		fx.Provide(
			repository.NewUserRepository,
			repository.NewPart2PracticeSetRepository,
			repository.NewPart2PracticeSetItemRepository,
			repository.NewPart2QuestionRepository,
			repository.NewPart3PracticeSetRepository,
			repository.NewPart3PracticeSetItemRepository,
			repository.NewPart3ConversationRepository,
			repository.NewPart4PracticeSetRepository,
			repository.NewPart4PracticeSetItemRepository,
			repository.NewPart4TalkRepository,
			repository.NewToeicAiPromptRepository,
		),

		// 服务层
		fx.Provide(
			service.NewPart2PracticeSetService,
			service.NewPart3PracticeSetService,
			service.NewPart4PracticeSetService,
			service.NewToeicAiPromptService,
		),

		// MCP处理器层
		fx.Provide(
			mcpHandler.NewMCPHandler,
		),

		// 构建HTTP Handler
		fx.Provide(
			func(h *mcpHandler.MCPHandler) http.Handler {
				return h.BuildHTTPHandler()
			},
		),

		// 启动服务器
		fx.Invoke(runHTTPServer),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start MCP application:", err)
	}

	<-app.Done()
	log.Println("MCP server stopped")
}
