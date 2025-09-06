package models

import (
	"time"
)

// LearningCategory 学习分类
type LearningCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:100"`        // 分类名称
	NameCN      string    `json:"name_cn" gorm:"not null;size:100"`     // 中文名称
	Description string    `json:"description" gorm:"size:500"`          // 描述
	Icon        string    `json:"icon" gorm:"size:100"`                 // 图标
	Color       string    `json:"color" gorm:"size:20;default:#409EFF"` // 主题色
	Sort        int       `json:"sort" gorm:"default:0"`                // 排序
	IsActive    bool      `json:"is_active" gorm:"default:true"`        // 是否激活
	CreatedBy   uint      `json:"created_by"`                           // 创建者ID
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 统计字段
	SongCount int64 `json:"song_count" gorm:"-"`  // 歌曲数量，不存储在数据库中

	// 关联
	Songs []Song `json:"songs,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName 指定表名
func (LearningCategory) TableName() string {
	return "learning_categories"
}

// Song 歌曲/学习材料
type Song struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null;size:200"`       // 标题
	TitleCN      string    `json:"title_cn" gorm:"size:200"`             // 中文标题
	Description  string    `json:"description" gorm:"size:1000"`         // 描述
	Lyrics       string    `json:"lyrics" gorm:"type:text"`              // 歌词
	LyricsCN     string    `json:"lyrics_cn" gorm:"type:text"`           // 中文歌词
	VideoURL     string    `json:"video_url" gorm:"size:500"`            // 视频链接
	AudioURL     string    `json:"audio_url" gorm:"size:500"`            // 音频链接
	CoverImage   string    `json:"cover_image" gorm:"size:500"`          // 封面图片
	Duration     int       `json:"duration" gorm:"default:0"`            // 时长(秒)
	Difficulty   int       `json:"difficulty" gorm:"default:1"`          // 难度等级(1-5)
	AgeRange     string    `json:"age_range" gorm:"size:50"`             // 适合年龄段
	Tags         string    `json:"tags" gorm:"size:500"`                 // 标签(JSON数组)
	ViewCount    int       `json:"view_count" gorm:"default:0"`          // 播放次数
	LikeCount    int       `json:"like_count" gorm:"default:0"`          // 点赞数
	IsPublished  bool      `json:"is_published" gorm:"default:false"`    // 是否发布
	Sort         int       `json:"sort" gorm:"default:0"`                // 排序
	CategoryID   *uint     `json:"category_id"`                          // 分类ID
	CreatedBy    uint      `json:"created_by"`                           // 创建者ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 临时字段，用于接收前端传递的分类名称
	CategoryName string `json:"category_name,omitempty" gorm:"-"`

	// 关联
	Category    *LearningCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Creator     *User             `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Vocabularies []Vocabulary     `json:"vocabularies,omitempty" gorm:"many2many:song_vocabularies;"`
	UserProgress []UserProgress   `json:"user_progress,omitempty" gorm:"foreignKey:SongID"`
}

// Vocabulary 词汇表
type Vocabulary struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Word         string    `json:"word" gorm:"not null;size:100;index"`     // 单词
	Pronunciation string   `json:"pronunciation" gorm:"size:200"`           // 发音(音标)
	PartOfSpeech string    `json:"part_of_speech" gorm:"size:50"`           // 词性
	Definition   string    `json:"definition" gorm:"size:500"`              // 英文释义
	DefinitionCN string    `json:"definition_cn" gorm:"size:500"`           // 中文释义
	Example      string    `json:"example" gorm:"size:1000"`                // 例句
	ExampleCN    string    `json:"example_cn" gorm:"size:1000"`             // 中文例句
	AudioURL     string    `json:"audio_url" gorm:"size:500"`               // 发音音频
	ImageURL     string    `json:"image_url" gorm:"size:500"`               // 配图
	Difficulty   int       `json:"difficulty" gorm:"default:1"`             // 难度等级
	Frequency    int       `json:"frequency" gorm:"default:0"`              // 使用频率
	Tags         string    `json:"tags" gorm:"size:500"`                    // 标签
	CreatedBy    uint      `json:"created_by"`                              // 创建者ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Creator *User `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Songs   []Song `json:"songs,omitempty" gorm:"many2many:song_vocabularies;"`
}

// UserProgress 用户学习进度
type UserProgress struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"index"`                    // 用户ID
	SongID           uint      `json:"song_id" gorm:"index"`                    // 歌曲ID
	Progress         int       `json:"progress" gorm:"default:0"`               // 进度百分比(0-100)
	IsCompleted      bool      `json:"is_completed" gorm:"default:false"`       // 是否完成
	IsLiked          bool      `json:"is_liked" gorm:"default:false"`           // 是否喜欢
	PlayCount        int       `json:"play_count" gorm:"default:0"`             // 播放次数
	StudyTimeMinutes int       `json:"study_time_minutes" gorm:"default:0"`     // 学习时长(分钟)
	LastStudiedAt    *time.Time `json:"last_studied_at"`                        // 最后学习时间
	Notes            string    `json:"notes" gorm:"type:text"`                  // 笔记 (支持多条笔记JSON格式)
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Song *Song `json:"song,omitempty" gorm:"foreignKey:SongID"`
}

// UserVocabulary 用户词汇掌握情况
type UserVocabulary struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index"`                      // 用户ID
	VocabularyID  uint      `json:"vocabulary_id" gorm:"index"`                // 词汇ID
	MasteryLevel  int       `json:"mastery_level" gorm:"default:0"`            // 掌握程度(0-5)
	ReviewCount   int       `json:"review_count" gorm:"default:0"`             // 复习次数
	CorrectCount  int       `json:"correct_count" gorm:"default:0"`            // 正确次数
	LastReviewAt  *time.Time `json:"last_review_at"`                           // 最后复习时间
	NextReviewAt  *time.Time `json:"next_review_at"`                           // 下次复习时间
	IsMarked      bool      `json:"is_marked" gorm:"default:false"`            // 是否标记
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	User       *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Vocabulary *Vocabulary `json:"vocabulary,omitempty" gorm:"foreignKey:VocabularyID"`
}

// LearningPlan 学习计划
type LearningPlan struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"index"`                    // 用户ID
	Name         string    `json:"name" gorm:"not null;size:200"`          // 计划名称
	Description  string    `json:"description" gorm:"size:1000"`           // 描述
	TargetLevel  int       `json:"target_level" gorm:"default:1"`          // 目标等级
	DailyGoal    int       `json:"daily_goal" gorm:"default:30"`           // 每日目标(分钟)
	WeeklyGoal   int       `json:"weekly_goal" gorm:"default:210"`         // 每周目标(分钟)
	StartDate    time.Time `json:"start_date"`                             // 开始日期
	EndDate      *time.Time `json:"end_date"`                              // 结束日期
	IsActive     bool      `json:"is_active" gorm:"default:true"`          // 是否激活
	Progress     int       `json:"progress" gorm:"default:0"`              // 进度百分比
	TotalMinutes int       `json:"total_minutes" gorm:"default:0"`         // 总学习时长
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// StudySession 学习会话
type StudySession struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"index"`                 // 用户ID
	SongID        *uint     `json:"song_id,omitempty"`                    // 歌曲ID(可选)
	SessionType   string    `json:"session_type" gorm:"size:50"`          // 会话类型(song|vocabulary|quiz)
	DurationMinutes int     `json:"duration_minutes"`                     // 时长
	StartTime     time.Time `json:"start_time"`                          // 开始时间
	EndTime       *time.Time `json:"end_time"`                           // 结束时间
	CompletionRate int      `json:"completion_rate" gorm:"default:0"`    // 完成率
	Score         int       `json:"score" gorm:"default:0"`              // 得分
	Notes         string    `json:"notes" gorm:"size:1000"`              // 笔记
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Song *Song `json:"song,omitempty" gorm:"foreignKey:SongID"`
}

// VideoSeries 视频系列
type VideoSeries struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null;size:200"`       // 系列标题
	TitleCN     string    `json:"title_cn" gorm:"size:200"`             // 中文标题
	Description string    `json:"description" gorm:"size:1000"`         // 描述
	CoverImage  string    `json:"cover_image" gorm:"size:500"`          // 封面图片
	Difficulty  int       `json:"difficulty" gorm:"default:1"`          // 难度等级(1-5)
	AgeRange    string    `json:"age_range" gorm:"size:50"`             // 适合年龄段
	Tags        string    `json:"tags" gorm:"size:500"`                 // 标签(JSON数组)
	ViewCount   int       `json:"view_count" gorm:"default:0"`          // 观看次数
	LikeCount   int       `json:"like_count" gorm:"default:0"`          // 点赞数
	IsPublished bool      `json:"is_published" gorm:"default:false"`    // 是否发布
	Sort        int       `json:"sort" gorm:"default:0"`                // 排序
	CreatedBy   uint      `json:"created_by"`                           // 创建者ID
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 统计字段
	EpisodeCount int64 `json:"episode_count" gorm:"-"`  // 剧集数量，不存储在数据库中
	IsLiked      bool  `json:"is_liked" gorm:"-"`       // 用户是否点赞，不存储在数据库中

	// 关联
	Creator  *User          `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Episodes []VideoEpisode `json:"episodes,omitempty" gorm:"foreignKey:SeriesID"`
}

// TableName 指定表名
func (VideoSeries) TableName() string {
	return "video_series"
}

// VideoEpisode 视频剧集
type VideoEpisode struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	SeriesID     uint      `json:"series_id" gorm:"not null;index"`      // 系列ID
	Title        string    `json:"title" gorm:"not null;size:200"`       // 剧集标题
	TitleCN      string    `json:"title_cn" gorm:"size:200"`             // 中文标题
	Description  string    `json:"description" gorm:"size:1000"`         // 描述
	VideoURL     string    `json:"video_url" gorm:"size:500"`            // 视频链接
	Thumbnail    string    `json:"thumbnail" gorm:"size:500"`            // 缩略图
	Duration     int       `json:"duration" gorm:"default:0"`            // 时长(秒)
	EpisodeNum   int       `json:"episode_num" gorm:"default:1"`         // 剧集编号
	Subtitles    string    `json:"subtitles" gorm:"type:text"`           // 字幕(JSON格式)
	Transcript   string    `json:"transcript" gorm:"type:text"`          // 文字稿
	ViewCount    int       `json:"view_count" gorm:"default:0"`          // 观看次数
	IsPublished  bool      `json:"is_published" gorm:"default:false"`    // 是否发布
	Sort         int       `json:"sort" gorm:"default:0"`                // 排序
	CreatedBy    uint      `json:"created_by"`                           // 创建者ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Series    *VideoSeries          `json:"series,omitempty" gorm:"foreignKey:SeriesID"`
	Creator   *User                 `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Progress  []VideoUserProgress   `json:"progress,omitempty" gorm:"foreignKey:EpisodeID"`
}

// TableName 指定表名
func (VideoEpisode) TableName() string {
	return "video_episodes"
}

// VideoUserProgress 用户视频观看进度
type VideoUserProgress struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	UserID           uint       `json:"user_id" gorm:"index"`                    // 用户ID
	SeriesID         uint       `json:"series_id" gorm:"index"`                  // 系列ID
	EpisodeID        uint       `json:"episode_id" gorm:"index"`                 // 剧集ID
	Progress         int        `json:"progress" gorm:"default:0"`               // 进度百分比(0-100)
	CurrentTime      int        `json:"current_time" gorm:"default:0"`           // 当前播放时间(秒)
	IsCompleted      bool       `json:"is_completed" gorm:"default:false"`       // 是否完成
	WatchTimeMinutes int        `json:"watch_time_minutes" gorm:"default:0"`     // 观看时长(分钟)
	LastWatchedAt    *time.Time `json:"last_watched_at"`                         // 最后观看时间
	Notes            string     `json:"notes" gorm:"type:text"`                  // 笔记
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// 计算字段（不存储在数据库中）
	WatchProgress    float64 `json:"watch_progress" gorm:"-"`    // 观看进度百分比(0.0-1.0)
	CurrentEpisode   int     `json:"current_episode" gorm:"-"`   // 当前观看的集数

	// 关联
	User    *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Series  *VideoSeries  `json:"series,omitempty" gorm:"foreignKey:SeriesID"`
	Episode *VideoEpisode `json:"episode,omitempty" gorm:"foreignKey:EpisodeID"`
}

// TableName 指定表名
func (VideoUserProgress) TableName() string {
	return "video_user_progress"
}

// VideoSeriesLike 用户视频系列点赞
type VideoSeriesLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`                      // 用户ID
	SeriesID  uint      `json:"series_id" gorm:"index"`                    // 系列ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User   *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Series *VideoSeries `json:"series,omitempty" gorm:"foreignKey:SeriesID"`
}

// TableName 指定表名
func (VideoSeriesLike) TableName() string {
	return "video_series_likes"
}