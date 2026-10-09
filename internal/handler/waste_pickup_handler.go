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

// Inisialisasi handler untuk pickup
func NewWastePickupHandler(s service.WastePickupService) *WastePickupHandler {
	return &WastePickupHandler{service: s}
}

// Struktur data payload untuk pengajuan pickup
type pickupSubmissionPayload struct {
	HouseholdID uuid.UUID `json:"household_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	SafetyCheck bool      `json:"safety_check"`
}

// Memproses pembuatan tiket pickup baru dan validasi safety check
func (h *WastePickupHandler) Create(c *gin.Context) {
	var payload pickupSubmissionPayload
	
	// Validasi kelengkapan body JSON
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memproses pengajuan: Pastikan filed terisi dengan benar."})
		return
	}

	// Teruskan ke service untuk pengecekan tunggakan dan safety check
	newPickup, srvErr := h.service.CreatePickup(payload.HouseholdID, payload.Type, payload.SafetyCheck)
	if srvErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": srvErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permintaan pickup berhasil dibuat.",
		"data":    newPickup,
	})
}

// Struktur payload khusus untuk penentuan jadwal
type schedulingPayload struct {
	PickupDate string `json:"pickup_date" binding:"required"`
}

// Menetapkan atau mengubah jadwal operasional pickup
func (h *WastePickupHandler) Schedule(c *gin.Context) {
	rawParam := c.Param("id")
	ticketID, parseErr := uuid.Parse(rawParam)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identitas tiket tidak valid. Harap gunakan format UUID."})
		return
	}

	var payload schedulingPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter jadwal (pickup_date) tidak ditemukan dalam request."})
		return
	}

	// Mekanisme fallback parsing tanggal agar fleksibel menerima berbagai format input
	var finalDate time.Time
	var timeErr error

	finalDate, timeErr = time.Parse("2006-01-02 15:04:05", payload.PickupDate)
	if timeErr != nil {
		finalDate, timeErr = time.Parse("2006-01-02 15:04", payload.PickupDate)
		if timeErr != nil {
			finalDate, timeErr = time.Parse("2006-01-02", payload.PickupDate)
			if timeErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Struktur tanggal tidak dikenali. Opsi format: 'YYYY-MM-DD', 'YYYY-MM-DD HH:MM', atau 'YYYY-MM-DD HH:MM:SS'",
				})
				return
			}
		}
	}

	updateErr := h.service.SchedulePickup(ticketID, finalDate)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Penjadwalan berhasil dikonfirmasi."})
}

// Menandai proses operasional selesai dan memicu pembuatan tagihan
func (h *WastePickupHandler) Complete(c *gin.Context) {
	rawParam := c.Param("id")
	ticketID, parseErr := uuid.Parse(rawParam)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format identitas tiket salah."})
		return
	}

	srvErr := h.service.CompletePickup(ticketID)
	if srvErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": srvErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operasional selesai. Invoice tagihan otomatis telah digenerate oleh sistem."})
}

// Menarik daftar semua aktivitas pengambilan
func (h *WastePickupHandler) GetAll(c *gin.Context) {
	pickupRecords, fetchErr := h.service.GetAllPickups()
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terdapat gangguan saat menarik riwayat operasional."})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar riwayat pickup berhasil dimuat.",
		"data": pickupRecords,
	})
}

// Mengambil detail satu tiket berdasarkan ID
func (h *WastePickupHandler) GetByID(c *gin.Context) {
	rawParam := c.Param("id")
	ticketID, parseErr := uuid.Parse(rawParam)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kriteria pencarian ID tidak valid."})
		return
	}

	record, findErr := h.service.GetPickupByID(ticketID)
	if findErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokumen tiket tidak ditemukan dalam arsip."})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail tiket ditemukan.",
		"data": record,
	})
}

// Membatalkan atau menghapus tiket yang masih berstatus pending
func (h *WastePickupHandler) Delete(c *gin.Context) {
	rawParam := c.Param("id")
	ticketID, parseErr := uuid.Parse(rawParam)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kriteria identitas UUID tidak valid."})
		return
	}

	delErr := h.service.DeletePickup(ticketID)
	if delErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Tiket pengambilan sukses dibatalkan dan dihapus dari sistem."})
}