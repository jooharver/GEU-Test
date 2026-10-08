package repository

import (
	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByID(id uuid.UUID) (*model.Payment, error)
	Update(payment *model.Payment) error
	FindByHouseholdID(householdID uuid.UUID) ([]model.Payment, error)
	FindAll() ([]model.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, "id = ?", id).Error
	return &payment, err
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	// Save akan melakukan UPDATE data yang sudah ada berdasarkan Primary Key (ID)
	return r.db.Save(payment).Error
}

func (r *paymentRepository) FindByHouseholdID(householdID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Where("household_id = ?", householdID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindAll() ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Find(&payments).Error
	return payments, err
}