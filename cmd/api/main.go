package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"khalif-alquran/internal/config"
	"khalif-alquran/internal/domain"
	"khalif-alquran/internal/handler"
	"khalif-alquran/pkg/database"
	"khalif-alquran/pkg/logger"

)

// @title           Khalif Al-Quran API
// @version         1.0
// @description     API Service for Khalif Al-Quran Application (Offline First)
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@khalifalquran.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type App struct {
	DB              *gorm.DB
	RDB             *redis.Client
	QuranHandler    *handler.QuranHandler
	BookmarkHandler *handler.BookmarkHandler
	Cfg             *config.Config // Tambahkan Config ke App struct jika perlu, tapi opsional
}

func NewApp(db *gorm.DB, rdb *redis.Client, qh *handler.QuranHandler, bh *handler.BookmarkHandler) *App {
	return &App{
		DB:              db,
		RDB:             rdb,
		QuranHandler:    qh,
		BookmarkHandler: bh,
	}
}

func main() {
	logger.Init()

	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()
	
	// Initialize App via Wire
	app, err := InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}

	if *refreshFlag {
		database.ResetSchema(app.DB)
		logger.Info("Database reset successfully")
	}

	// AutoMigrate
	app.DB.AutoMigrate(
		&domain.Surah{},
		&domain.Ayah{},
		&domain.Bookmark{},
	)

	// Seeding Data
	database.RunMigrations(app.DB)
	database.SeedQuran(app.DB)

	r := gin.New()
	r.Use(gin.Recovery())

	// Register Routes
	RegisterRoutes(r, app.QuranHandler, app.BookmarkHandler)

	// PERBAIKAN: Gunakan cfg.Port (sesuai config kamu), bukan cfg.AppPort
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Server start failed", zap.Error(err))
	}
}