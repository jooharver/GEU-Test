package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jooharver/geu-test/internal/service"
)

type ReportHandler struct {
	service service.ReportService
}

// Inisialisasi handler untuk rute pelaporan dashboard
func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{service: s}
}

// Endpoint untuk menarik statistik sampah berdasarkan tipe
func (h *ReportHandler) WasteSummary(c *gin.Context) {
	summaryData, fetchErr := h.service.GetWasteSummary()
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kendala internal saat menyusun laporan volume sampah."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data rekapitulasi volume sampah berhasil dimuat.",
		"data":    summaryData,
	})
}

// Endpoint untuk mengkalkulasi total pendapatan dari tagihan yang sudah lunas
func (h *ReportHandler) PaymentSummary(c *gin.Context) {
	revenueData, fetchErr := h.service.GetPaymentSummary()
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sistem gagal mengkalkulasi rekapitulasi pendapatan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data rekapitulasi pendapatan berhasil ditarik.",
		"data":    revenueData,
	})
}