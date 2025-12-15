package repository

import (
	"context"

	"gorm.io/gorm"

	"khalif-alquran/internal/domain"

)

type BookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) SaveBookmark(ctx context.Context, bookmark *domain.Bookmark) error {
	return r.db.WithContext(ctx).Save(bookmark).Error
}

func (r *BookmarkRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Bookmark, error) {
	var bookmarks []domain.Bookmark

	err := r.db.WithContext(ctx).
		Preload("Surah").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookmarks).Error

	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r *BookmarkRepository) DeleteBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND surah_id = ? AND ayah_number = ?", userID, surahID, ayahNumber).
		Delete(&domain.Bookmark{}).Error
}

func (r *BookmarkRepository) ClearAllBookmarks(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.Bookmark{}).Error
}