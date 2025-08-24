package service

import (
	"encoding/json"
	"fmt"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
	"time"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	UserID  uint                   `json:"user_id" binding:"required"`
	Type    string                 `json:"type" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data"`
}

type NotificationFilter struct {
	Type      string `json:"type"`
	IsRead    *bool  `json:"is_read"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type PaginatedNotifications struct {
	Notifications []*models.Notification `json:"notifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// CreateNotification 创建通知
func (s *NotificationService) CreateNotification(req *CreateNotificationRequest) (*models.Notification, error) {
	db := database.GetDB()

	dataJSON := ""
	if req.Data != nil {
		if data, err := json.Marshal(req.Data); err == nil {
			dataJSON = string(data)
		}
	}

	notification := &models.Notification{
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		IsRead:  false,
		Data:    dataJSON,
	}

	if err := db.Create(notification).Error; err != nil {
		return nil, fmt.Errorf("failed to create notification: %v", err)
	}

	return notification, nil
}

// GetNotifications 获取用户通知列表（分页）
func (s *NotificationService) GetNotifications(userID uint, filter *NotificationFilter) (*PaginatedNotifications, error) {
	db := database.GetDB()

	var notifications []*models.Notification
	var total int64

	query := db.Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}

	// 获取总数
	if err := query.Model(&models.Notification{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count notifications: %v", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %v", err)
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return &PaginatedNotifications{
		Notifications: notifications,
		Total:         total,
		Page:          filter.Page,
		Limit:         filter.Limit,
		TotalPages:    totalPages,
	}, nil
}

// GetUserNotifications 获取用户通知列表（向后兼容）
func (s *NotificationService) GetUserNotifications(userID uint, limit int) ([]models.Notification, error) {
	db := database.GetDB()

	var notifications []models.Notification
	query := db.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %v", err)
	}

	return notifications, nil
}

// MarkAsRead 标记通知为已读
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	db := database.GetDB()

	if err := db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to mark notification as read: %v", err)
	}

	return nil
}

// MarkAllAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	db := database.GetDB()

	if err := db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %v", err)
	}

	return nil
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(notificationID, userID uint) error {
	db := database.GetDB()

	if err := db.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf("failed to delete notification: %v", err)
	}

	return nil
}

// GetUnreadCount 获取未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get unread count: %v", err)
	}

	return count, nil
}

// CreateTaskNotification 创建任务相关通知
func (s *NotificationService) CreateTaskNotification(userID uint, taskTitle string, notificationType string) error {
	var title, message string

	switch notificationType {
	case "due_soon":
		title = "任务即将到期"
		message = fmt.Sprintf("任务\"%s\"即将到期，请及时处理", taskTitle)
	case "overdue":
		title = "任务已逾期"
		message = fmt.Sprintf("任务\"%s\"已逾期，请尽快完成", taskTitle)
	case "completed":
		title = "任务已完成"
		message = fmt.Sprintf("任务\"%s\"已完成，恭喜！", taskTitle)
	default:
		return fmt.Errorf("unknown notification type: %s", notificationType)
	}

	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		Data: map[string]interface{}{
			"task_title": taskTitle,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}

	_, err := s.CreateNotification(req)
	return err
}

// GetUnreadNotifications 获取未读通知
func (s *NotificationService) GetUnreadNotifications(userID uint) ([]*models.Notification, error) {
	db := database.GetDB()

	var notifications []*models.Notification
	if err := db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %v", err)
	}

	return notifications, nil
}

// CreateSystemNotification 创建系统通知
func (s *NotificationService) CreateSystemNotification(title, content string, userIDs []uint) error {
	for _, userID := range userIDs {
		req := &CreateNotificationRequest{
			UserID:  userID,
			Type:    "system",
			Title:   title,
			Message: content,
		}
		if _, err := s.CreateNotification(req); err != nil {
			return fmt.Errorf("failed to create system notification for user %d: %v", userID, err)
		}
	}
	return nil
}

// CreateBroadcastNotification 创建广播通知
func (s *NotificationService) CreateBroadcastNotification(title, content string) error {
	// TODO: 实现广播通知逻辑
	return fmt.Errorf("broadcast notification not implemented")
}

// CreateTaskReminderNotification 创建任务提醒通知
func (s *NotificationService) CreateTaskReminderNotification(userID uint, todoID uint) error {
	return s.CreateTaskNotification(userID, "任务提醒", "due_soon")
}

// CreateTaskOverdueNotification 创建任务逾期通知
func (s *NotificationService) CreateTaskOverdueNotification(userID uint, todoID uint) error {
	return s.CreateTaskNotification(userID, "任务逾期", "overdue")
}

// CreateTaskCompletedNotification 创建任务完成通知
func (s *NotificationService) CreateTaskCompletedNotification(userID uint, todoID uint) error {
	return s.CreateTaskNotification(userID, "任务完成", "completed")
}

// BatchMarkAsRead 批量标记为已读
func (s *NotificationService) BatchMarkAsRead(ids []uint, userID uint) error {
	db := database.GetDB()

	if err := db.Model(&models.Notification{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to batch mark notifications as read: %v", err)
	}

	return nil
}

// BatchDelete 批量删除通知
func (s *NotificationService) BatchDelete(ids []uint, userID uint) error {
	db := database.GetDB()

	if err := db.Where("id IN ? AND user_id = ?", ids, userID).
		Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf("failed to batch delete notifications: %v", err)
	}

	return nil
}

// CheckNotifications 检查通知（用于实时通知）
func (s *NotificationService) CheckNotifications(userID uint) error {
	// TODO: 实现实时通知检查逻辑
	return nil
}

// SendRealtimeNotification 发送实时通知
func (s *NotificationService) SendRealtimeNotification(userID uint, notification *models.Notification) {
	// TODO: 实现WebSocket实时通知
}
