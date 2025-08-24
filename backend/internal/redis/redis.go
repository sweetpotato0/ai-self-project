package redis

import (
	"context"
	"fmt"
	"time"

	"gin-web-framework/config"
	"gin-web-framework/pkg/logger"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端接口
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) error
	Close() error

	// 列表操作
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)

	// 集合操作
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SRem(ctx context.Context, key string, members ...interface{}) error

	// 哈希操作
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// 原生客户端
	GetNativeClient() interface{}
}

// redisClientImpl Redis客户端实现
type redisClientImpl struct {
	client *redis.Client
}

var (
	Client     *redis.Client
	clientImpl RedisClient
	ctx        = context.Background()
)

func Init() error {
	cfg := config.Get()

	redisConfig := cfg.GetRedis()

	// 检查Redis配置是否为空（表示禁用Redis）
	if redisConfig.Host == "" || redisConfig.Port == "" {
		logger.Info("Redis is disabled")
		return nil
	}

	Client = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password:           redisConfig.Password,
		DB:                 redisConfig.DB,
		PoolSize:           redisConfig.PoolSize,
		MinIdleConns:       redisConfig.MinIdleConns,
		MaxConnAge:         redisConfig.MaxConnAge,
		PoolTimeout:        redisConfig.PoolTimeout,
		IdleTimeout:        redisConfig.IdleTimeout,
		IdleCheckFrequency: redisConfig.IdleCheckFreq,
	})

	// 创建客户端实现
	clientImpl = &redisClientImpl{client: Client}

	// 测试连接
	err := clientImpl.Ping(ctx)
	if err != nil {
		logger.Warn("Failed to connect to Redis, continuing without Redis: %v", err)
		Client = nil
		clientImpl = nil
		return nil
	}

	logger.Info("Redis connected successfully")
	return nil
}

// 实现RedisClient接口
func (r *redisClientImpl) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("Redis is not available")
	}
	return r.client.Get(ctx, key).Result()
}

func (r *redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClientImpl) Del(ctx context.Context, keys ...string) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisClientImpl) Exists(ctx context.Context, keys ...string) (int64, error) {
	if r.client == nil {
		return 0, fmt.Errorf("Redis is not available")
	}
	return r.client.Exists(ctx, keys...).Result()
}

func (r *redisClientImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *redisClientImpl) Ping(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.Ping(ctx).Err()
}

func (r *redisClientImpl) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}

// 列表操作
func (r *redisClientImpl) LPush(ctx context.Context, key string, values ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.LPush(ctx, key, values...).Err()
}

func (r *redisClientImpl) RPush(ctx context.Context, key string, values ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.RPush(ctx, key, values...).Err()
}

func (r *redisClientImpl) LPop(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("Redis is not available")
	}
	return r.client.LPop(ctx, key).Result()
}

func (r *redisClientImpl) RPop(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("Redis is not available")
	}
	return r.client.RPop(ctx, key).Result()
}

func (r *redisClientImpl) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("Redis is not available")
	}
	return r.client.LRange(ctx, key, start, stop).Result()
}

func (r *redisClientImpl) LLen(ctx context.Context, key string) (int64, error) {
	if r.client == nil {
		return 0, fmt.Errorf("Redis is not available")
	}
	return r.client.LLen(ctx, key).Result()
}

// 集合操作
func (r *redisClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *redisClientImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("Redis is not available")
	}
	return r.client.SMembers(ctx, key).Result()
}

func (r *redisClientImpl) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("Redis is not available")
	}
	return r.client.SIsMember(ctx, key, member).Result()
}

func (r *redisClientImpl) SRem(ctx context.Context, key string, members ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.SRem(ctx, key, members...).Err()
}

// 哈希操作
func (r *redisClientImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *redisClientImpl) HGet(ctx context.Context, key, field string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("Redis is not available")
	}
	return r.client.HGet(ctx, key, field).Result()
}

func (r *redisClientImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("Redis is not available")
	}
	return r.client.HGetAll(ctx, key).Result()
}

func (r *redisClientImpl) HDel(ctx context.Context, key string, fields ...string) error {
	if r.client == nil {
		return fmt.Errorf("Redis is not available")
	}
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *redisClientImpl) GetNativeClient() interface{} {
	return r.client
}

// 全局函数，保持向后兼容
func Close() {
	if clientImpl != nil {
		clientImpl.Close()
	}
}

func Get(key string) (string, error) {
	if clientImpl == nil {
		return "", fmt.Errorf("Redis is not available")
	}
	return clientImpl.Get(ctx, key)
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if clientImpl == nil {
		return fmt.Errorf("Redis is not available")
	}
	return clientImpl.Set(ctx, key, value, expiration)
}

func Del(key string) error {
	if clientImpl == nil {
		return fmt.Errorf("Redis is not available")
	}
	return clientImpl.Del(ctx, key)
}

func Exists(key string) (bool, error) {
	if clientImpl == nil {
		return false, fmt.Errorf("Redis is not available")
	}
	result, err := clientImpl.Exists(ctx, key)
	return result > 0, err
}

func Expire(key string, expiration time.Duration) error {
	if clientImpl == nil {
		return fmt.Errorf("Redis is not available")
	}
	return clientImpl.Expire(ctx, key, expiration)
}

func GetClient() *redis.Client {
	return Client
}

// GetRedisClient 获取Redis客户端接口
func GetRedisClient() RedisClient {
	return clientImpl
}
