package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"khalif-alquran/internal/domain"

)

type QuranUC struct {
	surahRepo domain.SurahRepository
	ayahRepo  domain.AyahRepository
	redisRepo domain.RedisRepository
}

// NewQuranUseCase sekarang mengembalikan *QuranUC (Struct Pointer)
func NewQuranUseCase(surahRepo domain.SurahRepository, ayahRepo domain.AyahRepository, redisRepo domain.RedisRepository) *QuranUC {
	return &QuranUC{
		surahRepo: surahRepo,
		ayahRepo:  ayahRepo,
		redisRepo: redisRepo,
	}
}

func (uc *QuranUC) GetAllSurahs(ctx context.Context) ([]domain.Surah, error) {
	cacheKey := domain.CacheKeySurahAll

	if uc.redisRepo != nil {
		cachedData, err := uc.redisRepo.Get(ctx, cacheKey)
		if err == nil && cachedData != "" {
			var surahs []domain.Surah
			if err := json.Unmarshal([]byte(cachedData), &surahs); err == nil {
				return surahs, nil
			}
		}
	}

	surahs, err := uc.surahRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if uc.redisRepo != nil {
		if data, err := json.Marshal(surahs); err == nil {
			_ = uc.redisRepo.Set(ctx, cacheKey, data, 24*time.Hour)
		}
	}

	return surahs, nil
}

func (uc *QuranUC) GetSurahDetail(ctx context.Context, number int) (*domain.Surah, error) {
	cacheKey := fmt.Sprintf("%s%d", domain.CacheKeySurahPrefix, number)

	if uc.redisRepo != nil {
		cachedData, err := uc.redisRepo.Get(ctx, cacheKey)
		if err == nil && cachedData != "" {
			var surah domain.Surah
			if err := json.Unmarshal([]byte(cachedData), &surah); err == nil {
				return &surah, nil
			}
		}
	}

	surah, err := uc.surahRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	if uc.redisRepo != nil {
		if data, err := json.Marshal(surah); err == nil {
			_ = uc.redisRepo.Set(ctx, cacheKey, data, 24*time.Hour)
		}
	}

	return surah, nil
}

func (uc *QuranUC) GetAyahDetail(ctx context.Context, surahNumber, ayahNumber int) (*domain.Ayah, error) {
	return uc.ayahRepo.GetSpecificAyah(ctx, surahNumber, ayahNumber)
}

func (uc *QuranUC) Search(ctx context.Context, query string) (map[string]interface{}, error) {
	surahs, err := uc.surahRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	ayahs, err := uc.ayahRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"surahs": surahs,
		"ayahs":  ayahs,
	}

	return result, nil
}