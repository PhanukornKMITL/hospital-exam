package service

import (
	"errors"
	"strings"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
)

type HospitalService interface {
	GetHospitals() ([]entity.Hospital, error)
	CreateHospital(input HospitalCreateInput) (*entity.Hospital, error)
}

type hospitalService struct {
	repo repository.HospitalRepository
}

type HospitalCreateInput struct {
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
