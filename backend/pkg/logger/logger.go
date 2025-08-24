package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gin-web-framework/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 日志级别
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
	LevelPanic LogLevel = "panic"
)

// LoggerInterface 日志接口
type LoggerInterface interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	WithField(key string, value interface{}) LoggerInterface
	WithFields(fields map[string]interface{}) LoggerInterface
	WithContext(ctx context.Context) LoggerInterface
	WithError(err error) LoggerInterface

	SetLevel(level LogLevel)
	GetLevel() LogLevel
}

// Logger 增强的日志器
type Logger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
	mu     sync.RWMutex
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      LogLevel `json:"level"`
	Format     string   `json:"format"` // json, text
	Output     string   `json:"output"` // stdout, file, both
	FilePath   string   `json:"file_path"`
	MaxSize    int      `json:"max_size"` // MB
	MaxBackups int      `json:"max_backups"`
	MaxAge     int      `json:"max_age"` // days
	Compress   bool     `json:"compress"`

	// 额外配置
	EnableCaller   bool   `json:"enable_caller"`
	EnableColor    bool   `json:"enable_color"`
	TimeFormat     string `json:"time_format"`
	EnableRotation bool   `json:"enable_rotation"`
}

// DefaultLoggerConfig 默认配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:          LevelInfo,
		Format:         "json",
		Output:         "stdout",
		FilePath:       "logs/app.log",
		MaxSize:        100,
		MaxBackups:     5,
		MaxAge:         30,
		Compress:       true,
		EnableCaller:   true,
		EnableColor:    false,
		TimeFormat:     "2006-01-02 15:04:05",
		EnableRotation: true,
	}
}

var (
	globalLogger LoggerInterface
	once         sync.Once
)

// Init 初始化日志系统
func Init() {
	once.Do(func() {
		cfg := config.Get()
		appConfig := cfg.GetApp()

		logConfig := DefaultLoggerConfig()

		// 根据环境配置调整日志级别
		switch appConfig.Environment {
		case "development":
			logConfig.Level = LevelDebug
			logConfig.Format = "text"
			logConfig.EnableColor = true
		case "testing":
			logConfig.Level = LevelWarn
			logConfig.Output = "file"
		case "production":
			logConfig.Level = LevelInfo
			logConfig.Format = "json"
			logConfig.Output = "both"
		}

		if appConfig.Debug {
			logConfig.Level = LevelDebug
		}

		globalLogger = NewLogger(logConfig)
	})
}

// NewLogger 创建新的日志器
func NewLogger(config *LoggerConfig) LoggerInterface {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置输出
	setupLoggerOutput(logger, config)

	// 设置格式
	setupLoggerFormatter(logger, config)

	// 设置调用信息
	if config.EnableCaller {
		logger.SetReportCaller(true)
		if jsonFormatter, ok := logger.Formatter.(*logrus.JSONFormatter); ok {
			jsonFormatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			}
		}
	}

	return &Logger{
		logger: logger,
		entry:  logger.WithFields(logrus.Fields{}),
	}
}

// setupLoggerOutput 设置日志输出
func setupLoggerOutput(logger *logrus.Logger, config *LoggerConfig) {
	var writers []io.Writer

	switch config.Output {
	case "stdout":
		writers = append(writers, os.Stdout)
	case "file":
		writers = append(writers, getFileWriter(config))
	case "both":
		writers = append(writers, os.Stdout)
		writers = append(writers, getFileWriter(config))
	default:
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 1 {
		logger.SetOutput(writers[0])
	} else {
		logger.SetOutput(io.MultiWriter(writers...))
	}
}

// getFileWriter 获取文件写入器
func getFileWriter(config *LoggerConfig) io.Writer {
	// 确保日志目录存在
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return os.Stdout
	}

	if config.EnableRotation {
		return &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
			LocalTime:  true,
		}
	}

	file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return os.Stdout
	}

	return file
}

// setupLoggerFormatter 设置日志格式
func setupLoggerFormatter(logger *logrus.Logger, config *LoggerConfig) {
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
			PrettyPrint:     false,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: config.TimeFormat,
			FullTimestamp:   true,
			ForceColors:     config.EnableColor,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimeFormat,
		})
	}
}

// 实现LoggerInterface接口
func (l *Logger) Debug(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Panicf(format, args...)
}

func (l *Logger) WithField(key string, value interface{}) LoggerInterface {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return &Logger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

func (l *Logger) WithFields(fields map[string]interface{}) LoggerInterface {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return &Logger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

func (l *Logger) WithContext(ctx context.Context) LoggerInterface {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 从context中提取有用的信息
	fields := make(map[string]interface{})

	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}

	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}

	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}

	return &Logger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

func (l *Logger) WithError(err error) LoggerInterface {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return &Logger{
		logger: l.logger,
		entry:  l.entry.WithError(err),
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()

	logrusLevel, err := logrus.ParseLevel(string(level))
	if err != nil {
		logrusLevel = logrus.InfoLevel
	}
	l.logger.SetLevel(logrusLevel)
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return LogLevel(l.logger.GetLevel().String())
}

// 全局日志函数，保持向后兼容
func Debug(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debugf(format, args...)
	}
}

func Info(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Infof(format, args...)
	}
}

func Warn(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warnf(format, args...)
	}
}

func Error(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Errorf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatalf(format, args...)
	}
}

func WithField(key string, value interface{}) LoggerInterface {
	if globalLogger != nil {
		return globalLogger.WithField(key, value)
	}
	return nil
}

func WithFields(fields map[string]interface{}) LoggerInterface {
	if globalLogger != nil {
		return globalLogger.WithFields(fields)
	}
	return nil
}

func WithContext(ctx context.Context) LoggerInterface {
	if globalLogger != nil {
		return globalLogger.WithContext(ctx)
	}
	return nil
}

func WithError(err error) LoggerInterface {
	if globalLogger != nil {
		return globalLogger.WithError(err)
	}
	return nil
}

// GetLogger 获取全局日志器
func GetLogger() LoggerInterface {
	if globalLogger == nil {
		panic("logger not initialized, call Init() first")
	}
	return globalLogger
}

// SetLogger 设置全局日志器
func SetLogger(logger LoggerInterface) {
	globalLogger = logger
}

// LoggerMiddleware 中间件相关
func LoggerMiddleware(logger LoggerInterface) LoggerInterface {
	return logger
}

// 性能日志记录
func LogDuration(logger LoggerInterface, operation string, start time.Time) {
	duration := time.Since(start)
	logger.WithFields(map[string]interface{}{
		"operation":   operation,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
	}).Info("Operation completed")
}

// 结构化错误日志
func LogError(logger LoggerInterface, err error, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]interface{})
	}

	logger.WithFields(context).WithError(err).Error("Error occurred")
}

// SQL查询日志
func LogSQL(logger LoggerInterface, query string, duration time.Duration, args ...interface{}) {
	logger.WithFields(map[string]interface{}{
		"query":       query,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
		"args":        args,
	}).Debug("SQL query executed")
}

// HTTP请求日志
func LogHTTPRequest(logger LoggerInterface, method, url string, statusCode int, duration time.Duration, size int64) {
	level := "info"
	if statusCode >= 400 {
		level = "warn"
	}
	if statusCode >= 500 {
		level = "error"
	}

	fields := map[string]interface{}{
		"method":      method,
		"url":         url,
		"status_code": statusCode,
		"duration":    duration.String(),
		"duration_ms": duration.Milliseconds(),
		"size":        size,
	}

	switch level {
	case "error":
		logger.WithFields(fields).Error("HTTP request completed with error")
	case "warn":
		logger.WithFields(fields).Warn("HTTP request completed with warning")
	default:
		logger.WithFields(fields).Info("HTTP request completed")
	}
}

// 清理敏感信息
func SanitizeFields(fields map[string]interface{}) map[string]interface{} {
	sensitiveKeys := []string{
		"password", "token", "secret", "key", "auth", "credential",
		"passwd", "pwd", "pass", "authorization",
	}

	sanitized := make(map[string]interface{})
	for k, v := range fields {
		key := strings.ToLower(k)
		isSensitive := false

		for _, sensitive := range sensitiveKeys {
			if strings.Contains(key, sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			sanitized[k] = "[REDACTED]"
		} else {
			sanitized[k] = v
		}
	}

	return sanitized
}
