package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HouseholdID uuid.UUID  `gorm:"type:uuid;not null" json:"household_id"`
	Amount      float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Status      string     `gorm:"type:varchar(50);default:'unpaid'" json:"status"` // unpaid, paid
	ReceiptURL  string     `gorm:"type:varchar(255)" json:"receipt_url"`
	PaymentDate *time.Time `json:"payment_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relasi kembali ke Household
	Household   *Household `gorm:"foreignKey:HouseholdID" json:"-"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}