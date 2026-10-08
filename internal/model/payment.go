package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HouseholdID  uuid.UUID  `gorm:"type:uuid;not null" json:"household_id"`
	WasteID      uuid.UUID  `gorm:"type:uuid;not null" json:"waste_id"` 
	Amount       float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentDate  *time.Time `json:"payment_date"`
	Status       string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	ProofFileURL string     `gorm:"type:varchar(255)" json:"proof_file_url"` 
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Household *Household `gorm:"foreignKey:HouseholdID" json:"-"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}