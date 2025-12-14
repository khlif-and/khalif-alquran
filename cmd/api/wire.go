//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"khalif-alquran/internal/config"
	"khalif-alquran/internal/domain"
	"khalif-alquran/internal/handler"
	"khalif-alquran/internal/repository"
	"khalif-alquran/internal/usecase"

)

func InitializeApp() (*App, error) {
	wire.Build(
		// 1. Config & Infra
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,

		// 2. Repositories (Perhatikan nama struct sekarang diawali Huruf Besar)
		repository.NewSurahRepository,
		repository.NewAyahRepository,
		repository.NewRedisRepository,
		repository.NewBookmarkRepository,

		// 3. Bind Interface ke Struct (Public)
		wire.Bind(new(domain.SurahRepository), new(*repository.SurahRepository)),
		wire.Bind(new(domain.AyahRepository), new(*repository.AyahRepository)),
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepository)),
		wire.Bind(new(domain.BookmarkRepository), new(*repository.BookmarkRepository)),

		// 4. UseCases
		usecase.NewQuranUseCase,
		usecase.NewBookmarkUseCase,

		// 5. Bind UseCases
		wire.Bind(new(domain.QuranUseCase), new(*usecase.QuranUC)),
		wire.Bind(new(domain.BookmarkUseCase), new(*usecase.BookmarkUC)),

		// 6. Handlers
		handler.NewQuranHandler,
		handler.NewBookmarkHandler,

		// 7. App
		NewApp,
	)
	return &App{}, nil
}