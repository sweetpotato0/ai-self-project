package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP请求计数器
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP请求持续时间直方图
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 活跃请求数
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// 数据库连接池指标
	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// Redis连接指标
	redisConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_active",
			Help: "Number of active Redis connections",
		},
	)

	// 业务指标
	businessMetrics = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_operations_total",
			Help: "Total number of business operations",
		},
		[]string{"operation", "status"},
	)
)

// PrometheusMetrics 中间件
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 增加活跃请求数
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// 处理请求
		c.Next()

		// 记录请求持续时间
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)

		// 记录请求总数
		status := strconv.Itoa(c.Writer.Status())
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}

// BusinessMetrics 业务指标记录
func BusinessMetrics(operation string, status string) {
	businessMetrics.WithLabelValues(operation, status).Inc()
}

// UpdateDBMetrics 更新数据库指标
func UpdateDBMetrics(active, idle int) {
	dbConnectionsActive.Set(float64(active))
	dbConnectionsIdle.Set(float64(idle))
}

// UpdateRedisMetrics 更新Redis指标
func UpdateRedisMetrics(active int) {
	redisConnectionsActive.Set(float64(active))
}

// GetMetricsHandler 获取指标的处理函数
func GetMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加自定义指标
		c.Next()
	}
}
