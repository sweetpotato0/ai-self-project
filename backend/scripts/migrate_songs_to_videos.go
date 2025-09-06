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

	// Create video series based on categories
	categoryToSeries := make(map[uint]*models.VideoSeries)
	
	// Create series for each unique category
	var categories []models.LearningCategory
	if err := db.Find(&categories).Error; err != nil {
		log.Fatal("Failed to fetch categories:", err)
	}

	for _, category := range categories {
		series := &models.VideoSeries{
			Title:       "Super Simple Songs - " + category.Name,
			TitleCN:     "超级简单" + category.NameCN,
			Description: "A collection of educational songs for learning English through " + category.Name + ".",
			CoverImage:  "", // Will be set to first episode's cover
			Difficulty:  1,
			AgeRange:    "3-6",
			Tags:        `["儿歌", "启蒙", "互动"]`,
			IsPublished: true,
			Sort:        int(category.ID),
			CreatedBy:   1,
		}

		if err := db.Create(series).Error; err != nil {
			log.Printf("Failed to create video series for category %s: %v", category.Name, err)
			continue
		}

		categoryToSeries[category.ID] = series
		log.Printf("Created video series: %s (ID: %d)", series.Title, series.ID)
	}

	// Create video episodes from songs
	for i, song := range songs {
		var seriesID uint
		if song.CategoryID != nil && categoryToSeries[*song.CategoryID] != nil {
			seriesID = categoryToSeries[*song.CategoryID].ID
		} else {
			// Use first series as default
			for _, series := range categoryToSeries {
				seriesID = series.ID
				break
			}
		}

		if seriesID == 0 {
			log.Printf("Skipping song %s - no series available", song.Title)
			continue
		}

		// Parse tags
		tags := song.Tags
		if tags == "" {
			tags = `["儿歌", "启蒙"]`
		}

		episode := &models.VideoEpisode{
			Title:       song.Title,
			TitleCN:     song.TitleCN,
			Description: song.Description,
			VideoURL:    song.VideoURL,
			Thumbnail:   song.CoverImage,
			Duration:    song.Duration,
			EpisodeNum:  i + 1,
			Subtitles:   song.Lyrics,
			Transcript:  song.Lyrics,
			IsPublished: song.IsPublished,
			Sort:        song.Sort,
			SeriesID:    seriesID,
			CreatedBy:   song.CreatedBy,
		}

		if err := db.Create(episode).Error; err != nil {
			log.Printf("Failed to create video episode %s: %v", song.Title, err)
			continue
		}

		log.Printf("Created video episode: %s (Episode %d)", episode.Title, episode.EpisodeNum)

		// Update series cover image with first episode's thumbnail
		if episode.Thumbnail != "" {
			series := categoryToSeries[*song.CategoryID]
			if series != nil && series.CoverImage == "" {
				series.CoverImage = episode.Thumbnail
				db.Save(series)
			}
		}
	}

	log.Printf("Migration completed! Created %d video series and %d episodes", len(categoryToSeries), len(songs))
}