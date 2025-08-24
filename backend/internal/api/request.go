package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RequestContext 请求上下文
type RequestContext struct {
	RequestID  string    `json:"request_id"`
	UserID     uint      `json:"user_id,omitempty"`
	Username   string    `json:"username,omitempty"`
	APIVersion string    `json:"api_version"`
	Timestamp  time.Time `json:"timestamp"`
	ClientIP   string    `json:"client_ip"`
	UserAgent  string    `json:"user_agent"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	PaginationRequest
	Query     string `json:"query" form:"query"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc"`
}

// FilterRequest 过滤请求
type FilterRequest struct {
	SearchRequest
	Filters map[string]interface{} `json:"filters" form:"filters"`
}

// GetRequestContext 从Gin上下文获取请求上下文
func GetRequestContext(c *gin.Context) *RequestContext {
	ctx := &RequestContext{
		RequestID:  c.GetString("request_id"),
		APIVersion: GetVersionFromContext(c),
		Timestamp:  time.Now(),
		ClientIP:   c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		Method:     c.Request.Method,
		Path:       c.Request.URL.Path,
	}

	// 获取用户信息（如果已认证）
	if userID, exists := c.Get("user_id"); exists {
		ctx.UserID = userID.(uint)
	}
	if username, exists := c.Get("username"); exists {
		ctx.Username = username.(string)
	}

	return ctx
}

// ValidatePagination 验证分页参数
func ValidatePagination(req *PaginationRequest) error {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	return nil
}

// ValidateSearch 验证搜索参数
func ValidateSearch(req *SearchRequest) error {
	if err := ValidatePagination(&req.PaginationRequest); err != nil {
		return err
	}

	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	return nil
}

// GetOffset 获取分页偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取分页限制
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}
