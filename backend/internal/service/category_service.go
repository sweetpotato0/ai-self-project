package service

import (
	"context"
	"time"

	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/errors"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

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

// OptimizedCategoryService 优化的分类服务
type OptimizedCategoryService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

// NewOptimizedCategoryService 创建优化的分类服务实例
func NewOptimizedCategoryService(db *gorm.DB, logger logger.LoggerInterface) *OptimizedCategoryService {
	return &OptimizedCategoryService{
		db:     db,
		logger: logger,
	}
}

// CategoryWithChildren 带子分类的分类结构
type CategoryWithChildren struct {
	models.Category
	Children []*CategoryWithChildren `json:"children,omitempty"`
}

// GetAllCategories 获取用户的所有分类（扁平结构）
func (s *OptimizedCategoryService) GetAllCategories(ctx context.Context, userID uint) ([]models.Category, error) {
	var categories []models.Category

	if err := s.db.WithContext(ctx).Where("created_by = ?", userID).
		Preload("Parent").
		Order("level ASC, sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// CreateCategory 创建分类（带事务保护和优化验证）
func (s *OptimizedCategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	// 输入验证
	if category == nil {
		return errors.NewValidationError("分类信息不能为空", nil)
	}

	if category.Name == "" {
		return errors.NewValidationError("分类名称不能为空", nil)
	}

	// 使用事务确保数据一致性
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查名称是否重复
		var existing models.Category
		if err := tx.Where("name = ? AND created_by = ? AND deleted_at IS NULL",
			category.Name, category.CreatedBy).First(&existing).Error; err == nil {
			return errors.NewConflictError("分类名称已存在")
		} else if err != gorm.ErrRecordNotFound {
			return errors.NewDatabaseError("检查分类名称时发生错误", err)
		}

		// 验证并检查父分类（优化版）
		if category.ParentID != nil {
			if err := s.validateParentCategoryTx(ctx, tx, *category.ParentID, category.CreatedBy, category.ID); err != nil {
				return err
			}
		}

		// 创建分类
		if err := tx.Create(category).Error; err != nil {
			return errors.NewDatabaseError("创建分类失败", err)
		}

		s.logger.WithFields(map[string]any{
			"category_id":   category.ID,
			"category_name": category.Name,
			"user_id":       category.CreatedBy,
		}).Info("分类创建成功")

		return nil
	})
}

// validateParentCategoryTx 在事务中验证父分类（优化版，单查询获取所有层级）
func (s *OptimizedCategoryService) validateParentCategoryTx(ctx context.Context, tx *gorm.DB, parentID, userID, excludeID uint) error {
	// 一次性获取用户的所有分类，构建层级关系
	var allCategories []models.Category
	if err := tx.WithContext(ctx).Where("created_by = ? AND deleted_at IS NULL", userID).
		Find(&allCategories).Error; err != nil {
		return errors.NewDatabaseError("获取分类信息失败", err)
	}

	// 构建分类映射表
	categoryMap := make(map[uint]*models.Category)
	for i := range allCategories {
		categoryMap[allCategories[i].ID] = &allCategories[i]
	}

	// 检查父分类是否存在
	_, exists := categoryMap[parentID]
	if !exists {
		return errors.NewNotFoundError("父分类")
	}

	// 使用内存中的映射检查循环引用
	visited := make(map[uint]bool)
	currentID := parentID

	for currentID != 0 {
		if currentID == excludeID {
			return errors.NewValidationError("不能将分类设置为自己的子分类", nil)
		}

		if visited[currentID] {
			return errors.NewValidationError("检测到循环引用", nil)
		}
		visited[currentID] = true

		current, exists := categoryMap[currentID]
		if !exists || current.ParentID == nil {
			break
		}
		currentID = *current.ParentID
	}

	return nil
}

// GetCategoryByID 获取分类（优化预加载）
func (s *OptimizedCategoryService) GetCategoryByID(ctx context.Context, id, userID uint) (*models.Category, error) {
	var category models.Category

	// 使用索引优化的查询
	if err := s.db.WithContext(ctx).
		Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, userID).
		First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("分类")
		}
		return nil, errors.NewDatabaseError("查询分类失败", err)
	}

	return &category, nil
}

// GetCategoryTree 获取分类树（优化版，单查询构建树形结构）
func (s *OptimizedCategoryService) GetCategoryTree(ctx context.Context, userID uint) ([]*CategoryWithChildren, error) {
	// 单次查询获取所有分类
	var categories []models.Category
	if err := s.db.WithContext(ctx).
		Where("created_by = ? AND deleted_at IS NULL", userID).
		Order("sort ASC, created_at ASC").
		Find(&categories).Error; err != nil {
		return nil, errors.NewDatabaseError("获取分类列表失败", err)
	}

	return s.buildCategoryTree(categories), nil
}

// buildCategoryTree 构建分类树（内存中操作，避免递归查询）
func (s *OptimizedCategoryService) buildCategoryTree(categories []models.Category) []*CategoryWithChildren {
	// 创建映射表和结果切片
	categoryMap := make(map[uint]*CategoryWithChildren)
	var rootCategories []*CategoryWithChildren

	// 第一遍：创建所有节点
	for i := range categories {
		node := &CategoryWithChildren{
			Category: categories[i],
			Children: make([]*CategoryWithChildren, 0),
		}
		categoryMap[categories[i].ID] = node
	}

	// 第二遍：构建树形关系
	for _, node := range categoryMap {
		if node.ParentID == nil {
			// 根节点
			rootCategories = append(rootCategories, node)
		} else {
			// 子节点，添加到父节点的children中
			if parent, exists := categoryMap[*node.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return rootCategories
}

// UpdateCategory 更新分类（优化事务处理）
func (s *OptimizedCategoryService) UpdateCategory(ctx context.Context, id, userID uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.NewValidationError("没有要更新的字段", nil)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查分类是否存在
		var existing models.Category
		if err := tx.Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, userID).
			First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewNotFoundError("分类")
			}
			return errors.NewDatabaseError("查询分类失败", err)
		}

		// 如果要更新父分类，需要验证
		if parentID, exists := updates["parent_id"]; exists {
			if parentID != nil {
				parentIDUint := parentID.(uint)
				if err := s.validateParentCategoryTx(ctx, tx, parentIDUint, userID, id); err != nil {
					return err
				}
			}
		}

		// 如果要更新名称，检查重复
		if name, exists := updates["name"]; exists {
			var duplicate models.Category
			if err := tx.WithContext(ctx).Where("name = ? AND created_by = ? AND id != ? AND deleted_at IS NULL",
				name, userID, id).First(&duplicate).Error; err == nil {
				return errors.NewConflictError("分类名称已存在")
			} else if err != gorm.ErrRecordNotFound {
				return errors.NewDatabaseError("检查分类名称时发生错误", err)
			}
		}

		// 添加更新时间
		updates["updated_at"] = time.Now()

		// 执行更新
		if err := tx.WithContext(ctx).Model(&existing).Updates(updates).Error; err != nil {
			return errors.NewDatabaseError("更新分类失败", err)
		}

		s.logger.WithFields(map[string]any{
			"category_id": id,
			"user_id":     userID,
			"updates":     updates,
		}).Info("分类更新成功")

		return nil
	})
}

// DeleteCategory 删除分类（软删除，处理子分类）
func (s *OptimizedCategoryService) DeleteCategory(ctx context.Context, id, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查分类是否存在
		var category models.Category
		if err := tx.WithContext(ctx).Where("id = ? AND created_by = ? AND deleted_at IS NULL", id, userID).
			First(&category).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewNotFoundError("分类")
			}
			return errors.NewDatabaseError("查询分类失败", err)
		}

		// 检查是否有子分类
		var childCount int64
		if err := tx.WithContext(ctx).Model(&models.Category{}).
			Where("parent_id = ? AND created_by = ? AND deleted_at IS NULL", id, userID).
			Count(&childCount).Error; err != nil {
			return errors.NewDatabaseError("检查子分类失败", err)
		}

		if childCount > 0 {
			return errors.NewValidationError("存在子分类，无法删除", map[string]interface{}{
				"child_count": childCount,
			})
		}

		// 检查是否有关联的内容（这里可以根据业务需求添加）
		// 例如检查是否有文章使用此分类等

		// 执行软删除
		if err := tx.WithContext(ctx).Delete(&category).Error; err != nil {
			return errors.NewDatabaseError("删除分类失败", err)
		}

		s.logger.WithFields(map[string]any{
			"category_id":   id,
			"category_name": category.Name,
			"user_id":       userID,
		}).Info("分类删除成功")

		return nil
	})
}

// GetCategoryStats 获取分类统计信息（优化查询）
func (s *OptimizedCategoryService) GetCategoryStats(ctx context.Context, userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 使用聚合查询一次性获取统计信息
	var totalCount, rootCount int64

	// 总分类数
	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("created_by = ? AND deleted_at IS NULL", userID).
		Count(&totalCount).Error; err != nil {
		return nil, errors.NewDatabaseError("获取分类统计失败", err)
	}

	// 根分类数
	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("created_by = ? AND parent_id IS NULL AND deleted_at IS NULL", userID).
		Count(&rootCount).Error; err != nil {
		return nil, errors.NewDatabaseError("获取根分类统计失败", err)
	}

	// 最深层级（使用递归CTE或应用层计算）
	maxDepth := s.calculateMaxDepth(ctx, userID)

	stats["total_categories"] = totalCount
	stats["root_categories"] = rootCount
	stats["child_categories"] = totalCount - rootCount
	stats["max_depth"] = maxDepth

	return stats, nil
}

// calculateMaxDepth 计算分类的最大深度
func (s *OptimizedCategoryService) calculateMaxDepth(ctx context.Context, userID uint) int {
	// 获取所有分类
	var categories []models.Category
	if err := s.db.WithContext(ctx).
		Where("created_by = ? AND deleted_at IS NULL", userID).
		Find(&categories).Error; err != nil {
		return 0
	}

	if len(categories) == 0 {
		return 0
	}

	// 构建映射表
	categoryMap := make(map[uint]*models.Category)
	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	maxDepth := 0

	// 对每个根分类计算深度
	for _, category := range categories {
		if category.ParentID == nil {
			depth := s.calculateDepthRecursive(categoryMap, category.ID, 1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	return maxDepth
}

// calculateDepthRecursive 递归计算深度（在内存中操作）
func (s *OptimizedCategoryService) calculateDepthRecursive(categoryMap map[uint]*models.Category, currentID uint, currentDepth int) int {
	maxDepth := currentDepth

	// 查找所有子分类
	for _, category := range categoryMap {
		if category.ParentID != nil && *category.ParentID == currentID {
			depth := s.calculateDepthRecursive(categoryMap, category.ID, currentDepth+1)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
	}

	return maxDepth
}
