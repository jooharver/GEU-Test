package repository

import (
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"gorm.io/gorm"
)

// Menggunakan Interface adalah praktik standar Clean Architecture
type HouseholdRepository interface {
	Create(household *model.Household) error
	FindAll() ([]model.Household, error)
	FindByID(id uuid.UUID) (*model.Household, error)
	Delete(id uuid.UUID) error
}

type householdRepository struct {
	db *gorm.DB
}

func NewHouseholdRepository(db *gorm.DB) HouseholdRepository {
	return &householdRepository{db}
}

func (r *householdRepository) Create(household *model.Household) error {
	return r.db.Create(household).Error
}

func (r *householdRepository) FindAll() ([]model.Household, error) {
	var households []model.Household
	err := r.db.Find(&households).Error
	return households, err
}

func (r *householdRepository) FindByID(id uuid.UUID) (*model.Household, error) {
	var household model.Household
	// Preload digunakan agar saat mengambil data warga, riwayat sampah dan pembayarannya ikut ditarik
	err := r.db.Preload("Pickups").Preload("Payments").First(&household, "id = ?", id).Error
	return &household, err
}

func (r *householdRepository) Delete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&model.Household{}, "id = ?", id).Error
}