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

// Inisialisasi payment handler baru
func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// endpoint untuk memvalidasi pembayaran dengan upload struk
func (h *PaymentHandler) Confirm(c *gin.Context) {
	rawParam := c.Param("id")
	
	paymentID, parseErr := uuid.Parse(rawParam)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memuat tagihan. Pastikan formatnya menggunakan UUID."})
		return
	}

	// Ekstrak lampiran gambar bukti transfer
	uploadedFile, fileErr := c.FormFile("receipt")
	if fileErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lampiran 'receipt' tidak ditemukan. Mohon sertakan bukti transfer."})
		return
	}

	// Konstruksi nama unik menggunakan timestamp untuk mencegah redundan nama file
	timestampSuffix := time.Now().Unix()
	safeFileName := fmt.Sprintf("%d_%s", timestampSuffix, filepath.Base(uploadedFile.Filename))
	destinationPath := filepath.Join("uploads", safeFileName)

	// Tulis file gambar ke storage server lokal
	if saveErr := c.SaveUploadedFile(uploadedFile, destinationPath); saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terdapat kendala saat menyimpan berkas bukti bayar ke direktori."})
		return
	}

	// Update record pembayaran menggunakan service layer
	publicFileURL := "/" + destinationPath
	srvErr := h.service.ConfirmPayment(paymentID, publicFileURL)
	if srvErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": srvErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Verifikasi pembayaran berhasil diselesaikan.",
		"receipt_url": publicFileURL,
	})
}

// Endpoint rekapitulasi semua riwayat tagihan
func (h *PaymentHandler) GetAll(c *gin.Context) {
	paymentList, fetchErr := h.service.GetAllPayments()
	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sistem gagal memuat daftar transaksi pembayaran."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daftar transaksi tagihan berhasil ditarik.",
		"data":    paymentList,
	})
}