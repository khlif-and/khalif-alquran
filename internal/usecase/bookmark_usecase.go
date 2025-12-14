package usecase

import (
	"context"
	"time"

	"khalif-alquran/internal/domain"

)

type bookmarkUseCase struct {
	bookmarkRepo domain.BookmarkRepository
	timeout      time.Duration
}

// PERBAIKAN: Hapus parameter timeout dari sini agar cocok dengan wire_gen.go
func NewBookmarkUseCase(repo domain.BookmarkRepository) domain.BookmarkUseCase {
	return &bookmarkUseCase{
		bookmarkRepo: repo,
		timeout:      time.Second * 2, // Set default timeout 2 detik (Hardcoded)
	}
}

func (u *bookmarkUseCase) AddBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int, note string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	bookmark := &domain.Bookmark{
		UserID:     userID,
		SurahID:    surahID,
		AyahNumber: ayahNumber,
		Note:       note,
	}

	return u.bookmarkRepo.SaveBookmark(ctx, bookmark)
}

func (u *bookmarkUseCase) GetUserBookmarks(ctx context.Context, userID string) ([]domain.Bookmark, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.GetByUserID(ctx, userID)
}

func (u *bookmarkUseCase) RemoveBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.DeleteBookmark(ctx, userID, surahID, ayahNumber)
}

func (u *bookmarkUseCase) ClearBookmarks(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.ClearAllBookmarks(ctx, userID)
}