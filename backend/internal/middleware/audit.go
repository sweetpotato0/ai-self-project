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

type AuditMiddleware struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewAuditMiddleware(db *gorm.DB, logger *logger.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		db:     db,
		logger: logger,
	}
}

// AuditLog 审计日志中间件
func (a *AuditMiddleware) AuditLog() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 排除不需要审计的路径
		if a.shouldSkipAudit(c.Request.RequestURI) {
			c.Next()
			return
		}

		startTime := time.Now()
		
		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 使用自定义 ResponseWriter 来捕获响应
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 处理请求
		c.Next()

		// 记录审计日志
		a.logAuditRecord(c, requestBody, w.body.Bytes(), startTime)
	})
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (a *AuditMiddleware) logAuditRecord(c *gin.Context, requestBody, responseBody []byte, startTime time.Time) {
	// 获取用户信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	// 构建审计日志记录
	auditLog := &model.AuditLog{
		UserID:       getUserID(userID),
		Username:     getUsernameStr(username),
		Action:       a.getActionName(c.Request.Method),
		Resource:     a.getResourceName(c.Request.RequestURI),
		Method:       c.Request.Method,
		Path:         c.Request.RequestURI,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		RequestBody:  a.sanitizeRequestBody(requestBody),
		ResponseBody: a.sanitizeResponseBody(responseBody),
		StatusCode:   c.Writer.Status(),
		Duration:     time.Since(startTime),
		Timestamp:    startTime,
	}

	// 异步保存到数据库
	go func() {
		if err := a.db.Create(auditLog).Error; err != nil {
			a.logger.Errorf("Failed to save audit log: %v", err)
		}
	}()

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
	}).Info("Audit log recorded")
}

// shouldSkipAudit 判断是否跳过审计日志
func (a *AuditMiddleware) shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/api/v1/health",
		"/api/v1/metrics",
		"/api/v1/ws",
		"/uploads/",
		"/api/v1/docs",
	}

	for _, skipPath := range skipPaths {
		if strings.Contains(path, skipPath) {
			return true
		}
	}
	return false
}

// getActionName 根据HTTP方法生成操作名称
func (a *AuditMiddleware) getActionName(method string) string {
	// 简化操作类型为基本的CRUD操作
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

// getResourceName 获取资源名称
func (a *AuditMiddleware) getResourceName(path string) string {
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

// sanitizeRequestBody 清理请求体中的敏感信息
func (a *AuditMiddleware) sanitizeRequestBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// 如果不是JSON，直接返回字符串形式
		bodyStr := string(body)
		if len(bodyStr) > 1000 {
			return bodyStr[:1000] + "..."
		}
		return bodyStr
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
	bodyStr := string(sanitized)
	if len(bodyStr) > 1000 {
		return bodyStr[:1000] + "..."
	}
	return bodyStr
}

// sanitizeResponseBody 清理响应体中的敏感信息
func (a *AuditMiddleware) sanitizeResponseBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// 如果不是JSON，截断长度
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			return bodyStr[:500] + "..."
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
		return bodyStr[:500] + "..."
	}
	return bodyStr
}

// 辅助函数
func getUserID(userID interface{}) uint {
	if uid, ok := userID.(uint); ok {
		return uid
	}
	if uid, ok := userID.(float64); ok {
		return uint(uid)
	}
	return 0
}

func getUsernameStr(username interface{}) string {
	if name, ok := username.(string); ok {
		return name
	}
	return ""
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

