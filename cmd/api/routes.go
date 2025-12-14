package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"khalif-alquran/internal/handler"
	"khalif-alquran/pkg/middleware"

)

func RegisterRoutes(
	r *gin.Engine,
	quranHandler *handler.QuranHandler,
	bookmarkHandler *handler.BookmarkHandler,
) {
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		quran := api.Group("/quran")
		{
			quran.GET("/surahs", quranHandler.GetAllSurahs)
			quran.GET("/surahs/:number", quranHandler.GetSurahDetail)
			quran.GET("/search", quranHandler.Search)
		}

		bookmarks := api.Group("/bookmarks")
		{
			// Nanti ditambahkan middleware Auth di sini jika sudah ada user
			bookmarks.GET("/:user_id", bookmarkHandler.GetUserBookmarks)
			bookmarks.POST("/", bookmarkHandler.AddBookmark)
			bookmarks.DELETE("/", bookmarkHandler.RemoveBookmark)
		}
	}
}