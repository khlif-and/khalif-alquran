package domain

import "errors"

var (
	// Error Standar
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("data (surah/ayah) not found")
	ErrConflict            = errors.New("data already exists")
	ErrBadParamInput       = errors.New("invalid parameter input")
	
	// Error Spesifik Domain Al-Quran (Opsional, agar lebih jelas saat debugging)
	ErrInvalidSurahNumber  = errors.New("surah number must be between 1 and 114")
	ErrInvalidAyahNumber   = errors.New("ayah number is out of range for this surah")
)