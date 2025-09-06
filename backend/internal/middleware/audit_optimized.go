package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"gin-web-framework/internal/model"
	"gin-web-framework/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// 最大请求体大小限制 (1MB)
	maxRequestBodySize = 1024 * 1024
	// 最大响应体大小限制 (512KB)
	maxResponseBodySize = 512 * 1024
	// 截断提示
	truncationSuffix = "...[TRUNCATED]"
)

type OptimizedAuditMiddleware struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewOptimizedAuditMiddleware(db *gorm.DB, logger *logger.Logger) *OptimizedAuditMiddleware {
	return &OptimizedAuditMiddleware{
		db:     db,
		logger: logger,
	}
}

// AuditLog 优化的审计日志中间件
func (a *OptimizedAuditMiddleware) AuditLog() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 排除不需要审计的路径
		if a.shouldSkipAudit(c.Request.RequestURI) {
			c.Next()
			return
		}

		startTime := time.Now()
		
		// 安全地读取请求体（带大小限制）
		requestBody := a.safeReadRequestBody(c)

		// 使用优化的 ResponseWriter
		w := &optimizedResponseWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
			maxSize:       maxResponseBodySize,
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 异步记录审计日志（避免阻塞请求）
		go a.logAuditRecordAsync(c, requestBody, w.getBody(), startTime)
	})
}

// 优化的响应写入器
type optimizedResponseWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int
	written bool
}

func (w *optimizedResponseWriter) Write(b []byte) (int, error) {
	// 限制响应体缓存大小
	if w.body.Len()+len(b) <= w.maxSize {
		w.body.Write(b)
	} else if !w.written {
		// 只在第一次超限时写入截断信息
		remaining := w.maxSize - w.body.Len()
		if remaining > len(truncationSuffix) {
			w.body.Write(b[:remaining-len(truncationSuffix)])
			w.body.WriteString(truncationSuffix)
		}
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

func (w *optimizedResponseWriter) getBody() []byte {
	return w.body.Bytes()
}

// 安全地读取请求体
func (a *OptimizedAuditMiddleware) safeReadRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// 使用LimitReader防止内存耗尽
	limitReader := io.LimitReader(c.Request.Body, maxRequestBodySize+1)
	requestBody, err := io.ReadAll(limitReader)
	if err != nil {
		a.logger.Errorf("Failed to read request body: %v", err)
		return ""
	}

	// 重置请求体供后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	// 检查是否被截断
	bodyStr := string(requestBody)
	if len(requestBody) > maxRequestBodySize {
		bodyStr = bodyStr[:maxRequestBodySize-len(truncationSuffix)] + truncationSuffix
	}

	return a.sanitizeRequestBody(bodyStr)
}

// 异步记录审计日志
func (a *OptimizedAuditMiddleware) logAuditRecordAsync(c *gin.Context, requestBody string, responseBody []byte, startTime time.Time) {
	// 获取用户信息
	userID := a.getUserID(c.GetUint("user_id"))
	username := c.GetString("username")

	// 构建审计日志记录
	auditLog := &model.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       a.getActionName(c.Request.Method),
		Resource:     a.getResourceName(c.Request.RequestURI),
		Method:       c.Request.Method,
		Path:         c.Request.RequestURI,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		RequestBody:  requestBody,
		ResponseBody: a.sanitizeResponseBody(responseBody),
		StatusCode:   c.Writer.Status(),
		Duration:     time.Since(startTime),
		Timestamp:    startTime,
	}

	// 保存到数据库（带错误恢复）
	if err := a.db.Create(auditLog).Error; err != nil {
		a.logger.Errorf("Failed to save audit log: %v", err)
		// 这里可以添加降级策略，比如写入文件或发送到队列
	}

	// 记录结构化日志
	a.logger.WithFields(map[string]interface{}{
		"audit_id":    auditLog.ID,
		"user_id":     auditLog.UserID,
		"username":    auditLog.Username,
		"action":      auditLog.Action,
		"resource":    auditLog.Resource,
		"method":      auditLog.Method,
		"path":        auditLog.Path,
		"ip":          auditLog.IPAddress,
		"status":      auditLog.StatusCode,
		"duration_ms": auditLog.Duration.Milliseconds(),
		"req_size":    len(requestBody),
		"resp_size":   len(responseBody),
	}).Info("Audit log recorded")
}

// 辅助方法保持与原版本兼容
func (a *OptimizedAuditMiddleware) shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/api/v1/health",
		"/api/v1/metrics", 
		"/api/v1/ws",
		"/uploads/",
		"/api/v1/docs",
		"/favicon.ico",
		"/static/",
	}

	for _, skipPath := range skipPaths {
		if strings.Contains(path, skipPath) {
			return true
		}
	}
	return false
}

func (a *OptimizedAuditMiddleware) getActionName(method string) string {
	switch method {
	case "GET":
		return "查询"
	case "POST":
		return "创建"
	case "PUT", "PATCH":
		return "更新"
	case "DELETE":
		return "删除"
	default:
		return strings.ToUpper(method)
	}
}

func (a *OptimizedAuditMiddleware) getResourceName(path string) string {
	cleanPath := strings.TrimPrefix(path, "/api/v1")
	if strings.Contains(cleanPath, "?") {
		cleanPath = strings.Split(cleanPath, "?")[0]
	}
	
	segments := strings.Split(strings.Trim(cleanPath, "/"), "/")
	if len(segments) > 0 && segments[0] != "" {
		return segments[0]
	}
	return "unknown"
}

func (a *OptimizedAuditMiddleware) sanitizeRequestBody(body string) string {
	if body == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// 如果不是JSON，直接返回（已经截断过了）
		return body
	}

	// 清理敏感字段
	sensitiveFields := []string{
		"password", "token", "secret", "key", "auth",
		"oldPassword", "newPassword", "confirmPassword",
	}
	
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***REDACTED***"
		}
	}

	sanitized, _ := json.Marshal(data)
	return string(sanitized)
}

func (a *OptimizedAuditMiddleware) sanitizeResponseBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// 如果不是JSON，截断长度返回
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			return bodyStr[:500] + truncationSuffix
		}
		return bodyStr
	}

	// 清理敏感字段
	if dataField, ok := data["data"].(map[string]interface{}); ok {
		sensitiveFields := []string{"token", "password", "secret", "key"}
		for _, field := range sensitiveFields {
			if _, exists := dataField[field]; exists {
				dataField[field] = "***REDACTED***"
			}
		}
	}

	sanitized, _ := json.Marshal(data)
	bodyStr := string(sanitized)
	if len(bodyStr) > 500 {
		return bodyStr[:500] + truncationSuffix
	}
	return bodyStr
}

func (a *OptimizedAuditMiddleware) getUserID(userID uint) uint {
	return userID
}