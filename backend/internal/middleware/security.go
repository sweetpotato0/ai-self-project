package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"gin-web-framework/config"
	"gin-web-framework/pkg/logger"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 安全头设置
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss:;")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	})
}

// InputValidation 输入验证中间件
func InputValidation() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 验证Content-Type
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			// 对于PUT请求，如果没有Content-Type头部，允许通过
			if c.Request.Method == "PUT" && contentType == "" {
				// 允许没有Content-Type的PUT请求
			} else if !strings.Contains(contentType, "application/json") &&
				!strings.Contains(contentType, "multipart/form-data") &&
				!strings.Contains(contentType, "application/x-www-form-urlencoded") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Invalid content type",
					"code":  "INVALID_CONTENT_TYPE",
				})
				return
			}
		}

		// 验证User-Agent
		userAgent := c.GetHeader("User-Agent")
		if userAgent == "" {
			logger.Warnf("Request without User-Agent from IP: %s", c.ClientIP())
		}

		// 检查恶意请求模式
		if isMaliciousRequest(c) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Suspicious request detected",
				"code":  "SUSPICIOUS_REQUEST",
			})
			return
		}

		c.Next()
	})
}

// SQLInjectionProtection SQL注入防护中间件
func SQLInjectionProtection() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 检查查询参数
		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSQLInjection(value) {
					logger.Warnf("Potential SQL injection detected in query parameter %s: %s from IP: %s",
						key, value, c.ClientIP())
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "Invalid input detected",
						"code":  "INVALID_INPUT",
					})
					return
				}
			}
		}

		// 检查路径参数
		for _, param := range c.Params {
			if containsSQLInjection(param.Value) {
				logger.Warnf("Potential SQL injection detected in path parameter %s: %s from IP: %s",
					param.Key, param.Value, c.ClientIP())
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Invalid input detected",
					"code":  "INVALID_INPUT",
				})
				return
			}
		}

		c.Next()
	})
}

// isMaliciousRequest 检查是否为恶意请求
func isMaliciousRequest(c *gin.Context) bool {
	// 检查常见的恶意请求模式
	path := c.Request.URL.Path
	userAgent := c.GetHeader("User-Agent")

	// 检查路径遍历攻击
	if strings.Contains(path, "..") || strings.Contains(path, "\\") {
		return true
	}

	// 检查常见的恶意User-Agent
	maliciousPatterns := []string{
		"sqlmap", "nikto", "nmap", "w3af", "burp", "zap",
		"scanner", "crawler", "bot", "spider",
	}

	for _, pattern := range maliciousPatterns {
		if strings.Contains(strings.ToLower(userAgent), pattern) {
			return true
		}
	}

	return false
}

// RateLimiting 速率限制中间件
func RateLimiting() gin.HandlerFunc {
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
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": 60,
			})
			return
		}

		// 添加当前请求
		info.requests = append(info.requests, now)

		c.Next()
	})
}

// rateLimitInfo 速率限制信息
type rateLimitInfo struct {
	requests  []time.Time
	lastReset time.Time
}

// cleanup 清理过期的请求记录
func (r *rateLimitInfo) cleanup(now time.Time, limit int) {
	// 保留最近1秒内的请求
	cutoff := now.Add(-time.Second)

	var validRequests []time.Time
	for _, req := range r.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}

	r.requests = validRequests
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// containsSQLInjection 检查是否包含SQL注入模式
func containsSQLInjection(input string) bool {
	sqlPatterns := []string{
		"';", "--", "/*", "*/", "union", "select", "insert", "update", "delete", "drop", "create",
		"exec", "execute", "script", "javascript:", "vbscript:", "onload", "onerror",
	}

	input = strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}

	return false
}
