package repository

import (
	"context"

	"gorm.io/gorm"

	"khalif-alquran/internal/domain"

)

// bookmarkRepository struct (private agar dipaksa lewat Interface)
type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository mengembalikan domain.BookmarkRepository (Interface)
func NewBookmarkRepository(db *gorm.DB) domain.BookmarkRepository {
	return &bookmarkRepository{db: db}
}

// SaveBookmark menyimpan atau update penanda ayat untuk user
func (r *bookmarkRepository) SaveBookmark(ctx context.Context, bookmark *domain.Bookmark) error {
	// Menggunakan Save (Upsert) agar jika sudah ada bookmark di ID tersebut, di-update
	return r.db.WithContext(ctx).Save(bookmark).Error
}

// GetByUserID mengambil semua bookmark milik user tertentu
func (r *bookmarkRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Bookmark, error) {
	var bookmarks []domain.Bookmark

	// Preload Surah agar frontend tahu ini bookmark surat apa
	err := r.db.WithContext(ctx).
		Preload("Surah").
		Where("user_id = ?", userID).
		Order("created_at DESC"). // Yang terbaru ditaruh paling atas
		Find(&bookmarks).Error

	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

// DeleteBookmark menghapus bookmark spesifik
func (r *bookmarkRepository) DeleteBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND surah_id = ? AND ayah_number = ?", userID, surahID, ayahNumber).
		Delete(&domain.Bookmark{}).Error
}

// ClearAllBookmarks menghapus semua bookmark user (fitur reset)
func (r *bookmarkRepository) ClearAllBookmarks(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.Bookmark{}).Error
}