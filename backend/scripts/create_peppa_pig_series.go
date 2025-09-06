package main

import (
	"log"

	"gin-web-framework/config"
	"gin-web-framework/internal/database"
	"gin-web-framework/internal/models"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := database.GetDB()

	// Get all songs from English Learning
	var songs []models.Song
	if err := db.Preload("Category").Find(&songs).Error; err != nil {
		log.Fatal("Failed to fetch songs:", err)
	}

	log.Printf("Found %d songs to migrate", len(songs))

	// First, delete existing video series and episodes to clean up
	log.Println("Cleaning up existing video data...")
	if err := db.Exec("DELETE FROM video_user_progress").Error; err != nil {
		log.Printf("Warning: Failed to delete video user progress: %v", err)
	}
	if err := db.Exec("DELETE FROM video_series_like").Error; err != nil {
		log.Printf("Warning: Failed to delete video series likes: %v", err)
	}
	if err := db.Exec("DELETE FROM video_episode").Error; err != nil {
		log.Printf("Warning: Failed to delete video episodes: %v", err)
	}
	if err := db.Exec("DELETE FROM video_series").Error; err != nil {
		log.Printf("Warning: Failed to delete video series: %v", err)
	}

	// Create the Peppa Pig series
	peppaSeries := &models.VideoSeries{
		Title:       "小猪佩奇英语启蒙",
		TitleCN:     "小猪佩奇英语启蒙",
		Description: "跟着小猪佩奇一起学英语！通过有趣的儿歌和故事，让孩子们在快乐中学习英语。",
		CoverImage:  "https://supersimple.com/wp-content/uploads/2025/08/songs_doyouhaveacrayon_thumbnail_website_800x800_en-300x300.jpg",
		Difficulty:  1,
		AgeRange:    "3-8",
		Tags:        `["小猪佩奇", "英语启蒙", "儿歌", "互动学习", "卡通"]`,
		IsPublished: true,
		Sort:        1,
		CreatedBy:   1,
	}

	if err := db.Create(peppaSeries).Error; err != nil {
		log.Fatal("Failed to create Peppa Pig series:", err)
	}

	log.Printf("Created Peppa Pig series: %s (ID: %d)", peppaSeries.Title, peppaSeries.ID)

	// Create video episodes from all songs and put them under Peppa Pig series
	var episodeCount int
	for i, song := range songs {
		episodeCount++

		// Parse tags
		tags := song.Tags
		if tags == "" {
			tags = `["小猪佩奇", "英语启蒙", "儿歌"]`
		}

		// Create episode title with episode number
		episodeTitle := song.Title
		if song.TitleCN != "" {
			episodeTitle = song.TitleCN
		}

		episode := &models.VideoEpisode{
			Title:       episodeTitle,
			TitleCN:     song.TitleCN,
			Description: song.Description,
			VideoURL:    song.VideoURL,
			Thumbnail:   song.CoverImage,
			Duration:    song.Duration,
			EpisodeNum:  episodeCount,
			Subtitles:   song.Lyrics,
			Transcript:  song.Lyrics,
			IsPublished: song.IsPublished,
			Sort:        song.Sort,
			SeriesID:    peppaSeries.ID,
			CreatedBy:   song.CreatedBy,
		}

		if err := db.Create(episode).Error; err != nil {
			log.Printf("Failed to create video episode %s: %v", song.Title, err)
			continue
		}

		log.Printf("Created episode %d: %s", episodeCount, episode.Title)

		// Update series cover image with first episode's thumbnail if not set
		if i == 0 && episode.Thumbnail != "" {
			peppaSeries.CoverImage = episode.Thumbnail
			db.Save(peppaSeries)
			log.Printf("Updated series cover image to: %s", episode.Thumbnail)
		}
	}

	log.Printf("Migration completed! Created Peppa Pig series with %d episodes", episodeCount)
	log.Printf("Series ID: %d", peppaSeries.ID)
	log.Printf("All English songs have been migrated to the Peppa Pig series")
}