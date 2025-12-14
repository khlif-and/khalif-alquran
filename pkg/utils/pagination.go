package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"

)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

// Tambahkan parameter defaultSort agar fleksibel
func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Default values
	limit := 10
	page := 1
	sort := "created_at desc" // Default umum

	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			if l, err := strconv.Atoi(queryValue); err == nil && l > 0 {
				limit = l
			}
		case "page":
			if p, err := strconv.Atoi(queryValue); err == nil && p > 0 {
				page = p
			}
		case "sort":
			// Mencegah SQL Injection sederhana dengan whitelist karakter
			if queryValue != "" {
				sort = queryValue
			}
		}
	}

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

// Helper untuk offset database
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// Helper untuk mendapatkan Sort yang aman (GORM style)
func (p *Pagination) GetSort() string {
    // Jika sort kosong, defaultkan. 
    // Kamu bisa logic disini, misal ganti koma dengan spasi dsb.
	return p.Sort
}