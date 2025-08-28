package handler

import (
	"net/http"
	"time"

	"gin-web-framework/internal/middleware"
	"gin-web-framework/internal/models"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsService *service.SettingsService
	logger          logger.LoggerInterface
}

func NewSettingsHandler(settingsService *service.SettingsService, logger logger.LoggerInterface) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}

// GetSettings 获取用户设置
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	settings, err := h.settingsService.GetUserSettings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// UpdateProfile 更新个人资料
func (h *SettingsHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.UpdateProfile(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "个人资料更新成功"})
}

// ChangePassword 修改密码
func (h *SettingsHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.ChangePassword(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// UpdateNotificationSettings 更新通知设置
func (h *SettingsHandler) UpdateNotificationSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var req models.UpdateNotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.UpdateNotificationSettings(userID.(uint), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存通知设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知设置已保存"})
}

// UpdateInterfaceSettings 更新界面设置
func (h *SettingsHandler) UpdateInterfaceSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	var req models.UpdateInterfaceSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingsService.UpdateInterfaceSettings(userID.(uint), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存界面设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "界面设置已保存"})
}

// ExportData 导出用户数据
func (h *SettingsHandler) ExportData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	data, err := h.settingsService.GetUserData(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出数据失败"})
		return
	}

	data["exportTime"] = time.Now().Format(time.RFC3339)

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=taskmaster-export-"+time.Now().Format("2006-01-02")+".json")
	c.JSON(http.StatusOK, data)
}

// ClearCompletedTasks 清理已完成任务
func (h *SettingsHandler) ClearCompletedTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	if err := h.settingsService.ClearCompletedTasks(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清理任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已完成任务清理成功"})
}

// HealthCheck 健康检查
func (h *SettingsHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}

// GetMetrics 获取性能指标
func (h *SettingsHandler) GetMetrics(c *gin.Context) {
	metrics := middleware.GetPerformanceMetrics()

	c.JSON(http.StatusOK, gin.H{
		"performance": map[string]interface{}{
			"request_count":   metrics.RequestCount,
			"avg_duration_ms": float64(metrics.AvgDuration.Nanoseconds()) / 1e6,
			"max_duration_ms": float64(metrics.MaxDuration.Nanoseconds()) / 1e6,
			"min_duration_ms": float64(metrics.MinDuration.Nanoseconds()) / 1e6,
			"error_count":     metrics.ErrorCount,
			"active_requests": metrics.ActiveRequests,
			"memory_usage_mb": float64(metrics.MemoryUsage) / 1024 / 1024,
			"goroutine_count": metrics.GoroutineCount,
		},
		"timestamp": c.GetInt64("timestamp"),
	})
}