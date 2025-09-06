package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gin-web-framework/pkg/errors"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/security"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 增强的安全中间件
type SecurityMiddleware struct {
	validator     *security.InputValidator
	logger        logger.LoggerInterface
	rateLimiter   RateLimiter
	trustedProxies []string
}

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(clientIP string) bool
	Reset(clientIP string)
}

// NewSecurityMiddleware 创建安全中间件
func NewSecurityMiddleware(logger logger.LoggerInterface, trustedProxies []string) *SecurityMiddleware {
	validator := security.NewInputValidator()
	
	// 添加通用安全规则
	validator.AddRule("*", security.NewSQLInjectionRule(""))
	validator.AddRule("*", security.NewXSSRule(""))
	
	return &SecurityMiddleware{
		validator:      validator,
		logger:         logger,
		rateLimiter:    NewMemoryRateLimiter(100, time.Minute), // 每分钟100次请求
		trustedProxies: trustedProxies,
	}
}

// InputValidation 输入验证中间件
func (sm *SecurityMiddleware) InputValidation() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 只验证POST、PUT、PATCH请求的body
		if !isWriteMethod(c.Request.Method) {
			c.Next()
			return
		}

		// 读取请求体
		bodyBytes, err := sm.safeReadBody(c)
		if err != nil {
			sm.logger.WithFields(map[string]any{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
			}).Warn("读取请求体失败")
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// 重置请求体供后续处理
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 解析JSON并验证
		if len(bodyBytes) > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err == nil {
				if err := sm.validator.Validate(c.Request.Context(), data); err != nil {
					sm.logger.WithFields(map[string]any{
						"error":     err.Error(),
						"path":      c.Request.URL.Path,
						"client_ip": c.ClientIP(),
					}).Warn("输入验证失败")
					
					c.AbortWithStatusJSON(400, gin.H{
						"error":   "Input validation failed",
						"details": err.Error(),
					})
					return
				}
			}
		}

		c.Next()
	})
}

// RateLimit 限流中间件
func (sm *SecurityMiddleware) RateLimit() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := sm.getRealClientIP(c)
		
		if !sm.rateLimiter.Allow(clientIP) {
			sm.logger.WithFields(map[string]any{
				"client_ip": clientIP,
				"path":      c.Request.URL.Path,
				"method":    c.Request.Method,
			}).Warn("请求频率超限")
			
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	})
}

// SecurityHeaders 安全头中间件（增强版）
func (sm *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 基础安全头
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// 严格的CSP策略
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"child-src 'none'; " +
			"worker-src 'none'; " +
			"manifest-src 'self'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Header("Content-Security-Policy", csp)
		
		// HSTS头（仅HTTPS）
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// 权限策略
		c.Header("Permissions-Policy", 
			"accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")
		
		c.Next()
	})
}

// IPFiltering IP过滤中间件
func (sm *SecurityMiddleware) IPFiltering(blacklist, whitelist []string) gin.HandlerFunc {
	blacklistMap := make(map[string]bool)
	for _, ip := range blacklist {
		blacklistMap[ip] = true
	}
	
	whitelistMap := make(map[string]bool)
	for _, ip := range whitelist {
		whitelistMap[ip] = true
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := sm.getRealClientIP(c)
		
		// 检查白名单（如果存在）
		if len(whitelistMap) > 0 && !whitelistMap[clientIP] {
			sm.logger.WithFields(map[string]any{
				"client_ip": clientIP,
				"reason":    "not_in_whitelist",
			}).Warn("IP访问被拒绝")
			
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Access denied",
			})
			return
		}
		
		// 检查黑名单
		if blacklistMap[clientIP] {
			sm.logger.WithFields(map[string]any{
				"client_ip": clientIP,
				"reason":    "in_blacklist",
			}).Warn("IP访问被拒绝")
			
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Access denied",
			})
			return
		}

		c.Next()
	})
}

// RequestSizeLimit 请求大小限制中间件
func (sm *SecurityMiddleware) RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			sm.logger.WithFields(map[string]any{
				"content_length": c.Request.ContentLength,
				"max_size":       maxSize,
				"client_ip":      c.ClientIP(),
			}).Warn("请求体过大")
			
			c.AbortWithStatusJSON(413, gin.H{
				"error": "Request entity too large",
			})
			return
		}

		// 限制读取大小
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	})
}

// 辅助方法

func (sm *SecurityMiddleware) safeReadBody(c *gin.Context) ([]byte, error) {
	if c.Request.Body == nil {
		return []byte{}, nil
	}

	// 限制读取大小防止内存攻击
	const maxBodySize = 32 * 1024 * 1024 // 32MB
	limitReader := io.LimitReader(c.Request.Body, maxBodySize)
	
	bodyBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, errors.NewValidationError("Failed to read request body", err)
	}

	return bodyBytes, nil
}

func (sm *SecurityMiddleware) getRealClientIP(c *gin.Context) string {
	// 检查可信代理的头部
	if sm.isFromTrustedProxy(c) {
		// X-Forwarded-For 头部处理
		if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
			// 取第一个IP（客户端真实IP）
			if ips := strings.Split(xForwardedFor, ","); len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		
		// X-Real-IP 头部
		if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
			return xRealIP
		}
	}

	// 使用连接的远程地址
	return c.ClientIP()
}

func (sm *SecurityMiddleware) isFromTrustedProxy(c *gin.Context) bool {
	if len(sm.trustedProxies) == 0 {
		return false
	}
	
	remoteIP := c.ClientIP()
	for _, trustedIP := range sm.trustedProxies {
		if remoteIP == trustedIP {
			return true
		}
	}
	
	return false
}

func isWriteMethod(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}

// MemoryRateLimiter 内存限流器实现
type MemoryRateLimiter struct {
	requests   map[string][]time.Time
	limit      int
	window     time.Duration
	lastCleanup time.Time
}

func NewMemoryRateLimiter(limit int, window time.Duration) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		requests:   make(map[string][]time.Time),
		limit:      limit,
		window:     window,
		lastCleanup: time.Now(),
	}
}

func (rl *MemoryRateLimiter) Allow(clientIP string) bool {
	now := time.Now()
	
	// 定期清理过期数据
	if now.Sub(rl.lastCleanup) > rl.window {
		rl.cleanup(now)
		rl.lastCleanup = now
	}
	
	// 获取客户端的请求历史
	requests := rl.requests[clientIP]
	
	// 移除过期请求
	cutoff := now.Add(-rl.window)
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		return false
	}
	
	// 添加当前请求
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests
	
	return true
}

func (rl *MemoryRateLimiter) Reset(clientIP string) {
	delete(rl.requests, clientIP)
}

func (rl *MemoryRateLimiter) cleanup(now time.Time) {
	cutoff := now.Add(-rl.window)
	
	for clientIP, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		
		if len(validRequests) == 0 {
			delete(rl.requests, clientIP)
		} else {
			rl.requests[clientIP] = validRequests
		}
	}
}