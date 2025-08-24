package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gin-web-framework/config"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	db     *gorm.DB
	config config.ConfigInterface
	mu     sync.RWMutex
}

// DatabaseInterface 数据库接口
type DatabaseInterface interface {
	GetDB() *gorm.DB
	Close() error
	AutoMigrate() error
	SeedData() error
	HealthCheck() error
	GetConnectionStats() map[string]interface{}
	Transaction(fn func(*gorm.DB) error) error
	WithContext(ctx context.Context) *gorm.DB
}

var (
	dbManager *DatabaseManager
	once      sync.Once
)

// Init 初始化数据库连接
func Init() error {
	var initErr error
	once.Do(func() {
		cfg := config.Get()
		manager := &DatabaseManager{
			config: cfg,
		}

		if err := manager.connect(); err != nil {
			initErr = err
			return
		}

		// 自动迁移
		if cfg.GetDatabase().EnableMigration {
			if err := manager.AutoMigrate(); err != nil {
				logger.Error("Failed to auto migrate: " + err.Error())
				// 不阻断启动，只记录错误
			}
		}

		dbManager = manager
		logger.Info("Database connected successfully")
	})

	return initErr
}

// connect 建立数据库连接
func (dm *DatabaseManager) connect() error {
	dbConfig := dm.config.GetDatabase()

	// 配置GORM日志级别
	var logLevel gormLogger.LogLevel
	switch dbConfig.LogLevel {
	case "silent":
		logLevel = gormLogger.Silent
	case "error":
		logLevel = gormLogger.Error
	case "warn":
		logLevel = gormLogger.Warn
	case "info":
		logLevel = gormLogger.Info
	default:
		logLevel = gormLogger.Warn
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var db *gorm.DB
	var err error

	switch dbConfig.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			dbConfig.Host,
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.DBName,
			dbConfig.Port,
			dbConfig.SSLMode,
		)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.DBName,
		)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)

	case "sqlite":
		// 确保数据目录存在
		dbDir := filepath.Dir(dbConfig.DBPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		db, err = gorm.Open(sqlite.Open(dbConfig.DBPath), gormConfig)

	default:
		return fmt.Errorf("unsupported database driver: %s", dbConfig.Driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	if err := dm.configureConnectionPool(db, dbConfig); err != nil {
		return fmt.Errorf("failed to configure connection pool: %w", err)
	}

	dm.mu.Lock()
	dm.db = db
	dm.mu.Unlock()

	return nil
}

// configureConnectionPool 配置连接池
func (dm *DatabaseManager) configureConnectionPool(db *gorm.DB, dbConfig config.DatabaseConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	// SQLite特殊处理
	if dbConfig.Driver == "sqlite" {
		sqlDB.SetMaxOpenConns(1) // SQLite只允许一个写连接
	}

	return nil
}

// 实现DatabaseInterface接口
func (dm *DatabaseManager) GetDB() *gorm.DB {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.db
}

func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.db != nil {
		sqlDB, err := dm.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (dm *DatabaseManager) AutoMigrate() error {
	// 自动迁移模型
	if err := dm.db.AutoMigrate(
		&models.User{},
		&models.UserSettings{},
		&models.Product{},
		&models.TodoCategory{},
		&models.TodoPriority{},
		&models.Todo{},
		&models.Notification{},
		&models.Article{},
		&models.Category{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 初始化种子数据
	return dm.SeedData()
}

func (dm *DatabaseManager) SeedData() error {
	// 初始化优先级
	priorities := []models.TodoPriority{
		{Name: "低", Level: 1, Color: "#67C23A"},
		{Name: "中", Level: 2, Color: "#E6A23C"},
		{Name: "高", Level: 3, Color: "#F56C6C"},
		{Name: "紧急", Level: 4, Color: "#FF4949"},
		{Name: "立即", Level: 5, Color: "#FF0000"},
	}

	for _, priority := range priorities {
		var existing models.TodoPriority
		if err := dm.db.Where("level = ?", priority.Level).First(&existing).Error; err != nil {
			if err := dm.db.Create(&priority).Error; err != nil {
				logger.Error(fmt.Sprintf("Failed to create priority %s: %v", priority.Name, err))
			} else {
				logger.Info(fmt.Sprintf("Created priority: %s", priority.Name))
			}
		}
	}

	// 初始化分类
	categories := []models.TodoCategory{
		{Name: "工作", Color: "#409EFF", Description: "工作相关任务"},
		{Name: "学习", Color: "#67C23A", Description: "学习相关任务"},
		{Name: "生活", Color: "#E6A23C", Description: "日常生活任务"},
		{Name: "健康", Color: "#F56C6C", Description: "健康相关任务"},
		{Name: "娱乐", Color: "#909399", Description: "娱乐休闲任务"},
	}

	for _, category := range categories {
		var existing models.TodoCategory
		if err := dm.db.Where("name = ?", category.Name).First(&existing).Error; err != nil {
			if err := dm.db.Create(&category).Error; err != nil {
				logger.Error(fmt.Sprintf("Failed to create category %s: %v", category.Name, err))
			} else {
				logger.Info(fmt.Sprintf("Created category: %s", category.Name))
			}
		}
	}

	logger.Info("Seed data initialized")
	return nil
}

func (dm *DatabaseManager) HealthCheck() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func (dm *DatabaseManager) GetConnectionStats() map[string]interface{} {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

func (dm *DatabaseManager) Transaction(fn func(*gorm.DB) error) error {
	return dm.db.Transaction(fn)
}

func (dm *DatabaseManager) WithContext(ctx context.Context) *gorm.DB {
	return dm.db.WithContext(ctx)
}

// 全局函数，保持兼容性
func GetDB() *gorm.DB {
	if dbManager == nil {
		panic("database not initialized, call Init() first")
	}
	return dbManager.GetDB()
}

func GetDatabaseManager() DatabaseInterface {
	if dbManager == nil {
		panic("database not initialized, call Init() first")
	}
	return dbManager
}

func Close() error {
	if dbManager != nil {
		return dbManager.Close()
	}
	return nil
}

func AutoMigrate() error {
	if dbManager == nil {
		return fmt.Errorf("database not initialized")
	}
	return dbManager.AutoMigrate()
}

func SeedData() error {
	if dbManager == nil {
		return fmt.Errorf("database not initialized")
	}
	return dbManager.SeedData()
}
