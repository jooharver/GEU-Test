package repository

import(
	   "github.com/google/uuid"
	   "github.com/jooharver/geu-test/internal/model"
	   "gorm.io/gorm"
)

type WastePickupRepository interface{
	Create(pickup *model.WastePickup) error
	FindByHouseholdID(householdID uuid.UUID) ([]model.WastePickup, error)
}

type wastePickupRepository struct{
	db *gorm.DB
}

func NewWastePickupRepository(db *gorm.DB) WastePickupRepository{
	return &wastePickupRepository{db}
}

func (r *wastePickupRepository) Create(pickup *model.WastePickup) error {
	return r.db.Create(pickup).Error
}

func (r *wastePickupRepository) FindByHouseholdID(householdID uuid.UUID) ([]model.WastePickup, error) {
	var pickups []model.WastePickup
	// Mencari semua data sampah berdasarkan ID warga
	err := r.db.Where("household_id = ?", householdID).Find(&pickups).Error
	return pickups, err
}

