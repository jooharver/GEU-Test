package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jooharver/geu-test/internal/handler"
	"github.com/jooharver/geu-test/internal/repository"
	"github.com/jooharver/geu-test/internal/service"
	"github.com/jooharver/geu-test/pkg/database"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Info: file .env tidak ditemukan, menggunakan environment variable dari sistem.")
	}

	// konek database
	database.Connect()

	// buat seeder data
	database.Seed(database.DB)

	// setup layer repository untuk akses database
	householdRepo := repository.NewHouseholdRepository(database.DB)
	pickupRepo := repository.NewWastePickupRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)
	reportRepo := repository.NewReportRepository(database.DB)

	// setup layer service
	householdService := service.NewHouseholdService(householdRepo)
	pickupService := service.NewWastePickupService(pickupRepo, paymentRepo)
	paymentService := service.NewPaymentService(paymentRepo)
	reportService := service.NewReportService(reportRepo)

	// setup layer handler untuk request-response HTTP
	householdHandler := handler.NewHouseholdHandler(householdService)
	pickupHandler := handler.NewWastePickupHandler(pickupService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	reportHandler := handler.NewReportHandler(reportService)

	// inisialisasi router Gin
	r := gin.Default()

	// Buka akses publik folder uploads biar gambar bukti bayar bisa diakses dari browser
	r.Static("/uploads", "./uploads")

	// mendaftarkan semua routing API
	api := r.Group("/api")
	{
		// Warga/Household
		api.POST("/households", householdHandler.Create)
		api.GET("/households", householdHandler.GetAll)
		api.GET("/households/:id", householdHandler.GetByID)
		api.DELETE("/households/:id", householdHandler.Delete)

		// Sampah/Waste Pickup
		api.POST("/pickups", pickupHandler.Create)
		api.GET("/pickups", pickupHandler.GetAll)
		api.GET("/pickups/:id", pickupHandler.GetByID)      
		api.DELETE("/pickups/:id", pickupHandler.Delete)   
		api.PUT("/pickups/:id/schedule", pickupHandler.Schedule)
		api.PUT("/pickups/:id/complete", pickupHandler.Complete)

		// Pembayaran/Payment
		api.GET("/payments", paymentHandler.GetAll)
		api.PUT("/payments/:id/confirm", paymentHandler.Confirm)

		// Laporan/Reports
		api.GET("/reports/waste-summary", reportHandler.WasteSummary)
		api.GET("/reports/payment-summary", reportHandler.PaymentSummary)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server siap gas di port %s", port)
	r.Run(":" + port)
}