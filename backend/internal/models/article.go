package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Title      string         `json:"title" gorm:"not null"`
	Content    string         `json:"content" gorm:"type:text"`
	Summary    string         `json:"summary" gorm:"type:text"`
	CoverImage string         `json:"cover_image"`
	Status     string         `json:"status" gorm:"default:'draft'"` // draft, published, archived
	Tags       string         `json:"tags" gorm:"type:text"`         // JSON格式存储标签
	ViewCount  int            `json:"view_count" gorm:"default:0"`
	LikeCount  int            `json:"like_count" gorm:"default:0"`
	CreatedBy  uint           `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联
	User User `json:"user" gorm:"foreignKey:CreatedBy"`
}
