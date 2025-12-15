package usecase

import (
	"context"
	"time"

	"khalif-alquran/internal/domain"

)

type BookmarkUC struct {
	bookmarkRepo domain.BookmarkRepository
	timeout      time.Duration
}

func NewBookmarkUseCase(repo domain.BookmarkRepository) *BookmarkUC {
	return &BookmarkUC{
		bookmarkRepo: repo,
		timeout:      time.Second * 2,
	}
}

func (u *BookmarkUC) AddBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int, note string) error {
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

func (u *BookmarkUC) GetUserBookmarks(ctx context.Context, userID string) ([]domain.Bookmark, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.GetByUserID(ctx, userID)
}

func (u *BookmarkUC) RemoveBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.DeleteBookmark(ctx, userID, surahID, ayahNumber)
}

func (u *BookmarkUC) ClearBookmarks(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.bookmarkRepo.ClearAllBookmarks(ctx, userID)
}