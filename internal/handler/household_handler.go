package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/service"
)

// Ini DTO (Data Transfer Object) buat nangkep JSON dari Postman
// tag binding:"required" itu versi otomatisnya $request->validate() di Laravel
type createHouseholdRequest struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Address   string `json:"address" binding:"required"`
}

type HouseholdHandler struct {
	service service.HouseholdService
}

// Constructor buat inject service-nya
func NewHouseholdHandler(s service.HouseholdService) *HouseholdHandler {
	return &HouseholdHandler{service: s}
}

// POST /api/households
func (h *HouseholdHandler) Create(c *gin.Context) {
	var req createHouseholdRequest

	// Cek apakah JSON yang dikirim user sesuai sama struct di atas
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format inputan salah atau ada data yang kosong"})
		return
	}

	// Lempar ke service buat dieksekusi
	household, err := h.service.CreateHousehold(req.OwnerName, req.Address)
	if err != nil {
		// Kalo service balikin error (misal nama kosong), tampilin errornya
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kalo sukses, balikin response 201 (Created)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Data warga berhasil ditambahkan",
		"data":    household,
	})
}

// GET /api/households
func (h *HouseholdHandler) GetAll(c *gin.Context) {
	households, err := h.service.GetAllHouseholds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ngambil data warga"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil ngambil list warga",
		"data":    households,
	})
}

// GET /api/households/:id
func (h *HouseholdHandler) GetByID(c *gin.Context) {
	// Ambil id dari parameter URL
	idParam := c.Param("id")
	
	// Ubah string id jadi format UUID
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID gak valid nih, harus UUID"})
		return
	}

	household, err := h.service.GetHouseholdByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data warga gak ketemu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data warga ditemukan",
		"data":    household,
	})
}

// DELETE /api/households/:id
func (h *HouseholdHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID gak valid"})
		return
	}

	// Panggil fitur delete yang baru aja kita bikin di service tadi
	err = h.service.DeleteHousehold(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data warga berhasil dihapus selamanya",
	})
}