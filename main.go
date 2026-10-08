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
	// 1. Muat config environment
	if err := godotenv.Load(); err != nil {
		log.Println("Info: file .env gak ada, lanjut pake variabel environment dari system")
	}

	// 2. Konek ke database & jalanin AutoMigrate
	database.Connect()

	// 3. Setup Layer Repository (Tangan DB)
	householdRepo := repository.NewHouseholdRepository(database.DB)
	pickupRepo := repository.NewWastePickupRepository(database.DB)
	paymentRepo := repository.NewPaymentRepository(database.DB)

	// 4. Setup Layer Service (Otak Logika)
	householdService := service.NewHouseholdService(householdRepo)
	// perhatiin: pickup service butuh payment repo juga buat generate tagihan otomatis pas sampah selesai diambil
	pickupService := service.NewWastePickupService(pickupRepo, paymentRepo) 
	paymentService := service.NewPaymentService(paymentRepo)

	// 5. Setup Layer Handler (Penerima Request HTTP)
	householdHandler := handler.NewHouseholdHandler(householdService)
	pickupHandler := handler.NewWastePickupHandler(pickupService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// 6. Inisialisasi Router Gin
	r := gin.Default()

	// Buka akses publik ke folder uploads biar gambar bukti bayar bisa diakses dari browser
	r.Static("/uploads", "./uploads")

	// 7. Daftarin semua routing API
	api := r.Group("/api")
	{
		// Warga / Household
		api.POST("/households", householdHandler.Create)
		api.GET("/households", householdHandler.GetAll)
		api.GET("/households/:id", householdHandler.GetByID)
		api.DELETE("/households/:id", householdHandler.Delete)

		// Sampah / Waste Pickup
		api.POST("/pickups", pickupHandler.Create)
		api.PUT("/pickups/:id/schedule", pickupHandler.Schedule)
		api.PUT("/pickups/:id/complete", pickupHandler.Complete)

		// Pembayaran / Payment
		api.PUT("/payments/:id/confirm", paymentHandler.Confirm)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server siap gas di port %s", port)
	r.Run(":" + port)
}