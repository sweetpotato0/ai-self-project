package service

import (
	"errors"
	"fmt"
	"time"

	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"

	"gorm.io/gorm"
)

type TodoService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

func NewTodoService(db *gorm.DB, logger logger.LoggerInterface) *TodoService {
	return &TodoService{
		db:     db,
		logger: logger,
	}
}

// CreateTodoRequest 创建TODO请求
type CreateTodoRequest struct {
	Title          string     `json:"title" binding:"required"`
	Description    string     `json:"description"`
	PriorityID     uint       `json:"priority_id" binding:"required"`
	CategoryID     *uint      `json:"category_id"`     // 改为可选
	StartDate      *time.Time `json:"start_date"`      // 开始时间
	DueDate        *time.Time `json:"due_date"`        // 截止时间
	EstimatedHours float64    `json:"estimated_hours"` // 预估工时
}

// UpdateTodoRequest 更新TODO请求
type UpdateTodoRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	PriorityID     uint       `json:"priority_id"`
	CategoryID     *uint      `json:"category_id"`     // 改为可选
	StartDate      *time.Time `json:"start_date"`      // 开始时间
	DueDate        *time.Time `json:"due_date"`        // 截止时间
	EstimatedHours float64    `json:"estimated_hours"` // 预估工时
	ActualHours    float64    `json:"actual_hours"`    // 实际工时
}

type TodoFilter struct {
	Status     string     `json:"status"`
	Priority   string     `json:"priority"`
	CategoryID *uint      `json:"category_id"`
	DueDate    *time.Time `json:"due_date"`
	Overdue    *bool      `json:"overdue"`
	Search     string     `json:"search"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	SortBy     string     `json:"sort_by"`
	SortOrder  string     `json:"sort_order"`
}

type PaginatedTodos struct {
	Todos      []*models.Todo `json:"todos"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// CreateTodo 创建TODO
func (s *TodoService) CreateTodo(req CreateTodoRequest, userID uint) (*models.Todo, error) {

	todo := &models.Todo{
		Title:          req.Title,
		Description:    req.Description,
		Status:         "pending",
		PriorityID:     req.PriorityID,
		CategoryID:     req.CategoryID,
		StartDate:      req.StartDate,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
		CreatedBy:      userID,
	}

	if err := s.db.Create(todo).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %v", err)
	}

	return todo, nil
}

// GetTodoList 获取TODO列表
func (s *TodoService) GetTodoList(userID uint) ([]models.Todo, error) {

	var todos []models.Todo
	if err := s.db.Preload("Priority").Preload("Category").
		Where("created_by = ?", userID).
		Order("created_at DESC").
		Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get todos: %v", err)
	}

	return todos, nil
}

// GetTodoByID 根据ID获取TODO（实现TodoServiceInterface接口）
func (s *TodoService) GetTodoByID(id uint, userID uint) (*models.Todo, error) {

	var todo models.Todo
	err := s.db.Where("id = ? AND created_by = ?", id, userID).
		Preload("Priority").
		Preload("Category").
		First(&todo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &todo, nil
}

// UpdateTodo 更新TODO
func (s *TodoService) UpdateTodo(todoID, userID uint, req UpdateTodoRequest) (*models.Todo, error) {

	var todo models.Todo
	if err := s.db.Where("id = ? AND created_by = ?", todoID, userID).First(&todo).Error; err != nil {
		return nil, errors.New("todo not found")
	}

	// 更新字段
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Status != "" {
		todo.Status = req.Status
		if req.Status == "completed" && todo.CompletedAt == nil {
			now := time.Now()
			todo.CompletedAt = &now

			// 创建任务完成通知
			notificationManager := NewNotificationManager(s.logger)
			if err := notificationManager.CreateTaskCompletedNotification(&todo); err != nil {
				fmt.Printf("Failed to create completion notification: %v\n", err)
			}
		}
	}
	if req.PriorityID != 0 {
		todo.PriorityID = req.PriorityID
	}
	if req.CategoryID != nil {
		todo.CategoryID = req.CategoryID
	}
	if req.StartDate != nil {
		todo.StartDate = req.StartDate
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}
	if req.EstimatedHours > 0 {
		todo.EstimatedHours = req.EstimatedHours
	}
	if req.ActualHours > 0 {
		todo.ActualHours = req.ActualHours
	}

	if err := s.db.Save(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to update todo: %v", err)
	}

	return &todo, nil
}

// DeleteTodo 删除TODO
func (s *TodoService) DeleteTodo(todoID, userID uint) error {

	result := s.db.Where("id = ? AND created_by = ?", todoID, userID).Delete(&models.Todo{})
	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// BatchDelete 批量删除TODO（实现TodoServiceInterface接口）
func (s *TodoService) BatchDelete(ids []uint, userID uint) error {

	// 批量删除指定用户的TODO
	result := s.db.Where("id IN ? AND created_by = ?", ids, userID).Delete(&models.Todo{})
	if result.Error != nil {
		return fmt.Errorf("failed to batch delete todos: %v", result.Error)
	}

	// 检查是否有记录被删除
	if result.RowsAffected == 0 {
		return errors.New("no todos found to delete")
	}

	return nil
}

// BatchUpdateStatus 批量更新TODO状态（实现TodoServiceInterface接口）
func (s *TodoService) BatchUpdateStatus(ids []uint, userID uint, status string) error {

	// 批量更新指定用户的TODO状态
	result := s.db.Model(&models.Todo{}).
		Where("id IN ? AND created_by = ?", ids, userID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to batch update todo status: %v", result.Error)
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		return errors.New("no todos found to update")
	}

	// 如果状态是completed，设置完成时间
	if status == "completed" {
		now := time.Now()
		s.db.Model(&models.Todo{}).
			Where("id IN ? AND created_by = ?", ids, userID).
			Update("completed_at", now)
	}

	return nil
}

// GetOverdueTodos 获取逾期TODO（实现TodoServiceInterface接口）
func (s *TodoService) GetOverdueTodos(userID uint) ([]*models.Todo, error) {

	var todos []*models.Todo
	now := time.Now()

	err := s.db.Where("created_by = ? AND due_date < ? AND status != 'completed'", userID, now).
		Preload("Priority").
		Preload("Category").
		Order("due_date ASC").
		Find(&todos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get overdue todos: %v", err)
	}

	return todos, nil
}

// GetTodoStats 获取TODO统计信息（实现TodoServiceInterface接口）
func (s *TodoService) GetTodoStats(userID uint) (map[string]interface{}, error) {

	stats := make(map[string]interface{})

	// 统计总数
	var total int64
	if err := s.db.Model(&models.Todo{}).Where("created_by = ?", userID).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total todos: %v", err)
	}
	stats["total"] = total

	// 统计已完成
	var completed int64
	if err := s.db.Model(&models.Todo{}).Where("created_by = ? AND status = ?", userID, "completed").Count(&completed).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed todos: %v", err)
	}
	stats["completed"] = completed

	// 统计进行中
	var inProgress int64
	if err := s.db.Model(&models.Todo{}).Where("created_by = ? AND status = ?", userID, "in_progress").Count(&inProgress).Error; err != nil {
		return nil, fmt.Errorf("failed to count in-progress todos: %v", err)
	}
	stats["in_progress"] = inProgress

	// 统计待处理
	var pending int64
	if err := s.db.Model(&models.Todo{}).Where("created_by = ? AND status = ?", userID, "pending").Count(&pending).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending todos: %v", err)
	}
	stats["pending"] = pending

	// 统计逾期
	var overdue int64
	now := time.Now()
	if err := s.db.Model(&models.Todo{}).Where("created_by = ? AND due_date < ? AND status != ?", userID, now, "completed").Count(&overdue).Error; err != nil {
		return nil, fmt.Errorf("failed to count overdue todos: %v", err)
	}
	stats["overdue"] = overdue

	// 计算完成率
	if total > 0 {
		completionRate := float64(completed) / float64(total) * 100
		stats["completion_rate"] = completionRate
	} else {
		stats["completion_rate"] = 0.0
	}

	return stats, nil
}

// GetTodos 获取TODO列表（实现TodoServiceInterface接口）
func (s *TodoService) GetTodos(userID uint, filter TodoFilter) (*PaginatedTodos, error) {

	query := s.db.Where("created_by = ?", userID)

	// 应用过滤器
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Joins("JOIN priorities ON todos.priority_id = priorities.id").
			Where("priorities.name = ?", filter.Priority)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.DueDate != nil {
		query = query.Where("due_date <= ?", filter.DueDate)
	}
	if filter.Overdue != nil && *filter.Overdue {
		now := time.Now()
		query = query.Where("due_date < ? AND status != ?", now, "completed")
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}

	// 获取总数
	var total int64
	if err := query.Model(&models.Todo{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count todos: %v", err)
	}

	// 应用排序
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	// 获取数据
	var todos []*models.Todo
	if err := query.Preload("Priority").Preload("Category").Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get todos: %v", err)
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return &PaginatedTodos{
		Todos:      todos,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTodosByCategory 根据分类获取TODO（实现TodoServiceInterface接口）
func (s *TodoService) GetTodosByCategory(userID uint, categoryID uint) ([]*models.Todo, error) {

	var todos []*models.Todo
	err := s.db.Where("created_by = ? AND category_id = ?", userID, categoryID).
		Preload("Priority").
		Preload("Category").
		Order("created_at DESC").
		Find(&todos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get todos by category: %v", err)
	}

	return todos, nil
}

// GetTodosByPriority 根据优先级获取TODO（实现TodoServiceInterface接口）
func (s *TodoService) GetTodosByPriority(userID uint, priority string) ([]*models.Todo, error) {

	var todos []*models.Todo
	err := s.db.Joins("JOIN priorities ON todos.priority_id = priorities.id").
		Where("todos.created_by = ? AND priorities.name = ?", userID, priority).
		Preload("Priority").
		Preload("Category").
		Order("todos.created_at DESC").
		Find(&todos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get todos by priority: %v", err)
	}

	return todos, nil
}

// MarkCompleted 标记为已完成（实现TodoServiceInterface接口）
func (s *TodoService) MarkCompleted(id uint, userID uint) error {

	now := time.Now()
	result := s.db.Model(&models.Todo{}).
		Where("id = ? AND created_by = ?", id, userID).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to mark todo as completed: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// MarkInProgress 标记为进行中（实现TodoServiceInterface接口）
func (s *TodoService) MarkInProgress(id uint, userID uint) error {

	result := s.db.Model(&models.Todo{}).
		Where("id = ? AND created_by = ?", id, userID).
		Update("status", "in_progress")

	if result.Error != nil {
		return fmt.Errorf("failed to mark todo as in progress: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// MarkCancelled 标记为已取消（实现TodoServiceInterface接口）
func (s *TodoService) MarkCancelled(id uint, userID uint) error {

	result := s.db.Model(&models.Todo{}).
		Where("id = ? AND created_by = ?", id, userID).
		Update("status", "cancelled")

	if result.Error != nil {
		return fmt.Errorf("failed to mark todo as cancelled: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("todo not found")
	}

	return nil
}

// SearchTodos 搜索TODO（实现TodoServiceInterface接口）
func (s *TodoService) SearchTodos(userID uint, query string, filter TodoFilter) (*PaginatedTodos, error) {
	// 将搜索查询添加到过滤器中
	filter.Search = query
	return s.GetTodos(userID, filter)
}
