package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从gin上下文中安全地获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user not authenticated")
	}

	// 尝试多种类型转换
	switch v := userID.(type) {
	case uint:
		return v, nil
	case int:
		if v < 0 {
			return 0, errors.New("invalid user ID: negative value")
		}
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, errors.New("invalid user ID: negative value")
		}
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, errors.New("invalid user ID: negative value")
		}
		return uint(v), nil
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, errors.New("invalid user ID format")
		}
		return uint(id), nil
	default:
		return 0, errors.New("invalid user ID type")
	}
}

// GetUsernameFromContext 从gin上下文中安全地获取用户名
func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}

	if str, ok := username.(string); ok {
		return str, nil
	}
	return "", errors.New("invalid username type")
}

// GetRoleFromContext 从gin上下文中安全地获取用户角色
func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}

	if str, ok := role.(string); ok {
		return str, nil
	}
	return "", errors.New("invalid role type")
}