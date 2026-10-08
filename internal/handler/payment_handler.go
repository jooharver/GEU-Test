package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/service"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// PUT /api/payments/:id/confirm
func (h *PaymentHandler) Confirm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tagihan gak valid"})
		return
	}

	// 1. Nangkep file gambar dari form-data Postman (pake key "receipt")
	file, err := c.FormFile("receipt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File bukti bayar wajib diupload (pake key form-data 'receipt')"})
		return
	}

	// 2. Bikin nama file unik pake kombinasi angka timestamp waktu sekarang
	// Biar kalo ada warga ngupload file namanya "bukti.jpg" secara bersamaan, filenya gak saling nimpa
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	uploadPath := filepath.Join("uploads", fileName)

	// 3. Simpen file fisiknya beneran ke folder lokal /uploads/ (dibantu otomatis sama Gin)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Aduh, gagal nyimpen file ke server lokal"})
		return
	}

	// 4. Update status di database lewat Service
	fileURL := "/" + uploadPath
	err = h.service.ConfirmPayment(id, fileURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Pembayaran sukses dikonfirmasi!",
		"receipt_url": fileURL,
	})
}