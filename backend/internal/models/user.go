package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Nickname  string         `json:"nickname" gorm:""`
	Password  string         `json:"-" gorm:"not null"` // 密码不返回给前端
	Role      string         `json:"role" gorm:"default:user"`        // 用户角色: user, admin
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRole 定义用户角色常量
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
