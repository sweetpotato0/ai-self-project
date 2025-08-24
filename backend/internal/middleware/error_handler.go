package middleware

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			response.Error(c, http.StatusInternalServerError, err)
		} else if err, ok := recovered.(error); ok {
			// 根据错误类型返回不同的状态码
			switch err {
			// 用户相关错误
			case service.ErrUserNotFound:
				response.NotFound(c, err.Error())
			case service.ErrInvalidCredentials:
				response.Unauthorized(c, err.Error())
			case service.ErrEmailAlreadyExists:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidEmail:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidPassword, service.ErrWeakPassword:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidUsername:
				response.BadRequest(c, err.Error())
				
			// 任务相关错误
			case service.ErrTodoNotFound:
				response.NotFound(c, err.Error())
			case service.ErrInvalidPriority:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidCategory:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidStatus:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidTimeRange:
				response.BadRequest(c, err.Error())
				
			// 文章相关错误
			case service.ErrArticleNotFound:
				response.NotFound(c, err.Error())
				
			// 权限相关错误
			case service.ErrUnauthorized:
				response.Unauthorized(c, err.Error())
			case service.ErrForbidden:
				response.Forbidden(c, err.Error())
				
			// 请求验证错误
			case service.ErrInvalidRequest, service.ErrInvalidInput:
				response.BadRequest(c, err.Error())
			case service.ErrValidationFailed:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidStatisticsType:
				response.BadRequest(c, err.Error())
			case service.ErrInvalidDaysParameter:
				response.BadRequest(c, err.Error())
				
			// 数据库相关错误
			case service.ErrDatabaseOperation:
				response.Error(c, http.StatusInternalServerError, "数据库操作失败")
				
			default:
				response.Error(c, http.StatusInternalServerError, "内部服务器错误")
			}
		} else {
			response.Error(c, http.StatusInternalServerError, "未知错误")
		}
		c.Abort()
	})
}

// ValidationError 验证错误处理
func ValidationError(c *gin.Context, err error) {
	response.BadRequest(c, "数据验证失败: "+err.Error())
}

// DatabaseError 数据库错误处理
func DatabaseError(c *gin.Context, err error) {
	response.Error(c, http.StatusInternalServerError, "数据库操作失败: "+err.Error())
}

// AuthenticationError 认证错误处理
func AuthenticationError(c *gin.Context, err error) {
	response.Unauthorized(c, "认证失败: "+err.Error())
}
