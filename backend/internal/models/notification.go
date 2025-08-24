package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知模型
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"` // due_soon, overdue, completed, system
	Title     string         `json:"title" gorm:"not null"`
	Message   string         `json:"message" gorm:"not null"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	Data      string         `json:"data" gorm:"type:text"` // JSON格式的额外数据
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}
