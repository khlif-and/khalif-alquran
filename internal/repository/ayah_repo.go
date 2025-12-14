package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"khalif-alquran/internal/domain"

)

// AyahRepository diubah menjadi huruf besar (Public) agar bisa di-bind oleh Wire
type AyahRepository struct {
	db *gorm.DB
}

// NewAyahRepository mengembalikan *AyahRepository (Public Struct Pointer)
func NewAyahRepository(db *gorm.DB) *AyahRepository {
	return &AyahRepository{db: db}
}

func (r *AyahRepository) GetBySurahID(ctx context.Context, surahID uint) ([]domain.Ayah, error) {
	var ayahs []domain.Ayah
	err := r.db.WithContext(ctx).
		Where("surah_id = ?", surahID).
		Order("number ASC").
		Find(&ayahs).Error

	if err != nil {
		return nil, err
	}
	return ayahs, nil
}

func (r *AyahRepository) GetSpecificAyah(ctx context.Context, surahNumber, ayahNumber int) (*domain.Ayah, error) {
	var ayah domain.Ayah
	// Menggunakan Join untuk mencari berdasarkan Nomor Surat (bukan ID) dan Nomor Ayat
	err := r.db.WithContext(ctx).
		Joins("JOIN surahs ON surahs.number = ayahs.surah_id"). // Asumsi surah_id di table ayahs merujuk ke surahs.number
		Where("surahs.number = ? AND ayahs.number = ?", surahNumber, ayahNumber).
		First(&ayah).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &ayah, nil
}

func (r *AyahRepository) Search(ctx context.Context, query string) ([]domain.Ayah, error) {
	var ayahs []domain.Ayah
	searchQuery := "%" + query + "%"

	// Preload Surah agar frontend tahu ini ayat dari surat apa
	err := r.db.WithContext(ctx).
		Preload("Surah").
		Where("translation ILIKE ? OR text_latin ILIKE ?", searchQuery, searchQuery).
		Limit(20).
		Find(&ayahs).Error

	if err != nil {
		return nil, err
	}
	return ayahs, nil
}