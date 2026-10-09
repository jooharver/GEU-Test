package database

import (
	"log"

	"github.com/jooharver/geu-test/internal/model"
	"gorm.io/gorm"
)

// Seed akan memasukkan data awal jika database masih kosong
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&model.Household{}).Count(&count)

	// Jika tabel households kosong, masukkan data dummy
	if count == 0 {
		households := []model.Household{
			{
				OwnerName: "Budi Santoso",
				Email:     "budi.santoso@example.com",
				Phone:     "081234567890",
				Address:   "Jl. Mawar No. 12, Kota Malang",
			},
			{
				OwnerName: "Siti Aminah",
				Email:     "siti.aminah@example.com",
				Phone:     "081298765432",
				Address:   "Jl. Melati No. 45, Kota Malang",
			},
			{
				OwnerName: "Andi Wijaya",
				Email:     "andi.w@example.com",
				Phone:     "081345678901",
				Address:   "Jl. Kamboja No. 8, Kota Malang",
			},
		}

		for _, h := range households {
			if err := db.Create(&h).Error; err != nil {
				log.Printf("Gagal insert data seed %s: %v", h.OwnerName, err)
			}
		}
		log.Println("Data dummy warga berhasil ditambahkan!")
	} else {
		log.Println("Data seed sudah ada, melewati proses seeding.")
	}
}