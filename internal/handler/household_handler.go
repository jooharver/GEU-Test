package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/service"
)

// Menyimpan struktur payload untuk pendaftaran warga baru
type createHouseholdRequest struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type HouseholdHandler struct {
	service service.HouseholdService
}

// Inisialisasi handler untuk entitas household(warga)
func NewHouseholdHandler(s service.HouseholdService) *HouseholdHandler {
	return &HouseholdHandler{service: s}
}

// Menangani pembuatan data warga
func (h *HouseholdHandler) Create(c *gin.Context) {
	var payload createHouseholdRequest

	// Pengecekan payload inputan dari client
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format inputan salah atau ada data wajib yang masih kosong"})
		return
	}

	// Transfer pendaftaran ke layer service
	household, srvErr := h.service.CreateHousehold(payload.OwnerName, payload.Email, payload.Phone, payload.Address)
	if srvErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": srvErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi warga berhasil.",
		"data":    household,
	})
}

// Mengambil keseluruhan data warga yang terdaftar
func (h *HouseholdHandler) GetAll(c *gin.Context) {
	households, fetchErr := h.service.GetAllHouseholds()
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memuat data warga dari sistem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data warga berhasil dimuat.",
		"data":    households,
	})
}

// Menarik data satu warga spesifik berdasarkan ID
func (h *HouseholdHandler) GetByID(c *gin.Context) {
	rawID := c.Param("id")
	
	householdID, parseErr := uuid.Parse(rawID)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid. Harap periksa kembali."})
		return
	}

	household, findErr := h.service.GetHouseholdByID(householdID)
	if findErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data warga tidak ditemukan di sistem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Arsip warga ditemukan.",
		"data":    household,
	})
}

// Menghapus data warga berdasarkan ID
func (h *HouseholdHandler) Delete(c *gin.Context) {
	rawID := c.Param("id")
	
	householdID, parseErr := uuid.Parse(rawID)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format UUID tidak sesuai standar."})
		return
	}

	delErr := h.service.DeleteHousehold(householdID)
	if delErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": delErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data warga sudah dihapus dari sistem.",
	})
}