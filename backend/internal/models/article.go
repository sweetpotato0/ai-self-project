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
	User      User          `json:"user" gorm:"foreignKey:CreatedBy"`
	UserLikes []ArticleLike `json:"-" gorm:"foreignKey:ArticleID"`

	// 虚拟字段 - 不存储在数据库中
	IsLikedByUser bool `json:"is_liked_by_user" gorm:"-"`
}

// ArticleLike 文章点赞关系模型
type ArticleLike struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ArticleID uint           `json:"article_id" gorm:"not null;index:idx_article_user,unique"`
	UserID    uint           `json:"user_id" gorm:"not null;index:idx_article_user,unique"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联
	Article Article `json:"article" gorm:"foreignKey:ArticleID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}

// 添加唯一索引
func (ArticleLike) TableName() string {
	return "article_likes"
}
