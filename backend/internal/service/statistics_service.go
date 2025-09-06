package service

import (
	"fmt"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// StatisticsService 统计服务
type StatisticsService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(db *gorm.DB, logger logger.LoggerInterface) *StatisticsService {
	return &StatisticsService{
		db:     db,
		logger: logger,
	}
}

// StatisticsType 统计类型
type StatisticsType string

const (
	StatisticsTypeTodo    StatisticsType = "todo"
	StatisticsTypeArticle StatisticsType = "article"
)

// TodoStatistics 任务统计
type TodoStatistics struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	InProgress int64 `json:"inProgress"`
	Completed  int64 `json:"completed"`
	Cancelled  int64 `json:"cancelled"`
	Overdue    int64 `json:"overdue"`
}

// ArticleStatistics 文章统计
type ArticleStatistics struct {
	Total      int64 `json:"total"`
	Published  int64 `json:"published"`
	Draft      int64 `json:"draft"`
	Archived   int64 `json:"archived"`
	TotalViews int64 `json:"totalViews"`
	TotalLikes int64 `json:"totalLikes"`
}

// TrendData 趋势数据
type TrendData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
	Type  string `json:"type"`
}

// GetStatistics 获取统计数据
func (s *StatisticsService) GetStatistics(statType StatisticsType, userID uint) (interface{}, error) {
	switch statType {
	case StatisticsTypeTodo:
		return s.GetTodoStatistics(userID)
	case StatisticsTypeArticle:
		return s.GetArticleStatistics(userID)
	default:
		return nil, ErrInvalidStatisticsType
	}
}

// GetTrends 获取趋势数据
func (s *StatisticsService) GetTrends(statType StatisticsType, userID uint, days int) ([]TrendData, error) {
	switch statType {
	case StatisticsTypeTodo:
		return s.GetTodoTrends(userID, days)
	case StatisticsTypeArticle:
		return s.GetArticleTrends(userID, days)
	default:
		return nil, ErrInvalidStatisticsType
	}
}

// getTodoStatistics 获取任务统计
func (s *StatisticsService) GetTodoStatistics(userID uint) (*TodoStatistics, error) {
	var stats TodoStatistics

	// 总任务数
	if err := s.db.Model(&models.Todo{}).Where("created_by = ?", userID).Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total todos: %w", err)
	}

	// 各状态任务数 - 使用批量查询优化性能
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := s.db.Model(&models.Todo{}).
		Select("status, COUNT(*) as count").
		Where("created_by = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	// 映射状态计数
	for _, sc := range statusCounts {
		switch sc.Status {
		case "pending":
			stats.Pending = sc.Count
		case "in_progress":
			stats.InProgress = sc.Count
		case "completed":
			stats.Completed = sc.Count
		case "cancelled":
			stats.Cancelled = sc.Count
		}
	}

	// 逾期任务数
	now := time.Now()
	if err := s.db.Model(&models.Todo{}).
		Where("created_by = ? AND due_date < ? AND status NOT IN (?)",
			userID, now, []string{"completed", "cancelled"}).
		Count(&stats.Overdue).Error; err != nil {
		return nil, fmt.Errorf("failed to count overdue todos: %w", err)
	}

	return &stats, nil
}

// getArticleStatistics 获取文章统计
func (s *StatisticsService) GetArticleStatistics(userID uint) (*ArticleStatistics, error) {
	var stats ArticleStatistics

	// 总文章数
	if err := s.db.Model(&models.Article{}).Where("created_by = ?", userID).Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total articles: %w", err)
	}

	// 各状态文章数 - 批量查询
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := s.db.Model(&models.Article{}).
		Select("status, COUNT(*) as count").
		Where("created_by = ?", userID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get article status counts: %w", err)
	}

	// 映射状态计数
	for _, sc := range statusCounts {
		switch sc.Status {
		case "published":
			stats.Published = sc.Count
		case "draft":
			stats.Draft = sc.Count
		case "archived":
			stats.Archived = sc.Count
		}
	}

	// 总浏览量和点赞数
	var result struct {
		TotalViews int64 `json:"total_views"`
		TotalLikes int64 `json:"total_likes"`
	}
	if err := s.db.Model(&models.Article{}).
		Where("created_by = ?", userID).
		Select("COALESCE(SUM(view_count), 0) as total_views, COALESCE(SUM(like_count), 0) as total_likes").
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get article views and likes: %w", err)
	}
	stats.TotalViews = result.TotalViews
	stats.TotalLikes = result.TotalLikes

	return &stats, nil
}

// getTodoTrends 获取任务趋势
func (s *StatisticsService) GetTodoTrends(userID uint, days int) ([]TrendData, error) {
	var trends []TrendData

	// 生成日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	// 按日期统计任务完成数
	rows, err := s.db.Model(&models.Todo{}).
		Select("DATE(completed_at) as date, COUNT(*) as count").
		Where("created_by = ? AND completed_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("DATE(completed_at)").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建日期映射
	dateMap := make(map[string]int64)
	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		dateMap[date] = count
	}

	// 填充所有日期
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		count := dateMap[date]
		trends = append(trends, TrendData{
			Date:  date,
			Count: count,
			Type:  "todo",
		})
	}

	return trends, nil
}

// getArticleTrends 获取文章趋势
func (s *StatisticsService) GetArticleTrends(userID uint, days int) ([]TrendData, error) {
	var trends []TrendData

	// 生成日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	// 按日期统计文章浏览量
	rows, err := s.db.Model(&models.Article{}).
		Select("DATE(created_at) as date, SUM(view_count) as count").
		Where("created_by = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("DATE(created_at)").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建日期映射
	dateMap := make(map[string]int64)
	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		dateMap[date] = count
	}

	// 填充所有日期
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		count := dateMap[date]
		trends = append(trends, TrendData{
			Date:  date,
			Count: count,
			Type:  "article",
		})
	}

	return trends, nil
}

// ValidateStatisticsType 验证统计类型
func ValidateStatisticsType(statType string) (StatisticsType, error) {
	switch StatisticsType(statType) {
	case StatisticsTypeTodo, StatisticsTypeArticle:
		return StatisticsType(statType), nil
	default:
		return "", ErrInvalidStatisticsType
	}
}

// ValidateDays 验证天数参数
func ValidateDays(daysStr string) (int, error) {
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 7, ErrInvalidDaysParameter
	}

	// 限制天数范围
	if days < 1 || days > 365 {
		return 7, ErrInvalidDaysParameter
	}

	return days, nil
}

// GetActiveUsersCount 获取活跃用户数量
func (s *StatisticsService) GetActiveUsersCount(days int) (int64, error) {

	// 计算指定天数前的时间
	startDate := time.Now().AddDate(0, 0, -days)

	var count int64
	if err := s.db.Model(&models.User{}).
		Where("updated_at >= ?", startDate).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get active users count: %v", err)
	}

	return count, nil
}

// GetDailyActiveUsers 获取每日活跃用户数趋势
func (s *StatisticsService) GetDailyActiveUsers(days int) ([]TrendData, error) {

	var trends []TrendData

	// 生成过去几天的数据
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var count int64
		if err := s.db.Model(&models.User{}).
			Where("updated_at >= ? AND updated_at < ?", startOfDay, endOfDay).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to get daily active users: %v", err)
		}

		trends = append(trends, TrendData{
			Date:  startOfDay.Format("2006-01-02"),
			Count: count,
		})
	}

	return trends, nil
}

// GetSystemStats 获取系统统计信息
func (s *StatisticsService) GetSystemStats() (map[string]interface{}, error) {

	stats := make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	if err := s.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to get total users: %v", err)
	}
	stats["total_users"] = totalUsers

	// 总任务数
	var totalTodos int64
	if err := s.db.Model(&models.Todo{}).Count(&totalTodos).Error; err != nil {
		return nil, fmt.Errorf("failed to get total todos: %v", err)
	}
	stats["total_todos"] = totalTodos

	// 总文章数
	var totalArticles int64
	if err := s.db.Model(&models.Article{}).Count(&totalArticles).Error; err != nil {
		return nil, fmt.Errorf("failed to get total articles: %v", err)
	}
	stats["total_articles"] = totalArticles

	// 活跃用户数（最近7天）
	activeUsers, err := s.GetActiveUsersCount(7)
	if err != nil {
		return nil, err
	}
	stats["active_users_7d"] = activeUsers

	return stats, nil
}

// GetTodosByStatus 获取各状态任务数量统计
func (s *StatisticsService) GetTodosByStatus(userID uint) (map[string]int64, error) {
	result := make(map[string]int64)

	statuses := []string{"pending", "in_progress", "completed", "cancelled"}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.Todo{}).
			Where("created_by = ? AND status = ?", userID, status).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count todos by status %s: %v", status, err)
		}
		result[status] = count
	}

	return result, nil
}

// GetTodosByPriority 获取各优先级任务数量统计
func (s *StatisticsService) GetTodosByPriority(userID uint) (map[string]int64, error) {
	result := make(map[string]int64)

	priorities := []string{"low", "medium", "high", "urgent"}

	for _, priority := range priorities {
		var count int64
		if err := s.db.Model(&models.Todo{}).
			Where("created_by = ? AND priority = ?", userID, priority).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count todos by priority %s: %v", priority, err)
		}
		result[priority] = count
	}

	return result, nil
}

// GetTodosByCategory 获取各分类任务数量统计
func (s *StatisticsService) GetTodosByCategory(userID uint) (map[string]int64, error) {
	result := make(map[string]int64)

	// 获取用户的所有分类
	var categories []models.Category
	if err := s.db.Where("created_by = ?", userID).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}

	// 统计每个分类的任务数量
	for _, category := range categories {
		var count int64
		if err := s.db.Model(&models.Todo{}).
			Where("created_by = ? AND category_id = ?", userID, category.ID).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count todos for category %d: %v", category.ID, err)
		}
		result[category.Name] = count
	}

	// 统计未分类的任务
	var uncategorizedCount int64
	if err := s.db.Model(&models.Todo{}).
		Where("created_by = ? AND (category_id IS NULL OR category_id = 0)", userID).
		Count(&uncategorizedCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count uncategorized todos: %v", err)
	}
	result["uncategorized"] = uncategorizedCount

	return result, nil
}

// GetArticlesByStatus 获取各状态文章数量统计
func (s *StatisticsService) GetArticlesByStatus(userID uint) (map[string]int64, error) {
	result := make(map[string]int64)

	statuses := []string{"draft", "published", "archived"}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.Article{}).
			Where("created_by = ? AND status = ?", userID, status).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count articles by status %s: %v", status, err)
		}
		result[status] = count
	}

	return result, nil
}

// GetTopArticles 获取热门文章（按浏览量排序）
func (s *StatisticsService) GetTopArticles(userID uint, limit int) ([]*models.Article, error) {

	var articles []*models.Article
	if err := s.db.Where("created_by = ?", userID).
		Order("view_count DESC").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("failed to get top articles: %v", err)
	}

	return articles, nil
}

// GetUserActivityStats 获取用户活动统计
func (s *StatisticsService) GetUserActivityStats(userID uint, days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	startDate := time.Now().AddDate(0, 0, -days)

	// 创建的任务数
	var createdTodos int64
	if err := s.db.Model(&models.Todo{}).
		Where("created_by = ? AND created_at >= ?", userID, startDate).
		Count(&createdTodos).Error; err != nil {
		return nil, fmt.Errorf("failed to count created todos: %v", err)
	}
	stats["created_todos"] = createdTodos

	// 完成的任务数
	var completedTodos int64
	if err := s.db.Model(&models.Todo{}).
		Where("created_by = ? AND status = 'completed' AND updated_at >= ?", userID, startDate).
		Count(&completedTodos).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed todos: %v", err)
	}
	stats["completed_todos"] = completedTodos

	// 创建的文章数
	var createdArticles int64
	if err := s.db.Model(&models.Article{}).
		Where("created_by = ? AND created_at >= ?", userID, startDate).
		Count(&createdArticles).Error; err != nil {
		return nil, fmt.Errorf("failed to count created articles: %v", err)
	}
	stats["created_articles"] = createdArticles

	// 发布的文章数
	var publishedArticles int64
	if err := s.db.Model(&models.Article{}).
		Where("created_by = ? AND status = 'published' AND updated_at >= ?", userID, startDate).
		Count(&publishedArticles).Error; err != nil {
		return nil, fmt.Errorf("failed to count published articles: %v", err)
	}
	stats["published_articles"] = publishedArticles

	return stats, nil
}
