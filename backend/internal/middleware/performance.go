package middleware

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"gin-web-framework/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	RequestCount   int64
	TotalDuration  time.Duration
	AvgDuration    time.Duration
	MaxDuration    time.Duration
	MinDuration    time.Duration
	ErrorCount     int64
	ActiveRequests int64
	MemoryUsage    uint64
	GoroutineCount int
}

var (
	metrics = &PerformanceMetrics{
		MinDuration: time.Hour, // 初始化为一个很大的值
	}
)

// PerformanceMonitor 性能监控中间件
func PerformanceMonitor() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// 增加活跃请求数
		metrics.ActiveRequests++

		// 获取内存统计
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		metrics.MemoryUsage = m.Alloc
		metrics.GoroutineCount = runtime.NumGoroutine()

		// 添加性能相关的请求头
		c.Header("X-Request-Start", strconv.FormatInt(start.UnixNano(), 10))

		// 处理请求
		c.Next()

		// 减少活跃请求数
		metrics.ActiveRequests--

		// 计算请求处理时间
		duration := time.Since(start)

		// 更新指标
		updateMetrics(duration, c.Writer.Status())

		// 添加响应头
		c.Header("X-Response-Time", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6))
		c.Header("X-Memory-Usage", fmt.Sprintf("%d", metrics.MemoryUsage))
		c.Header("X-Goroutines", fmt.Sprintf("%d", metrics.GoroutineCount))

		// 记录慢请求
		if duration > 1*time.Second {
			logger.Warnf("Slow request detected: %s %s took %v",
				c.Request.Method, c.Request.URL.Path, duration)
		}

		// 记录详细性能日志
		logPerformanceMetrics(c, duration, start)
	})
}

// updateMetrics 更新性能指标
func updateMetrics(duration time.Duration, statusCode int) {
	metrics.RequestCount++
	metrics.TotalDuration += duration

	// 更新平均响应时间
	metrics.AvgDuration = time.Duration(int64(metrics.TotalDuration) / metrics.RequestCount)

	// 更新最大响应时间
	if duration > metrics.MaxDuration {
		metrics.MaxDuration = duration
	}

	// 更新最小响应时间
	if duration < metrics.MinDuration {
		metrics.MinDuration = duration
	}

	// 更新错误计数
	if statusCode >= 400 {
		metrics.ErrorCount++
	}
}

// logPerformanceMetrics 记录性能指标
func logPerformanceMetrics(c *gin.Context, duration time.Duration, start time.Time) {
	ctx := context.WithValue(c.Request.Context(), "request_id", c.GetString("request_id"))

	fields := map[string]interface{}{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"status":      c.Writer.Status(),
		"duration_ms": float64(duration.Nanoseconds()) / 1e6,
		"size":        c.Writer.Size(),
		"ip":          c.ClientIP(),
		"user_agent":  c.Request.UserAgent(),
		"timestamp":   start.Unix(),
		"memory_mb":   float64(metrics.MemoryUsage) / 1024 / 1024,
		"goroutines":  metrics.GoroutineCount,
		"active_reqs": metrics.ActiveRequests,
	}

	// 添加用户ID（如果存在）
	if userID, exists := c.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	logger.WithContext(ctx).WithFields(fields).Info("Request completed")
}

// GetPerformanceMetrics 获取性能指标
func GetPerformanceMetrics() *PerformanceMetrics {
	// 创建副本避免并发问题
	return &PerformanceMetrics{
		RequestCount:   metrics.RequestCount,
		TotalDuration:  metrics.TotalDuration,
		AvgDuration:    metrics.AvgDuration,
		MaxDuration:    metrics.MaxDuration,
		MinDuration:    metrics.MinDuration,
		ErrorCount:     metrics.ErrorCount,
		ActiveRequests: metrics.ActiveRequests,
		MemoryUsage:    metrics.MemoryUsage,
		GoroutineCount: metrics.GoroutineCount,
	}
}

// ResetMetrics 重置性能指标
func ResetMetrics() {
	metrics.RequestCount = 0
	metrics.TotalDuration = 0
	metrics.AvgDuration = 0
	metrics.MaxDuration = 0
	metrics.MinDuration = time.Hour
	metrics.ErrorCount = 0
	// 不重置 ActiveRequests，因为可能有正在处理的请求
}

// MemoryCleanup 内存清理中间件
func MemoryCleanup() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// 在请求完成后，如果内存使用过高，触发GC
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// 如果内存使用超过100MB，触发垃圾回收
		if m.Alloc > 100*1024*1024 {
			runtime.GC()
			logger.Debugf("Memory cleanup triggered, usage: %.2f MB",
				float64(m.Alloc)/1024/1024)
		}
	})
}

// RequestSizeLimit 请求大小限制中间件
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{
				"error":    "Request entity too large",
				"max_size": fmt.Sprintf("%d bytes", maxSize),
			})
			c.Abort()
			return
		}
		c.Next()
	})
}

// TimeoutMiddleware 请求超时中间件
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			c.JSON(408, gin.H{
				"error":   "Request timeout",
				"timeout": timeout.String(),
			})
			c.Abort()
			logger.Warnf("Request timeout: %s %s after %v",
				c.Request.Method, c.Request.URL.Path, timeout)
		}
	})
}

// CircuitBreaker 熔断器结构
type CircuitBreaker struct {
	maxFailures     int
	resetTimeout    time.Duration
	currentFailures int
	lastFailureTime time.Time
	state           string // "closed", "open", "half-open"
}

// NewCircuitBreaker 创建新的熔断器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

// CircuitBreakerMiddleware 熔断器中间件
func (cb *CircuitBreaker) CircuitBreakerMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 检查熔断器状态
		if cb.shouldTrip() {
			c.JSON(503, gin.H{
				"error":  "Service temporarily unavailable",
				"reason": "Circuit breaker is open",
			})
			c.Abort()
			return
		}

		c.Next()

		// 根据响应状态更新熔断器
		if c.Writer.Status() >= 500 {
			cb.recordFailure()
		} else {
			cb.recordSuccess()
		}
	})
}

// shouldTrip 检查是否应该触发熔断
func (cb *CircuitBreaker) shouldTrip() bool {
	now := time.Now()

	switch cb.state {
	case "open":
		// 如果重置时间已过，切换到半开状态
		if now.Sub(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = "half-open"
			return false
		}
		return true
	case "half-open":
		return false
	default: // closed
		return false
	}
}

// recordFailure 记录失败
func (cb *CircuitBreaker) recordFailure() {
	cb.currentFailures++
	cb.lastFailureTime = time.Now()

	if cb.currentFailures >= cb.maxFailures {
		cb.state = "open"
		logger.Warnf("Circuit breaker opened due to %d failures", cb.currentFailures)
	}
}

// recordSuccess 记录成功
func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == "half-open" {
		cb.state = "closed"
		cb.currentFailures = 0
		logger.Infof("Circuit breaker closed after successful request")
	} else if cb.state == "closed" {
		cb.currentFailures = 0
	}
}
