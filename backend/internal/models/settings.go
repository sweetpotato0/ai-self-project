package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSettings struct {
	ID                     uint           `json:"id" gorm:"primarykey"`
	UserID                 uint           `json:"user_id" gorm:"not null;index"`
	DueReminder            bool           `json:"due_reminder" gorm:"default:true"`
	CompletionNotification bool           `json:"completion_notification" gorm:"default:true"`
	NewTaskNotification    bool           `json:"new_task_notification" gorm:"default:true"`
	EmailNotification      bool           `json:"email_notification" gorm:"default:false"`
	Theme                  string         `json:"theme" gorm:"default:'light'"`
	Language               string         `json:"language" gorm:"default:'zh-CN'"`
	Timezone               string         `json:"timezone" gorm:"default:'Asia/Shanghai'"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UpdateNotificationSettingsRequest struct {
	DueReminder            bool `json:"due_reminder"`
	CompletionNotification bool `json:"completion_notification"`
	NewTaskNotification    bool `json:"new_task_notification"`
	EmailNotification      bool `json:"email_notification"`
}

type UpdateInterfaceSettingsRequest struct {
	Theme    string `json:"theme" binding:"required,oneof=light dark auto"`
	Language string `json:"language" binding:"required,oneof=zh-CN en-US"`
	Timezone string `json:"timezone" binding:"required"`
}
