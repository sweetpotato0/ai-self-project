package service

import (
	"fmt"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// StatisticsService 统计服务
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

// getDB 获取数据库连接
func (s *StatisticsService) getDB() *gorm.DB {
	// 每次都重新获取数据库连接，确保连接是活跃的
	db := database.GetDB()
	if db == nil {
		panic("Database connection is nil")
	}
	return db
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
		return s.getTodoStatistics(userID)
	case StatisticsTypeArticle:
		return s.getArticleStatistics(userID)
	default:
		return nil, ErrInvalidStatisticsType
	}
}

// GetTrends 获取趋势数据
func (s *StatisticsService) GetTrends(statType StatisticsType, userID uint, days int) ([]TrendData, error) {
	switch statType {
	case StatisticsTypeTodo:
		return s.getTodoTrends(userID, days)
	case StatisticsTypeArticle:
		return s.getArticleTrends(userID, days)
	default:
		return nil, ErrInvalidStatisticsType
	}
}

// getTodoStatistics 获取任务统计
func (s *StatisticsService) getTodoStatistics(userID uint) (*TodoStatistics, error) {
	var stats TodoStatistics
	db := s.getDB()

	// 总任务数
	if err := db.Model(&models.Todo{}).Where("created_by = ?", userID).Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total todos: %w", err)
	}

	// 各状态任务数 - 使用批量查询优化性能
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := db.Model(&models.Todo{}).
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
	if err := db.Model(&models.Todo{}).
		Where("created_by = ? AND due_date < ? AND status NOT IN (?)",
			userID, now, []string{"completed", "cancelled"}).
		Count(&stats.Overdue).Error; err != nil {
		return nil, fmt.Errorf("failed to count overdue todos: %w", err)
	}

	return &stats, nil
}

// getArticleStatistics 获取文章统计
func (s *StatisticsService) getArticleStatistics(userID uint) (*ArticleStatistics, error) {
	var stats ArticleStatistics
	db := s.getDB()

	// 总文章数
	if err := db.Model(&models.Article{}).Where("created_by = ?", userID).Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total articles: %w", err)
	}

	// 各状态文章数 - 批量查询
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := db.Model(&models.Article{}).
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
	if err := db.Model(&models.Article{}).
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
func (s *StatisticsService) getTodoTrends(userID uint, days int) ([]TrendData, error) {
	var trends []TrendData

	// 生成日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	// 按日期统计任务完成数
	rows, err := s.getDB().Model(&models.Todo{}).
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
func (s *StatisticsService) getArticleTrends(userID uint, days int) ([]TrendData, error) {
	var trends []TrendData

	// 生成日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	// 按日期统计文章浏览量
	rows, err := s.getDB().Model(&models.Article{}).
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
