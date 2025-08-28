package service

import (
	"fmt"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"
	"time"
)

// NotificationManager 通知管理器
type NotificationManager struct {
	notificationService *NotificationService
}

func NewNotificationManager(logger logger.LoggerInterface) *NotificationManager {
	return &NotificationManager{
		notificationService: NewNotificationService(logger),
	}
}

// CheckTaskDueNotifications 检查任务到期通知
func (nm *NotificationManager) CheckTaskDueNotifications() error {
	db := database.GetDB()
	now := time.Now()

	// 查找即将到期的任务（24小时内）
	dueSoon := now.Add(24 * time.Hour)
	var dueSoonTasks []models.Todo

	if err := db.Where("due_date <= ? AND due_date > ? AND status != ?",
		dueSoon, now, "completed").Find(&dueSoonTasks).Error; err != nil {
		return fmt.Errorf("failed to query due soon tasks: %v", err)
	}

	for _, task := range dueSoonTasks {
		// 检查是否已经发送过通知
		var existingNotification models.Notification
		if err := db.Where("user_id = ? AND type = ? AND data LIKE ?",
			task.CreatedBy, "due_soon", "%"+task.Title+"%").First(&existingNotification).Error; err == nil {
			// 已经发送过通知，跳过
			continue
		}

		// 创建即将到期通知
		req := &CreateNotificationRequest{
			UserID: task.CreatedBy,
			Type:   "due_soon",
			Title:  "任务即将到期提醒",
			Message: fmt.Sprintf("任务「%s」将在 %s 到期，请及时处理",
				task.Title, task.DueDate.Format("2006-01-02 15:04")),
			Data: map[string]interface{}{
				"task_id":    task.ID,
				"task_title": task.Title,
				"due_date":   task.DueDate.Format(time.RFC3339),
				"hours_left": int(time.Until(*task.DueDate).Hours()),
			},
		}

		if _, err := nm.notificationService.CreateNotification(*req); err != nil {
			fmt.Printf("Failed to create due soon notification for task %d: %v\n", task.ID, err)
		}
	}

	// 查找已逾期的任务
	var overdueTasks []models.Todo
	if err := db.Where("due_date < ? AND status != ?", now, "completed").Find(&overdueTasks).Error; err != nil {
		return fmt.Errorf("failed to query overdue tasks: %v", err)
	}

	for _, task := range overdueTasks {
		// 检查是否已经发送过逾期通知
		var existingNotification models.Notification
		if err := db.Where("user_id = ? AND type = ? AND data LIKE ?",
			task.CreatedBy, "overdue", "%"+task.Title+"%").First(&existingNotification).Error; err == nil {
			// 已经发送过通知，跳过
			continue
		}

		// 创建逾期通知
		req := &CreateNotificationRequest{
			UserID: task.CreatedBy,
			Type:   "overdue",
			Title:  "任务已逾期",
			Message: fmt.Sprintf("任务「%s」已逾期 %d 天，请尽快完成",
				task.Title, int(time.Since(*task.DueDate).Hours()/24)),
			Data: map[string]interface{}{
				"task_id":      task.ID,
				"task_title":   task.Title,
				"due_date":     task.DueDate.Format(time.RFC3339),
				"days_overdue": int(time.Since(*task.DueDate).Hours() / 24),
			},
		}

		if _, err := nm.notificationService.CreateNotification(*req); err != nil {
			fmt.Printf("Failed to create overdue notification for task %d: %v\n", task.ID, err)
		}
	}

	return nil
}

// CreateTaskCompletedNotification 创建任务完成通知
func (nm *NotificationManager) CreateTaskCompletedNotification(task *models.Todo) error {
	req := &CreateNotificationRequest{
		UserID:  task.CreatedBy,
		Type:    "completed",
		Title:   "任务已完成",
		Message: fmt.Sprintf("恭喜！任务「%s」已完成", task.Title),
		Data: map[string]interface{}{
			"task_id":      task.ID,
			"task_title":   task.Title,
			"completed_at": task.CompletedAt.Format(time.RFC3339),
		},
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}

// CreateArticlePublishedNotification 创建文章发布通知
func (nm *NotificationManager) CreateArticlePublishedNotification(article *models.Article) error {
	req := &CreateNotificationRequest{
		UserID:  article.CreatedBy,
		Type:    "article_published",
		Title:   "文章已发布",
		Message: fmt.Sprintf("文章「%s」已成功发布", article.Title),
		Data: map[string]interface{}{
			"article_id":    article.ID,
			"article_title": article.Title,
			"published_at":  time.Now().Format(time.RFC3339),
		},
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}

// CreateSystemNotification 创建系统通知
func (nm *NotificationManager) CreateSystemNotification(userID uint, title, message string, data map[string]interface{}) error {
	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    "system",
		Title:   title,
		Message: message,
		Data:    data,
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}

// CreateWelcomeNotification 创建欢迎通知
func (nm *NotificationManager) CreateWelcomeNotification(userID uint, username string) error {
	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    "welcome",
		Title:   "欢迎使用任务管理系统",
		Message: fmt.Sprintf("欢迎 %s！开始创建您的第一个任务吧", username),
		Data: map[string]interface{}{
			"username":   username,
			"welcome_at": time.Now().Format(time.RFC3339),
			"tips": []string{
				"创建任务时设置合理的截止日期",
				"定期查看任务进度",
				"及时处理即将到期的任务",
			},
		},
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}

// RunNotificationChecks 运行通知检查（定时任务）
func (nm *NotificationManager) RunNotificationChecks() {
	// 检查任务到期通知
	if err := nm.CheckTaskDueNotifications(); err != nil {
		fmt.Printf("Failed to check task due notifications: %v\n", err)
	}

	// 这里可以添加其他定时检查逻辑
	// 比如每日总结、周报等
}

// CreateDailySummaryNotification 创建每日总结通知
func (nm *NotificationManager) CreateDailySummaryNotification(userID uint) error {
	db := database.GetDB()

	// 获取今日完成的任务数
	var completedToday int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND status = ? AND DATE(completed_at) = DATE(?)",
			userID, "completed", today).Count(&completedToday).Error; err != nil {
		return fmt.Errorf("failed to count completed tasks: %v", err)
	}

	// 获取待处理的任务数
	var pendingTasks int64
	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND status = ?", userID, "pending").Count(&pendingTasks).Error; err != nil {
		return fmt.Errorf("failed to count pending tasks: %v", err)
	}

	// 获取即将到期的任务数
	dueSoon := time.Now().Add(24 * time.Hour)
	var dueSoonTasks int64
	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND due_date <= ? AND status != ?",
			userID, dueSoon, "completed").Count(&dueSoonTasks).Error; err != nil {
		return fmt.Errorf("failed to count due soon tasks: %v", err)
	}

	message := fmt.Sprintf("今日完成 %d 个任务，还有 %d 个待处理任务", completedToday, pendingTasks)
	if dueSoonTasks > 0 {
		message += fmt.Sprintf("，其中 %d 个即将到期", dueSoonTasks)
	}

	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    "daily_summary",
		Title:   "今日任务总结",
		Message: message,
		Data: map[string]interface{}{
			"completed_today": completedToday,
			"pending_tasks":   pendingTasks,
			"due_soon_tasks":  dueSoonTasks,
			"summary_date":    today.Format("2006-01-02"),
		},
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}

// CreateWeeklyReportNotification 创建周报通知
func (nm *NotificationManager) CreateWeeklyReportNotification(userID uint) error {
	db := database.GetDB()

	// 获取本周完成的任务数
	var completedThisWeek int64
	weekStart := time.Now().Truncate(24 * time.Hour)
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.Add(-24 * time.Hour)
	}

	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND status = ? AND completed_at >= ?",
			userID, "completed", weekStart).Count(&completedThisWeek).Error; err != nil {
		return fmt.Errorf("failed to count weekly completed tasks: %v", err)
	}

	// 获取本周创建的任务数
	var createdThisWeek int64
	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND created_at >= ?", userID, weekStart).Count(&createdThisWeek).Error; err != nil {
		return fmt.Errorf("failed to count weekly created tasks: %v", err)
	}

	completionRate := 0.0
	if createdThisWeek > 0 {
		completionRate = float64(completedThisWeek) / float64(createdThisWeek) * 100
	}

	message := fmt.Sprintf("本周完成 %d 个任务，完成率 %.1f%%", completedThisWeek, completionRate)

	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    "weekly_report",
		Title:   "本周任务报告",
		Message: message,
		Data: map[string]interface{}{
			"completed_this_week": completedThisWeek,
			"created_this_week":   createdThisWeek,
			"completion_rate":     completionRate,
			"week_start":          weekStart.Format("2006-01-02"),
			"week_end":            time.Now().Format("2006-01-02"),
		},
	}

	_, err := nm.notificationService.CreateNotification(*req)
	return err
}
