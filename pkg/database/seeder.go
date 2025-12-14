package database

import (
	"encoding/json"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"khalif-alquran/internal/domain"
	"khalif-alquran/pkg/logger"

)

func SeedQuran(db *gorm.DB) {
	var count int64
	db.Model(&domain.Surah{}).Count(&count)
	if count > 0 {
		logger.Info("Surahs already seeded, skipping...")
		return
	}

	// Sesuaikan path dengan lokasi file di Docker
	pattern := "pkg/database/seeds/data/*.json"
	files, err := filepath.Glob(pattern)
	if err != nil {
		logger.Error("Failed to list seed files", zap.Error(err))
		return
	}

	if len(files) == 0 {
		// PERBAIKAN: Ganti logger.Warn jadi logger.Info
		logger.Info("No seed files found in " + pattern) 
		return
	}

	logger.Info("Start seeding Surahs...", zap.Int("files_found", len(files)))

	for _, filename := range files {
		fileData, err := os.ReadFile(filename)
		if err != nil {
			logger.Error("Failed to read seed file", zap.String("file", filename), zap.Error(err))
			continue
		}

		var surah domain.Surah
		if err := json.Unmarshal(fileData, &surah); err != nil {
			logger.Error("Failed to parse json", zap.String("file", filename), zap.Error(err))
			continue
		}

		if err := db.Create(&surah).Error; err != nil {
			logger.Error("Failed to insert surah to DB", zap.String("surah", surah.Name), zap.Error(err))
			return
		}

		logger.Info("Seeded: " + surah.LatinName)
	}

	logger.Info("Database seeding completed.")
}