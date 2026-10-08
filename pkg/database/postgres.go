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

	// Konfigurasi DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	// Membuka koneksi menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL!")

	// Jalankan AutoMigrate untuk semua entitas
	err = db.AutoMigrate(&model.Household{}, &model.WastePickup{}, &model.Payment{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	log.Println("Migrasi seluruh tabel berhasil!")

	DB = db
}