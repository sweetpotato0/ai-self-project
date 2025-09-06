package service

import (
	"fmt"
	"gin-web-framework/internal/models"
	"gin-web-framework/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type EnglishLearningService struct {
	db     *gorm.DB
	logger logger.LoggerInterface
}

func NewEnglishLearningService(db *gorm.DB, logger logger.LoggerInterface) *EnglishLearningService {
	return &EnglishLearningService{
		db:     db,
		logger: logger,
	}
}

// ====== 分类管理 ======

type CategoryFilter struct {
	IsActive *bool `json:"is_active,omitempty"`
	Page     int   `json:"page,omitempty"`
	Limit    int   `json:"limit,omitempty"`
}

type PaginatedCategories struct {
	Categories []*models.LearningCategory `json:"categories"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

func (s *EnglishLearningService) GetCategories(filter *CategoryFilter) (*PaginatedCategories, error) {
	var categories []*models.LearningCategory
	var total int64

	query := s.db.Model(&models.LearningCategory{})

	// 应用过滤器
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count categories: %v", err)
	}

	// 分页
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Order("sort ASC, created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}

	// 计算每个分类的歌曲数量
	for _, category := range categories {
		var songCount int64
		s.db.Model(&models.Song{}).Where("category_id = ?", category.ID).Count(&songCount)
		category.SongCount = songCount
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return &PaginatedCategories{
		Categories: categories,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *EnglishLearningService) CreateCategory(category *models.LearningCategory) error {
	if err := s.db.Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}
	return nil
}

func (s *EnglishLearningService) UpdateCategory(id uint, updates map[string]interface{}) error {
	if err := s.db.Model(&models.LearningCategory{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update category: %v", err)
	}
	return nil
}

func (s *EnglishLearningService) DeleteCategory(id uint) error {
	// 检查是否有关联的歌曲
	var count int64
	s.db.Model(&models.Song{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with existing songs")
	}

	if err := s.db.Delete(&models.LearningCategory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}
	return nil
}

// ====== 歌曲/学习材料管理 ======

type SongFilter struct {
	CategoryID  *uint  `json:"category_id,omitempty"`
	Difficulty  *int   `json:"difficulty,omitempty"`
	AgeRange    string `json:"age_range,omitempty"`
	IsPublished *bool  `json:"is_published,omitempty"`
	Search      string `json:"search,omitempty"`
	Tags        string `json:"tags,omitempty"`
	SortBy      string `json:"sort_by,omitempty"`
	SortOrder   string `json:"sort_order,omitempty"`
	Page        int    `json:"page,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

type PaginatedSongs struct {
	Songs      []*models.Song `json:"songs"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

func (s *EnglishLearningService) GetSongs(filter *SongFilter) (*PaginatedSongs, error) {
	var songs []*models.Song
	var total int64

	query := s.db.Model(&models.Song{}).Preload("Category").Preload("Creator")

	// 应用过滤器
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Difficulty != nil {
		query = query.Where("difficulty = ?", *filter.Difficulty)
	}
	if filter.AgeRange != "" {
		query = query.Where("age_range = ?", filter.AgeRange)
	}
	if filter.IsPublished != nil {
		query = query.Where("is_published = ?", *filter.IsPublished)
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR lyrics LIKE ?",
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.Tags != "" {
		query = query.Where("tags LIKE ?", "%\""+filter.Tags+"\"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count songs: %v", err)
	}

	// 排序
	orderBy := "sort ASC, created_at DESC"
	if filter.SortBy != "" {
		validSortFields := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"title":      true,
			"view_count": true,
			"like_count": true,
			"difficulty": true,
			"sort":       true,
		}
		if validSortFields[filter.SortBy] {
			direction := "DESC"
			if filter.SortOrder == "asc" {
				direction = "ASC"
			}
			orderBy = filter.SortBy + " " + direction
		}
	}

	// 分页
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Order(orderBy).Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("failed to get songs: %v", err)
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	return &PaginatedSongs{
		Songs:      songs,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *EnglishLearningService) GetSongByID(id uint) (*models.Song, error) {
	var song models.Song
	if err := s.db.Preload("Category").Preload("Creator").Preload("Vocabularies").First(&song, id).Error; err != nil {
		return nil, fmt.Errorf("song not found: %v", err)
	}

	// 增加浏览量
	s.db.Model(&models.Song{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &song, nil
}

func (s *EnglishLearningService) CreateSong(song *models.Song) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 如果没有分类ID但有分类名称，先创建分类
		if song.CategoryID == nil && song.CategoryName != "" {
			// 检查分类是否已存在
			var existingCategory models.LearningCategory
			err := tx.Where("name_cn = ? OR name = ?", song.CategoryName, song.CategoryName).First(&existingCategory).Error

			if err == gorm.ErrRecordNotFound {
				// 分类不存在，创建新分类
				newCategory := &models.LearningCategory{
					Name:      song.CategoryName,
					NameCN:    song.CategoryName,
					Icon:      "📚",       // 默认图标
					Color:     "#667eea", // 默认颜色
					IsActive:  true,
					Sort:      0,
					CreatedBy: song.CreatedBy,
				}

				if err := tx.Create(newCategory).Error; err != nil {
					return fmt.Errorf("failed to create category: %v", err)
				}

				song.CategoryID = &newCategory.ID
			} else if err != nil {
				return fmt.Errorf("failed to check existing category: %v", err)
			} else {
				// 分类已存在，使用现有分类
				song.CategoryID = &existingCategory.ID
			}
		}

		// 创建歌曲
		if err := tx.Create(song).Error; err != nil {
			return fmt.Errorf("failed to create song: %v", err)
		}

		return nil
	})
}

func (s *EnglishLearningService) UpdateSong(id uint, updates map[string]interface{}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 如果更新数据包含category_name但没有category_id，处理分类创建
		if categoryName, hasCategoryName := updates["category_name"]; hasCategoryName && categoryName != "" {
			if _, hasCategoryID := updates["category_id"]; !hasCategoryID || updates["category_id"] == nil {
				// 检查分类是否已存在
				var existingCategory models.LearningCategory
				categoryNameStr := categoryName.(string)
				err := tx.Where("name_cn = ? OR name = ?", categoryNameStr, categoryNameStr).First(&existingCategory).Error

				if err == gorm.ErrRecordNotFound {
					// 获取当前歌曲信息以获取创建者
					var currentSong models.Song
					if err := tx.First(&currentSong, id).Error; err != nil {
						return fmt.Errorf("failed to find song: %v", err)
					}

					// 分类不存在，创建新分类
					newCategory := &models.LearningCategory{
						Name:      categoryNameStr,
						NameCN:    categoryNameStr,
						Icon:      "📚",       // 默认图标
						Color:     "#667eea", // 默认颜色
						IsActive:  true,
						Sort:      0,
						CreatedBy: currentSong.CreatedBy,
					}

					if err := tx.Create(newCategory).Error; err != nil {
						return fmt.Errorf("failed to create category: %v", err)
					}

					updates["category_id"] = newCategory.ID
				} else if err != nil {
					return fmt.Errorf("failed to check existing category: %v", err)
				} else {
					// 分类已存在，使用现有分类
					updates["category_id"] = existingCategory.ID
				}
			}

			// 删除category_name字段，避免更新到数据库
			delete(updates, "category_name")
		}

		// 更新歌曲
		if err := tx.Model(&models.Song{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update song: %v", err)
		}

		return nil
	})
}

func (s *EnglishLearningService) DeleteSong(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的用户进度
		if err := tx.Where("song_id = ?", id).Delete(&models.UserProgress{}).Error; err != nil {
			return err
		}

		// 删除歌曲词汇关联
		if err := tx.Exec("DELETE FROM song_vocabularies WHERE song_id = ?", id).Error; err != nil {
			return err
		}

		// 删除歌曲
		if err := tx.Delete(&models.Song{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *EnglishLearningService) LikeSong(songID, userID uint) error {
	// 检查是否已存在进度记录
	var progress models.UserProgress
	err := s.db.Where("song_id = ? AND user_id = ?", songID, userID).First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新的进度记录
		progress = models.UserProgress{
			SongID:  songID,
			UserID:  userID,
			IsLiked: true,
		}
		if err := s.db.Create(&progress).Error; err != nil {
			return fmt.Errorf("failed to create user progress: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to query user progress: %v", err)
	} else {
		// 更新现有记录
		if progress.IsLiked {
			return fmt.Errorf("already liked")
		}
		if err := s.db.Model(&progress).UpdateColumn("is_liked", true).Error; err != nil {
			return fmt.Errorf("failed to like song: %v", err)
		}
	}

	// 增加歌曲点赞数
	if err := s.db.Model(&models.Song{}).Where("id = ?", songID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to update like count: %v", err)
	}

	return nil
}

func (s *EnglishLearningService) UnlikeSong(songID, userID uint) error {
	var progress models.UserProgress
	if err := s.db.Where("song_id = ? AND user_id = ? AND is_liked = ?", songID, userID, true).First(&progress).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("not liked yet")
		}
		return fmt.Errorf("failed to query user progress: %v", err)
	}

	// 更新点赞状态
	if err := s.db.Model(&progress).UpdateColumn("is_liked", false).Error; err != nil {
		return fmt.Errorf("failed to unlike song: %v", err)
	}

	// 减少歌曲点赞数
	if err := s.db.Model(&models.Song{}).Where("id = ? AND like_count > 0", songID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to update like count: %v", err)
	}

	return nil
}

// ====== 用户学习进度 ======

func (s *EnglishLearningService) UpdateUserProgress(userID, songID uint, updates map[string]interface{}) error {
	var progress models.UserProgress
	err := s.db.Where("user_id = ? AND song_id = ?", userID, songID).First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		progress = models.UserProgress{
			UserID: userID,
			SongID: songID,
		}

		// 设置默认值和更新字段
		for key, value := range updates {
			switch key {
			case "progress":
				if v, ok := value.(int); ok {
					progress.Progress = v
				}
			case "is_completed":
				if v, ok := value.(bool); ok {
					progress.IsCompleted = v
				}
			case "play_count":
				if v, ok := value.(int); ok {
					progress.PlayCount = v
				}
			case "study_time_minutes":
				if v, ok := value.(int); ok {
					progress.StudyTimeMinutes = v
				}
			case "notes":
				if v, ok := value.(string); ok {
					progress.Notes = v
				}
			}
		}

		now := time.Now()
		progress.LastStudiedAt = &now

		if err := s.db.Create(&progress).Error; err != nil {
			return fmt.Errorf("failed to create user progress: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to query user progress: %v", err)
	} else {
		// 更新现有记录
		now := time.Now()
		updates["last_studied_at"] = &now

		if err := s.db.Model(&progress).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update user progress: %v", err)
		}
	}

	return nil
}

func (s *EnglishLearningService) GetUserProgress(userID uint, songID *uint) ([]models.UserProgress, error) {
	var progress []models.UserProgress
	query := s.db.Preload("Song").Preload("Song.Category").Where("user_id = ?", userID)

	if songID != nil {
		query = query.Where("song_id = ?", *songID)
	}

	if err := query.Find(&progress).Error; err != nil {
		return nil, fmt.Errorf("failed to get user progress: %v", err)
	}

	return progress, nil
}

// ====== 推荐系统 ======

func (s *EnglishLearningService) GetRecommendedSongs(userID uint, limit int) ([]*models.Song, error) {
	var songs []*models.Song

	// 简单的推荐算法：基于用户喜欢的分类和未完成的歌曲
	query := s.db.Model(&models.Song{}).
		Preload("Category").
		Where("is_published = ?", true).
		Where("id NOT IN (?)",
			s.db.Table("user_progress").
				Select("song_id").
				Where("user_id = ? AND is_completed = ?", userID, true))

	if err := query.Order("view_count DESC, like_count DESC").
		Limit(limit).
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("failed to get recommended songs: %v", err)
	}

	return songs, nil
}

// ====== 统计信息 ======

type LearningStats struct {
	TotalSongs        int64  `json:"total_songs"`
	CompletedSongs    int64  `json:"completed_songs"`
	TotalStudyMinutes int64  `json:"total_study_minutes"`
	CurrentStreak     int    `json:"current_streak"`
	FavoriteCategory  string `json:"favorite_category"`
	Level             int    `json:"level"`
}

func (s *EnglishLearningService) GetUserStats(userID uint) (*LearningStats, error) {
	stats := &LearningStats{}

	// 总歌曲数
	s.db.Model(&models.Song{}).Where("is_published = ?", true).Count(&stats.TotalSongs)

	// 已完成歌曲数
	s.db.Model(&models.UserProgress{}).Where("user_id = ? AND is_completed = ?", userID, true).Count(&stats.CompletedSongs)

	// 总学习时长
	s.db.Model(&models.UserProgress{}).Where("user_id = ?", userID).Select("COALESCE(SUM(study_time_minutes), 0)").Scan(&stats.TotalStudyMinutes)

	// 计算等级（基于学习时长）
	stats.Level = int(stats.TotalStudyMinutes/60) + 1 // 每60分钟升一级

	// 获取最喜欢的分类
	var categoryName string
	s.db.Table("user_progress").
		Select("lc.name").
		Joins("JOIN songs s ON user_progress.song_id = s.id").
		Joins("JOIN learning_categories lc ON s.category_id = lc.id").
		Where("user_progress.user_id = ?", userID).
		Group("lc.name").
		Order("COUNT(*) DESC").
		Limit(1).
		Scan(&categoryName)
	stats.FavoriteCategory = categoryName

	return stats, nil
}
