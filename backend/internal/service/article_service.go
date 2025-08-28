package service

import (
	"encoding/json"
	"fmt"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/utils"

	"gorm.io/gorm"
)

type ArticleService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

func NewArticleService(db *gorm.DB, logger logger.LoggerInterface) *ArticleService {
	return &ArticleService{
		db:     db,
		logger: logger,
	}
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Summary    string   `json:"summary"`
	CoverImage string   `json:"cover_image"`
	Status     string   `json:"status"`
	Tags       []string `json:"tags"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Summary    string   `json:"summary"`
	CoverImage string   `json:"cover_image"`
	Status     string   `json:"status"`
	Tags       []string `json:"tags"`
}

type ArticleFilter struct {
	Status      string `json:"status"`
	IsPublished *bool  `json:"is_published"`
	Search      string `json:"search"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	SortBy      string `json:"sort_by"`
	SortOrder   string `json:"sort_order"`
}

type PaginatedArticles struct {
	Articles   []*models.Article `json:"articles"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req CreateArticleRequest, userID uint) (*models.Article, error) {
	// 将标签转换为JSON字符串
	tagsJSON := ""
	if req.Tags != nil {
		if tags, err := json.Marshal(req.Tags); err == nil {
			tagsJSON = string(tags)
		}
	}

	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		Tags:       tagsJSON,
		CreatedBy:  userID,
	}

	if err := s.db.Create(article).Error; err != nil {
		return nil, fmt.Errorf("failed to create article: %v", err)
	}

	return article, nil
}

// GetUserArticles 获取用户文章列表
func (s *ArticleService) GetUserArticles(userID uint, filter *ArticleFilter) (*PaginatedArticles, error) {
	var articles []*models.Article
	var total int64

	query := s.db.Where("created_by = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.IsPublished != nil {
		query = query.Where("status = ?", "published")
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 获取总数
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count articles: %v", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("failed to get articles: %v", err)
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return &PaginatedArticles{
		Articles:   articles,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(id uint, userID uint) (*models.Article, error) {
	var article models.Article
	if err := s.db.Where("id = ? AND created_by = ?", id, userID).First(&article).Error; err != nil {
		return nil, fmt.Errorf("article not found: %v", err)
	}
	return &article, nil
}

// GetArticles 获取文章列表
func (s *ArticleService) GetArticles(userID uint, filter ArticleFilter) (*PaginatedArticles, error) {
	var articles []*models.Article
	var total int64

	query := s.db.Where("created_by = ?", userID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 获取总数
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	pagination := utils.NewPaginationInfo(filter.Page, filter.Limit, total)

	// 排序
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	if err := query.Order(sortBy + " " + sortOrder).Offset(pagination.Offset).Limit(pagination.Limit).Find(&articles).Error; err != nil {
		return nil, err
	}

	return &PaginatedArticles{
		Articles:   articles,
		Total:      pagination.Total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: pagination.TotalPages,
	}, nil
}

// PublishArticle 发布文章
func (s *ArticleService) PublishArticle(id uint, userID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ? AND created_by = ?", id, userID).Update("status", "published").Error
}

// ArchiveArticle 归档文章
func (s *ArticleService) ArchiveArticle(id uint, userID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ? AND created_by = ?", id, userID).Update("status", "archived").Error
}

// RestoreArticle 恢复文章
func (s *ArticleService) RestoreArticle(id uint, userID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ? AND created_by = ?", id, userID).Update("status", "draft").Error
}

// UpdateContent 更新文章内容
func (s *ArticleService) UpdateContent(id uint, userID uint, content string) error {
	return s.db.Model(&models.Article{}).Where("id = ? AND created_by = ?", id, userID).Update("content", content).Error
}

// IncrementViewCount 增加浏览次数
func (s *ArticleService) IncrementViewCount(id uint) error {
	return s.db.Model(&models.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// LikeArticle 点赞文章
func (s *ArticleService) LikeArticle(id uint, userID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}

// UnlikeArticle 取消点赞文章
func (s *ArticleService) UnlikeArticle(id uint, userID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count - 1")).Error
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id uint, userID uint, req UpdateArticleRequest) (*models.Article, error) {
	var article models.Article
	if err := s.db.Where("id = ? AND created_by = ?", id, userID).First(&article).Error; err != nil {
		return nil, fmt.Errorf("article not found: %v", err)
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	if req.Status != "" {
		article.Status = req.Status
	}
	if req.Tags != nil {
		if tags, err := json.Marshal(req.Tags); err == nil {
			article.Tags = string(tags)
		}
	}

	if err := s.db.Save(&article).Error; err != nil {
		return nil, fmt.Errorf("failed to update article: %v", err)
	}

	return &article, nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(articleID, userID uint) error {
	if err := s.db.Where("id = ? AND created_by = ?", articleID, userID).Delete(&models.Article{}).Error; err != nil {
		return fmt.Errorf("failed to delete article: %v", err)
	}

	return nil
}

// GetArticleStatistics 获取文章统计信息
func (s *ArticleService) GetArticleStatistics(userID uint) (*ArticleStatistics, error) {
	var stats ArticleStatistics

	// 总文章数
	if err := s.db.Model(&models.Article{}).Where("created_by = ?", userID).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// 已发布文章数
	if err := s.db.Model(&models.Article{}).Where("created_by = ? AND status = ?", userID, "published").Count(&stats.Published).Error; err != nil {
		return nil, err
	}

	// 草稿文章数
	if err := s.db.Model(&models.Article{}).Where("created_by = ? AND status = ?", userID, "draft").Count(&stats.Draft).Error; err != nil {
		return nil, err
	}

	// 总浏览量
	if err := s.db.Model(&models.Article{}).Where("created_by = ?", userID).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	// 总点赞数
	if err := s.db.Model(&models.Article{}).Where("created_by = ?", userID).Select("COALESCE(SUM(like_count), 0)").Scan(&stats.TotalLikes).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// IncrementArticleViewCount 增加文章浏览次数
func (s *ArticleService) IncrementArticleViewCount(articleID uint) error {
	return s.db.Model(&models.Article{}).Where("id = ?", articleID).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// GetArticleTrends 获取文章趋势数据
func (s *ArticleService) GetArticleTrends(userID uint, days int) ([]TrendData, error) {
	var trends []TrendData

	query := `
		SELECT
			DATE(created_at) as date,
			SUM(view_count) as count
		FROM articles
		WHERE created_by = ?
		AND created_at >= DATE('now', '-? days')
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	rows, err := s.db.Raw(query, userID, days).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get article trends: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			return nil, fmt.Errorf("failed to scan article trend row: %v", err)
		}

		trends = append(trends, TrendData{
			Date:  date,
			Count: count,
			Type:  "article",
		})
	}

	return trends, nil
}

// GetArticleStats 获取文章统计
func (s *ArticleService) GetArticleStats(userID uint) (map[string]interface{}, error) {
	stats, err := s.GetArticleStatistics(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":       stats.Total,
		"published":   stats.Published,
		"draft":       stats.Draft,
		"total_views": stats.TotalViews,
		"total_likes": stats.TotalLikes,
	}, nil
}

// GetPopularArticles 获取热门文章
func (s *ArticleService) GetPopularArticles(userID uint, limit int) ([]*models.Article, error) {
	var articles []*models.Article
	if err := s.db.Where("created_by = ?", userID).Order("view_count DESC").Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// GetRecentArticles 获取最近文章
func (s *ArticleService) GetRecentArticles(userID uint, limit int) ([]*models.Article, error) {
	var articles []*models.Article
	if err := s.db.Where("created_by = ?", userID).Order("created_at DESC").Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// SearchArticles 搜索文章
func (s *ArticleService) SearchArticles(userID uint, query string, filter ArticleFilter) (*PaginatedArticles, error) {
	// 设置搜索查询
	filter.Search = query
	return s.GetArticles(userID, filter)
}
