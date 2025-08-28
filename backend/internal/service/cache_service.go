package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gin-web-framework/internal/redis"
	"gin-web-framework/pkg/logger"

	redisClient "github.com/go-redis/redis/v8"
)

// CacheService 缓存服务
type CacheService struct {
	redisClient redis.RedisClient
	localCache  *sync.Map // 本地缓存，用于热点数据
	mutex       sync.RWMutex
	logger      logger.LoggerInterface
}

// NewCacheService 创建缓存服务实例
func NewCacheService(redisClient redis.RedisClient, logger logger.LoggerInterface) *CacheService {
	return &CacheService{
		redisClient: redisClient,
		logger:      logger,
		localCache:  &sync.Map{},
	}
}

// CacheKey 缓存键生成器
type CacheKey struct {
	Prefix string
	UserID uint
	Type   string
	Params map[string]interface{}
}

// GenerateKey 生成缓存键
func (ck *CacheKey) GenerateKey() string {
	key := fmt.Sprintf("%s:user:%d:type:%s", ck.Prefix, ck.UserID, ck.Type)

	if ck.Params != nil {
		for k, v := range ck.Params {
			key += fmt.Sprintf(":%s:%v", k, v)
		}
	}

	return key
}

// LocalCacheItem 本地缓存项
type LocalCacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired 检查是否过期
func (item *LocalCacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// 实现CacheServiceInterface接口
func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	duration := time.Duration(expiration) * time.Second
	err = cs.redisClient.Set(ctx, key, data, duration)
	if err != nil {
		logger.Errorf("Failed to set cache key %s: %v", key, err)
		return err
	}

	// 热点数据同时存储到本地缓存
	if cs.isHotKey(key) {
		cs.setLocalCache(key, value, duration)
	}

	return nil
}

func (cs *CacheService) Get(ctx context.Context, key string) (string, error) {
	if cs.redisClient == nil {
		return "", fmt.Errorf("Redis client is not available")
	}

	// 先检查本地缓存
	if value, found := cs.getLocalCache(key); found {
		if strValue, ok := value.(string); ok {
			return strValue, nil
		}
	}

	result, err := cs.redisClient.Get(ctx, key)
	if err != nil {
		logger.Debugf("Cache miss for key %s: %v", key, err)
		return "", err
	}

	// 热点数据存储到本地缓存
	if cs.isHotKey(key) {
		cs.setLocalCache(key, result, 5*time.Minute)
	}

	return result, nil
}

func (cs *CacheService) GetObject(ctx context.Context, key string, dest interface{}) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	// 先检查本地缓存
	if value, found := cs.getLocalCache(key); found {
		if data, err := json.Marshal(value); err == nil {
			return json.Unmarshal(data, dest)
		}
	}

	result, err := cs.redisClient.Get(ctx, key)
	if err != nil {
		logger.Debugf("Cache miss for key %s: %v", key, err)
		return err
	}

	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	// 热点数据存储到本地缓存
	if cs.isHotKey(key) {
		cs.setLocalCache(key, dest, 5*time.Minute)
	}

	return nil
}

func (cs *CacheService) Delete(ctx context.Context, key string) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	// 删除本地缓存
	cs.localCache.Delete(key)

	return cs.redisClient.Del(ctx, key)
}

func (cs *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	if cs.redisClient == nil {
		return false, fmt.Errorf("Redis client is not available")
	}

	// 先检查本地缓存
	if _, found := cs.getLocalCache(key); found {
		return true, nil
	}

	result, err := cs.redisClient.Exists(ctx, key)
	return result > 0, err
}

// 批量操作
func (cs *CacheService) MSet(ctx context.Context, pairs map[string]interface{}, expiration int) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	pipeline := cs.getNativeClient().Pipeline()
	duration := time.Duration(expiration) * time.Second

	for key, value := range pairs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipeline.Set(ctx, key, data, duration)
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (cs *CacheService) MGet(ctx context.Context, keys []string) (map[string]string, error) {
	if cs.redisClient == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}

	results := make(map[string]string)
	pipeline := cs.getNativeClient().Pipeline()

	// 创建管道获取所有键
	cmds := make(map[string]*redisClient.StringCmd)
	for _, key := range keys {
		cmds[key] = pipeline.Get(ctx, key)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil && err != redisClient.Nil {
		return nil, err
	}

	// 收集结果
	for key, cmd := range cmds {
		if value, err := cmd.Result(); err == nil {
			results[key] = value
		}
	}

	return results, nil
}

func (cs *CacheService) MDelete(ctx context.Context, keys []string) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	// 删除本地缓存
	for _, key := range keys {
		cs.localCache.Delete(key)
	}

	return cs.redisClient.Del(ctx, keys...)
}

// 列表操作
func (cs *CacheService) ListPush(ctx context.Context, key string, values ...interface{}) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.LPush(ctx, key, values...)
}

func (cs *CacheService) ListPop(ctx context.Context, key string) (string, error) {
	if cs.redisClient == nil {
		return "", fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.LPop(ctx, key)
}

func (cs *CacheService) ListRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if cs.redisClient == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.LRange(ctx, key, start, stop)
}

func (cs *CacheService) ListLength(ctx context.Context, key string) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.LLen(ctx, key)
}

// 集合操作
func (cs *CacheService) SetAdd(ctx context.Context, key string, members ...interface{}) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.SAdd(ctx, key, members...)
}

func (cs *CacheService) SetMembers(ctx context.Context, key string) ([]string, error) {
	if cs.redisClient == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.SMembers(ctx, key)
}

func (cs *CacheService) SetIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if cs.redisClient == nil {
		return false, fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.SIsMember(ctx, key, member)
}

func (cs *CacheService) SetRemove(ctx context.Context, key string, members ...interface{}) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.SRem(ctx, key, members...)
}

// 过期和TTL
func (cs *CacheService) SetExpiration(ctx context.Context, key string, expiration int) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}
	duration := time.Duration(expiration) * time.Second
	return cs.redisClient.Expire(ctx, key, duration)
}

func (cs *CacheService) GetTTL(ctx context.Context, key string) (int, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().TTL(ctx, key)
	duration, err := result.Result()
	if err != nil {
		return 0, err
	}
	return int(duration.Seconds()), nil
}

// 模式匹配
func (cs *CacheService) Keys(ctx context.Context, pattern string) ([]string, error) {
	if cs.redisClient == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().Keys(ctx, pattern)
	return result.Result()
}

func (cs *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := cs.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cs.MDelete(ctx, keys)
	}
	return nil
}

// 原子操作
func (cs *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().Incr(ctx, key)
	return result.Result()
}

func (cs *CacheService) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().IncrBy(ctx, key, value)
	return result.Result()
}

func (cs *CacheService) Decrement(ctx context.Context, key string) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().Decr(ctx, key)
	return result.Result()
}

func (cs *CacheService) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().DecrBy(ctx, key, value)
	return result.Result()
}

// 分布式锁
func (cs *CacheService) Lock(ctx context.Context, key string, expiration int) (bool, error) {
	if cs.redisClient == nil {
		return false, fmt.Errorf("Redis client is not available")
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	lockValue := cs.generateLockValue()
	duration := time.Duration(expiration) * time.Second

	result := cs.getNativeClient().SetNX(ctx, lockKey, lockValue, duration)
	success, err := result.Result()
	if err != nil {
		return false, err
	}

	if success {
		// 存储锁值用于解锁验证
		cs.setLocalCache(fmt.Sprintf("lockval:%s", key), lockValue, duration)
	}

	return success, nil
}

func (cs *CacheService) Unlock(ctx context.Context, key string) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	lockKey := fmt.Sprintf("lock:%s", key)
	lockValKey := fmt.Sprintf("lockval:%s", key)

	// 获取锁值进行验证
	lockValue, found := cs.getLocalCache(lockValKey)
	if !found {
		return fmt.Errorf("lock not found or expired")
	}

	// Lua脚本确保原子性解锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result := cs.getNativeClient().Eval(ctx, script, []string{lockKey}, lockValue)
	deleted, err := result.Result()
	if err != nil {
		return err
	}

	if deleted == int64(1) {
		cs.localCache.Delete(lockValKey)
		return nil
	}

	return fmt.Errorf("failed to unlock: lock not owned or expired")
}

// 健康检查
func (cs *CacheService) Ping(ctx context.Context) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}
	return cs.redisClient.Ping(ctx)
}

func (cs *CacheService) FlushDB(ctx context.Context) error {
	if cs.redisClient == nil {
		return fmt.Errorf("Redis client is not available")
	}

	// 清空本地缓存
	cs.localCache = &sync.Map{}

	result := cs.getNativeClient().FlushDB(ctx)
	return result.Err()
}

// 统计
func (cs *CacheService) Info(ctx context.Context) (map[string]string, error) {
	if cs.redisClient == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}

	result := cs.getNativeClient().Info(ctx)
	infoStr, err := result.Result()
	if err != nil {
		return nil, err
	}

	info := make(map[string]string)
	lines := strings.Split(infoStr, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				info[parts[0]] = parts[1]
			}
		}
	}

	return info, nil
}

// 私有方法
func (cs *CacheService) getNativeClient() *redisClient.Client {
	return cs.redisClient.GetNativeClient().(*redisClient.Client)
}

func (cs *CacheService) isHotKey(key string) bool {
	// 判断是否为热点键，可以根据业务需求调整
	hotPrefixes := []string{"statistics:", "trends:", "user:", "session:"}
	for _, prefix := range hotPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func (cs *CacheService) setLocalCache(key string, value interface{}, duration time.Duration) {
	item := &LocalCacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
	cs.localCache.Store(key, item)
}

func (cs *CacheService) getLocalCache(key string) (interface{}, bool) {
	value, found := cs.localCache.Load(key)
	if !found {
		return nil, false
	}

	item, ok := value.(*LocalCacheItem)
	if !ok {
		cs.localCache.Delete(key)
		return nil, false
	}

	if item.IsExpired() {
		cs.localCache.Delete(key)
		return nil, false
	}

	return item.Value, true
}

func (cs *CacheService) generateLockValue() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CleanupExpiredLocalCache 清理过期的本地缓存
func (cs *CacheService) CleanupExpiredLocalCache() {
	cs.localCache.Range(func(key, value interface{}) bool {
		if item, ok := value.(*LocalCacheItem); ok && item.IsExpired() {
			cs.localCache.Delete(key)
		}
		return true
	})
}

// StatisticsCache 统计缓存服务
type StatisticsCache struct {
	cacheService CacheServiceInterface
}

// NewStatisticsCache 创建统计缓存服务
func NewStatisticsCache(cacheService CacheServiceInterface) *StatisticsCache {
	return &StatisticsCache{
		cacheService: cacheService,
	}
}

// GetStatisticsCacheKey 获取统计缓存键
func (sc *StatisticsCache) GetStatisticsCacheKey(userID uint, statType StatisticsType) string {
	key := &CacheKey{
		Prefix: "statistics",
		UserID: userID,
		Type:   string(statType),
	}
	return key.GenerateKey()
}

// GetTrendsCacheKey 获取趋势缓存键
func (sc *StatisticsCache) GetTrendsCacheKey(userID uint, statType StatisticsType, days int) string {
	key := &CacheKey{
		Prefix: "trends",
		UserID: userID,
		Type:   string(statType),
		Params: map[string]interface{}{
			"days": days,
		},
	}
	return key.GenerateKey()
}

// GetCachedStatistics 获取缓存的统计数据
func (sc *StatisticsCache) GetCachedStatistics(ctx context.Context, userID uint, statType StatisticsType, dest interface{}) error {
	key := sc.GetStatisticsCacheKey(userID, statType)
	return sc.cacheService.GetObject(ctx, key, dest)
}

// SetCachedStatistics 设置统计缓存
func (sc *StatisticsCache) SetCachedStatistics(ctx context.Context, userID uint, statType StatisticsType, data interface{}) error {
	key := sc.GetStatisticsCacheKey(userID, statType)
	return sc.cacheService.Set(ctx, key, data, 300) // 缓存5分钟
}

// GetCachedTrends 获取缓存的趋势数据
func (sc *StatisticsCache) GetCachedTrends(ctx context.Context, userID uint, statType StatisticsType, days int, dest interface{}) error {
	key := sc.GetTrendsCacheKey(userID, statType, days)
	return sc.cacheService.GetObject(ctx, key, dest)
}

// SetCachedTrends 设置趋势缓存
func (sc *StatisticsCache) SetCachedTrends(ctx context.Context, userID uint, statType StatisticsType, days int, data interface{}) error {
	key := sc.GetTrendsCacheKey(userID, statType, days)
	return sc.cacheService.Set(ctx, key, data, 600) // 缓存10分钟
}

// InvalidateUserCache 清除用户相关缓存
func (sc *StatisticsCache) InvalidateUserCache(ctx context.Context, userID uint) error {
	patterns := []string{
		fmt.Sprintf("statistics:user:%d:*", userID),
		fmt.Sprintf("trends:user:%d:*", userID),
	}

	for _, pattern := range patterns {
		if err := sc.cacheService.DeletePattern(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

// WarmupCache 预热缓存
func (sc *StatisticsCache) WarmupCache(ctx context.Context, userID uint) error {
	// 预热常用的统计数据
	statTypes := []StatisticsType{StatisticsTypeTodo, StatisticsTypeArticle}
	trendDays := []int{7, 30, 90}

	// 使用 goroutine 并发预热
	errCh := make(chan error, len(statTypes)*(1+len(trendDays)))
	defer close(errCh)

	// 预热统计数据
	for _, statType := range statTypes {
		go func(st StatisticsType) {
			key := sc.GetStatisticsCacheKey(userID, st)
			// 检查是否已缓存
			if exists, _ := sc.cacheService.Exists(ctx, key); !exists {
				// 这里可以调用实际的统计服务获取数据并缓存
				logger.Debugf("Cache warmup needed for statistics key: %s", key)
			}
			errCh <- nil
		}(statType)

		// 预热趋势数据
		for _, days := range trendDays {
			go func(st StatisticsType, d int) {
				key := sc.GetTrendsCacheKey(userID, st, d)
				if exists, _ := sc.cacheService.Exists(ctx, key); !exists {
					logger.Debugf("Cache warmup needed for trends key: %s", key)
				}
				errCh <- nil
			}(statType, days)
		}
	}

	// 等待所有预热任务完成
	for i := 0; i < len(statTypes)*(1+len(trendDays)); i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}

// GetCacheHitRate 获取缓存命中率
func (sc *StatisticsCache) GetCacheHitRate(ctx context.Context) (float64, error) {
	info, err := sc.cacheService.Info(ctx)
	if err != nil {
		return 0, err
	}

	hits, _ := strconv.ParseFloat(info["keyspace_hits"], 64)
	misses, _ := strconv.ParseFloat(info["keyspace_misses"], 64)

	if hits+misses == 0 {
		return 0, nil
	}

	return hits / (hits + misses), nil
}

// DatabaseSize 获取数据库大小（实现CacheServiceInterface接口）
func (cs *CacheService) DatabaseSize(ctx context.Context) (int64, error) {
	if cs.redisClient == nil {
		return 0, fmt.Errorf("Redis client is not available")
	}

	// 由于Redis客户端接口没有DBSize方法，我们返回一个默认值
	// 在实际生产环境中，可以通过Redis的INFO命令获取更多信息
	logger.Debug("DatabaseSize method called - returning default value")
	return 0, nil
}
