package router

import (
	"time"

	"gin-web-framework/config"
	"gin-web-framework/internal/api"
	"gin-web-framework/internal/container"
	"gin-web-framework/internal/middleware"
	"gin-web-framework/pkg/logger"

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
	
	// 审计日志中间件
	auditMiddleware := middleware.NewAuditMiddleware(container.GetDB(), logger.GetLogger().(*logger.Logger))
	r.Use(auditMiddleware.AuditLog())

	// 获取handler实例
	userHandler := container.GetUserHandler()
	todoHandler := container.GetTodoHandler()
	articleHandler := container.GetArticleHandler()
	notificationHandler := container.GetNotificationHandler()
	statisticsHandler := container.GetStatisticsHandler()
	categoryHandler := container.GetCategoryHandler()
	settingsHandler := container.GetSettingsHandler()
	networkHandler := container.GetNetworkHandler()
	auditHandler := container.GetAuditHandler()
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

		// 认证相关路由
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/refresh", userHandler.RefreshToken)
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
			articles.POST("/:id/like", middleware.AuthMiddleware(), articleHandler.LikeArticle)
			articles.DELETE("/:id/like", middleware.AuthMiddleware(), articleHandler.UnlikeArticle)
			articles.POST("/:id/view", articleHandler.IncrementViewCount) // 浏览量不需要认证
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
			tools.POST("/network/port-scan", middleware.AuthMiddleware(), networkHandler.PortScan)
			tools.POST("/network/dns-lookup", middleware.AuthMiddleware(), networkHandler.DNSLookup)
		}

		// 审计日志相关路由
		auditLogs := apiGroup.Group("/audit-logs")
		{
			auditLogs.GET("", middleware.AuthMiddleware(), auditHandler.GetAuditLogs)
			auditLogs.GET("/:id", middleware.AuthMiddleware(), auditHandler.GetAuditLogByID)
			auditLogs.GET("/stats", middleware.AuthMiddleware(), auditHandler.GetAuditLogStats)
			auditLogs.DELETE("", middleware.AuthMiddleware(), auditHandler.DeleteAuditLogs)
		}

		// 英文学习相关路由
		englishLearningHandler := container.GetEnglishLearningHandler()
		englishVideoHandler := container.GetEnglishVideoHandler()
		learning := apiGroup.Group("/learning")
		{
			// 学习分类管理
			categories := learning.Group("/categories")
			{
				categories.GET("", englishLearningHandler.GetCategories)
				categories.POST("", middleware.AuthMiddleware(), englishLearningHandler.CreateCategory)
				categories.PUT("/:id", middleware.AuthMiddleware(), englishLearningHandler.UpdateCategory)
				categories.DELETE("/:id", middleware.AuthMiddleware(), englishLearningHandler.DeleteCategory)
			}

			// 歌曲/学习材料管理
			songs := learning.Group("/songs")
			{
				songs.GET("", englishLearningHandler.GetSongs)
				songs.GET("/:id", englishLearningHandler.GetSongByID)
				songs.POST("", middleware.AuthMiddleware(), englishLearningHandler.CreateSong)
				songs.PUT("/:id", middleware.AuthMiddleware(), englishLearningHandler.UpdateSong)
				songs.DELETE("/:id", middleware.AuthMiddleware(), englishLearningHandler.DeleteSong)
				songs.POST("/:id/like", middleware.AuthMiddleware(), englishLearningHandler.LikeSong)
				songs.DELETE("/:id/like", middleware.AuthMiddleware(), englishLearningHandler.UnlikeSong)
				songs.PUT("/:id/progress", middleware.AuthMiddleware(), englishLearningHandler.UpdateProgress)
			}

			// 用户学习相关
			user := learning.Group("/user")
			{
				user.GET("/progress", middleware.AuthMiddleware(), englishLearningHandler.GetProgress)
				user.GET("/recommendations", middleware.AuthMiddleware(), englishLearningHandler.GetRecommendations)
				user.GET("/stats", middleware.AuthMiddleware(), englishLearningHandler.GetStats)
			}
		}

		// 英文视频相关路由
		englishVideos := apiGroup.Group("/english-videos")
		{
			// 根级别的统计和进度路由
			englishVideos.GET("/stats", englishVideoHandler.GetVideoStats)
			englishVideos.GET("/progress", middleware.AuthMiddleware(), englishVideoHandler.GetUserProgress)

			// 视频系列
			series := englishVideos.Group("/series")
			{
				series.GET("", englishVideoHandler.GetVideoSeries)
				series.GET("/search", englishVideoHandler.SearchVideoSeries)
				series.GET("/recommended", englishVideoHandler.GetRecommendedSeries)
				series.GET("/:seriesId", englishVideoHandler.GetVideoSeriesDetail)
				series.POST("/:seriesId/toggle-like", middleware.AuthMiddleware(), englishVideoHandler.ToggleSeriesLike)
				series.GET("/:seriesId/episodes", englishVideoHandler.GetEpisodes)
				
				// 管理员功能
				admin := series.Group("/admin")
				{
					admin.POST("", middleware.AuthMiddleware(), englishVideoHandler.CreateVideoSeries)
					admin.PUT("/:seriesId", middleware.AuthMiddleware(), englishVideoHandler.UpdateVideoSeries)
					admin.DELETE("/:seriesId", middleware.AuthMiddleware(), englishVideoHandler.DeleteVideoSeries)
					admin.POST("/:seriesId/episodes", middleware.AuthMiddleware(), englishVideoHandler.CreateEpisode)
					admin.POST("/:seriesId/episodes/batch", middleware.AuthMiddleware(), englishVideoHandler.BatchImportEpisodes)
				}
			}

			// 剧集
			episodes := englishVideos.Group("/episodes")
			{
				episodes.GET("/:episodeId", englishVideoHandler.GetEpisodeDetail)
				episodes.GET("/:episodeId/progress", middleware.AuthMiddleware(), englishVideoHandler.GetEpisodeProgress)
				episodes.POST("/:episodeId/progress", middleware.AuthMiddleware(), englishVideoHandler.UpdateEpisodeProgress)
				
				// 管理员功能
				admin := episodes.Group("/admin")
				{
					admin.PUT("/:episodeId", middleware.AuthMiddleware(), englishVideoHandler.UpdateEpisode)
					admin.DELETE("/:episodeId", middleware.AuthMiddleware(), englishVideoHandler.DeleteEpisode)
					admin.POST("/uncategorized", middleware.AuthMiddleware(), englishVideoHandler.CreateUncategorizedEpisode)
				}
			}

			// 用户统计
			user := englishVideos.Group("/user")
			{
				user.GET("/stats", middleware.AuthMiddleware(), englishVideoHandler.GetUserVideoStats)
			}
		}
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	return r
}
