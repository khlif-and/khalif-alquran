package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-alquran/internal/domain"
	"khalif-alquran/pkg/utils" // PENTING: Menggunakan utils, bukan response

)

type QuranHandler struct {
	quranUC domain.QuranUseCase
}

func NewQuranHandler(quranUC domain.QuranUseCase) *QuranHandler {
	return &QuranHandler{
		quranUC: quranUC,
	}
}

// GetAllSurahs godoc
// @Summary      Get All Surahs
// @Description  Get a list of all 114 Surahs (metadata only)
// @Tags         Quran
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /quran/surahs [get]
func (h *QuranHandler) GetAllSurahs(c *gin.Context) {
	surahs, err := h.quranUC.GetAllSurahs(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch surahs: "+err.Error())
		return
	}

	// Utils SuccessResponse hanya menerima (ctx, code, data)
	utils.SuccessResponse(c, http.StatusOK, surahs)
}

// GetSurahDetail godoc
// @Summary      Get Surah Detail
// @Description  Get specific Surah details including all Ayahs
// @Tags         Quran
// @Accept       json
// @Produce      json
// @Param        number   path      int  true  "Surah Number (1-114)"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /quran/surahs/{number} [get]
func (h *QuranHandler) GetSurahDetail(c *gin.Context) {
	numberStr := c.Param("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid surah number")
		return
	}

	surah, err := h.quranUC.GetSurahDetail(c.Request.Context(), number)
	if err != nil {
		if err == domain.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Surah not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch surah detail: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, surah)
}

// Search godoc
// @Summary      Search Quran
// @Description  Search for Surah names or Ayah texts/translations
// @Tags         Quran
// @Accept       json
// @Produce      json
// @Param        q    query     string  true  "Search Query"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /quran/search [get]
func (h *QuranHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Query param 'q' is required")
		return
	}

	result, err := h.quranUC.Search(c.Request.Context(), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Search failed: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}