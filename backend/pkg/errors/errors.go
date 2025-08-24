package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// ErrorType 错误类型
type ErrorType string

const (
	// Business errors
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeBusiness     ErrorType = "BUSINESS_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeTooManyReqs  ErrorType = "TOO_MANY_REQUESTS"

	// System errors
	ErrorTypeInternal    ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase    ErrorType = "DATABASE_ERROR"
	ErrorTypeNetwork     ErrorType = "NETWORK_ERROR"
	ErrorTypeTimeout     ErrorType = "TIMEOUT_ERROR"
	ErrorTypeUnavailable ErrorType = "SERVICE_UNAVAILABLE"
)

// ErrorCode 错误码
type ErrorCode int

const (
	// 2xx Success codes (for reference)
	CodeSuccess ErrorCode = 200

	// 4xx Client error codes
	CodeBadRequest       ErrorCode = 400
	CodeUnauthorized     ErrorCode = 401
	CodeForbidden        ErrorCode = 403
	CodeNotFound         ErrorCode = 404
	CodeConflict         ErrorCode = 409
	CodeValidationFailed ErrorCode = 422
	CodeTooManyRequests  ErrorCode = 429

	// 5xx Server error codes
	CodeInternalError      ErrorCode = 500
	CodeDatabaseError      ErrorCode = 501
	CodeServiceUnavailable ErrorCode = 503
	CodeTimeout            ErrorCode = 504
)

// AppError 应用错误结构
type AppError struct {
	Type      ErrorType   `json:"type"`
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	UserID    uint        `json:"user_id,omitempty"`

	// Internal fields
	Cause    error  `json:"-"`
	Stack    string `json:"-"`
	File     string `json:"-"`
	Line     int    `json:"-"`
	Function string `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap 支持error wrapping
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ToJSON 转换为JSON格式
func (e *AppError) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// GetHTTPStatus 获取HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeTooManyReqs:
		return http.StatusTooManyRequests
	case ErrorTypeDatabase, ErrorTypeInternal, ErrorTypeNetwork, ErrorTypeTimeout:
		return http.StatusInternalServerError
	case ErrorTypeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithUserID 添加用户ID
func (e *AppError) WithUserID(userID uint) *AppError {
	e.UserID = userID
	return e
}

// 错误构造函数
func newAppError(errType ErrorType, code ErrorCode, message string, cause error) *AppError {
	// 获取调用栈信息
	pc, file, line, ok := runtime.Caller(2)
	var function string
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}

	return &AppError{
		Type:      errType,
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now().UTC(),
		File:      file,
		Line:      line,
		Function:  function,
		Stack:     getStackTrace(),
	}
}

// 业务错误构造函数
func NewValidationError(message string, details interface{}) *AppError {
	return newAppError(ErrorTypeValidation, CodeValidationFailed, message, nil).WithDetails(details)
}

func NewBusinessError(message string) *AppError {
	return newAppError(ErrorTypeBusiness, CodeBadRequest, message, nil)
}

func NewNotFoundError(resource string) *AppError {
	return newAppError(ErrorTypeNotFound, CodeNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return newAppError(ErrorTypeUnauthorized, CodeUnauthorized, message, nil)
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return newAppError(ErrorTypeForbidden, CodeForbidden, message, nil)
}

func NewConflictError(message string) *AppError {
	return newAppError(ErrorTypeConflict, CodeConflict, message, nil)
}

func NewTooManyRequestsError(message string) *AppError {
	if message == "" {
		message = "Too many requests"
	}
	return newAppError(ErrorTypeTooManyReqs, CodeTooManyRequests, message, nil)
}

// 系统错误构造函数
func NewInternalError(message string, cause error) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return newAppError(ErrorTypeInternal, CodeInternalError, message, cause)
}

func NewDatabaseError(message string, cause error) *AppError {
	if message == "" {
		message = "Database operation failed"
	}
	return newAppError(ErrorTypeDatabase, CodeDatabaseError, message, cause)
}

func NewNetworkError(message string, cause error) *AppError {
	if message == "" {
		message = "Network operation failed"
	}
	return newAppError(ErrorTypeNetwork, CodeInternalError, message, cause)
}

func NewTimeoutError(message string, cause error) *AppError {
	if message == "" {
		message = "Operation timeout"
	}
	return newAppError(ErrorTypeTimeout, CodeTimeout, message, cause)
}

func NewServiceUnavailableError(message string) *AppError {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	return newAppError(ErrorTypeUnavailable, CodeServiceUnavailable, message, nil)
}

// Wrap 包装现有错误
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，则更新消息
	if appErr, ok := err.(*AppError); ok {
		appErr.Message = message
		return appErr
	}

	return newAppError(ErrorTypeInternal, CodeInternalError, message, err)
}

// WrapWithType 包装现有错误并指定类型
func WrapWithType(err error, errType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	var code ErrorCode
	switch errType {
	case ErrorTypeValidation:
		code = CodeValidationFailed
	case ErrorTypeNotFound:
		code = CodeNotFound
	case ErrorTypeUnauthorized:
		code = CodeUnauthorized
	case ErrorTypeForbidden:
		code = CodeForbidden
	case ErrorTypeConflict:
		code = CodeConflict
	case ErrorTypeTooManyReqs:
		code = CodeTooManyRequests
	case ErrorTypeDatabase:
		code = CodeDatabaseError
	case ErrorTypeUnavailable:
		code = CodeServiceUnavailable
	case ErrorTypeTimeout:
		code = CodeTimeout
	default:
		code = CodeInternalError
	}

	return newAppError(errType, code, message, err)
}

// Is 检查错误类型
func Is(err error, errType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errType
	}
	return false
}

// IsCode 检查错误码
func IsCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// getStackTrace 获取调用栈
func getStackTrace() string {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return string(buf)
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	HandleError(err error, requestID string) *AppError
	LogError(err *AppError)
	ShouldRetry(err *AppError) bool
	GetRetryDelay(err *AppError, attempt int) time.Duration
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct{}

// HandleError 处理错误
func (h *DefaultErrorHandler) HandleError(err error, requestID string) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		if appErr.RequestID == "" {
			appErr.RequestID = requestID
		}
		return appErr
	}

	// 包装为内部错误
	return NewInternalError("An unexpected error occurred", err).WithRequestID(requestID)
}

// LogError 记录错误日志
func (h *DefaultErrorHandler) LogError(err *AppError) {
	// 这里可以集成具体的日志系统
	fmt.Printf("[ERROR] %s: %s (Request: %s, User: %d)\n",
		err.Type, err.Message, err.RequestID, err.UserID)
}

// ShouldRetry 判断是否应该重试
func (h *DefaultErrorHandler) ShouldRetry(err *AppError) bool {
	switch err.Type {
	case ErrorTypeTimeout, ErrorTypeNetwork, ErrorTypeUnavailable:
		return true
	case ErrorTypeDatabase:
		// 某些数据库错误可以重试
		return true
	default:
		return false
	}
}

// GetRetryDelay 获取重试延迟
func (h *DefaultErrorHandler) GetRetryDelay(err *AppError, attempt int) time.Duration {
	// 指数退避算法
	baseDelay := time.Second
	maxDelay := time.Minute

	delay := time.Duration(attempt) * baseDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// 全局错误处理器
var GlobalErrorHandler ErrorHandler = &DefaultErrorHandler{}

// SetGlobalErrorHandler 设置全局错误处理器
func SetGlobalErrorHandler(handler ErrorHandler) {
	GlobalErrorHandler = handler
}
