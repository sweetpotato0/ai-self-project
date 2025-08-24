package service

import (
	"context"
	"gin-web-framework/internal/models"
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	// 用户认证
	Register(req RegisterRequest) (*models.User, error)
	Login(req LoginRequest) (*LoginResponse, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateProfile(userID uint, req UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID uint, req ChangePasswordRequest) error

	// 用户管理
	ListUsers(page, limit int) (*PaginatedUsers, error)
	DeleteUser(id uint) error
	UpdateUserStatus(id uint, status string) error

	// 验证和工具
	ValidateUser(user *models.User) error
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) bool
}

// TodoServiceInterface 任务服务接口
type TodoServiceInterface interface {
	// 基本CRUD
	CreateTodo(req CreateTodoRequest, userID uint) (*models.Todo, error)
	GetTodoByID(id uint, userID uint) (*models.Todo, error)
	GetTodos(userID uint, filter TodoFilter) (*PaginatedTodos, error)
	UpdateTodo(id uint, userID uint, req UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(id uint, userID uint) error

	// 状态管理
	MarkCompleted(id uint, userID uint) error
	MarkInProgress(id uint, userID uint) error
	MarkCancelled(id uint, userID uint) error

	// 批量操作
	BatchUpdateStatus(ids []uint, userID uint, status string) error
	BatchDelete(ids []uint, userID uint) error

	// 统计和查询
	GetOverdueTodos(userID uint) ([]*models.Todo, error)
	GetTodosByCategory(userID uint, categoryID uint) ([]*models.Todo, error)
	GetTodosByPriority(userID uint, priority string) ([]*models.Todo, error)
	GetTodoStats(userID uint) (map[string]interface{}, error)

	// 搜索
	SearchTodos(userID uint, query string, filter TodoFilter) (*PaginatedTodos, error)
}

// ArticleServiceInterface 文章服务接口
type ArticleServiceInterface interface {
	// 基本CRUD
	CreateArticle(req CreateArticleRequest, userID uint) (*models.Article, error)
	GetArticleByID(id uint, userID uint) (*models.Article, error)
	GetArticles(userID uint, filter ArticleFilter) (*PaginatedArticles, error)
	UpdateArticle(id uint, userID uint, req UpdateArticleRequest) (*models.Article, error)
	DeleteArticle(id uint, userID uint) error

	// 状态管理
	PublishArticle(id uint, userID uint) error
	ArchiveArticle(id uint, userID uint) error
	RestoreArticle(id uint, userID uint) error

	// 内容管理
	UpdateContent(id uint, userID uint, content string) error
	IncrementViewCount(id uint) error
	LikeArticle(id uint, userID uint) error
	UnlikeArticle(id uint, userID uint) error

	// 统计和查询
	GetArticleStats(userID uint) (map[string]interface{}, error)
	GetPopularArticles(userID uint, limit int) ([]*models.Article, error)
	GetRecentArticles(userID uint, limit int) ([]*models.Article, error)

	// 搜索
	SearchArticles(userID uint, query string, filter ArticleFilter) (*PaginatedArticles, error)
}

// NotificationServiceInterface 通知服务接口
type NotificationServiceInterface interface {
	// 基本CRUD
	CreateNotification(req CreateNotificationRequest) (*models.Notification, error)
	GetNotifications(userID uint, filter NotificationFilter) (*PaginatedNotifications, error)
	GetUnreadNotifications(userID uint) ([]*models.Notification, error)
	MarkAsRead(id uint, userID uint) error
	MarkAllAsRead(userID uint) error
	DeleteNotification(id uint, userID uint) error

	// 统计
	GetUnreadCount(userID uint) (int64, error)

	// 系统通知
	CreateSystemNotification(title, content string, userIDs []uint) error
	CreateBroadcastNotification(title, content string) error

	// 任务相关通知
	CreateTaskReminderNotification(userID uint, todoID uint) error
	CreateTaskOverdueNotification(userID uint, todoID uint) error
	CreateTaskCompletedNotification(userID uint, todoID uint) error

	// 批量操作
	BatchMarkAsRead(ids []uint, userID uint) error
	BatchDelete(ids []uint, userID uint) error
}

// StatisticsServiceInterface 统计服务接口
type StatisticsServiceInterface interface {
	// 基础统计
	GetStatistics(statType StatisticsType, userID uint) (interface{}, error)
	GetTrends(statType StatisticsType, userID uint, days int) ([]TrendData, error)

	// 任务统计
	GetTodoStatistics(userID uint) (*TodoStatistics, error)
	GetTodoTrends(userID uint, days int) ([]TrendData, error)
	GetTodosByStatus(userID uint) (map[string]int64, error)
	GetTodosByPriority(userID uint) (map[string]int64, error)
	GetTodosByCategory(userID uint) (map[string]int64, error)

	// 文章统计
	GetArticleStatistics(userID uint) (*ArticleStatistics, error)
	GetArticleTrends(userID uint, days int) ([]TrendData, error)
	GetArticlesByStatus(userID uint) (map[string]int64, error)
	GetTopArticles(userID uint, limit int) ([]*models.Article, error)

	// 用户活跃度统计
	GetUserActivityStats(userID uint, days int) (map[string]interface{}, error)
	GetActiveUsersCount(days int) (int64, error)

	// 系统统计
	GetSystemStats() (map[string]interface{}, error)
	GetDailyActiveUsers(days int) ([]TrendData, error)
}

// CategoryServiceInterface 分类服务接口
type CategoryServiceInterface interface {
	// 基本CRUD
	CreateCategory(req CreateCategoryRequest) (*models.Category, error)
	GetCategoryByID(id uint, userID uint) (*models.Category, error)
	GetCategories(userID uint) ([]*models.Category, error)
	GetCategoryTree(userID uint) ([]*models.Category, error)
	UpdateCategory(id uint, userID uint, req UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(id uint, userID uint) error

	// 层级管理
	GetChildCategories(parentID uint, userID uint) ([]*models.Category, error)
	GetCategoryPath(id uint, userID uint) ([]*models.Category, error)
	MoveCategory(id uint, newParentID *uint, userID uint) error

	// 统计
	GetCategoryStats(userID uint) (map[uint]interface{}, error)
	GetCategoriesWithTodoCount(userID uint) ([]*models.Category, error)

	// 验证
	ValidateCategoryHierarchy(id uint, parentID *uint, userID uint) error
	IsCategoryOwner(id uint, userID uint) bool
}

// CacheServiceInterface 缓存服务接口
type CacheServiceInterface interface {
	// 基本操作
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Get(ctx context.Context, key string) (string, error)
	GetObject(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 批量操作
	MSet(ctx context.Context, pairs map[string]interface{}, expiration int) error
	MGet(ctx context.Context, keys []string) (map[string]string, error)
	MDelete(ctx context.Context, keys []string) error

	// 列表操作
	ListPush(ctx context.Context, key string, values ...interface{}) error
	ListPop(ctx context.Context, key string) (string, error)
	ListRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ListLength(ctx context.Context, key string) (int64, error)

	// 集合操作
	SetAdd(ctx context.Context, key string, members ...interface{}) error
	SetMembers(ctx context.Context, key string) ([]string, error)
	SetIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SetRemove(ctx context.Context, key string, members ...interface{}) error

	// 过期和TTL
	SetExpiration(ctx context.Context, key string, expiration int) error
	GetTTL(ctx context.Context, key string) (int, error)

	// 模式匹配
	Keys(ctx context.Context, pattern string) ([]string, error)
	DeletePattern(ctx context.Context, pattern string) error

	// 原子操作
	Increment(ctx context.Context, key string) (int64, error)
	IncrementBy(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	DecrementBy(ctx context.Context, key string, value int64) (int64, error)

	// 锁
	Lock(ctx context.Context, key string, expiration int) (bool, error)
	Unlock(ctx context.Context, key string) error

	// 健康检查
	Ping(ctx context.Context) error
	FlushDB(ctx context.Context) error

	// 统计
	Info(ctx context.Context) (map[string]string, error)
	DatabaseSize(ctx context.Context) (int64, error)
}

// SchedulerServiceInterface 调度服务接口
type SchedulerServiceInterface interface {
	// 调度器管理
	Start() error
	Stop() error
	IsRunning() bool

	// 任务管理
	AddJob(name string, spec string, job func()) error
	RemoveJob(name string) error
	UpdateJob(name string, spec string, job func()) error
	GetJobStatus(name string) (interface{}, error)
	ListJobs() ([]interface{}, error)

	// 任务执行
	RunJobNow(name string) error
	GetJobHistory(name string, limit int) ([]interface{}, error)

	// 系统任务
	ScheduleTaskReminders() error
	ScheduleOverdueChecks() error
	ScheduleCleanupJobs() error
	ScheduleStatisticsUpdate() error
}

// WebSocketServiceInterface WebSocket服务接口
type WebSocketServiceInterface interface {
	// 连接管理
	AddConnection(userID uint, conn interface{}) error
	RemoveConnection(userID uint) error
	GetConnection(userID uint) (interface{}, bool)
	GetActiveConnections() map[uint]interface{}
	GetConnectionCount() int

	// 消息发送
	SendToUser(userID uint, message interface{}) error
	SendToUsers(userIDs []uint, message interface{}) error
	BroadcastToAll(message interface{}) error
	BroadcastToRoom(room string, message interface{}) error

	// 房间管理
	JoinRoom(userID uint, room string) error
	LeaveRoom(userID uint, room string) error
	GetRoomMembers(room string) []uint
	GetUserRooms(userID uint) []string

	// 事件处理
	HandleConnection(userID uint) error
	HandleDisconnection(userID uint) error
	HandleMessage(userID uint, messageType string, data interface{}) error

	// 通知推送
	PushNotification(userID uint, notification *models.Notification) error
	PushTaskReminder(userID uint, todo *models.Todo) error
	PushSystemMessage(message string) error
}
