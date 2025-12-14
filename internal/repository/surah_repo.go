package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"khalif-alquran/internal/domain"

)

// SurahRepository diubah menjadi huruf besar (Public)
type SurahRepository struct {
	db *gorm.DB
}

// NewSurahRepository mengembalikan *SurahRepository (Public Struct Pointer)
func NewSurahRepository(db *gorm.DB) *SurahRepository {
	return &SurahRepository{db: db}
}

func (r *SurahRepository) GetAll(ctx context.Context) ([]domain.Surah, error) {
	var surahs []domain.Surah
	err := r.db.WithContext(ctx).
		Order("number ASC").
		Find(&surahs).Error
	return surahs, err
}

func (r *SurahRepository) GetByNumber(ctx context.Context, number int) (*domain.Surah, error) {
	var surah domain.Surah

	err := r.db.WithContext(ctx).
		Preload("Ayahs", func(db *gorm.DB) *gorm.DB {
			return db.Order("ayahs.number ASC")
		}).
		Where("number = ?", number).
		First(&surah).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &surah, nil
}

func (r *SurahRepository) Search(ctx context.Context, query string) ([]domain.Surah, error) {
	var surahs []domain.Surah
	searchQuery := "%" + query + "%"

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR latin_name ILIKE ? OR english_name ILIKE ?", searchQuery, searchQuery, searchQuery).
		Order("number ASC").
		Find(&surahs).Error

	return surahs, err
}