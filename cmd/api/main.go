package main

import (
	"context" // Tambahkan import context
	"flag"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"khalif-alquran/internal/config"
	"khalif-alquran/internal/domain"
	"khalif-alquran/internal/handler"
	grpcHandler "khalif-alquran/internal/handler/grpc" // Alias untuk membedakan dengan handler HTTP
	"khalif-alquran/pkg/database"
	"khalif-alquran/pkg/logger"
	"khalif-alquran/pkg/pb"

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
	DB               *gorm.DB
	RDB              *redis.Client
	QuranHandler     *handler.QuranHandler
	BookmarkHandler  *handler.BookmarkHandler
	GrpcQuranHandler *grpcHandler.QuranHandler // Field baru untuk gRPC Handler
	Cfg              *config.Config
}

// NewApp diperbarui untuk menerima gRPC Handler dari Wire
func NewApp(
	db *gorm.DB,
	rdb *redis.Client,
	qh *handler.QuranHandler,
	bh *handler.BookmarkHandler,
	gqh *grpcHandler.QuranHandler, // Parameter baru
) *App {
	return &App{
		DB:               db,
		RDB:              rdb,
		QuranHandler:     qh,
		BookmarkHandler:  bh,
		GrpcQuranHandler: gqh, // Assign ke struct
	}
}

func main() {
	logger.Init()

	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()

	// Initialize App via Wire (dependency injection)
	app, err := InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}

	if *refreshFlag {
		// 1. Reset Database PostgreSQL
		database.ResetSchema(app.DB)
		logger.Info("Database reset successfully")

		// 2. OTOMATIS: Hapus Semua Cache Redis (Tambahan Baru)
		if app.RDB != nil {
			if err := app.RDB.FlushAll(context.Background()).Err(); err != nil {
				logger.Error("Failed to flush redis", zap.Error(err))
			} else {
				logger.Info("Redis cache cleared successfully (Auto-Flush)")
			}
		}
	}

	// AutoMigrate Database Tables
	app.DB.AutoMigrate(
		&domain.Surah{},
		&domain.Ayah{},
		&domain.Bookmark{},
	)

	// Seeding Data
	database.RunMigrations(app.DB)
	database.SeedQuran(app.DB)

	// --- Jalankan gRPC Server (Concurrent) ---
	go func() {
		grpcPort := ":50051" // Port khusus untuk gRPC
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			logger.Fatal("Failed to listen grpc", zap.Error(err))
		}

		s := grpc.NewServer()

		// Register Service gRPC ke Server
		pb.RegisterQuranServiceServer(s, app.GrpcQuranHandler)

		logger.Info("gRPC Server starting", zap.String("port", grpcPort))
		if err := s.Serve(lis); err != nil {
			logger.Fatal("Failed to serve grpc", zap.Error(err))
		}
	}()

	// --- Jalankan HTTP Server (Gin) ---
	r := gin.New()
	r.Use(gin.Recovery())

	// Register Routes HTTP
	RegisterRoutes(r, app.QuranHandler, app.BookmarkHandler)

	// Tentukan Port HTTP
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	logger.Info("HTTP Server starting", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Server start failed", zap.Error(err))
	}
}