package handler

import (
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *service.ArticleService
	logger         logger.LoggerInterface
}

func NewArticleHandler(articleService *service.ArticleService, logger logger.LoggerInterface) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		logger:         logger,
	}
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req service.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	userID := getUserIDFromContext(c)
	article, err := h.articleService.CreateArticle(req, userID)
	if err != nil {
		response.InternalServerError(c, "Failed to create article")
		return
	}

	response.Success(c, article)
}

// GetUserArticles 获取用户文章列表
func (h *ArticleHandler) GetUserArticles(c *gin.Context) {
	userID := getUserIDFromContext(c)
	h.logger.Debugf("GetUserArticles called for user %d", userID)

	// 获取查询参数并记录日志
	status := c.Query("status")
	search := c.Query("search")
	tag := c.Query("tag")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	h.logger.Debugf("Raw query params - status: '%s', search: '%s', tag: '%s', sort_by: '%s', sort_order: '%s', page: %d, limit: %d",
		status, search, tag, sortBy, sortOrder, page, limit)

	// 解析前端传来的sort参数（格式：field:direction）
	sortParam := c.Query("sort")
	if sortParam != "" {
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			sortBy = parts[0]
			sortOrder = parts[1]
			h.logger.Debugf("Parsed sort param: sortBy='%s', sortOrder='%s'", sortBy, sortOrder)
		}
	}

	filter := &service.ArticleFilter{
		Status:    status,
		Search:    search,
		Tags:      tag,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		Limit:     limit,
	}

	h.logger.Debugf("Filter object: %+v", filter)

	result, err := h.articleService.GetUserArticles(userID, filter)
	if err != nil {
		h.logger.Errorf("Failed to get articles: %v", err)
		response.InternalServerError(c, "Failed to get articles")
		return
	}

	h.logger.Debugf("Retrieved %d articles out of %d total", len(result.Articles), result.Total)
	response.Success(c, result)
}

// GetArticleByID 根据ID获取文章
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	userID := getUserIDFromContext(c)

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	article, err := h.articleService.GetArticleByID(uint(articleID), userID)
	if err != nil {
		response.NotFound(c, "Article not found")
		return
	}

	response.Success(c, article)
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID := getUserIDFromContext(c)

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	var req service.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	article, err := h.articleService.UpdateArticle(uint(articleID), userID, req)
	if err != nil {
		response.InternalServerError(c, "Failed to update article")
		return
	}

	response.Success(c, article)
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID := getUserIDFromContext(c)

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	if err := h.articleService.DeleteArticle(uint(articleID), userID); err != nil {
		response.InternalServerError(c, "Failed to delete article")
		return
	}

	response.Success(c, gin.H{"message": "Article deleted successfully"})
}

// LikeArticle 点赞文章
func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	userID := getUserIDFromContext(c)
	h.logger.Debugf("LikeArticle: userID from context: %d", userID)

	if userID == 0 {
		h.logger.Errorf("LikeArticle: No user ID found in context")
		response.BadRequest(c, "Authentication required")
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Errorf("LikeArticle: Invalid article ID: %s, error: %v", c.Param("id"), err)
		response.BadRequest(c, "Invalid article ID")
		return
	}

	h.logger.Debugf("LikeArticle: Attempting to like article %d for user %d", articleID, userID)

	if err := h.articleService.LikeArticle(uint(articleID), userID); err != nil {
		if err.Error() == "already liked" {
			h.logger.Debugf("LikeArticle: User %d already liked article %d", userID, articleID)
			response.BadRequest(c, "Already liked this article")
			return
		}
		h.logger.Errorf("LikeArticle: Service error: %v", err)
		response.InternalServerError(c, "Failed to like article")
		return
	}

	h.logger.Infof("LikeArticle: Successfully liked article %d for user %d", articleID, userID)
	response.Success(c, gin.H{"message": "Article liked successfully"})
}

// UnlikeArticle 取消点赞文章
func (h *ArticleHandler) UnlikeArticle(c *gin.Context) {
	userID := getUserIDFromContext(c)

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	if err := h.articleService.UnlikeArticle(uint(articleID), userID); err != nil {
		response.InternalServerError(c, "Failed to unlike article")
		return
	}

	response.Success(c, gin.H{"message": "Article unliked successfully"})
}

// IncrementViewCount 增加文章浏览量
func (h *ArticleHandler) IncrementViewCount(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID")
		return
	}

	if err := h.articleService.IncrementViewCount(uint(articleID)); err != nil {
		response.InternalServerError(c, "Failed to increment view count")
		return
	}

	response.Success(c, gin.H{"message": "View count incremented successfully"})
}

// GetArticleStats 获取文章统计信息
func (h *ArticleHandler) GetArticleStats(c *gin.Context) {
	userID := getUserIDFromContext(c)

	stats, err := h.articleService.GetArticleStatistics(userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get article stats")
		return
	}

	response.Success(c, stats)
}
