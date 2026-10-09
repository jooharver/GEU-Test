package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jooharver/geu-test/internal/model"
	"github.com/jooharver/geu-test/internal/repository"
)

type HouseholdService interface {
	CreateHousehold(ownerName, email, phone, address string) (*model.Household, error)
	GetAllHouseholds() ([]model.Household, error)
	GetHouseholdByID(id uuid.UUID) (*model.Household, error)
	DeleteHousehold(id uuid.UUID) error
}

type householdService struct {
	repo repository.HouseholdRepository
}

// Constructor
func NewHouseholdService(repo repository.HouseholdRepository) HouseholdService {
	return &householdService{repo}
}

// Fungsi nambah warga baru
func (s *householdService) CreateHousehold(ownerName, email, phone, address string) (*model.Household, error) {
	if ownerName == "" {
		return nil, errors.New("nama pemilik wajib diisi")
	}
	if address == "" {
		return nil, errors.New("alamat wajib diisi")
	}

	// prepare struct modelnya
	household := &model.Household{
		OwnerName: ownerName,
		Email:     email,
		Phone:     phone,
		Address:   address,
	}

	// lempar ke repo buat di save ke db
	err := s.repo.Create(household)
	if err != nil {
		return nil, err
	}

	return household, nil
}

// Fungsi ambil semua data list warga
func (s *householdService) GetAllHouseholds() ([]model.Household, error) {
	return s.repo.FindAll()
}

// Fungsi search warga by id
func (s *householdService) GetHouseholdByID(id uuid.UUID) (*model.Household, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID warga tidak valid")
	}
	
	return s.repo.FindByID(id)
}

// Fungsi hapus warga by id
func (s *householdService) DeleteHousehold(id uuid.UUID) error {
	// pastikan orangnya ada di db
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("Tidak bisa dihapus, data warga tidak ditemukan")
	}

	// kalo ketemu, langsung eksekusi hapus
	return s.repo.Delete(id)
}