package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/middleware"
	"vibe-coding-starter/pkg/logger"
)

// Server HTTP 服务器
type Server struct {
	config         *config.Config
	logger         logger.Logger
	httpServer     *http.Server
	middleware     *middleware.Middleware
	userHandler    *handler.UserHandler
	articleHandler *handler.ArticleHandler
	healthHandler  *handler.HealthHandler
	dictHandler    *handler.DictHandler
	productcategoryHandler *handler.ProductCategoryHandler
	departmentHandler *handler.DepartmentHandler
}

// New 创建新的服务器实例
func New(
	config *config.Config,
	logger logger.Logger,
	middleware *middleware.Middleware,
	userHandler *handler.UserHandler,
	articleHandler *handler.ArticleHandler,
	healthHandler *handler.HealthHandler,
	dictHandler *handler.DictHandler,
	productcategoryHandler *handler.ProductCategoryHandler,
	departmentHandler *handler.DepartmentHandler,
) *Server {
	return &Server{
		config:         config,
		logger:         logger,
		middleware:     middleware,
		userHandler:    userHandler,
		articleHandler: articleHandler,
		healthHandler:  healthHandler,
		dictHandler:    dictHandler,
		productcategoryHandler: productcategoryHandler,
		departmentHandler: departmentHandler,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 设置 Gin 模式
	gin.SetMode(s.config.Server.Mode)

	// 创建 Gin 引擎
	engine := gin.New()

	// 设置中间件
	s.setupMiddleware(engine)

	// 设置路由
	s.setupRoutes(engine)

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:         s.config.Server.GetAddress(),
		Handler:      engine,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}

	s.logger.Info("Starting HTTP server",
		"address", s.config.Server.GetAddress(),
		"mode", s.config.Server.Mode,
	)

	// 启动服务器
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware(engine *gin.Engine) {
	// 使用新的中间件系统
	s.middleware.SetupGlobalMiddleware(engine)
}

// setupRoutes 设置路由
func (s *Server) setupRoutes(engine *gin.Engine) {
	// 健康检查路由（使用专用中间件）
	health := engine.Group("/")
	health.Use(s.middleware.HealthCheckAPI()...)

	// 健康检查路由直接注册到引擎上
	s.healthHandler.RegisterRoutes(engine)

	// API 路由组
	api := engine.Group("/api")
	api.Use(s.middleware.SetupAPIMiddleware()...)
	{
		v1 := api.Group("/v1")
		{
			// 公共路由（不需要认证）
			public := v1.Group("")
			public.Use(s.middleware.PublicAPI()...)
			{
				// 用户注册和登录路由
				users := public.Group("/users")
				users.POST("/register", s.userHandler.Register)
				users.POST("/login", s.userHandler.Login)

				// 文章公共路由（查看文章列表和详情）
				articles := public.Group("/articles")
				{
					articles.GET("", s.articleHandler.List)
					articles.GET("/search", s.articleHandler.Search)
					articles.GET("/:id", s.articleHandler.GetByID)
				}

				// 数据字典路由（不需要认证，便于测试）
				s.dictHandler.RegisterRoutes(public)
			}

			// 受保护的路由（需要认证）
			protected := v1.Group("")
			protected.Use(s.middleware.ProtectedAPI()...)
			{
				// 用户路由
				s.userHandler.RegisterRoutes(protected)

				// 用户文章管理路由（只能操作自己的文章）
				userArticles := protected.Group("/user/articles")
				{
					userArticles.GET("", s.articleHandler.ListUserArticles) // 用户专用：只返回当前用户的文章
					userArticles.POST("", s.articleHandler.Create)
					userArticles.PUT("/:id", s.articleHandler.Update)
					userArticles.DELETE("/:id", s.articleHandler.Delete)
				}
			}

			// 管理员路由
			admin := v1.Group("/admin")
			admin.Use(s.middleware.AdminAPI()...)
			{
				// 管理员专用路由
				s.userHandler.RegisterRoutes(admin)

				// 管理员文章管理路由（可以操作所有文章）
				adminArticles := admin.Group("/articles")
				{
					adminArticles.GET("", s.articleHandler.ListAllArticles) // 管理员专用：返回所有文章
					adminArticles.POST("", s.articleHandler.Create)
					adminArticles.PUT("/:id", s.articleHandler.Update)
					adminArticles.DELETE("/:id", s.articleHandler.Delete)
				}

				// ProductCategory管理路由
				s.productcategoryHandler.RegisterRoutes(admin)

				// Department管理路由
				s.departmentHandler.RegisterRoutes(admin)

			}
		}
	}

	// Swagger 文档路由 (开发环境)
	if s.config.Server.Mode == "debug" {
		s.setupSwaggerRoutes(engine)
	}
}

// setupSwaggerRoutes 设置 Swagger 文档路由
func (s *Server) setupSwaggerRoutes(engine *gin.Engine) {
	// Swagger 文档路由
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 重定向根路径到 Swagger 文档
	engine.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
