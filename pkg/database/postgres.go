package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jooharver/geu-test/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah instance global untuk dipanggil dari repository
var DB *gorm.DB

// Connect melakukan inisialisasi koneksi ke PostgreSQL dan menjalankan migrasi
func Connect() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// config data source name (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	// Membuka koneksi GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database gagal terhubung: %v", err)
	}

	log.Println("Database berhasil terhubung")

	// Jalankan AutoMigrate untuk entitas yang ada
	err = db.AutoMigrate(&model.Household{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	log.Println("Migrasi tabel Household berhasil")

	DB = db
}