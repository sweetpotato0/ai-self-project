package service

import (
	"encoding/json"
	"fmt"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RealtimeNotificationService 实时通知服务
type RealtimeNotificationService struct {
	clients    map[uint]*websocket.Conn // 用户ID -> WebSocket连接
	clientsMux sync.RWMutex
	stopChan   chan bool
}

// NewRealtimeNotificationService 创建实时通知服务
func NewRealtimeNotificationService() *RealtimeNotificationService {
	return &RealtimeNotificationService{
		clients:  make(map[uint]*websocket.Conn),
		stopChan: make(chan bool),
	}
}

// AddClient 添加WebSocket客户端
func (s *RealtimeNotificationService) AddClient(userID uint, conn *websocket.Conn) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	// 关闭旧连接
	if oldConn, exists := s.clients[userID]; exists {
		oldConn.Close()
	}

	s.clients[userID] = conn
	log.Printf("用户 %d 已连接实时通知服务", userID)
}

// RemoveClient 移除WebSocket客户端
func (s *RealtimeNotificationService) RemoveClient(userID uint) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	if conn, exists := s.clients[userID]; exists {
		conn.Close()
		delete(s.clients, userID)
		log.Printf("用户 %d 已断开实时通知服务", userID)
	}
}

// SendNotification 发送实时通知
func (s *RealtimeNotificationService) SendNotification(userID uint, notification *models.Notification) error {
	s.clientsMux.RLock()
	conn, exists := s.clients[userID]
	s.clientsMux.RUnlock()

	if !exists {
		return fmt.Errorf("用户 %d 未连接", userID)
	}

	message := map[string]interface{}{
		"type": "notification",
		"data": notification,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化通知失败: %v", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// BroadcastNotification 广播通知给所有在线用户
func (s *RealtimeNotificationService) BroadcastNotification(notification *models.Notification) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	message := map[string]interface{}{
		"type": "notification",
		"data": notification,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化广播通知失败: %v", err)
		return
	}

	for userID, conn := range s.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("发送通知给用户 %d 失败: %v", userID, err)
			// 移除失效连接
			go s.RemoveClient(userID)
		}
	}
}

// StartTaskMonitoring 启动任务监控
func (s *RealtimeNotificationService) StartTaskMonitoring() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer ticker.Stop()

	log.Println("启动任务监控服务")

	for {
		select {
		case <-ticker.C:
			s.checkOverdueTasks()
			s.checkDueSoonTasks()
		case <-s.stopChan:
			log.Println("停止任务监控服务")
			return
		}
	}
}

// StopTaskMonitoring 停止任务监控
func (s *RealtimeNotificationService) StopTaskMonitoring() {
	close(s.stopChan)
}

// checkOverdueTasks 检查逾期任务
func (s *RealtimeNotificationService) checkOverdueTasks() {
	db := database.GetDB()
	now := time.Now()

	var overdueTasks []models.Todo
	if err := db.Where("due_date < ? AND status != ?", now, "completed").Find(&overdueTasks).Error; err != nil {
		log.Printf("查询逾期任务失败: %v", err)
		return
	}

	notificationService := NewNotificationService()

	for _, task := range overdueTasks {
		// 检查是否已经发送过逾期通知（避免重复通知）
		var existingNotification models.Notification
		if err := db.Where("user_id = ? AND type = ? AND data LIKE ?",
			task.CreatedBy, "overdue", "%"+task.Title+"%").First(&existingNotification).Error; err == nil {
			// 已经发送过通知，跳过
			continue
		}

		// 创建逾期通知
		notification, err := notificationService.CreateNotification(&CreateNotificationRequest{
			UserID:  task.CreatedBy,
			Type:    "overdue",
			Title:   "任务已逾期",
			Message: fmt.Sprintf("任务\"%s\"已逾期，请尽快完成", task.Title),
			Data: map[string]interface{}{
				"task_id":    task.ID,
				"task_title": task.Title,
				"due_date":   task.DueDate.Format(time.RFC3339),
				"overdue_at": now.Format(time.RFC3339),
			},
		})

		if err != nil {
			log.Printf("创建逾期通知失败: %v", err)
			continue
		}

		// 发送实时通知
		if err := s.SendNotification(task.CreatedBy, notification); err != nil {
			log.Printf("发送逾期实时通知失败: %v", err)
		}

		log.Printf("已发送逾期通知给用户 %d，任务: %s", task.CreatedBy, task.Title)
	}
}

// checkDueSoonTasks 检查即将到期的任务
func (s *RealtimeNotificationService) checkDueSoonTasks() {
	db := database.GetDB()
	now := time.Now()
	dueSoonTime := now.Add(1 * time.Hour) // 1小时内到期

	var dueSoonTasks []models.Todo
	if err := db.Where("due_date BETWEEN ? AND ? AND status NOT IN (?)",
		now, dueSoonTime, []string{"completed", "cancelled"}).Find(&dueSoonTasks).Error; err != nil {
		log.Printf("查询即将到期任务失败: %v", err)
		return
	}

	notificationService := NewNotificationService()

	for _, task := range dueSoonTasks {
		// 检查是否已经发送过即将到期通知
		var existingNotification models.Notification
		if err := db.Where("user_id = ? AND type = ? AND data LIKE ?",
			task.CreatedBy, "due_soon", "%"+task.Title+"%").First(&existingNotification).Error; err == nil {
			// 已经发送过通知，跳过
			continue
		}

		// 计算剩余时间
		remainingTime := task.DueDate.Sub(now)
		remainingMinutes := int(remainingTime.Minutes())

		// 创建即将到期通知
		notification, err := notificationService.CreateNotification(&CreateNotificationRequest{
			UserID:  task.CreatedBy,
			Type:    "due_soon",
			Title:   "任务即将到期",
			Message: fmt.Sprintf("任务\"%s\"将在 %d 分钟后到期", task.Title, remainingMinutes),
			Data: map[string]interface{}{
				"task_id":           task.ID,
				"task_title":        task.Title,
				"due_date":          task.DueDate.Format(time.RFC3339),
				"remaining_minutes": remainingMinutes,
			},
		})

		if err != nil {
			log.Printf("创建即将到期通知失败: %v", err)
			continue
		}

		// 发送实时通知
		if err := s.SendNotification(task.CreatedBy, notification); err != nil {
			log.Printf("发送即将到期实时通知失败: %v", err)
		}

		log.Printf("已发送即将到期通知给用户 %d，任务: %s，剩余 %d 分钟", task.CreatedBy, task.Title, remainingMinutes)
	}
}

// GetOnlineUsers 获取在线用户数量
func (s *RealtimeNotificationService) GetOnlineUsers() int {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()
	return len(s.clients)
}
