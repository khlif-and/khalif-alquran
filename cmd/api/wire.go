//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"khalif-alquran/internal/config"
	"khalif-alquran/internal/domain"
	"khalif-alquran/internal/handler"
	grpcHandler "khalif-alquran/internal/handler/grpc"
	"khalif-alquran/internal/repository"
	"khalif-alquran/internal/usecase"

)

func InitializeApp() (*App, error) {
	wire.Build(
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,

		repository.NewSurahRepository,
		repository.NewAyahRepository,
		repository.NewRedisRepository,
		repository.NewBookmarkRepository,

		wire.Bind(new(domain.SurahRepository), new(*repository.SurahRepository)),
		wire.Bind(new(domain.AyahRepository), new(*repository.AyahRepository)),
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepository)),
		wire.Bind(new(domain.BookmarkRepository), new(*repository.BookmarkRepository)),

		usecase.NewQuranUseCase,
		usecase.NewBookmarkUseCase,

		wire.Bind(new(domain.QuranUseCase), new(*usecase.QuranUC)),
		wire.Bind(new(domain.BookmarkUseCase), new(*usecase.BookmarkUC)),

		handler.NewQuranHandler,
		handler.NewBookmarkHandler,
		grpcHandler.NewQuranHandler,

		NewApp,
	)
	return &App{}, nil
}