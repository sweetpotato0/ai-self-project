package models

import (
	"time"

	"gorm.io/gorm"
)

// TodoCategory TODO分类
type TodoCategory struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex"`
	Color       string         `json:"color" gorm:"default:'#3B82F6'"` // 分类颜色
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TodoPriority TODO优先级
type TodoPriority struct {
	ID    uint   `json:"id" gorm:"primarykey"`
	Name  string `json:"name" gorm:"not null;uniqueIndex"`
	Level int    `json:"level" gorm:"not null;uniqueIndex"` // 优先级等级 1-5
	Color string `json:"color" gorm:"not null"`             // 优先级颜色
}

// Todo TODO项目
type Todo struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	Title          string         `json:"title" gorm:"not null"`
	Description    string         `json:"description"`
	Status         string         `json:"status" gorm:"default:'pending'"` // pending, in_progress, completed, cancelled
	PriorityID     uint           `json:"priority_id" gorm:"not null"`
	CategoryID     *uint          `json:"category_id"`     // 改为可选，因为可能没有分类
	StartDate      *time.Time     `json:"start_date"`      // 开始时间
	DueDate        *time.Time     `json:"due_date"`        // 截止时间
	CompletedAt    *time.Time     `json:"completed_at"`    // 完成时间
	EstimatedHours float64        `json:"estimated_hours"` // 预估工时（小时）
	ActualHours    float64        `json:"actual_hours"`    // 实际工时（小时）
	CreatedBy      uint           `json:"created_by" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Priority TodoPriority `json:"priority" gorm:"foreignKey:PriorityID"`
	Category *Category    `json:"category" gorm:"foreignKey:CategoryID"`
}

// TodoNotification TODO通知
type TodoNotification struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	TodoID    uint           `json:"todo_id" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"` // due_soon, overdue, completed
	Message   string         `json:"message" gorm:"not null"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Todo Todo `json:"todo" gorm:"foreignKey:TodoID"`
}

// TableName 指定表名
func (TodoCategory) TableName() string {
	return "todo_categories"
}

func (TodoPriority) TableName() string {
	return "todo_priorities"
}

func (Todo) TableName() string {
	return "todos"
}

func (TodoNotification) TableName() string {
	return "todo_notifications"
}
