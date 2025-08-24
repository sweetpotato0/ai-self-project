package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseStatus 响应状态
type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
	StatusWarning ResponseStatus = "warning"
	StatusInfo    ResponseStatus = "info"
)

// ErrorCode 错误码
type ErrorCode string

const (
	// 通用错误码
	ErrorCodeValidation         ErrorCode = "VALIDATION_ERROR"
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest         ErrorCode = "BAD_REQUEST"
	ErrorCodeConflict           ErrorCode = "CONFLICT"
	ErrorCodeTooManyRequests    ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// 业务错误码
	ErrorCodeUserNotFound            ErrorCode = "USER_NOT_FOUND"
	ErrorCodeInvalidCredentials      ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeEmailExists             ErrorCode = "EMAIL_EXISTS"
	ErrorCodeUsernameExists          ErrorCode = "USERNAME_EXISTS"
	ErrorCodeInvalidToken            ErrorCode = "INVALID_TOKEN"
	ErrorCodeTokenExpired            ErrorCode = "TOKEN_EXPIRED"
	ErrorCodeInsufficientPermissions ErrorCode = "INSUFFICIENT_PERMISSIONS"
)

// APIResponse 标准API响应结构
type APIResponse struct {
	Status     ResponseStatus `json:"status"`
	Code       int            `json:"code"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data,omitempty"`
	Error      *APIError      `json:"error,omitempty"`
	Meta       *ResponseMeta  `json:"meta,omitempty"`
	Timestamp  time.Time      `json:"timestamp"`
	RequestID  string         `json:"request_id,omitempty"`
	APIVersion string         `json:"api_version,omitempty"`
}

// APIError 错误详情
type APIError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ResponseMeta 响应元数据
type ResponseMeta struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Sorting    *SortingMeta    `json:"sorting,omitempty"`
	Filtering  *FilteringMeta  `json:"filtering,omitempty"`
	Cache      *CacheMeta      `json:"cache,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// SortingMeta 排序元数据
type SortingMeta struct {
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// FilteringMeta 过滤元数据
type FilteringMeta struct {
	AppliedFilters   map[string]interface{} `json:"applied_filters"`
	AvailableFilters []string               `json:"available_filters"`
}

// CacheMeta 缓存元数据
type CacheMeta struct {
	FromCache bool      `json:"from_cache"`
	TTL       int       `json:"ttl,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// 便捷响应函数

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	response := &APIResponse{
		Status:     StatusSuccess,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       data,
		Timestamp:  time.Now(),
		RequestID:  c.GetString("request_id"),
		APIVersion: GetVersionFromContext(c),
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, errorCode ErrorCode, message string) {
	response := &APIResponse{
		Status:  StatusError,
		Code:    statusCode,
		Message: message,
		Error: &APIError{
			Code:    errorCode,
			Message: message,
		},
		Timestamp:  time.Now(),
		RequestID:  c.GetString("request_id"),
		APIVersion: GetVersionFromContext(c),
	}
	c.JSON(statusCode, response)
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, ErrorCodeBadRequest, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, ErrorCodeUnauthorized, message)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, ErrorCodeForbidden, message)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, ErrorCodeNotFound, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, ErrorCodeInternalServer, message)
}

// ValidationErrorResponse 验证错误响应
func ValidationErrorResponse(c *gin.Context, message string, details interface{}) {
	response := &APIResponse{
		Status:  StatusError,
		Code:    http.StatusBadRequest,
		Message: message,
		Error: &APIError{
			Code:    ErrorCodeValidation,
			Message: message,
			Details: details,
		},
		Timestamp:  time.Now(),
		RequestID:  c.GetString("request_id"),
		APIVersion: GetVersionFromContext(c),
	}
	c.JSON(http.StatusBadRequest, response)
}
