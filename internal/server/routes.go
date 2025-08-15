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
	fileService service.FileService,
	dictService service.DictService,
	productService service.ProductService,
	productCategoryService service.ProductCategoryService,
	departmentService service.DepartmentService,
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

	// 文件路由
	fileHandler := handler.NewFileHandler(fileService, logger)
	fileHandler.RegisterRoutes(v1)

	// 数据字典路由
	dictHandler := handler.NewDictHandler(dictService, logger)
	dictHandler.RegisterRoutes(v1)

	// 产品分类路由
	productCategoryHandler := handler.NewProductCategoryHandler(productCategoryService, logger)
	productCategoryHandler.RegisterRoutes(v1)

	// 产品路由
	productHandler := handler.NewProductHandler(productService, logger)
	productHandler.RegisterRoutes(v1)

	// 部门路由
	departmentHandler := handler.NewDepartmentHandler(departmentService, logger)
	departmentHandler.RegisterRoutes(v1)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
