package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WastePickup struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HouseholdID uuid.UUID  `gorm:"type:uuid;not null" json:"household_id"`
	Type        string     `gorm:"type:varchar(50);not null" json:"type"` 
	Status      string     `gorm:"type:varchar(50);default:'pending'" json:"status"` 
	PickupDate  *time.Time `json:"pickup_date"`
	SafetyCheck bool       `gorm:"default:false" json:"safety_check"` 
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Household *Household `gorm:"foreignKey:HouseholdID" json:"-"`
}

func (w *WastePickup) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}