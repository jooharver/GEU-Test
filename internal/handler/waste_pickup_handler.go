package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/service"
)

type WastePickupHandler struct {
	service service.WastePickupService
}

func NewWastePickupHandler(s service.WastePickupService) *WastePickupHandler {
	return &WastePickupHandler{service: s}
}

// struct buat nangkep JSON dari body request pas bikin jadwal baru
type createPickupRequest struct {
	HouseholdID uuid.UUID `json:"household_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	SafetyCheck bool      `json:"safety_check"`
}

// POST /api/pickups
func (h *WastePickupHandler) Create(c *gin.Context) {
	var req createPickupRequest
	// bind JSON ke struct, kalo ada yang kurang otomatis error
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format inputan salah atau ada data wajib yang kosong"})
		return
	}

	pickup, err := h.service.CreatePickup(req.HouseholdID, req.Type, req.SafetyCheck)
	if err != nil {
		// balikin error dari aturan bisnis (misal: masih ada utang)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permintaan pickup berhasil dibuat",
		"data":    pickup,
	})
}

// struct buat nangkep JSON tanggal jadwal
type schedulePickupRequest struct {
	PickupDate time.Time `json:"pickup_date" binding:"required"`
}

// PUT /api/pickups/:id/schedule
func (h *WastePickupHandler) Schedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID pickup gak valid"})
		return
	}

	var req schedulePickupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah. Pake format RFC3339 ya (contoh: 2026-10-15T10:00:00Z)"})
		return
	}

	err = h.service.SchedulePickup(id, req.PickupDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pickup berhasil diubah jadi scheduled"})
}

// PUT /api/pickups/:id/complete
func (h *WastePickupHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID pickup gak valid"})
		return
	}

	err = h.service.CompletePickup(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pickup kelar! Tagihan pembayaran otomatis udah digenerate ke database."})
}