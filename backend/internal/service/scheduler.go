package service

import (
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	notificationManager *NotificationManager
	stopChan            chan bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		notificationManager: NewNotificationManager(),
		stopChan:            make(chan bool),
	}
}

// Start 启动定时任务
func (s *Scheduler) Start() {
	go s.runNotificationChecks()
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// runNotificationChecks 运行通知检查定时任务
func (s *Scheduler) runNotificationChecks() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.notificationManager.RunNotificationChecks()
		case <-s.stopChan:
			return
		}
	}
}
