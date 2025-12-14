package main

import (
	"log"
	"time"
	"khalif-alquran/internal/config"
	"khalif-alquran/pkg/database"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

)

func ProvideDB(cfg *config.Config) *gorm.DB {
	// Memastikan database fisik ada (CREATE DATABASE if not exists)
	database.EnsureDBExists(cfg.DBUrl)

	// Set logger ke mode error agar log tidak terlalu berisik
	dbLogger := logger.Default.LogMode(logger.Error)

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true, // Cache prepared statement untuk performa
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}

	// Konfigurasi Connection Pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		// Password: "", // Set jika ada password di config
		// DB:       0,  // Gunakan DB default
	})
}