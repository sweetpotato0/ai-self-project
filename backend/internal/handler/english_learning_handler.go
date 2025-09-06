package handler

import (
	"gin-web-framework/internal/models"
	"gin-web-framework/internal/service"
	"gin-web-framework/pkg/logger"
	"gin-web-framework/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EnglishLearningHandler struct {
	englishLearningService *service.EnglishLearningService
	logger                 logger.LoggerInterface
}

func NewEnglishLearningHandler(englishLearningService *service.EnglishLearningService, logger logger.LoggerInterface) *EnglishLearningHandler {
	return &EnglishLearningHandler{
		englishLearningService: englishLearningService,
		logger:                 logger,
	}
}

// ====== 分类管理 ======

// GetCategories 获取学习分类列表
func (h *EnglishLearningHandler) GetCategories(c *gin.Context) {
	var filter service.CategoryFilter

	// 解析查询参数
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter.Page = page
	filter.Limit = limit

	result, err := h.englishLearningService.GetCategories(&filter)
	if err != nil {
		h.logger.Errorf("Failed to get categories: %v", err)
		response.InternalServerError(c, "Failed to get categories")
		return
	}

	response.Success(c, result)
}

// CreateCategory 创建学习分类
func (h *EnglishLearningHandler) CreateCategory(c *gin.Context) {
	var category models.LearningCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.englishLearningService.CreateCategory(&category); err != nil {
		h.logger.Errorf("Failed to create category: %v", err)
		response.InternalServerError(c, "Failed to create category")
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新学习分类
func (h *EnglishLearningHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.englishLearningService.UpdateCategory(uint(id), updates); err != nil {
		h.logger.Errorf("Failed to update category: %v", err)
		response.InternalServerError(c, "Failed to update category")
		return
	}

	response.Success(c, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory 删除学习分类
func (h *EnglishLearningHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	if err := h.englishLearningService.DeleteCategory(uint(id)); err != nil {
		h.logger.Errorf("Failed to delete category: %v", err)
		if strings.Contains(err.Error(), "cannot delete category with existing songs") {
			response.BadRequest(c, err.Error())
		} else {
			response.InternalServerError(c, "Failed to delete category")
		}
		return
	}

	response.Success(c, gin.H{"message": "Category deleted successfully"})
}

// ====== 歌曲/学习材料管理 ======

// GetSongs 获取歌曲列表
func (h *EnglishLearningHandler) GetSongs(c *gin.Context) {
	var filter service.SongFilter

	// 解析查询参数
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryIDUint := uint(categoryID)
			filter.CategoryID = &categoryIDUint
		}
	}

	if difficultyStr := c.Query("difficulty"); difficultyStr != "" {
		if difficulty, err := strconv.Atoi(difficultyStr); err == nil {
			filter.Difficulty = &difficulty
		}
	}

	if isPublishedStr := c.Query("is_published"); isPublishedStr != "" {
		isPublished := isPublishedStr == "true"
		filter.IsPublished = &isPublished
	}

	filter.AgeRange = c.Query("age_range")
	filter.Search = c.Query("search")
	filter.Tags = c.Query("tags")
	filter.SortBy = c.Query("sort_by")
	filter.SortOrder = c.DefaultQuery("sort_order", "desc")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter.Page = page
	filter.Limit = limit

	// 解析前端传来的sort参数（格式：field:direction）
	if sortParam := c.Query("sort"); sortParam != "" {
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			filter.SortBy = parts[0]
			filter.SortOrder = parts[1]
		}
	}

	result, err := h.englishLearningService.GetSongs(&filter)
	if err != nil {
		h.logger.Errorf("Failed to get songs: %v", err)
		response.InternalServerError(c, "Failed to get songs")
		return
	}

	response.Success(c, result)
}

// GetSongByID 获取单个歌曲详情
func (h *EnglishLearningHandler) GetSongByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	song, err := h.englishLearningService.GetSongByID(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get song: %v", err)
		response.NotFound(c, "Song not found")
		return
	}

	response.Success(c, song)
}

// CreateSong 创建歌曲
func (h *EnglishLearningHandler) CreateSong(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.Errorf("Failed to bind JSON data: %v", err)
		response.BadRequest(c, "Invalid request data")
		return
	}

	song.CreatedBy = userID

	if err := h.englishLearningService.CreateSong(&song); err != nil {
		h.logger.Errorf("Failed to create song: %v", err)
		response.InternalServerError(c, "Failed to create song")
		return
	}

	response.Success(c, song)
}

// UpdateSong 更新歌曲
func (h *EnglishLearningHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.englishLearningService.UpdateSong(uint(id), updates); err != nil {
		h.logger.Errorf("Failed to update song: %v", err)
		response.InternalServerError(c, "Failed to update song")
		return
	}

	response.Success(c, gin.H{"message": "Song updated successfully"})
}

// DeleteSong 删除歌曲
func (h *EnglishLearningHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	if err := h.englishLearningService.DeleteSong(uint(id)); err != nil {
		h.logger.Errorf("Failed to delete song: %v", err)
		response.InternalServerError(c, "Failed to delete song")
		return
	}

	response.Success(c, gin.H{"message": "Song deleted successfully"})
}

// LikeSong 点赞歌曲
func (h *EnglishLearningHandler) LikeSong(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	songID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	if err := h.englishLearningService.LikeSong(uint(songID), userID); err != nil {
		h.logger.Errorf("Failed to like song: %v", err)
		if strings.Contains(err.Error(), "already liked") {
			response.BadRequest(c, "Already liked this song")
		} else {
			response.InternalServerError(c, "Failed to like song")
		}
		return
	}

	response.Success(c, gin.H{"message": "Song liked successfully"})
}

// UnlikeSong 取消点赞歌曲
func (h *EnglishLearningHandler) UnlikeSong(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	songID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	if err := h.englishLearningService.UnlikeSong(uint(songID), userID); err != nil {
		h.logger.Errorf("Failed to unlike song: %v", err)
		if strings.Contains(err.Error(), "not liked yet") {
			response.BadRequest(c, "Song not liked yet")
		} else {
			response.InternalServerError(c, "Failed to unlike song")
		}
		return
	}

	response.Success(c, gin.H{"message": "Song unliked successfully"})
}

// ====== 用户学习进度 ======

// UpdateProgress 更新学习进度
func (h *EnglishLearningHandler) UpdateProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	songID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid song ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.englishLearningService.UpdateUserProgress(userID, uint(songID), updates); err != nil {
		h.logger.Errorf("Failed to update progress: %v", err)
		response.InternalServerError(c, "Failed to update progress")
		return
	}

	response.Success(c, gin.H{"message": "Progress updated successfully"})
}

// GetProgress 获取学习进度
func (h *EnglishLearningHandler) GetProgress(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	var songID *uint
	if songIDStr := c.Query("song_id"); songIDStr != "" {
		if id, err := strconv.ParseUint(songIDStr, 10, 32); err == nil {
			songIDUint := uint(id)
			songID = &songIDUint
		}
	}

	progress, err := h.englishLearningService.GetUserProgress(userID, songID)
	if err != nil {
		h.logger.Errorf("Failed to get progress: %v", err)
		response.InternalServerError(c, "Failed to get progress")
		return
	}

	response.Success(c, progress)
}

// ====== 推荐和统计 ======

// GetRecommendations 获取推荐歌曲
func (h *EnglishLearningHandler) GetRecommendations(c *gin.Context) {
	userID := getUserIDFromContext(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	songs, err := h.englishLearningService.GetRecommendedSongs(userID, limit)
	if err != nil {
		h.logger.Errorf("Failed to get recommendations: %v", err)
		response.InternalServerError(c, "Failed to get recommendations")
		return
	}

	response.Success(c, gin.H{"songs": songs})
}

// GetStats 获取用户学习统计
func (h *EnglishLearningHandler) GetStats(c *gin.Context) {
	userID := getUserIDFromContext(c)

	stats, err := h.englishLearningService.GetUserStats(userID)
	if err != nil {
		h.logger.Errorf("Failed to get user stats: %v", err)
		response.InternalServerError(c, "Failed to get user stats")
		return
	}

	response.Success(c, stats)
}