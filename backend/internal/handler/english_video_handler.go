package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gin-web-framework/internal/api"
	"gin-web-framework/internal/service"
)

type EnglishVideoHandler struct {
	videoService *service.EnglishVideoService
}

func NewEnglishVideoHandler(videoService *service.EnglishVideoService) *EnglishVideoHandler {
	return &EnglishVideoHandler{
		videoService: videoService,
	}
}

// GetVideoSeries 获取视频系列列表
func (h *EnglishVideoHandler) GetVideoSeries(c *gin.Context) {
	var req api.VideoSeriesListRequest
	
	// 手动绑定查询参数，跳过验证
	req.Search = c.Query("search")
	req.Difficulty, _ = strconv.Atoi(c.Query("difficulty"))
	req.AgeRange = c.Query("age_range")
	req.Tag = c.Query("tag")
	req.SortBy = c.Query("sort_by")
	req.SortOrder = c.Query("sort_order")
	
	// 手动处理分页参数
	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		req.Page = page
	} else {
		req.Page = 1
	}
	
	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		req.PageSize = pageSize
	} else {
		req.PageSize = 20
	}
	
	// 设置默认分页参数并确保有效性
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	// 确保排序顺序的有效性
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	result, err := h.videoService.GetVideoSeries(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetVideoSeriesDetail 获取视频系列详情
func (h *EnglishVideoHandler) GetVideoSeriesDetail(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	userID := getUserID(c)
	series, err := h.videoService.GetVideoSeriesDetail(uint(seriesID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": series})
}

// GetEpisodes 获取系列的剧集列表
func (h *EnglishVideoHandler) GetEpisodes(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	var req api.EpisodeListRequest
	
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 设置默认分页参数并确保有效性
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	// 确保排序顺序的有效性
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	episodes, err := h.videoService.GetEpisodes(uint(seriesID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": episodes})
}

// GetEpisodeDetail 获取剧集详情
func (h *EnglishVideoHandler) GetEpisodeDetail(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	episode, err := h.videoService.GetEpisodeDetail(uint(episodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": episode})
}

// GetEpisodeProgress 获取用户观看进度
func (h *EnglishVideoHandler) GetEpisodeProgress(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	progress, err := h.videoService.GetEpisodeProgress(userID, uint(episodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}

// UpdateEpisodeProgress 更新观看进度
func (h *EnglishVideoHandler) UpdateEpisodeProgress(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req api.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.videoService.UpdateEpisodeProgress(userID, uint(episodeID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}

// ToggleSeriesLike 收藏/取消收藏系列
func (h *EnglishVideoHandler) ToggleSeriesLike(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	result, err := h.videoService.ToggleSeriesLike(userID, uint(seriesID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetVideoStats 获取视频统计信息（不需要认证）
func (h *EnglishVideoHandler) GetVideoStats(c *gin.Context) {
	stats, err := h.videoService.GetVideoStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUserProgress 获取用户观看进度
func (h *EnglishVideoHandler) GetUserProgress(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	progress, err := h.videoService.GetUserProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}

// GetUserVideoStats 获取用户观看统计
func (h *EnglishVideoHandler) GetUserVideoStats(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := h.videoService.GetUserVideoStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// SearchVideoSeries 搜索视频系列
func (h *EnglishVideoHandler) SearchVideoSeries(c *gin.Context) {
	var req api.SearchVideoSeriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserID(c)
	result, err := h.videoService.SearchVideoSeries(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetRecommendedSeries 获取推荐视频系列
func (h *EnglishVideoHandler) GetRecommendedSeries(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	userID := getUserID(c)
	series, err := h.videoService.GetRecommendedSeries(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": series})
}

// 管理员接口

// CreateVideoSeries 创建视频系列
func (h *EnglishVideoHandler) CreateVideoSeries(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req api.CreateVideoSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	series, err := h.videoService.CreateVideoSeries(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": series})
}

// UpdateVideoSeries 更新视频系列
func (h *EnglishVideoHandler) UpdateVideoSeries(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	var req api.UpdateVideoSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.videoService.UpdateVideoSeries(uint(seriesID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video series updated successfully"})
}

// DeleteVideoSeries 删除视频系列
func (h *EnglishVideoHandler) DeleteVideoSeries(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	err = h.videoService.DeleteVideoSeries(uint(seriesID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video series deleted successfully"})
}

// CreateEpisode 创建剧集
func (h *EnglishVideoHandler) CreateEpisode(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req api.CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	episode, err := h.videoService.CreateEpisode(uint(seriesID), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": episode})
}

// UpdateEpisode 更新剧集
func (h *EnglishVideoHandler) UpdateEpisode(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	var req api.UpdateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.videoService.UpdateEpisode(uint(episodeID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode updated successfully"})
}

// DeleteEpisode 删除剧集
func (h *EnglishVideoHandler) DeleteEpisode(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	err = h.videoService.DeleteEpisode(uint(episodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}

// BatchImportEpisodes 批量导入剧集
func (h *EnglishVideoHandler) BatchImportEpisodes(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("seriesId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req api.BatchImportEpisodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.videoService.BatchImportEpisodes(uint(seriesID), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// CreateUncategorizedEpisode 创建未分类剧集
func (h *EnglishVideoHandler) CreateUncategorizedEpisode(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req api.CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	episode, err := h.videoService.CreateUncategorizedEpisode(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": episode})
}

// getUserID 从上下文中获取用户ID
func getUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}