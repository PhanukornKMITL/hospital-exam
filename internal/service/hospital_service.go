package service

import (
	"errors"
	"strings"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/google/uuid"
)

type HospitalService interface {
	GetHospitals() ([]entity.Hospital, error)
	CreateHospital(input HospitalCreateInput) (*entity.Hospital, error)
	UpdateHospital(id uuid.UUID, input HospitalUpdateInput) (*entity.Hospital, error)
	DeleteHospital(id uuid.UUID) error
}

type hospitalService struct {
	repo repository.HospitalRepository
}

type HospitalCreateInput struct {
	Name    string
	Address string
}

type HospitalUpdateInput struct {
	Name    string
	Address string
}

func NewHospitalService(repo repository.HospitalRepository) HospitalService {
	return &hospitalService{repo: repo}
}

func (s *hospitalService) GetHospitals() ([]entity.Hospital, error) {
	return s.repo.FindAll()
}

func (s *hospitalService) CreateHospital(input HospitalCreateInput) (*entity.Hospital, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("hospital name is required")
	}

	h := &entity.Hospital{
		Name:    input.Name,
		Address: input.Address,
	}
	return s.repo.Create(h)
}

func (s *hospitalService) UpdateHospital(id uuid.UUID, input HospitalUpdateInput) (*entity.Hospital, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New("hospital name is required")
	}

	hospital, err := s.repo.FindByID(id)
	if err != nil || hospital == nil {
		return nil, errors.New("hospital not found")
	}

	hospital.Name = input.Name
	hospital.Address = input.Address

	return s.repo.Update(hospital)
}

func (s *hospitalService) DeleteHospital(id uuid.UUID) error {
	hospital, err := s.repo.FindByID(id)
	if err != nil || hospital == nil {
		return errors.New("hospital not found")
	}
	return s.repo.Delete(id)
}
