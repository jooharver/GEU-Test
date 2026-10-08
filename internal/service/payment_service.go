package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"github.com/jooharver/geu-test/internal/repository"
)

type PaymentService interface {
	ConfirmPayment(id uuid.UUID, proofFileURL string) error
	GetPaymentsByHousehold(householdID uuid.UUID) ([]model.Payment, error)
	GetAllPayments() ([]model.Payment, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (s *paymentService) ConfirmPayment(id uuid.UUID, proofFileURL string) error {
	// ambil data tagihan
	payment, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data tagihan tidak ditemukan")
	}

	// validasi tagihan (lunas/belum)
	if payment.Status == "paid" {
		return errors.New("tagihan ini sudah lunas")
	}

	// cek aturan #5: wajib upload bukti bayar
	if proofFileURL == "" {
		return errors.New("wajib upload bukti bayar")
	}

	// update status dan simpan path file
	payment.Status = "paid"
	payment.ProofFileURL = proofFileURL

	// simpan perubahan ke database
	return s.repo.Update(payment)
}

func (s *paymentService) GetPaymentsByHousehold(householdID uuid.UUID) ([]model.Payment, error) {
	return s.repo.FindByHouseholdID(householdID)
}

func (s *paymentService) GetAllPayments() ([]model.Payment, error) {
	return s.repo.FindAll()
}