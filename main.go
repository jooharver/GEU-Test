package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jooharver/geu-test/pkg/database"
)

func main() {
	// Muat file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan environment system")
	}

	// Inisialisasi koneksi database dan jalankan migrasi
	database.Connect()

	// Inisialisasi router Gin
	r := gin.Default()

	// Endpoint sederhana untuk tes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "API Pengumpulan Sampah Green Energi Utama Menyala!",
		})
	})

	// Ambil PORT dari .env, default ke 8080 jika tidak ada
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}