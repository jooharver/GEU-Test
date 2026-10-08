package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WastePickup struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HouseholdID uuid.UUID  `gorm:"type:uuid;not null" json:"household_id"`
	Weight      float64    `gorm:"type:decimal(10,2);not null" json:"weight"`
	WasteType   string     `gorm:"type:varchar(50);not null" json:"waste_type"` 
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"` 
	PickupDate  *time.Time `json:"pickup_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi kembali ke Household
	Household   *Household `gorm:"foreignKey:HouseholdID" json:"-"`
}

func (w *WastePickup) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}