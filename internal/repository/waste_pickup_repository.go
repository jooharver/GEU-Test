package repository

import(
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"gorm.io/gorm"
)

type WastePickupRepository interface {
	Create(pickup *model.WastePickup) error
	FindByHouseholdID(householdID uuid.UUID) ([]model.WastePickup, error)
	FindByID(id uuid.UUID) (*model.WastePickup, error)
	Update(pickup *model.WastePickup) error
	FindAll() ([]model.WastePickup, error)
	Delete(id uuid.UUID) error
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
	// Mencari semua data sampah berdasarkan id household
	err := r.db.Where("household_id = ?", householdID).Find(&pickups).Error
	return pickups, err
}

func (r *wastePickupRepository) FindByID(id uuid.UUID) (*model.WastePickup, error) {
	var pickup model.WastePickup
	err := r.db.First(&pickup, "id = ?", id).Error
	return &pickup, err
}

func (r *wastePickupRepository) Update(pickup *model.WastePickup) error {
	return r.db.Save(pickup).Error
}

func (r *wastePickupRepository) FindAll() ([]model.WastePickup, error) {
	var pickups []model.WastePickup
	// Mengambil semua data pickup sampah dari database
	err := r.db.Find(&pickups).Error
	return pickups, err
}

func (r *wastePickupRepository) Delete(id uuid.UUID) error {
	// Menghapus data berdasarkan ID
	return r.db.Delete(&model.WastePickup{}, "id = ?", id).Error
}