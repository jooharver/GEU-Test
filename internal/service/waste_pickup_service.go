package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"github.com/jooharver/geu-test/internal/repository"
)

//interface
type WastePickupService interface {
	CreatePickup(householdID uuid.UUID, wasteType string, safetyCheck bool) (*model.WastePickup, error)
	SchedulePickup(id uuid.UUID, pickupDate time.Time) error
	CompletePickup(id uuid.UUID) error
	GetAllPickups() ([]model.WastePickup, error)
	GetPickupByID(id uuid.UUID) (*model.WastePickup, error)
	DeletePickup(id uuid.UUID) error
}

type wastePickupService struct {
	pickupRepo  repository.WastePickupRepository
	paymentRepo repository.PaymentRepository
}

// butuh 2 repo karena saat pickup selesai, harus otomatis bikin payment
func NewWastePickupService(pr repository.WastePickupRepository, payR repository.PaymentRepository) WastePickupService {
	return &wastePickupService{
		pickupRepo:  pr,
		paymentRepo: payR,
	}
}

func (s *wastePickupService) CreatePickup(householdID uuid.UUID, wasteType string, safetyCheck bool) (*model.WastePickup, error) {
	// cek aturan 1:gak boleh ada utang pending payment
	payments, err := s.paymentRepo.FindByHouseholdID(householdID)
	if err != nil {
		return nil, errors.New("gagal mengecek data pembayaran")
	}

	for _, p := range payments {
		if p.Status == "pending" {
			return nil, errors.New("tidak bisa request pickup, masih ada tagihan pending")
		}
	}

	pickup := &model.WastePickup{
		HouseholdID: householdID,
		Type:        wasteType,
		Status:      "pending",
		SafetyCheck: safetyCheck,
	}

	err = s.pickupRepo.Create(pickup)
	if err != nil {
		return nil, err
	}

	return pickup, nil
}

func (s *wastePickupService) SchedulePickup(id uuid.UUID, pickupDate time.Time) error {
	// ambil data pickup nya
	pickup, err := s.pickupRepo.FindByID(id)
	if err != nil {
		return errors.New("data pickup tidak ditemukan")
	}

	// cek aturan 2: harus dari status pending
	if pickup.Status != "pending" {
		return errors.New("pickup cuma bisa dijadwalkan kalau statusnya masih pending")
	}

	// cek aturan 3: kalo elektronik, safety check wajib true
	if pickup.Type == "electronic" && !pickup.SafetyCheck {
		return errors.New("sampah elektronik wajib lolos safety check dulu")
	}

	pickup.Status = "scheduled"
	pickup.PickupDate = &pickupDate

	return s.pickupRepo.Update(pickup)
}

func (s *wastePickupService) CompletePickup(id uuid.UUID) error {
	pickup, err := s.pickupRepo.FindByID(id)
	if err != nil {
		return errors.New("data pickup tidak ditemukan")
	}

	pickup.Status = "completed"
	
	// update status pickup nya
	err = s.pickupRepo.Update(pickup)
	if err != nil {
		return err
	}

	// cek aturan 4: hitung tarif buat digenerate jadi payment
	var amount float64
	if pickup.Type == "electronic" {
		amount = 100000
	} else {
		// organic, plastic, paper
		amount = 50000 
	}

	// generate payment
	newPayment := &model.Payment{
		HouseholdID: pickup.HouseholdID,
		WasteID:     pickup.ID,
		Amount:      amount,
		Status:      "pending",
	}

	return s.paymentRepo.Create(newPayment)
}

func (s *wastePickupService) GetAllPickups() ([]model.WastePickup, error) {
	return s.pickupRepo.FindAll()
}

func (s *wastePickupService) GetPickupByID(id uuid.UUID) (*model.WastePickup, error) {
	return s.pickupRepo.FindByID(id)
}

func (s *wastePickupService) DeletePickup(id uuid.UUID) error {
	// 1. Cari data pickup-nya dulu
	pickup, err := s.pickupRepo.FindByID(id)
	if err != nil {
		return errors.New("data pickup tidak ditemukan")
	}

	// 2. Jangan izinkan hapus kalau statusnya sudah completed
	if pickup.Status == "completed" {
		return errors.New("tidak dapat menghapus permintaan yang sudah diselesaikan dan menjadi tagihan")
	}

	// 3. Hapus jika aman
	return s.pickupRepo.Delete(id)
}