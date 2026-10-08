package model

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

//membuat household struct
type Household struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    OwnerName  string    `gorm:"type:varchar(100);not null" json:"owner_name"`
    Address    string    `gorm:"type:text;not null" json:"address"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

	// Relasi
	Pickups  []WastePickup `gorm:"foreignKey:HouseholdID" json:"pickups,omitempty"`
	Payments []Payment     `gorm:"foreignKey:HouseholdID" json:"payments,omitempty"`
}

// tangani UUID jika kosong
func (h *Household) BeforeCreate(tx *gorm.DB) (err error) {
    if h.ID == uuid.Nil {
        h.ID = uuid.New()
    }
    return
}