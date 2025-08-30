package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	UserID       uint           `json:"user_id" gorm:"index;comment:操作用户ID"`
	Username     string         `json:"username" gorm:"size:100;comment:用户名"`
	Action       string         `json:"action" gorm:"size:100;index;comment:操作类型"`
	Resource     string         `json:"resource" gorm:"size:100;index;comment:资源类型"`
	Method       string         `json:"method" gorm:"size:20;comment:HTTP方法"`
	Path         string         `json:"path" gorm:"size:500;comment:请求路径"`
	IPAddress    string         `json:"ip_address" gorm:"size:45;index;comment:客户端IP"`
	UserAgent    string         `json:"user_agent" gorm:"size:500;comment:用户代理"`
	RequestBody  string         `json:"request_body" gorm:"type:text;comment:请求体"`
	ResponseBody string         `json:"response_body" gorm:"type:text;comment:响应体"`
	StatusCode   int            `json:"status_code" gorm:"index;comment:响应状态码"`
	Duration     time.Duration  `json:"duration" gorm:"comment:请求耗时(纳秒)"`
	Timestamp    time.Time      `json:"timestamp" gorm:"index;comment:操作时间"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogQuery 审计日志查询参数
type AuditLogQuery struct {
	UserID     uint      `form:"user_id"`
	Username   string    `form:"username"`
	Action     string    `form:"action"`
	Resource   string    `form:"resource"`
	Method     string    `form:"method"`
	IPAddress  string    `form:"ip_address"`
	StatusCode int       `form:"status_code"`
	StartTime  time.Time `form:"start_time"`
	EndTime    time.Time `form:"end_time"`
	Page       int       `form:"page" binding:"min=1"`
	Limit      int       `form:"limit" binding:"min=1,max=100"`
	OrderBy    string    `form:"order_by"` // timestamp, duration, status_code
	Order      string    `form:"order"`    // asc, desc
}

// AuditLogStats 审计日志统计
type AuditLogStats struct {
	TotalCount        int64                    `json:"total_count"`
	TodayCount        int64                    `json:"today_count"`
	SuccessCount      int64                    `json:"success_count"`
	ErrorCount        int64                    `json:"error_count"`
	TopUsers          []UserActionCount        `json:"top_users"`
	TopActions        []ActionCount            `json:"top_actions"`
	TopResources      []ResourceCount          `json:"top_resources"`
	HourlyDistribution []HourlyCount           `json:"hourly_distribution"`
	StatusDistribution []StatusCount           `json:"status_distribution"`
}

// UserActionCount 用户操作统计
type UserActionCount struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Count    int64  `json:"count"`
}

// ActionCount 操作类型统计
type ActionCount struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// ResourceCount 资源类型统计
type ResourceCount struct {
	Resource string `json:"resource"`
	Count    int64  `json:"count"`
}

// HourlyCount 每小时统计
type HourlyCount struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

// StatusCount 状态码统计
type StatusCount struct {
	StatusCode int   `json:"status_code"`
	Count      int64 `json:"count"`
}