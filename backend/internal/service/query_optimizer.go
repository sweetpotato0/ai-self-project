package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

// QueryOptimizer 查询优化器
type QueryOptimizer struct {
	db           *gorm.DB
	cache        CacheServiceInterface
	queryCache   sync.Map // 查询结果缓存
	preparedStmt sync.Map // 预编译语句缓存
	mutex        sync.RWMutex
}

// QueryResult 查询结果
type QueryResult struct {
	Data      interface{}
	CreatedAt time.Time
	TTL       time.Duration
}

// IsExpired 检查查询结果是否过期
func (qr *QueryResult) IsExpired() bool {
	return time.Since(qr.CreatedAt) > qr.TTL
}

// NewQueryOptimizer 创建查询优化器
func NewQueryOptimizer(db *gorm.DB, cache CacheServiceInterface) *QueryOptimizer {
	return &QueryOptimizer{
		db:    db,
		cache: cache,
	}
}

// OptimizedQuery 优化查询配置
type OptimizedQuery struct {
	SQL           string
	Args          []interface{}
	CacheKey      string
	CacheTTL      time.Duration
	UseLocalCache bool
	UseRedisCache bool
	PrepareStmt   bool
}

// ExecuteOptimizedQuery 执行优化查询
func (qo *QueryOptimizer) ExecuteOptimizedQuery(ctx context.Context, config *OptimizedQuery, dest interface{}) error {
	// 生成缓存键
	cacheKey := config.CacheKey
	if cacheKey == "" {
		cacheKey = qo.generateCacheKey(config.SQL, config.Args)
	}

	// 1. 检查本地缓存
	if config.UseLocalCache {
		if result, found := qo.getFromLocalCache(cacheKey); found {
			return qo.copyResult(result.Data, dest)
		}
	}

	// 2. 检查Redis缓存
	if config.UseRedisCache {
		if err := qo.cache.GetObject(ctx, cacheKey, dest); err == nil {
			// Redis缓存命中，同时更新本地缓存
			if config.UseLocalCache {
				qo.setToLocalCache(cacheKey, dest, config.CacheTTL)
			}
			return nil
		}
	}

	// 3. 执行数据库查询
	start := time.Now()

	var err error
	if config.PrepareStmt {
		err = qo.executeWithPreparedStmt(config.SQL, config.Args, dest)
	} else {
		err = qo.db.WithContext(ctx).Raw(config.SQL, config.Args...).Scan(dest).Error
	}

	if err != nil {
		logger.Errorf("Database query failed: %v", err)
		return err
	}

	duration := time.Since(start)

	// 记录查询性能
	qo.logQueryPerformance(config.SQL, duration, len(config.Args))

	// 4. 缓存查询结果
	if config.UseRedisCache && config.CacheTTL > 0 {
		go func() {
			if err := qo.cache.Set(ctx, cacheKey, dest, int(config.CacheTTL.Seconds())); err != nil {
				logger.Errorf("Failed to cache query result: %v", err)
			}
		}()
	}

	if config.UseLocalCache && config.CacheTTL > 0 {
		qo.setToLocalCache(cacheKey, dest, config.CacheTTL)
	}

	return nil
}

// BatchQuery 批量查询
type BatchQueryItem struct {
	SQL      string
	Args     []interface{}
	Dest     interface{}
	CacheKey string
	CacheTTL time.Duration
}

// ExecuteBatchQueries 执行批量查询
func (qo *QueryOptimizer) ExecuteBatchQueries(ctx context.Context, queries []*BatchQueryItem) error {
	// 使用事务执行批量查询
	return qo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, query := range queries {
			config := &OptimizedQuery{
				SQL:           query.SQL,
				Args:          query.Args,
				CacheKey:      query.CacheKey,
				CacheTTL:      query.CacheTTL,
				UseLocalCache: true,
				UseRedisCache: true,
				PrepareStmt:   true,
			}

			if err := qo.ExecuteOptimizedQuery(ctx, config, query.Dest); err != nil {
				return err
			}
		}
		return nil
	})
}

// generateCacheKey 生成缓存键
func (qo *QueryOptimizer) generateCacheKey(sql string, args []interface{}) string {
	// 标准化SQL（移除多余空格）
	normalizedSQL := strings.TrimSpace(strings.ReplaceAll(sql, "\n", " "))
	normalizedSQL = strings.ReplaceAll(normalizedSQL, "\t", " ")
	for strings.Contains(normalizedSQL, "  ") {
		normalizedSQL = strings.ReplaceAll(normalizedSQL, "  ", " ")
	}

	// 创建包含SQL和参数的字符串
	keyData := fmt.Sprintf("%s|%v", normalizedSQL, args)

	// 生成MD5哈希
	hasher := md5.New()
	hasher.Write([]byte(keyData))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("query:%s", hash[:16]) // 使用前16位作为缓存键
}

// getFromLocalCache 从本地缓存获取
func (qo *QueryOptimizer) getFromLocalCache(key string) (*QueryResult, bool) {
	value, found := qo.queryCache.Load(key)
	if !found {
		return nil, false
	}

	result, ok := value.(*QueryResult)
	if !ok || result.IsExpired() {
		qo.queryCache.Delete(key)
		return nil, false
	}

	return result, true
}

// setToLocalCache 设置到本地缓存
func (qo *QueryOptimizer) setToLocalCache(key string, data interface{}, ttl time.Duration) {
	result := &QueryResult{
		Data:      data,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}
	qo.queryCache.Store(key, result)
}

// executeWithPreparedStmt 使用预编译语句执行查询
func (qo *QueryOptimizer) executeWithPreparedStmt(sql string, args []interface{}, dest interface{}) error {
	// 检查预编译语句缓存
	stmtKey := qo.generateCacheKey(sql, nil)

	value, found := qo.preparedStmt.Load(stmtKey)
	if found {
		if stmt, ok := value.(*gorm.DB); ok {
			return stmt.Raw(sql, args...).Scan(dest).Error
		}
	}

	// 创建新的预编译语句
	stmt := qo.db.Session(&gorm.Session{PrepareStmt: true})
	qo.preparedStmt.Store(stmtKey, stmt)

	return stmt.Raw(sql, args...).Scan(dest).Error
}

// copyResult 复制查询结果
func (qo *QueryOptimizer) copyResult(src, dest interface{}) error {
	// 这里可以使用反射或者JSON序列化/反序列化来复制数据
	// 为了简单起见，这里返回nil，实际实现中需要根据具体需求来处理
	return nil
}

// logQueryPerformance 记录查询性能
func (qo *QueryOptimizer) logQueryPerformance(sql string, duration time.Duration, argCount int) {
	if duration > 100*time.Millisecond {
		logger.Warnf("Slow query detected: %s (took %v, %d args)",
			qo.truncateSQL(sql), duration, argCount)
	} else {
		logger.Debugf("Query executed: %s (took %v, %d args)",
			qo.truncateSQL(sql), duration, argCount)
	}
}

// truncateSQL 截断SQL用于日志
func (qo *QueryOptimizer) truncateSQL(sql string) string {
	if len(sql) > 100 {
		return sql[:97] + "..."
	}
	return sql
}

// CleanExpiredCache 清理过期的本地缓存
func (qo *QueryOptimizer) CleanExpiredCache() {
	qo.queryCache.Range(func(key, value interface{}) bool {
		if result, ok := value.(*QueryResult); ok && result.IsExpired() {
			qo.queryCache.Delete(key)
		}
		return true
	})
}

// GetCacheStats 获取缓存统计
func (qo *QueryOptimizer) GetCacheStats() map[string]interface{} {
	localCacheCount := 0
	qo.queryCache.Range(func(key, value interface{}) bool {
		localCacheCount++
		return true
	})

	preparedStmtCount := 0
	qo.preparedStmt.Range(func(key, value interface{}) bool {
		preparedStmtCount++
		return true
	})

	return map[string]interface{}{
		"local_cache_entries":   localCacheCount,
		"prepared_stmt_entries": preparedStmtCount,
	}
}

// InvalidateCache 清除指定模式的缓存
func (qo *QueryOptimizer) InvalidateCache(pattern string) error {
	ctx := context.Background()

	// 清除Redis缓存
	if err := qo.cache.DeletePattern(ctx, pattern); err != nil {
		return err
	}

	// 清除本地缓存
	qo.queryCache.Range(func(key, value interface{}) bool {
		if keyStr, ok := key.(string); ok && strings.Contains(keyStr, pattern) {
			qo.queryCache.Delete(key)
		}
		return true
	})

	return nil
}

// DatabaseConnectionPool 数据库连接池管理器
type DatabaseConnectionPool struct {
	db     *gorm.DB
	config *ConnectionPoolConfig
	mutex  sync.RWMutex
}

// ConnectionPoolConfig 连接池配置
type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewDatabaseConnectionPool 创建数据库连接池管理器
func NewDatabaseConnectionPool(db *gorm.DB, config *ConnectionPoolConfig) *DatabaseConnectionPool {
	return &DatabaseConnectionPool{
		db:     db,
		config: config,
	}
}

// OptimizeConnectionPool 优化连接池配置
func (dcp *DatabaseConnectionPool) OptimizeConnectionPool() error {
	dcp.mutex.Lock()
	defer dcp.mutex.Unlock()

	sqlDB, err := dcp.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(dcp.config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dcp.config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dcp.config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dcp.config.ConnMaxIdleTime)

	logger.Infof("Database connection pool optimized: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		dcp.config.MaxOpenConns, dcp.config.MaxIdleConns,
		dcp.config.ConnMaxLifetime, dcp.config.ConnMaxIdleTime)

	return nil
}

// GetConnectionPoolStats 获取连接池统计信息
func (dcp *DatabaseConnectionPool) GetConnectionPoolStats() (map[string]interface{}, error) {
	dcp.mutex.RLock()
	defer dcp.mutex.RUnlock()

	sqlDB, err := dcp.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	stats := sqlDB.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}

// HealthCheck 连接池健康检查
func (dcp *DatabaseConnectionPool) HealthCheck() error {
	sqlDB, err := dcp.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	stats := sqlDB.Stats()

	// 检查连接池是否健康
	if stats.OpenConnections >= stats.MaxOpenConnections {
		logger.Warnf("Database connection pool is at maximum capacity: %d/%d",
			stats.OpenConnections, stats.MaxOpenConnections)
	}

	if stats.WaitCount > 100 {
		logger.Warnf("High database connection wait count: %d", stats.WaitCount)
	}

	return nil
}
