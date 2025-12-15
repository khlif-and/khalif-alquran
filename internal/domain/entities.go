package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

)

// --- Helper Structs untuk Tajwid (JSONB Support) ---

// TajwidRule merepresentasikan satu aturan tajwid dalam array
type TajwidRule struct {
	Rule    string `json:"rule"`
	Segment string `json:"segment"`
}

// TajwidList adalah slice custom agar bisa disimpan sebagai JSONB di Postgres
type TajwidList []TajwidRule

// Value: Mengubah struct Go menjadi JSON string untuk disimpan ke Database (GORM interface)
func (t TajwidList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan: Mengubah data JSON dari Database menjadi struct Go (GORM interface)
func (t *TajwidList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &t)
}

// --- Entities ---

type Surah struct {
	ID             uint      `gorm:"primaryKey" json:"-"`
	Number         int       `gorm:"uniqueIndex" json:"number"`
	Name           string    `json:"name"`
	LatinName      string    `json:"latin_name"`
	EnglishName    string    `json:"english_name"`
	IndonesianName string    `json:"indonesian_name"` // FIELD BARU
	RevelationType string    `json:"revelation_type"`
	
	// Tag json disesuaikan dengan key di file seed ("ayah_count")
	TotalAyahs     int       `json:"ayah_count"` 
	
	Ayahs          []Ayah    `gorm:"foreignKey:SurahID" json:"ayahs,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Ayah struct {
	ID           uint       `gorm:"primaryKey" json:"-"`
	SurahID      uint       `gorm:"index" json:"surah_id"`
	Surah        Surah      `gorm:"foreignKey:SurahID" json:"-"`
	Number       int        `json:"number"`
	TextArabic   string     `gorm:"type:text" json:"text_arabic"`
	TextLatin    string     `gorm:"type:text" json:"text_latin"`
	
	// Tag json disesuaikan dengan key di file seed ("translation_id")
	Translation  string     `gorm:"type:text" json:"translation_id"`
	
	Tafsir       string     `gorm:"type:text" json:"tafsir"`
	AsbabunNuzul string     `gorm:"type:text" json:"asbabun_nuzul"`
	
	// Menggunakan tipe custom TajwidList dan tag "tajwid_info"
	// Tipe gorm:jsonb agar tersimpan efisien di Postgres
	TajwidInfo   TajwidList `gorm:"type:jsonb" json:"tajwid_info"` 
	
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Bookmark struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"index" json:"user_id"`
	SurahID    uint      `gorm:"index" json:"surah_id"`
	Surah      Surah     `gorm:"foreignKey:SurahID" json:"surah,omitempty"`
	AyahNumber int       `json:"ayah_number"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// --- Interfaces ---

type SurahRepository interface {
	GetAll(ctx context.Context) ([]Surah, error)
	GetByNumber(ctx context.Context, number int) (*Surah, error)
	Search(ctx context.Context, query string) ([]Surah, error)
}

type AyahRepository interface {
	GetBySurahID(ctx context.Context, surahID uint) ([]Ayah, error)
	GetSpecificAyah(ctx context.Context, surahNumber, ayahNumber int) (*Ayah, error)
	Search(ctx context.Context, query string) ([]Ayah, error)
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

type BookmarkRepository interface {
	SaveBookmark(ctx context.Context, bookmark *Bookmark) error
	GetByUserID(ctx context.Context, userID string) ([]Bookmark, error)
	DeleteBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error
	ClearAllBookmarks(ctx context.Context, userID string) error
}

type QuranUseCase interface {
	GetAllSurahs(ctx context.Context) ([]Surah, error)
	GetSurahDetail(ctx context.Context, number int) (*Surah, error)
	GetAyahDetail(ctx context.Context, surahNumber, ayahNumber int) (*Ayah, error)
	Search(ctx context.Context, query string) (map[string]interface{}, error)
	ClearCache(ctx context.Context) error
}

type BookmarkUseCase interface {
	AddBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int, note string) error
	GetUserBookmarks(ctx context.Context, userID string) ([]Bookmark, error)
	RemoveBookmark(ctx context.Context, userID string, surahID uint, ayahNumber int) error
	ClearBookmarks(ctx context.Context, userID string) error
}