package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"gin-web-framework/pkg/jwt"
	"gin-web-framework/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  "INVALID_AUTH_FORMAT",
			})
			return
		}

		token := parts[1]

		// 验证JWT token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			logger.Warnf("Invalid JWT token: %v from IP: %s", err, c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", "user") // 默认角色，实际项目中应该从数据库获取

		logger.Debugf("User authenticated: %s (ID: %d) from IP: %s",
			claims.Username, claims.UserID, c.ClientIP())

		c.Next()
	})
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有认证头，继续处理
			c.Next()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 格式错误，但不阻止请求
			logger.Warnf("Invalid authorization header format from IP: %s", c.ClientIP())
			c.Next()
			return
		}

		token := parts[1]

		// 验证JWT token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			// token无效，但不阻止请求
			logger.Debugf("Invalid JWT token in optional auth: %v from IP: %s", err, c.ClientIP())
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", "user") // 默认角色，实际项目中应该从数据库获取

		logger.Debugf("User authenticated (optional): %s (ID: %d) from IP: %s",
			claims.Username, claims.UserID, c.ClientIP())

		c.Next()
	})
}

// RoleMiddleware 角色授权中间件
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 检查用户是否已认证
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTHENTICATION_REQUIRED",
			})
			return
		}

		// 获取用户角色
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "User role not found",
				"code":  "ROLE_NOT_FOUND",
			})
			return
		}

		role := userRole.(string)

		// 检查用户角色是否在所需角色列表中
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			logger.Warnf("Access denied: user %v with role %s attempted to access resource requiring roles %v from IP: %s",
				userID, role, requiredRoles, c.ClientIP())
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":          "Insufficient permissions",
				"code":           "INSUFFICIENT_PERMISSIONS",
				"required_roles": requiredRoles,
				"user_role":      role,
			})
			return
		}

		logger.Debugf("Access granted: user %v with role %s accessing resource from IP: %s",
			userID, role, c.ClientIP())

		c.Next()
	})
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// UserMiddleware 用户权限中间件
func UserMiddleware() gin.HandlerFunc {
	return RoleMiddleware("user", "admin")
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// GetCurrentUserRole 获取当前用户角色
func GetCurrentUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// RequireAuth 要求认证的辅助函数
func RequireAuth(c *gin.Context) (uint, error) {
	userID, exists := GetCurrentUserID(c)
	if !exists {
		return 0, fmt.Errorf("authentication required")
	}
	return userID, nil
}
