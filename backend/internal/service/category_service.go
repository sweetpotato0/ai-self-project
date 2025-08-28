package service

import (
	"errors"
	"fmt"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

// CategoryService 分类服务
type CategoryService struct{
	logger logger.LoggerInterface
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(logger logger.LoggerInterface) *CategoryService {
	return &CategoryService{
		logger: logger,
	}
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Sort        int    `json:"sort"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	Sort        int    `json:"sort"`
}

// getDB 获取数据库连接
func (s *CategoryService) getDB() *gorm.DB {
	db := database.GetDB()
	if db == nil {
		panic("Database connection is nil")
	}
	return db
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(category *models.Category) error {
	db := s.getDB()

	// 验证父分类是否存在
	if category.ParentID != nil {
		var parent models.Category
		if err := db.First(&parent, *category.ParentID).Error; err != nil {
			return fmt.Errorf("父分类不存在: %w", err)
		}
		// 检查是否创建循环引用
		if err := s.checkCircularReference(*category.ParentID, category.CreatedBy); err != nil {
			return err
		}
	}

	return db.Create(category).Error
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(id uint, userID uint) (*models.Category, error) {
	db := s.getDB()
	var category models.Category

	if err := db.Where("id = ? AND created_by = ?", id, userID).
		Preload("Parent").
		Preload("Children").
		First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// GetCategoryTree 获取用户的分类树
func (s *CategoryService) GetCategoryTree(userID uint) ([]models.Category, error) {
	db := s.getDB()
	var categories []models.Category

	// 获取所有顶级分类（ParentID为null）
	if err := db.Where("created_by = ? AND parent_id IS NULL", userID).
		Preload("Children").
		Order("sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	// 递归加载子分类
	for i := range categories {
		if err := s.loadChildrenRecursively(&categories[i], userID); err != nil {
			return nil, err
		}
	}

	return categories, nil
}

// GetAllCategories 获取用户的所有分类（扁平结构）
func (s *CategoryService) GetAllCategories(userID uint) ([]models.Category, error) {
	db := s.getDB()
	var categories []models.Category

	if err := db.Where("created_by = ?", userID).
		Preload("Parent").
		Order("level ASC, sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, userID uint, updates map[string]interface{}) error {
	db := s.getDB()

	// 检查分类是否存在且属于当前用户
	var category models.Category
	if err := db.Where("id = ? AND created_by = ?", id, userID).First(&category).Error; err != nil {
		return fmt.Errorf("分类不存在或无权限: %w", err)
	}

	// 如果要更新父分类，检查循环引用
	if parentID, exists := updates["parent_id"]; exists {
		if parentID != nil {
			parentIDUint := parentID.(uint)
			if parentIDUint == id {
				return errors.New("不能将自己设为父分类")
			}
			if err := s.checkCircularReference(parentIDUint, userID); err != nil {
				return err
			}
		}
	}

	return db.Model(&category).Updates(updates).Error
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint, userID uint) error {
	db := s.getDB()

	// 检查分类是否存在且属于当前用户
	var category models.Category
	if err := db.Where("id = ? AND created_by = ?", id, userID).First(&category).Error; err != nil {
		return fmt.Errorf("分类不存在或无权限: %w", err)
	}

	// 检查是否有子分类
	var childCount int64
	if err := db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("请先删除子分类")
	}

	// 检查是否有关联的任务
	var todoCount int64
	if err := db.Model(&models.Todo{}).Where("category_id = ?", id).Count(&todoCount).Error; err != nil {
		return err
	}
	if todoCount > 0 {
		return errors.New("该分类下还有任务，无法删除")
	}

	return db.Delete(&category).Error
}

// loadChildrenRecursively 递归加载子分类
func (s *CategoryService) loadChildrenRecursively(category *models.Category, userID uint) error {
	db := s.getDB()

	if err := db.Where("parent_id = ? AND created_by = ?", category.ID, userID).
		Preload("Children").
		Order("sort ASC, id ASC").
		Find(&category.Children).Error; err != nil {
		return err
	}

	// 递归加载每个子分类的子分类
	for i := range category.Children {
		if err := s.loadChildrenRecursively(&category.Children[i], userID); err != nil {
			return err
		}
	}

	return nil
}

// checkCircularReference 检查循环引用
func (s *CategoryService) checkCircularReference(parentID uint, userID uint) error {
	db := s.getDB()

	// 检查父分类是否存在
	var parent models.Category
	if err := db.Where("id = ? AND created_by = ?", parentID, userID).First(&parent).Error; err != nil {
		return fmt.Errorf("父分类不存在: %w", err)
	}

	// 检查是否形成循环引用（通过递归检查父分类的父分类）
	visited := make(map[uint]bool)
	currentID := parentID

	for currentID != 0 {
		if visited[currentID] {
			return errors.New("检测到循环引用")
		}
		visited[currentID] = true

		var current models.Category
		if err := db.Where("id = ? AND created_by = ?", currentID, userID).First(&current).Error; err != nil {
			break
		}

		if current.ParentID == nil {
			break
		}
		currentID = *current.ParentID
	}

	return nil
}

// GetCategoryStats 获取分类统计信息
func (s *CategoryService) GetCategoryStats(userID uint) (map[string]interface{}, error) {
	db := s.getDB()

	stats := make(map[string]interface{})

	// 总分类数
	var totalCount int64
	if err := db.Model(&models.Category{}).Where("created_by = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total"] = totalCount

	// 顶级分类数
	var rootCount int64
	if err := db.Model(&models.Category{}).Where("created_by = ? AND parent_id IS NULL", userID).Count(&rootCount).Error; err != nil {
		return nil, err
	}
	stats["root_count"] = rootCount

	// 子分类数
	var childCount int64
	if err := db.Model(&models.Category{}).Where("created_by = ? AND parent_id IS NOT NULL", userID).Count(&childCount).Error; err != nil {
		return nil, err
	}
	stats["child_count"] = childCount

	// 有任务的分类数
	var activeCount int64
	if err := db.Model(&models.Category{}).
		Joins("JOIN todos ON categories.id = todos.category_id").
		Where("categories.created_by = ?", userID).
		Distinct("categories.id").
		Count(&activeCount).Error; err != nil {
		return nil, err
	}
	stats["active_count"] = activeCount

	return stats, nil
}
