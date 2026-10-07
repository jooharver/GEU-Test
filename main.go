package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load file .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan")
	}

	// Inisialisasi router Gin
	r := gin.Default()

	// Endpoint untuk tes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "API sukses berjalan",
		})
	})

	// ambil port dari .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}