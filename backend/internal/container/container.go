package container

import (
	"context"
	"fmt"
	"sync"

	"gin-web-framework/config"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/handler"
	"gin-web-framework/internal/redis"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	mu       sync.RWMutex
	services map[string]interface{}
	config   config.ConfigInterface
	db       *gorm.DB
	redis    redis.RedisClient
}

// ContainerInterface 容器接口
type ContainerInterface interface {
	// 配置相关
	GetConfig() config.ConfigInterface

	// 数据库相关
	GetDB() *gorm.DB
	GetRedis() redis.RedisClient

	// 服务层 - 返回接口类型
	GetUserService() service.UserServiceInterface
	GetTodoService() service.TodoServiceInterface
	GetArticleService() service.ArticleServiceInterface
	GetNotificationService() service.NotificationServiceInterface
	GetStatisticsService() service.StatisticsServiceInterface
	GetCategoryService() service.CategoryServiceInterface
	GetCacheService() service.CacheServiceInterface
	GetQueryOptimizer() *service.QueryOptimizer
	GetStatisticsCache() *service.StatisticsCache

	// 处理器层
	GetUserHandler() *handler.UserHandler
	GetTodoHandler() *handler.TodoHandler
	GetArticleHandler() *handler.ArticleHandler
	GetNotificationHandler() *handler.NotificationHandler
	GetStatisticsHandler() *handler.StatisticsHandler
	GetCategoryHandler() *handler.CategoryHandler
	GetSettingsHandler() *handler.SettingsHandler
	GetToolsHandler() *handler.ToolsHandler
	GetUploadHandler() *handler.UploadHandler

	// 容器管理
	Register(name string, service interface{})
	Get(name string) (interface{}, error)
	Has(name string) bool
	Shutdown(ctx context.Context) error
}

// NewContainer 创建新的容器实例
func NewContainer(cfg config.ConfigInterface, db *gorm.DB, redisClient redis.RedisClient) ContainerInterface {
	container := &Container{
		services: make(map[string]interface{}),
		config:   cfg,
		db:       db,
		redis:    redisClient,
	}

	// 在创建时直接初始化所有服务，避免运行时死锁
	container.initializeAllServices()

	return container
}

// initializeAllServices 初始化所有服务，避免运行时死锁
func (c *Container) initializeAllServices() {
	// 获取全局logger
	globalLogger := logger.GetLogger()

	// 创建所有服务实例 - 逐步添加logger
	userService := service.NewUserService(globalLogger)
	todoService := service.NewTodoService(globalLogger)
	articleService := service.NewArticleService(c.db, globalLogger)
	notificationService := service.NewNotificationService(globalLogger)
	statisticsService := service.NewStatisticsService(globalLogger)
	categoryService := service.NewCategoryService(globalLogger)
	cacheService := service.NewCacheService(c.redis, globalLogger)
	settingsService := service.NewSettingsService(c.db, globalLogger)

	// 创建依赖缓存服务的组件
	queryOptimizer := service.NewQueryOptimizer(c.db, cacheService)
	statisticsCache := service.NewStatisticsCache(cacheService)

	// 创建所有处理器实例 - 逐步统一依赖注入模式
	userHandler := handler.NewUserHandler(userService, globalLogger)
	todoHandler := handler.NewTodoHandler(todoService, globalLogger)
	articleHandler := handler.NewArticleHandler(articleService, globalLogger)
	notificationHandler := handler.NewNotificationHandler(notificationService, globalLogger)
	statisticsHandler := handler.NewStatisticsHandler(statisticsService, globalLogger)
	categoryHandler := handler.NewCategoryHandler(categoryService, globalLogger)
	settingsHandler := handler.NewSettingsHandler(settingsService, globalLogger)
	toolsHandler := handler.NewToolsHandler(globalLogger)
	uploadHandler := handler.NewUploadHandler()

	// 注册所有服务
	c.services["user_service"] = userService
	c.services["todo_service"] = todoService
	c.services["article_service"] = articleService
	c.services["notification_service"] = notificationService
	c.services["statistics_service"] = statisticsService
	c.services["category_service"] = categoryService
	c.services["cache_service"] = cacheService
	c.services["settings_service"] = settingsService
	c.services["query_optimizer"] = queryOptimizer
	c.services["statistics_cache"] = statisticsCache

	// 注册所有处理器
	c.services["user_handler"] = userHandler
	c.services["todo_handler"] = todoHandler
	c.services["article_handler"] = articleHandler
	c.services["notification_handler"] = notificationHandler
	c.services["statistics_handler"] = statisticsHandler
	c.services["category_handler"] = categoryHandler
	c.services["settings_handler"] = settingsHandler
	c.services["tools_handler"] = toolsHandler
	c.services["upload_handler"] = uploadHandler

	logger.Info("All services initialized successfully")
}

// GetConfig 获取配置
func (c *Container) GetConfig() config.ConfigInterface {
	return c.config
}

// GetDB 获取数据库连接
func (c *Container) GetDB() *gorm.DB {
	return c.db
}

// GetRedis 获取Redis连接
func (c *Container) GetRedis() redis.RedisClient {
	return c.redis
}

// 服务层实现 - 直接从已初始化的服务中获取
func (c *Container) GetUserService() service.UserServiceInterface {
	return c.services["user_service"].(service.UserServiceInterface)
}

func (c *Container) GetTodoService() service.TodoServiceInterface {
	return c.services["todo_service"].(service.TodoServiceInterface)
}

func (c *Container) GetArticleService() service.ArticleServiceInterface {
	return c.services["article_service"].(service.ArticleServiceInterface)
}

func (c *Container) GetNotificationService() service.NotificationServiceInterface {
	return c.services["notification_service"].(service.NotificationServiceInterface)
}

func (c *Container) GetStatisticsService() service.StatisticsServiceInterface {
	return c.services["statistics_service"].(service.StatisticsServiceInterface)
}

func (c *Container) GetCategoryService() service.CategoryServiceInterface {
	return c.services["category_service"].(service.CategoryServiceInterface)
}

func (c *Container) GetCacheService() service.CacheServiceInterface {
	return c.services["cache_service"].(service.CacheServiceInterface)
}

func (c *Container) GetQueryOptimizer() *service.QueryOptimizer {
	return c.services["query_optimizer"].(*service.QueryOptimizer)
}

func (c *Container) GetStatisticsCache() *service.StatisticsCache {
	return c.services["statistics_cache"].(*service.StatisticsCache)
}

// 处理器层实现 - 直接从已初始化的处理器中获取
func (c *Container) GetUserHandler() *handler.UserHandler {
	return c.services["user_handler"].(*handler.UserHandler)
}

func (c *Container) GetTodoHandler() *handler.TodoHandler {
	return c.services["todo_handler"].(*handler.TodoHandler)
}

func (c *Container) GetArticleHandler() *handler.ArticleHandler {
	return c.services["article_handler"].(*handler.ArticleHandler)
}

func (c *Container) GetNotificationHandler() *handler.NotificationHandler {
	return c.services["notification_handler"].(*handler.NotificationHandler)
}

func (c *Container) GetStatisticsHandler() *handler.StatisticsHandler {
	return c.services["statistics_handler"].(*handler.StatisticsHandler)
}

func (c *Container) GetCategoryHandler() *handler.CategoryHandler {
	return c.services["category_handler"].(*handler.CategoryHandler)
}

func (c *Container) GetSettingsHandler() *handler.SettingsHandler {
	return c.services["settings_handler"].(*handler.SettingsHandler)
}

func (c *Container) GetToolsHandler() *handler.ToolsHandler {
	return c.services["tools_handler"].(*handler.ToolsHandler)
}

// Register 注册服务
func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[name] = service
	logger.Infof("Service '%s' registered successfully", name)
}

// Get 获取服务
func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if service, exists := c.services[name]; exists {
		return service, nil
	}

	return nil, fmt.Errorf("service '%s' not found", name)
}

// Has 检查服务是否存在
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.services[name]
	return exists
}

// Shutdown 关闭容器，清理资源
func (c *Container) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down container...")

	c.mu.Lock()
	defer c.mu.Unlock()

	// 清理服务
	for name, service := range c.services {
		if shutdownable, ok := service.(interface{ Shutdown(context.Context) error }); ok {
			if err := shutdownable.Shutdown(ctx); err != nil {
				logger.Error(fmt.Sprintf("Failed to shutdown service %s: %v", name, err))
			}
		}
	}

	// 清空服务映射
	c.services = make(map[string]interface{})

	logger.Info("Container shutdown completed")
	return nil
}

// 全局容器实例
var globalContainer ContainerInterface

// SetGlobalContainer 设置全局容器
func SetGlobalContainer(container ContainerInterface) {
	globalContainer = container
}

// GetGlobalContainer 获取全局容器
func GetGlobalContainer() ContainerInterface {
	if globalContainer == nil {
		panic("global container not initialized")
	}
	return globalContainer
}

// InitializeContainer 初始化容器
func InitializeContainer() error {
	// 获取配置
	cfg := config.Get()

	// 获取数据库连接
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// 获取Redis连接
	redisClient := redis.GetRedisClient()

	// 创建容器（会自动初始化所有服务）
	container := NewContainer(cfg, db, redisClient)

	// 设置为全局容器
	SetGlobalContainer(container)

	logger.Info("Container initialized successfully with all services pre-loaded")
	return nil
}


// GetUploadHandler 获取上传处理器
func (c *Container) GetUploadHandler() *handler.UploadHandler {
	return c.services["upload_handler"].(*handler.UploadHandler)
}
