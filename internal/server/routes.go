package server

import (
	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/handler"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// SetupRoutes 设置路由
func SetupRoutes(
	r *gin.Engine,
	userService service.UserService,
	articleService service.ArticleService,
	logger logger.Logger,
) {
	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 用户路由
	userHandler := handler.NewUserHandler(userService, logger)
	userHandler.RegisterRoutes(v1)

	// 文章路由
	articleHandler := handler.NewArticleHandler(articleService, logger)
	articleHandler.RegisterRoutes(v1)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
