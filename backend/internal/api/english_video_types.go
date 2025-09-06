package api

import "gin-web-framework/internal/models"

// Video Series Types

// VideoSeriesListRequest 视频系列列表请求
type VideoSeriesListRequest struct {
	PaginationRequest
	Search     string `json:"search" form:"search"`
	Difficulty int    `json:"difficulty" form:"difficulty"`
	AgeRange   string `json:"age_range" form:"age_range"`
	Tag        string `json:"tag" form:"tag"`
	SortBy     string `json:"sort_by" form:"sort_by"`
	SortOrder  string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc"`
}

// VideoSeriesListResponse 视频系列列表响应
type VideoSeriesListResponse struct {
	List     []models.VideoSeries `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// CreateVideoSeriesRequest 创建视频系列请求
type CreateVideoSeriesRequest struct {
	Title       string   `json:"title" binding:"required"`
	TitleCN     string   `json:"title_cn"`
	Description string   `json:"description"`
	CoverImage  string   `json:"cover_image"`
	Difficulty  int      `json:"difficulty" binding:"min=1,max=5"`
	AgeRange    string   `json:"age_range"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
	Sort        int      `json:"sort"`
}

// UpdateVideoSeriesRequest 更新视频系列请求
type UpdateVideoSeriesRequest struct {
	Title       string   `json:"title"`
	TitleCN     string   `json:"title_cn"`
	Description string   `json:"description"`
	CoverImage  string   `json:"cover_image"`
	Difficulty  int      `json:"difficulty" binding:"omitempty,min=1,max=5"`
	AgeRange    string   `json:"age_range"`
	Tags        []string `json:"tags"`
	IsPublished *bool    `json:"is_published"`
	Sort        int      `json:"sort"`
}

// SearchVideoSeriesRequest 搜索视频系列请求
type SearchVideoSeriesRequest struct {
	Query      string `json:"query" form:"q"`
	Difficulty int    `json:"difficulty" form:"difficulty"`
	AgeRange   string `json:"age_range" form:"age_range"`
	Tag        string `json:"tag" form:"tag"`
}

// Episode Types

// EpisodeListRequest 剧集列表请求
type EpisodeListRequest struct {
	PaginationRequest
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc"`
}

// EpisodeListResponse 剧集列表响应
type EpisodeListResponse struct {
	List     []models.VideoEpisode `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Series   *models.VideoSeries   `json:"series"`
}

// CreateEpisodeRequest 创建剧集请求
type CreateEpisodeRequest struct {
	Title       string `json:"title" binding:"required"`
	TitleCN     string `json:"title_cn"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	Duration    int    `json:"duration"`
	EpisodeNum  int    `json:"episode_num" binding:"required,min=1"`
	Subtitles   string `json:"subtitles"`
	Transcript  string `json:"transcript"`
	IsPublished bool   `json:"is_published"`
	Sort        int    `json:"sort"`
}

// UpdateEpisodeRequest 更新剧集请求
type UpdateEpisodeRequest struct {
	Title       string `json:"title"`
	TitleCN     string `json:"title_cn"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url"`
	Thumbnail   string `json:"thumbnail"`
	Duration    int    `json:"duration"`
	EpisodeNum  int    `json:"episode_num" binding:"omitempty,min=1"`
	Subtitles   string `json:"subtitles"`
	Transcript  string `json:"transcript"`
	IsPublished *bool  `json:"is_published"`
	Sort        int    `json:"sort"`
}

// BatchImportEpisodesRequest 批量导入剧集请求
type BatchImportEpisodesRequest struct {
	Episodes []CreateEpisodeRequest `json:"episodes" binding:"required,dive"`
}

// BatchImportResult 批量导入结果
type BatchImportResult struct {
	SuccessCount int      `json:"success_count"`
	FailCount    int      `json:"fail_count"`
	Errors       []string `json:"errors,omitempty"`
}

// Progress Types

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	Progress         int   `json:"progress" binding:"min=0,max=100"`
	CurrentTime      int   `json:"current_time" binding:"min=0"`
	WatchTimeMinutes int   `json:"watch_time_minutes" binding:"min=0"`
	IsCompleted      *bool `json:"is_completed"`
	Notes            string `json:"notes"`
}

// ToggleLikeResponse 切换点赞响应
type ToggleLikeResponse struct {
	IsLiked   bool `json:"is_liked"`
	LikeCount int  `json:"like_count"`
}

// VideoStatsResponse 视频统计响应
type VideoStatsResponse struct {
	TotalSeries   int64 `json:"total_series"`   // 总系列数
	TotalEpisodes int64 `json:"total_episodes"` // 总集数
	TotalViews    int64 `json:"total_views"`    // 总观看次数
	WatchedSeries int64 `json:"watched_series"` // 已看完系列数
}

// UserVideoStatsResponse 用户视频统计响应
type UserVideoStatsResponse struct {
	TotalWatchTime     int   `json:"total_watch_time"`     // 总观看时长(分钟)
	CompletedEpisodes  int64 `json:"completed_episodes"`   // 完成的剧集数
	FavoriteSeries     int64 `json:"favorite_series"`      // 收藏的系列数
	LearningStreak     int   `json:"learning_streak"`      // 学习连续天数
}