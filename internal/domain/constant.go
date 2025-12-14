package domain

const (
	// Roles (Tetap dipertahankan untuk sekuritas endpoint admin/seeding)
	RoleAdmin = "Admin"
	RoleUser  = "User"

	// Cache Keys khusus Al-Quran
	CacheKeySurahAll    = "quran:surahs:all"   // Untuk list semua surah
	CacheKeySurahPrefix = "quran:surah:"       // Untuk detail per surah (misal: quran:surah:1)
)