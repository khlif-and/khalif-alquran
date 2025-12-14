package domain

import (
	"context"
	"time"

)

// --- Entities ---

type Surah struct {
	ID             uint      `gorm:"primaryKey" json:"-"`
	Number         int       `gorm:"uniqueIndex" json:"number"`       // Nomor Surah (1-114)
	Name           string    `json:"name"`                            // Teks Arab (misal: الفاتحة)
	LatinName      string    `json:"latin_name"`                      // Teks Latin (misal: Al-Fatihah)
	EnglishName    string    `json:"english_name"`                    // Teks Inggris (misal: The Opening)
	RevelationType string    `json:"revelation_type"`                 // Mekkah/Madinah
	TotalAyahs     int       `json:"total_ayahs"`                     // Jumlah Ayat
	Ayahs          []Ayah    `gorm:"foreignKey:SurahID" json:"ayahs,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Ayah struct {
	ID           uint      `gorm:"primaryKey" json:"-"`
	SurahID      uint      `gorm:"index" json:"surah_id"`
	Surah        Surah     `gorm:"foreignKey:SurahID" json:"-"`
	Number       int       `json:"number"`                   // Nomor Ayat dalam Surah
	TextArabic   string    `gorm:"type:text" json:"text_arabic"` // Teks Utsmani
	TextLatin    string    `gorm:"type:text" json:"text_latin"`  // Transliterasi
	Translation  string    `gorm:"type:text" json:"translation"` // Terjemahan Indo
	Tafsir       string    `gorm:"type:text" json:"tafsir"`      // Tafsir Ringkas/Kemenag
	AsbabunNuzul string    `gorm:"type:text" json:"asbabun_nuzul"` // Nullable
	TajwidMeta   string    `gorm:"type:text" json:"tajwid_meta"`   // JSON String posisi tajwid
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Bookmark struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"index" json:"user_id"` // Dari JWT
	SurahID    uint      `gorm:"index" json:"surah_id"`
	Surah      Surah     `gorm:"foreignKey:SurahID" json:"surah,omitempty"`
	AyahNumber int       `json:"ayah_number"`
	Note       string    `json:"note"` // Opsional: Catatan user
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// --- Interfaces ---

// SurahRepository menangani akses data level Surah
type SurahRepository interface {
	GetAll(ctx context.Context) ([]Surah, error)
	GetByNumber(ctx context.Context, number int) (*Surah, error)
	Search(ctx context.Context, query string) ([]Surah, error)
}

// AyahRepository menangani akses data level Ayat
type AyahRepository interface {
	GetBySurahID(ctx context.Context, surahID uint) ([]Ayah, error)
	GetSpecificAyah(ctx context.Context, surahNumber, ayahNumber int) (*Ayah, error)
	Search(ctx context.Context, query string) ([]Ayah, error)
}

// RedisRepository untuk caching data Quran yang statis
type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

// BookmarkRepository menangani penyimpanan penanda ayat
type BookmarkRepository interface {
	SaveBookmark(ctx context.Context, bookmark *Bookmark) error
	GetByUserID(ctx context.Context, userID string) ([]Bookmark, error)
	DeleteBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error
	ClearAllBookmarks(ctx context.Context, userID string) error
}

// QuranUseCase menggabungkan logika bisnis untuk Surah dan Ayat
type QuranUseCase interface {
	GetAllSurahs(ctx context.Context) ([]Surah, error)
	GetSurahDetail(ctx context.Context, number int) (*Surah, error) // Include Ayahs
	GetAyahDetail(ctx context.Context, surahNumber, ayahNumber int) (*Ayah, error)
	Search(ctx context.Context, query string) (map[string]interface{}, error) // Bisa return Surah & Ayat
}

// BookmarkUseCase menangani logika bisnis bookmark
type BookmarkUseCase interface {
	AddBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int, note string) error
	GetUserBookmarks(ctx context.Context, userID string) ([]Bookmark, error)
	RemoveBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error
	ClearBookmarks(ctx context.Context, userID string) error
}