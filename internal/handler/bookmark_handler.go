package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-alquran/internal/domain"
	"khalif-alquran/pkg/utils"

)

type BookmarkHandler struct {
	bookmarkUC domain.BookmarkUseCase
}

func NewBookmarkHandler(bookmarkUC domain.BookmarkUseCase) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkUC: bookmarkUC,
	}
}

// GetUserBookmarks godoc
// @Summary      Get User Bookmarks
// @Description  Get all bookmarks for a specific user
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true  "User ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /bookmarks/{user_id} [get]
func (h *BookmarkHandler) GetUserBookmarks(c *gin.Context) {
	userID := c.Param("user_id")

	// PERBAIKAN: Hapus pagination, sesuai interface UseCase terbaru
	bookmarks, err := h.bookmarkUC.GetUserBookmarks(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch bookmarks: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, bookmarks)
}

// AddBookmark godoc
// @Summary      Add Bookmark
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        request body domain.Bookmark true "Bookmark Data"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /bookmarks [post]
func (h *BookmarkHandler) AddBookmark(c *gin.Context) {
	var req domain.Bookmark
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// PERBAIKAN: Gunakan AddBookmark dan pecah parameter sesuai UseCase
	if err := h.bookmarkUC.AddBookmark(c.Request.Context(), req.UserID, req.SurahID, req.AyahNumber, req.Note); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create bookmark: "+err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusCreated, "Bookmark created successfully")
}

// RemoveBookmark godoc
// @Summary      Remove Bookmark
// @Tags         Bookmarks
// @Param        user_id    query     string  true  "User ID"
// @Param        surah_id   query     int     true  "Surah ID"
// @Param        ayah_number query    int     true  "Ayah Number"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /bookmarks [delete]
func (h *BookmarkHandler) RemoveBookmark(c *gin.Context) {
	userID := c.Query("user_id")
	surahIDStr := c.Query("surah_id")
	ayahNumberStr := c.Query("ayah_number")

	if userID == "" || surahIDStr == "" || ayahNumberStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "user_id, surah_id, and ayah_number are required")
		return
	}

	surahID, err := strconv.Atoi(surahIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid surah_id format")
		return
	}

	ayahNumber, err := strconv.Atoi(ayahNumberStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ayah_number format")
		return
	}

	// PERBAIKAN: Gunakan RemoveBookmark (ganti nama dari Delete)
	if err := h.bookmarkUC.RemoveBookmark(c.Request.Context(), userID, uint(surahID), ayahNumber); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete bookmark: "+err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusOK, "Bookmark deleted successfully")
}