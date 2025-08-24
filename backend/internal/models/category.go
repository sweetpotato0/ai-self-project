package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 任务分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`          // 分类名称
	Description string         `json:"description" gorm:"size:500"`            // 分类描述
	Color       string         `json:"color" gorm:"size:20;default:'#409eff'"` // 分类颜色
	Icon        string         `json:"icon" gorm:"size:50"`                    // 分类图标
	ParentID    *uint          `json:"parent_id" gorm:"index"`                 // 父分类ID
	Level       int            `json:"level" gorm:"default:1"`                 // 分类层级
	Sort        int            `json:"sort" gorm:"default:0"`                  // 排序
	IsActive    bool           `json:"is_active" gorm:"default:true"`          // 是否激活
	CreatedBy   uint           `json:"created_by" gorm:"not null;index"`       // 创建者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// 关联关系
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`   // 父分类
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"` // 子分类
	Todos    []Todo     `json:"todos,omitempty" gorm:"foreignKey:CategoryID"`  // 关联的任务
	User     User       `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`    // 创建者
}

// TableName 指定表名
func (Category) TableName() string {
	return "category"
}

// BeforeCreate GORM钩子：创建前设置默认值
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Color == "" {
		c.Color = "#409eff"
	}
	if c.Level == 0 {
		c.Level = 1
	}
	return nil
}

// BeforeSave GORM钩子：保存前计算层级
func (c *Category) BeforeSave(tx *gorm.DB) error {
	if c.ParentID != nil {
		var parent Category
		if err := tx.First(&parent, *c.ParentID).Error; err == nil {
			c.Level = parent.Level + 1
		}
	} else {
		c.Level = 1
	}
	return nil
}
