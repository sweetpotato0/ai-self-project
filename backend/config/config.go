package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// ConfigInterface 配置接口
type ConfigInterface interface {
	GetServer() ServerConfig
	GetDatabase() DatabaseConfig
	GetRedis() RedisConfig
	GetJWT() JWTConfig
	GetApp() AppConfig
	GetTelemetry() TelemetryConfig
	Validate() error
	Reload() error
}

// Config 应用配置
type Config struct {
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	Redis     RedisConfig     `json:"redis"`
	JWT       JWTConfig       `json:"jwt"`
	App       AppConfig       `json:"app"`
	Telemetry TelemetryConfig `json:"telemetry"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            string        `json:"port" validate:"required"`
	Mode            string        `json:"mode" validate:"oneof=debug release test"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	MaxHeaderBytes  int           `json:"max_header_bytes"`
	TrustedProxies  []string      `json:"trusted_proxies"`
	EnableMetrics   bool          `json:"enable_metrics"`
	EnablePprof     bool          `json:"enable_pprof"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `json:"driver" validate:"required,oneof=mysql postgres sqlite"`
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	DBName          string        `json:"db_name"`
	SSLMode         string        `json:"ssl_mode"`
	DBPath          string        `json:"db_path"` // SQLite specific
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	EnableMigration bool          `json:"enable_migration"`
	LogLevel        string        `json:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host          string        `json:"host"`
	Port          string        `json:"port"`
	Password      string        `json:"password"`
	DB            int           `json:"db"`
	PoolSize      int           `json:"pool_size"`
	MinIdleConns  int           `json:"min_idle_conns"`
	MaxConnAge    time.Duration `json:"max_conn_age"`
	PoolTimeout   time.Duration `json:"pool_timeout"`
	IdleTimeout   time.Duration `json:"idle_timeout"`
	IdleCheckFreq time.Duration `json:"idle_check_freq"`
	EnableCluster bool          `json:"enable_cluster"`
	ClusterAddrs  []string      `json:"cluster_addrs"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string           `json:"secret" validate:"required,min=32"`
	ExpireHours   int              `json:"expire_hours" validate:"min=1,max=720"`
	RefreshExpire int              `json:"refresh_expire"`
	Issuer        string           `json:"issuer"`
	SigningMethod string           `json:"signing_method"`
	TokenLookup   string           `json:"token_lookup"`
	TokenHeadName string           `json:"token_head_name"`
	TimeFunc      func() time.Time `json:"-"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name               string   `json:"name" validate:"required"`
	Version            string   `json:"version" validate:"required"`
	Environment        string   `json:"environment" validate:"oneof=development production testing"`
	Debug              bool     `json:"debug"`
	TimeZone           string   `json:"timezone"`
	EnableCORS         bool     `json:"enable_cors"`
	CORSOrigins        []string `json:"cors_origins"`
	CORSMethods        []string `json:"cors_methods"`
	CORSHeaders        []string `json:"cors_headers"`
	RateLimitEnabled   bool     `json:"rate_limit_enabled"`
	RateLimitRPS       int      `json:"rate_limit_rps"`
	RateLimitBurst     int      `json:"rate_limit_burst"`
	UploadMaxSize      int64    `json:"upload_max_size"`
	UploadAllowedTypes []string `json:"upload_allowed_types"`
}

// TelemetryConfig OpenTelemetry配置
type TelemetryConfig struct {
	Enabled        bool   `json:"enabled"`
	OTLPEndpoint   string `json:"otlp_endpoint"`
	EnableTracing  bool   `json:"enable_tracing"`
	EnableMetrics  bool   `json:"enable_metrics"`
	ServiceName    string `json:"service_name"`
	ServiceVersion string `json:"service_version"`
	Environment    string `json:"environment"`
}

var instance *Config

// Load 加载配置
func Load() error {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		// 如果.env文件不存在，不报错，使用默认值
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Mode:            getEnv("GIN_MODE", "debug"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", "120s"),
			MaxHeaderBytes:  getIntEnv("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
			TrustedProxies:  getSliceEnv("SERVER_TRUSTED_PROXIES", []string{"127.0.0.1"}),
			EnableMetrics:   getBoolEnv("SERVER_ENABLE_METRICS", false),
			EnablePprof:     getBoolEnv("SERVER_ENABLE_PPROF", false),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", "30s"),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			Username:        getEnv("DB_USERNAME", "postgres"),
			Password:        getEnv("DB_PASSWORD", "password"),
			DBName:          getEnv("DB_NAME", "gin_web_framework"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			DBPath:          getEnv("DB_PATH", "./todo.db"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 50),             // 增加最大连接数
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),             // 优化空闲连接数
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", "30m"),  // 增加连接生存时间
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", "10m"), // 增加空闲时间
			EnableMigration: getBoolEnv("DB_ENABLE_MIGRATION", true),
			LogLevel:        getEnv("DB_LOG_LEVEL", "warn"),
		},
		Redis: RedisConfig{
			Host:          getEnv("REDIS_HOST", "localhost"),
			Port:          getEnv("REDIS_PORT", "6379"),
			Password:      getEnv("REDIS_PASSWORD", ""),
			DB:            getIntEnv("REDIS_DB", 0),
			PoolSize:      getIntEnv("REDIS_POOL_SIZE", 10),
			MinIdleConns:  getIntEnv("REDIS_MIN_IDLE_CONNS", 3),
			MaxConnAge:    getDurationEnv("REDIS_MAX_CONN_AGE", "30m"),
			PoolTimeout:   getDurationEnv("REDIS_POOL_TIMEOUT", "30s"),
			IdleTimeout:   getDurationEnv("REDIS_IDLE_TIMEOUT", "5m"),
			IdleCheckFreq: getDurationEnv("REDIS_IDLE_CHECK_FREQ", "1m"),
			EnableCluster: getBoolEnv("REDIS_ENABLE_CLUSTER", false),
			ClusterAddrs:  getSliceEnv("REDIS_CLUSTER_ADDRS", []string{}),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", generateRandomSecret()),
			ExpireHours:   getIntEnv("JWT_EXPIRE_HOURS", 24),
			RefreshExpire: getIntEnv("JWT_REFRESH_EXPIRE", 168), // 7 days
			Issuer:        getEnv("JWT_ISSUER", "gin-web-framework"),
			SigningMethod: getEnv("JWT_SIGNING_METHOD", "HS256"),
			TokenLookup:   getEnv("JWT_TOKEN_LOOKUP", "header:Authorization"),
			TokenHeadName: getEnv("JWT_TOKEN_HEAD_NAME", "Bearer"),
			TimeFunc:      time.Now,
		},
		App: AppConfig{
			Name:               getEnv("APP_NAME", "Gin Web Framework"),
			Version:            getEnv("APP_VERSION", "1.0.0"),
			Environment:        getEnv("APP_ENV", "development"),
			Debug:              getBoolEnv("APP_DEBUG", true),
			TimeZone:           getEnv("APP_TIMEZONE", "UTC"),
			EnableCORS:         getBoolEnv("APP_ENABLE_CORS", true),
			CORSOrigins:        getSliceEnv("APP_CORS_ORIGINS", []string{"*"}),
			CORSMethods:        getSliceEnv("APP_CORS_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			CORSHeaders:        getSliceEnv("APP_CORS_HEADERS", []string{"Origin", "Content-Type", "Authorization"}),
			RateLimitEnabled:   getBoolEnv("APP_RATE_LIMIT_ENABLED", false),
			RateLimitRPS:       getIntEnv("APP_RATE_LIMIT_RPS", 100),
			RateLimitBurst:     getIntEnv("APP_RATE_LIMIT_BURST", 200),
			UploadMaxSize:      getInt64Env("APP_UPLOAD_MAX_SIZE", 10<<20), // 10MB
			UploadAllowedTypes: getSliceEnv("APP_UPLOAD_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/gif"}),
		},
		Telemetry: TelemetryConfig{
			Enabled:        getBoolEnv("TELEMETRY_ENABLED", true),
			OTLPEndpoint:   getEnv("OTLP_ENDPOINT", "http://localhost:4318/v1/traces"),
			EnableTracing:  getBoolEnv("ENABLE_TRACING", true),
			EnableMetrics:  getBoolEnv("ENABLE_METRICS", true),
			ServiceName:    getEnv("TELEMETRY_SERVICE_NAME", "ai-self-project-backend"),
			ServiceVersion: getEnv("TELEMETRY_SERVICE_VERSION", "1.0.0"),
			Environment:    getEnv("TELEMETRY_ENVIRONMENT", "development"),
		},
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	instance = cfg
	return nil
}

// Get 获取配置实例
func Get() ConfigInterface {
	if instance == nil {
		panic("config not loaded, call Load() first")
	}
	return instance
}

// 实现ConfigInterface接口
func (c *Config) GetServer() ServerConfig       { return c.Server }
func (c *Config) GetDatabase() DatabaseConfig   { return c.Database }
func (c *Config) GetRedis() RedisConfig         { return c.Redis }
func (c *Config) GetJWT() JWTConfig             { return c.JWT }
func (c *Config) GetApp() AppConfig             { return c.App }
func (c *Config) GetTelemetry() TelemetryConfig { return c.Telemetry }

// Validate 验证配置
func (c *Config) Validate() error {
	var errs []string

	// 验证服务器配置
	if c.Server.Port == "" {
		errs = append(errs, "server port is required")
	}
	if !contains([]string{"debug", "release", "test"}, c.Server.Mode) {
		errs = append(errs, "server mode must be debug, release, or test")
	}

	// 验证数据库配置
	if !contains([]string{"mysql", "postgres", "sqlite"}, c.Database.Driver) {
		errs = append(errs, "database driver must be mysql, postgres, or sqlite")
	}
	if c.Database.Driver != "sqlite" && (c.Database.Host == "" || c.Database.DBName == "") {
		errs = append(errs, "database host and name are required for non-sqlite drivers")
	}

	// 验证JWT配置
	if len(c.JWT.Secret) < 32 {
		errs = append(errs, "JWT secret must be at least 32 characters")
	}
	if c.JWT.ExpireHours < 1 || c.JWT.ExpireHours > 720 {
		errs = append(errs, "JWT expire hours must be between 1 and 720")
	}

	// 验证应用配置
	if c.App.Name == "" {
		errs = append(errs, "app name is required")
	}
	if c.App.Version == "" {
		errs = append(errs, "app version is required")
	}
	if !contains([]string{"development", "production", "testing"}, c.App.Environment) {
		errs = append(errs, "app environment must be development, production, or testing")
	}

	if len(errs) > 0 {
		return errors.New("configuration validation errors: " + strings.Join(errs, "; "))
	}

	return nil
}

// Reload 重新加载配置
func (c *Config) Reload() error {
	return Load()
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// 如果解析失败，尝试解析默认值
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 0
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateRandomSecret() string {
	// 在生产环境中，应该生成真正的随机密钥
	// 这里只是一个默认值，建议在.env中设置
	// 警告：生产环境必须设置JWT_SECRET环境变量！
	defaultSecret := "DEVELOPMENT_ONLY_JWT_SECRET_PLEASE_CHANGE_IN_PRODUCTION_123456789012345678901234567890"

	// 如果是生产环境但没有设置JWT_SECRET，应该报错
	if os.Getenv("APP_ENV") == "production" && os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET environment variable is required in production")
	}

	return defaultSecret
}
