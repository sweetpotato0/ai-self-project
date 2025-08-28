package router

import (
	"time"

	"gin-web-framework/config"
	"gin-web-framework/internal/api"
	"gin-web-framework/internal/container"
	"gin-web-framework/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(container container.ContainerInterface) *gin.Engine {
	cfg := config.Get()

	// 设置Gin模式
	gin.SetMode(cfg.GetServer().Mode)

	r := gin.New()

	// 使用中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.PerformanceMonitor())
	r.Use(middleware.MemoryCleanup())
	r.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB限制

	// 可观测性中间件
	r.Use(middleware.PrometheusMetrics()) // Prometheus指标收集
	r.Use(middleware.Tracing())           // 分布式追踪

	// CORS配置（必须在安全头之前）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 安全中间件
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.InputValidation())
	r.Use(middleware.SQLInjectionProtection())

	// 认证中间件（可选认证）
	r.Use(middleware.OptionalAuthMiddleware())

	// 获取handler实例
	userHandler := container.GetUserHandler()
	todoHandler := container.GetTodoHandler()
	articleHandler := container.GetArticleHandler()
	notificationHandler := container.GetNotificationHandler()
	statisticsHandler := container.GetStatisticsHandler()
	categoryHandler := container.GetCategoryHandler()
	settingsHandler := container.GetSettingsHandler()
	toolsHandler := container.GetToolsHandler()
	uploadHandler := container.GetUploadHandler()

	// API路由组
	apiGroup := r.Group("/api/v1")
	{
		// 健康检查和监控
		apiGroup.GET("/health", settingsHandler.HealthCheck)
		apiGroup.GET("/metrics", settingsHandler.GetMetrics)

		// API文档
		apiGroup.GET("/docs", api.GenerateSwaggerUIHandler())
		apiGroup.GET("/docs/json", api.GenerateDocsHandler())

		// 用户相关路由
		users := apiGroup.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.GET("/profile", middleware.AuthMiddleware(), userHandler.GetProfile)
			users.PUT("/profile", middleware.AuthMiddleware(), userHandler.UpdateProfile)
		}

		// TODO相关路由
		todos := apiGroup.Group("/todos")
		{
			todos.GET("", middleware.AuthMiddleware(), todoHandler.GetTodos)
			todos.POST("", middleware.AuthMiddleware(), todoHandler.CreateTodo)
			todos.GET("/:id", middleware.AuthMiddleware(), todoHandler.GetTodo)
			todos.PUT("/:id", middleware.AuthMiddleware(), todoHandler.UpdateTodo)
			todos.DELETE("/:id", middleware.AuthMiddleware(), todoHandler.DeleteTodo)
		}

		// 通知相关路由
		notifications := apiGroup.Group("/notifications")
		{
			notifications.GET("", middleware.AuthMiddleware(), notificationHandler.GetNotifications)
			notifications.GET("/unread-count", middleware.AuthMiddleware(), notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", middleware.AuthMiddleware(), notificationHandler.MarkAsRead)
			notifications.PUT("/mark-all-read", middleware.AuthMiddleware(), notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", middleware.AuthMiddleware(), notificationHandler.DeleteNotification)
		}

		// WebSocket路由
		apiGroup.GET("/ws", container.GetWebSocketHandler().WebSocket)

		// 文章相关路由
		articles := apiGroup.Group("/articles")
		{
			articles.POST("", middleware.AuthMiddleware(), articleHandler.CreateArticle)
			articles.GET("", middleware.AuthMiddleware(), articleHandler.GetUserArticles)
			articles.GET("/stats", middleware.AuthMiddleware(), articleHandler.GetArticleStats)
			articles.GET("/:id", middleware.AuthMiddleware(), articleHandler.GetArticleByID)
			articles.PUT("/:id", middleware.AuthMiddleware(), articleHandler.UpdateArticle)
			articles.DELETE("/:id", middleware.AuthMiddleware(), articleHandler.DeleteArticle)
		}

		// 上传相关路由
		upload := apiGroup.Group("/upload")
		{
			upload.POST("/image", middleware.AuthMiddleware(), uploadHandler.UploadImage)
		}

		// 统计相关路由
		statistics := apiGroup.Group("/statistics")
		{
			statistics.GET("", middleware.AuthMiddleware(), statisticsHandler.GetStatistics)
			statistics.GET("/trends", middleware.AuthMiddleware(), statisticsHandler.GetTrends)
		}

		// 分类相关路由
		categories := apiGroup.Group("/categories")
		{
			categories.POST("", middleware.AuthMiddleware(), categoryHandler.CreateCategory)
			categories.GET("/tree", middleware.AuthMiddleware(), categoryHandler.GetCategoryTree)
			categories.GET("", middleware.AuthMiddleware(), categoryHandler.GetAllCategories)
			categories.GET("/stats", middleware.AuthMiddleware(), categoryHandler.GetCategoryStats)
			categories.GET("/:id", middleware.AuthMiddleware(), categoryHandler.GetCategoryByID)
			categories.PUT("/:id", middleware.AuthMiddleware(), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.AuthMiddleware(), categoryHandler.DeleteCategory)
		}

		// 设置相关路由
		settings := apiGroup.Group("/settings")
		{
			settings.GET("", middleware.AuthMiddleware(), settingsHandler.GetSettings)
			settings.PUT("/profile", middleware.AuthMiddleware(), settingsHandler.UpdateProfile)
			settings.PUT("/password", middleware.AuthMiddleware(), settingsHandler.ChangePassword)
			settings.PUT("/notifications", middleware.AuthMiddleware(), settingsHandler.UpdateNotificationSettings)
			settings.PUT("/interface", middleware.AuthMiddleware(), settingsHandler.UpdateInterfaceSettings)
			settings.GET("/export", middleware.AuthMiddleware(), settingsHandler.ExportData)
			settings.DELETE("/completed-tasks", middleware.AuthMiddleware(), settingsHandler.ClearCompletedTasks)
		}

		// 工具相关路由
		tools := apiGroup.Group("/tools")
		{
			// 网络工具
			tools.POST("/network/port-scan", middleware.AuthMiddleware(), toolsHandler.PortScan)
		}
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	return r
}
