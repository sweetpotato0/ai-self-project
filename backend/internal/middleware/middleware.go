package middleware

import (
	"strings"
	"time"

	"gin-web-framework/config"
	"gin-web-framework/pkg/jwt"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 为每个请求添加唯一ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Logger 自定义日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
		)
		return ""
	})
}

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证JWT token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	cfg := config.Get()
	appCfg := cfg.GetApp()

	if !appCfg.RateLimitEnabled {
		return gin.HandlerFunc(func(c *gin.Context) {
			c.Next()
		})
	}

	// 简单的内存速率限制器
	limiter := make(map[string]*rateLimitInfo)

	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// 获取或创建速率限制信息
		info, exists := limiter[clientIP]
		if !exists {
			info = &rateLimitInfo{
				requests:  make([]time.Time, 0),
				lastReset: now,
			}
			limiter[clientIP] = info
		}

		// 清理过期的请求记录
		info.cleanup(now, appCfg.RateLimitRPS)

		// 检查是否超过限制
		if len(info.requests) >= appCfg.RateLimitRPS {
			response.TooManyRequests(c, "Rate limit exceeded")
			return
		}

		// 添加当前请求
		info.requests = append(info.requests, now)

		c.Next()
	})
}

// CORS CORS中间件
func CORS() gin.HandlerFunc {
	cfg := config.Get()
	appCfg := cfg.GetApp()

	if !appCfg.EnableCORS {
		return gin.HandlerFunc(func(c *gin.Context) {
			c.Next()
		})
	}

	corsConfig := cors.Config{
		AllowOrigins:     appCfg.CORSOrigins,
		AllowMethods:     appCfg.CORSMethods,
		AllowHeaders:     appCfg.CORSHeaders,
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}

