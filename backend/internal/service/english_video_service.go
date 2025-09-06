package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gin-web-framework/internal/api"
	"gin-web-framework/internal/models"

	"gorm.io/gorm"
)

type EnglishVideoService struct {
	db *gorm.DB
}

func NewEnglishVideoService(db *gorm.DB) *EnglishVideoService {
	return &EnglishVideoService{
		db: db,
	}
}

// GetVideoSeries 获取视频系列列表
func (s *EnglishVideoService) GetVideoSeries(req *api.VideoSeriesListRequest) (*api.VideoSeriesListResponse, error) {
	var series []models.VideoSeries
	var total int64

	query := s.db.Model(&models.VideoSeries{}).Where("is_published = ?", true)

	// 搜索
	if req.Search != "" {
		query = query.Where("title LIKE ? OR title_cn LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 难度过滤
	if req.Difficulty > 0 {
		query = query.Where("difficulty = ?", req.Difficulty)
	}

	// 年龄段过滤
	if req.AgeRange != "" {
		query = query.Where("age_range = ?", req.AgeRange)
	}

	// 标签过滤
	if req.Tag != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, req.Tag))
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	orderBy := "sort DESC, created_at DESC"
	if req.SortBy != "" {
		switch req.SortBy {
		case "title":
			orderBy = "title " + getSortOrder(req.SortOrder)
		case "difficulty":
			orderBy = "difficulty " + getSortOrder(req.SortOrder)
		case "view_count":
			orderBy = "view_count " + getSortOrder(req.SortOrder)
		case "like_count":
			orderBy = "like_count " + getSortOrder(req.SortOrder)
		case "created_at":
			orderBy = "created_at " + getSortOrder(req.SortOrder)
		}
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order(orderBy).
		Offset(offset).
		Limit(req.PageSize).
		Preload("Creator").
		Find(&series).Error; err != nil {
		return nil, err
	}

	// 添加统计信息
	for i := range series {
		var episodeCount int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id = ? AND is_published = ?", series[i].ID, true).Count(&episodeCount)
		series[i].EpisodeCount = episodeCount
	}

	// 检查是否有未分类的视频（series_id为NULL的视频）
	var uncategorizedCount int64
	s.db.Model(&models.VideoEpisode{}).Where("series_id IS NULL AND is_published = ?", true).Count(&uncategorizedCount)

	if uncategorizedCount > 0 {
		// 创建虚拟的未分类系列
		uncategorizedSeries := models.VideoSeries{
			ID:           0, // 特殊ID标识虚拟系列
			Title:        "Uncategorized",
			TitleCN:      "未分类",
			Description:  "未分类的视频内容",
			CoverImage:   "",
			Difficulty:   1,
			AgeRange:     "全年龄",
			Tags:         "[]",
			ViewCount:    0,
			LikeCount:    0,
			IsPublished:  true,
			Sort:         999999, // 排在最后
			EpisodeCount: uncategorizedCount,
		}

		// 计算未分类视频的总观看次数
		var totalViews int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id IS NULL AND is_published = ?", true).
			Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
		uncategorizedSeries.ViewCount = int(totalViews)

		// 将虚拟系列添加到结果中
		series = append(series, uncategorizedSeries)
		total += 1
	}

	return &api.VideoSeriesListResponse{
		List:     series,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetVideoSeriesDetail 获取视频系列详情
func (s *EnglishVideoService) GetVideoSeriesDetail(seriesID, userID uint) (*models.VideoSeries, error) {
	// 处理虚拟未分类系列（ID = 0）
	if seriesID == 0 {
		var uncategorizedCount int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id IS NULL AND is_published = ?", true).Count(&uncategorizedCount)

		if uncategorizedCount == 0 {
			return nil, errors.New("video series not found")
		}

		// 计算未分类视频的总观看次数
		var totalViews int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id IS NULL AND is_published = ?", true).
			Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)

		return &models.VideoSeries{
			ID:           0,
			Title:        "Uncategorized",
			TitleCN:      "未分类",
			Description:  "未分类的视频内容",
			CoverImage:   "",
			Difficulty:   1,
			AgeRange:     "全年龄",
			Tags:         "[]",
			ViewCount:    int(totalViews),
			LikeCount:    0,
			IsPublished:  true,
			Sort:         999999,
			EpisodeCount: uncategorizedCount,
			IsLiked:      false, // 虚拟系列不支持点赞
		}, nil
	}

	var series models.VideoSeries
	if err := s.db.Where("id = ? AND is_published = ?", seriesID, true).
		Preload("Creator").
		First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video series not found")
		}
		return nil, err
	}

	// 获取剧集数量
	var episodeCount int64
	s.db.Model(&models.VideoEpisode{}).Where("series_id = ? AND is_published = ?", seriesID, true).Count(&episodeCount)
	series.EpisodeCount = episodeCount

	// 检查用户是否点赞
	if userID > 0 {
		var like models.VideoSeriesLike
		if err := s.db.Where("user_id = ? AND series_id = ?", userID, seriesID).First(&like).Error; err == nil {
			series.IsLiked = true
		}
	}

	// 增加观看次数
	s.db.Model(&series).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	return &series, nil
}

// GetEpisodes 获取系列的剧集列表
func (s *EnglishVideoService) GetEpisodes(seriesID uint, req *api.EpisodeListRequest) (*api.EpisodeListResponse, error) {
	var episodes []models.VideoEpisode
	var total int64
	var series models.VideoSeries

	// 处理虚拟未分类系列（ID = 0）
	if seriesID == 0 {
		// 查询未分类的视频（series_id为NULL）
		query := s.db.Model(&models.VideoEpisode{}).Where("series_id IS NULL AND is_published = ?", true)

		// 获取总数
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}

		// 排序
		orderBy := "created_at DESC" // 未分类视频按创建时间倒序
		if req.SortBy != "" {
			switch req.SortBy {
			case "title":
				orderBy = "title " + getSortOrder(req.SortOrder)
			case "episode_num":
				orderBy = "episode_num " + getSortOrder(req.SortOrder)
			case "duration":
				orderBy = "duration " + getSortOrder(req.SortOrder)
			case "created_at":
				orderBy = "created_at " + getSortOrder(req.SortOrder)
			}
		}

		// 分页
		offset := (req.Page - 1) * req.PageSize
		if err := query.Order(orderBy).
			Offset(offset).
			Limit(req.PageSize).
			Preload("Creator").
			Find(&episodes).Error; err != nil {
			return nil, err
		}

		// 构造虚拟系列信息
		series = models.VideoSeries{
			ID:           0,
			Title:        "Uncategorized",
			TitleCN:      "未分类",
			Description:  "未分类的视频内容",
			CoverImage:   "",
			Difficulty:   1,
			AgeRange:     "全年龄",
			Tags:         "[]",
			ViewCount:    0,
			LikeCount:    0,
			IsPublished:  true,
			Sort:         999999,
			EpisodeCount: total,
		}
	} else {
		// 检查系列是否存在
		if err := s.db.Where("id = ? AND is_published = ?", seriesID, true).First(&series).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("video series not found")
			}
			return nil, err
		}

		query := s.db.Model(&models.VideoEpisode{}).Where("series_id = ? AND is_published = ?", seriesID, true)

		// 获取总数
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}

		// 排序
		orderBy := "episode_num ASC"
		if req.SortBy != "" {
			switch req.SortBy {
			case "title":
				orderBy = "title " + getSortOrder(req.SortOrder)
			case "episode_num":
				orderBy = "episode_num " + getSortOrder(req.SortOrder)
			case "duration":
				orderBy = "duration " + getSortOrder(req.SortOrder)
			case "created_at":
				orderBy = "created_at " + getSortOrder(req.SortOrder)
			}
		}

		// 分页
		offset := (req.Page - 1) * req.PageSize
		if err := query.Order(orderBy).
			Offset(offset).
			Limit(req.PageSize).
			Preload("Creator").
			Find(&episodes).Error; err != nil {
			return nil, err
		}
	}

	return &api.EpisodeListResponse{
		List:     episodes,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Series:   &series,
	}, nil
}

// GetEpisodeDetail 获取剧集详情
func (s *EnglishVideoService) GetEpisodeDetail(episodeID uint) (*models.VideoEpisode, error) {
	var episode models.VideoEpisode
	if err := s.db.Where("id = ? AND is_published = ?", episodeID, true).
		Preload("Series").
		Preload("Creator").
		First(&episode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("episode not found")
		}
		return nil, err
	}

	// 增加观看次数
	s.db.Model(&episode).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	return &episode, nil
}

// GetEpisodeProgress 获取用户观看进度
func (s *EnglishVideoService) GetEpisodeProgress(userID, episodeID uint) (*models.VideoUserProgress, error) {
	var progress models.VideoUserProgress
	if err := s.db.Where("user_id = ? AND episode_id = ?", userID, episodeID).First(&progress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.VideoUserProgress{
				UserID:    userID,
				EpisodeID: episodeID,
				Progress:  0,
			}, nil
		}
		return nil, err
	}

	return &progress, nil
}

// UpdateEpisodeProgress 更新观看进度
func (s *EnglishVideoService) UpdateEpisodeProgress(userID, episodeID uint, req *api.UpdateProgressRequest) error {
	// 获取剧集信息
	var episode models.VideoEpisode
	if err := s.db.Where("id = ?", episodeID).First(&episode).Error; err != nil {
		return errors.New("episode not found")
	}

	var progress models.VideoUserProgress
	err := s.db.Where("user_id = ? AND episode_id = ?", userID, episodeID).First(&progress).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新的进度记录
		progress = models.VideoUserProgress{
			UserID:    userID,
			SeriesID:  episode.SeriesID,
			EpisodeID: episodeID,
		}
	}

	// 更新进度数据
	if req.Progress >= 0 {
		progress.Progress = req.Progress
	}
	if req.CurrentTime >= 0 {
		progress.CurrentTime = req.CurrentTime
	}
	if req.WatchTimeMinutes > 0 {
		progress.WatchTimeMinutes += req.WatchTimeMinutes
	}
	if req.IsCompleted != nil {
		progress.IsCompleted = *req.IsCompleted
	}
	if req.Notes != "" {
		progress.Notes = req.Notes
	}

	now := time.Now()
	progress.LastWatchedAt = &now

	return s.db.Save(&progress).Error
}

// ToggleSeriesLike 收藏/取消收藏系列
func (s *EnglishVideoService) ToggleSeriesLike(userID, seriesID uint) (*api.ToggleLikeResponse, error) {
	// 虚拟未分类系列不支持点赞
	if seriesID == 0 {
		return nil, errors.New("uncategorized series does not support like operation")
	}

	// 检查系列是否存在
	var series models.VideoSeries
	if err := s.db.Where("id = ?", seriesID).First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video series not found")
		}
		return nil, err
	}

	var like models.VideoSeriesLike
	err := s.db.Where("user_id = ? AND series_id = ?", userID, seriesID).First(&like).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var isLiked bool
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建点赞
		like = models.VideoSeriesLike{
			UserID:   userID,
			SeriesID: seriesID,
		}
		if err := s.db.Create(&like).Error; err != nil {
			return nil, err
		}
		s.db.Model(&series).UpdateColumn("like_count", gorm.Expr("like_count + 1"))
		isLiked = true
	} else {
		// 取消点赞
		if err := s.db.Delete(&like).Error; err != nil {
			return nil, err
		}
		s.db.Model(&series).UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		isLiked = false
	}

	// 获取最新点赞数
	s.db.Where("id = ?", seriesID).First(&series)

	return &api.ToggleLikeResponse{
		IsLiked:   isLiked,
		LikeCount: series.LikeCount,
	}, nil
}

// GetVideoStats 获取视频统计信息（全局统计）
func (s *EnglishVideoService) GetVideoStats() (*api.VideoStatsResponse, error) {
	stats := &api.VideoStatsResponse{}

	// 总系列数
	s.db.Model(&models.VideoSeries{}).Where("is_published = ?", true).Count(&stats.TotalSeries)

	// 总集数
	s.db.Model(&models.VideoEpisode{}).Where("is_published = ?", true).Count(&stats.TotalEpisodes)

	// 总观看次数
	s.db.Model(&models.VideoSeries{}).Where("is_published = ?", true).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews)

	// 已看完的系列数（简单统计：有用户完成了该系列所有集数的系列）
	s.db.Raw(`
		SELECT COUNT(DISTINCT vs.id)
		FROM video_series vs
		JOIN video_episodes ve ON vs.id = ve.series_id
		JOIN video_user_progress vup ON ve.id = vup.episode_id
		WHERE vs.is_published = ? AND ve.is_published = ? AND vup.is_completed = ?
		GROUP BY vs.id
		HAVING COUNT(DISTINCT ve.id) = (
			SELECT COUNT(*)
			FROM video_episodes ve2
			WHERE ve2.series_id = vs.id AND ve2.is_published = ?
		)
	`, true, true, true, true).Scan(&stats.WatchedSeries)

	return stats, nil
}

// GetUserProgress 获取用户观看进度列表
func (s *EnglishVideoService) GetUserProgress(userID uint) ([]models.VideoUserProgress, error) {
	var progress []models.VideoUserProgress

	if err := s.db.Where("user_id = ?", userID).
		Preload("Series").
		Order("last_watched_at DESC").
		Find(&progress).Error; err != nil {
		return nil, err
	}

	// 计算每个系列的观看进度百分比
	for i := range progress {
		if progress[i].SeriesID > 0 {
			// 获取系列总集数
			var totalEpisodes int64
			s.db.Model(&models.VideoEpisode{}).
				Where("series_id = ? AND is_published = ?", progress[i].SeriesID, true).
				Count(&totalEpisodes)

			// 获取用户已完成的集数
			var completedEpisodes int64
			s.db.Model(&models.VideoUserProgress{}).
				Where("user_id = ? AND series_id = ? AND is_completed = ?", userID, progress[i].SeriesID, true).
				Count(&completedEpisodes)

			if totalEpisodes > 0 {
				progress[i].WatchProgress = float64(completedEpisodes) / float64(totalEpisodes)
			}

			// 获取当前观看的集数
			var currentEpisode models.VideoEpisode
			if err := s.db.Where("id = ?", progress[i].EpisodeID).First(&currentEpisode).Error; err == nil {
				progress[i].CurrentEpisode = currentEpisode.EpisodeNum
			}
		}
	}

	return progress, nil
}

// GetUserVideoStats 获取用户观看统计
func (s *EnglishVideoService) GetUserVideoStats(userID uint) (*api.UserVideoStatsResponse, error) {
	stats := &api.UserVideoStatsResponse{}

	// 总观看时长
	s.db.Model(&models.VideoUserProgress{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(watch_time_minutes), 0)").
		Scan(&stats.TotalWatchTime)

	// 完成的剧集数
	s.db.Model(&models.VideoUserProgress{}).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Count(&stats.CompletedEpisodes)

	// 收藏的系列数
	s.db.Model(&models.VideoSeriesLike{}).
		Where("user_id = ?", userID).
		Count(&stats.FavoriteSeries)

	// 学习连续天数（简单实现：最近7天内有观看记录的天数）
	s.db.Raw(`
		SELECT COUNT(DISTINCT DATE(last_watched_at))
		FROM video_user_progress
		WHERE user_id = ? AND last_watched_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
	`, userID).Scan(&stats.LearningStreak)

	return stats, nil
}

// SearchVideoSeries 搜索视频系列
func (s *EnglishVideoService) SearchVideoSeries(req *api.SearchVideoSeriesRequest, userID uint) ([]models.VideoSeries, error) {
	var series []models.VideoSeries

	query := s.db.Model(&models.VideoSeries{}).Where("is_published = ?", true)

	if req.Query != "" {
		query = query.Where("title LIKE ? OR title_cn LIKE ? OR description LIKE ?",
			"%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	// 应用过滤条件
	if req.Difficulty > 0 {
		query = query.Where("difficulty = ?", req.Difficulty)
	}
	if req.AgeRange != "" {
		query = query.Where("age_range = ?", req.AgeRange)
	}
	if req.Tag != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, req.Tag))
	}

	if err := query.Order("view_count DESC, created_at DESC").
		Limit(50). // 限制搜索结果数量
		Preload("Creator").
		Find(&series).Error; err != nil {
		return nil, err
	}

	// 添加统计信息
	for i := range series {
		var episodeCount int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id = ? AND is_published = ?", series[i].ID, true).Count(&episodeCount)
		series[i].EpisodeCount = episodeCount

		// 检查用户点赞状态
		if userID > 0 {
			var like models.VideoSeriesLike
			if err := s.db.Where("user_id = ? AND series_id = ?", userID, series[i].ID).First(&like).Error; err == nil {
				series[i].IsLiked = true
			}
		}
	}

	return series, nil
}

// GetRecommendedSeries 获取推荐视频系列
func (s *EnglishVideoService) GetRecommendedSeries(userID uint, limit int) ([]models.VideoSeries, error) {
	var series []models.VideoSeries

	// 简单的推荐算法：按观看次数和点赞数排序
	if err := s.db.Where("is_published = ?", true).
		Order("(view_count * 0.7 + like_count * 0.3) DESC, created_at DESC").
		Limit(limit).
		Preload("Creator").
		Find(&series).Error; err != nil {
		return nil, err
	}

	// 添加统计信息
	for i := range series {
		var episodeCount int64
		s.db.Model(&models.VideoEpisode{}).Where("series_id = ? AND is_published = ?", series[i].ID, true).Count(&episodeCount)
		series[i].EpisodeCount = episodeCount

		// 检查用户点赞状态
		if userID > 0 {
			var like models.VideoSeriesLike
			if err := s.db.Where("user_id = ? AND series_id = ?", userID, series[i].ID).First(&like).Error; err == nil {
				series[i].IsLiked = true
			}
		}
	}

	return series, nil
}

// 管理员功能

// CreateVideoSeries 创建视频系列
func (s *EnglishVideoService) CreateVideoSeries(req *api.CreateVideoSeriesRequest, userID uint) (*models.VideoSeries, error) {
	tagsJSON, _ := json.Marshal(req.Tags)

	series := models.VideoSeries{
		Title:       req.Title,
		TitleCN:     req.TitleCN,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Difficulty:  req.Difficulty,
		AgeRange:    req.AgeRange,
		Tags:        string(tagsJSON),
		IsPublished: req.IsPublished,
		Sort:        req.Sort,
		CreatedBy:   userID,
	}

	if err := s.db.Create(&series).Error; err != nil {
		return nil, err
	}

	return &series, nil
}

// UpdateVideoSeries 更新视频系列
func (s *EnglishVideoService) UpdateVideoSeries(seriesID uint, req *api.UpdateVideoSeriesRequest) error {
	var series models.VideoSeries
	if err := s.db.Where("id = ?", seriesID).First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("video series not found")
		}
		return err
	}

	updates := map[string]interface{}{
		"title":        req.Title,
		"title_cn":     req.TitleCN,
		"description":  req.Description,
		"cover_image":  req.CoverImage,
		"difficulty":   req.Difficulty,
		"age_range":    req.AgeRange,
		"tags":         req.Tags,
		"is_published": req.IsPublished,
	}

	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}

	return s.db.Model(&series).Updates(updates).Error
}

// DeleteVideoSeries 删除视频系列
func (s *EnglishVideoService) DeleteVideoSeries(seriesID uint) error {
	// 检查是否存在剧集
	var episodeCount int64
	s.db.Model(&models.VideoEpisode{}).Where("series_id = ?", seriesID).Count(&episodeCount)
	if episodeCount > 0 {
		return errors.New("cannot delete series with existing episodes")
	}

	// 删除相关的点赞记录
	s.db.Where("series_id = ?", seriesID).Delete(&models.VideoSeriesLike{})

	// 删除系列
	return s.db.Delete(&models.VideoSeries{}, seriesID).Error
}

// CreateEpisode 创建剧集
func (s *EnglishVideoService) CreateEpisode(seriesID uint, req *api.CreateEpisodeRequest, userID uint) (*models.VideoEpisode, error) {
	// 检查系列是否存在
	var series models.VideoSeries
	if err := s.db.Where("id = ?", seriesID).First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video series not found")
		}
		return nil, err
	}

	episode := models.VideoEpisode{
		SeriesID:    seriesID,
		Title:       req.Title,
		TitleCN:     req.TitleCN,
		Description: req.Description,
		VideoURL:    req.VideoURL,
		Thumbnail:   req.Thumbnail,
		Duration:    req.Duration,
		EpisodeNum:  req.EpisodeNum,
		Subtitles:   req.Subtitles,
		Transcript:  req.Transcript,
		IsPublished: req.IsPublished,
		Sort:        req.Sort,
		CreatedBy:   userID,
	}

	if err := s.db.Create(&episode).Error; err != nil {
		return nil, err
	}

	return &episode, nil
}

// UpdateEpisode 更新剧集
func (s *EnglishVideoService) UpdateEpisode(episodeID uint, req *api.UpdateEpisodeRequest) error {
	var episode models.VideoEpisode
	if err := s.db.Where("id = ?", episodeID).First(&episode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("episode not found")
		}
		return err
	}

	updates := map[string]interface{}{
		"title":        req.Title,
		"title_cn":     req.TitleCN,
		"description":  req.Description,
		"thumbnail":    req.Thumbnail,
		"subtitles":    req.Subtitles,
		"transcript":   req.Transcript,
		"video_url":    req.VideoURL,
		"duration":     req.Duration,
		"episode_num":  req.EpisodeNum,
		"is_published": req.IsPublished,
	}

	if req.Sort > 0 {
		updates["sort"] = req.Sort
	}

	return s.db.Model(&episode).Updates(updates).Error
}

// DeleteEpisode 删除剧集
func (s *EnglishVideoService) DeleteEpisode(episodeID uint) error {
	// 删除相关的进度记录
	s.db.Where("episode_id = ?", episodeID).Delete(&models.VideoUserProgress{})

	// 删除剧集
	return s.db.Delete(&models.VideoEpisode{}, episodeID).Error
}

// BatchImportEpisodes 批量导入剧集
func (s *EnglishVideoService) BatchImportEpisodes(seriesID uint, req *api.BatchImportEpisodesRequest, userID uint) (*api.BatchImportResult, error) {
	// 检查系列是否存在
	var series models.VideoSeries
	if err := s.db.Where("id = ?", seriesID).First(&series).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video series not found")
		}
		return nil, err
	}

	var successCount, failCount int
	var errors []string

	for i, episodeReq := range req.Episodes {
		episode := models.VideoEpisode{
			SeriesID:    seriesID,
			Title:       episodeReq.Title,
			TitleCN:     episodeReq.TitleCN,
			Description: episodeReq.Description,
			VideoURL:    episodeReq.VideoURL,
			Thumbnail:   episodeReq.Thumbnail,
			Duration:    episodeReq.Duration,
			EpisodeNum:  episodeReq.EpisodeNum,
			Subtitles:   episodeReq.Subtitles,
			Transcript:  episodeReq.Transcript,
			IsPublished: episodeReq.IsPublished,
			Sort:        episodeReq.Sort,
			CreatedBy:   userID,
		}

		if err := s.db.Create(&episode).Error; err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("Episode %d: %s", i+1, err.Error()))
		} else {
			successCount++
		}
	}

	return &api.BatchImportResult{
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
	}, nil
}

// CreateUncategorizedEpisode 创建未分类剧集
func (s *EnglishVideoService) CreateUncategorizedEpisode(req *api.CreateEpisodeRequest, userID uint) (*models.VideoEpisode, error) {
	episode := models.VideoEpisode{
		// 注意：这里不设置SeriesID，让其为NULL表示未分类
		Title:       req.Title,
		TitleCN:     req.TitleCN,
		Description: req.Description,
		VideoURL:    req.VideoURL,
		Thumbnail:   req.Thumbnail,
		Duration:    req.Duration,
		EpisodeNum:  req.EpisodeNum,
		Subtitles:   req.Subtitles,
		Transcript:  req.Transcript,
		IsPublished: req.IsPublished,
		Sort:        req.Sort,
		CreatedBy:   userID,
	}

	if err := s.db.Create(&episode).Error; err != nil {
		return nil, err
	}

	return &episode, nil
}

// 辅助函数

func getSortOrder(order string) string {
	if order == "asc" {
		return "ASC"
	}
	return "DESC"
}
