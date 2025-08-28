package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// getUserIDFromContext 从上下文中获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

type NotificationHandler struct {
	notificationService service.NotificationServiceInterface
	logger             logger.LoggerInterface
}

func NewNotificationHandler(notificationService service.NotificationServiceInterface, logger logger.LoggerInterface) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:             logger,
	}
}

// GetNotifications 获取用户通知列表
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	filter := service.NotificationFilter{
		Limit: limit,
	}
	result, err := h.notificationService.GetNotifications(userID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get notifications")
		return
	}

	response.Success(c, gin.H{
		"notifications": result.Notifications,
		"total": result.Total,
		"page": result.Page,
		"limit": result.Limit,
	})
}

// MarkAsRead 标记通知为已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.MarkAsRead(uint(notificationID), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	response.Success(c, gin.H{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead 标记所有通知为已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	response.Success(c, gin.H{
		"message": "All notifications marked as read",
	})
}

// CheckNotifications 手动触发通知检查
func (h *NotificationHandler) CheckNotifications(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 创建通知管理器并运行检查
	notificationManager := service.NewNotificationManager(h.logger)
	notificationManager.RunNotificationChecks()

	response.Success(c, gin.H{
		"message": "Notification check completed",
	})
}

// DeleteNotification 删除通知
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.DeleteNotification(uint(notificationID), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	response.Success(c, gin.H{
		"message": "Notification deleted successfully",
	})
}

// GetUnreadCount 获取未读通知数量
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	response.Success(c, gin.H{
		"unread_count": count,
	})
}
