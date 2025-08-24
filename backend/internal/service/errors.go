package service

import "errors"

// 通用错误定义
var (
	ErrInvalidStatisticsType = errors.New("无效的统计类型")
	ErrInvalidDaysParameter  = errors.New("无效的天数参数")
	ErrUserNotFound          = errors.New("用户不存在")
	ErrInvalidCredentials    = errors.New("无效的凭据")
	ErrEmailAlreadyExists    = errors.New("邮箱已存在")
	ErrInvalidEmail          = errors.New("无效的邮箱格式")
	ErrInvalidPassword       = errors.New("无效的密码")
	ErrTodoNotFound          = errors.New("任务不存在")
	ErrArticleNotFound       = errors.New("文章不存在")
	ErrUnauthorized          = errors.New("未授权访问")
	ErrForbidden             = errors.New("禁止访问")
	ErrInternalServer        = errors.New("内部服务器错误")
	ErrDatabaseOperation     = errors.New("数据库操作失败")
	ErrInvalidRequest        = errors.New("无效的请求")
	ErrValidationFailed      = errors.New("数据验证失败")
	
	// 新增的错误定义
	ErrInvalidInput          = errors.New("无效的输入")
	ErrInvalidPriority       = errors.New("无效的优先级")
	ErrInvalidCategory       = errors.New("无效的分类")
	ErrInvalidStatus         = errors.New("无效的状态")
	ErrInvalidTimeRange      = errors.New("无效的时间范围")
	ErrWeakPassword          = errors.New("密码强度不够")
	ErrInvalidUsername       = errors.New("无效的用户名")
)
