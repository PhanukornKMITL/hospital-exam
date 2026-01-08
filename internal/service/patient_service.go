package service

import (
	"strings"
	"time"

	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/google/uuid"
)

type PatientService interface {
	GetPatients() ([]entity.Patient, error)
	GetPatientsByHospital(hospitalID uuid.UUID) ([]entity.Patient, error)
	CreatePatient(input PatientCreateInput) (*entity.Patient, error)
}

type patientService struct {
	repo repository.PatientRepository
}

type PatientCreateInput struct {
	HospitalID   uuid.UUID
	FirstNameTH  string
	MiddleNameTH string
	LastNameTH   string

	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string

	DateOfBirth *time.Time

	NationalID  string
	PassportID  string
	PhoneNumber string
	Email       string
	Gender      string
}

func NewPatientService(repo repository.PatientRepository) PatientService {
	return &patientService{repo: repo}
}

func (s *patientService) GetPatients() ([]entity.Patient, error) {
	return s.repo.FindAll()
}

func (s *patientService) GetPatientsByHospital(hospitalID uuid.UUID) ([]entity.Patient, error) {
	return s.repo.FindByHospitalID(hospitalID)
}

func (s *patientService) CreatePatient(input PatientCreateInput) (*entity.Patient, error) {
	patient := &entity.Patient{
		HospitalID:   input.HospitalID,
		FirstNameTH:  input.FirstNameTH,
		MiddleNameTH: input.MiddleNameTH,
		LastNameTH:   input.LastNameTH,
		FirstNameEN:  input.FirstNameEN,
		MiddleNameEN: input.MiddleNameEN,
		LastNameEN:   input.LastNameEN,
		DateOfBirth:  input.DateOfBirth,
		NationalID:   normalizeOptionalPtr(input.NationalID),
		PassportID:   normalizeOptionalPtr(input.PassportID),
		PhoneNumber:  input.PhoneNumber,
		Email:        input.Email,
		Gender:       input.Gender,
	}
	return s.repo.CreateWithGeneratedHN(patient)
}

// normalizeOptionalPtr trims and returns nil if empty/whitespace, else pointer.
func normalizeOptionalPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := strings.TrimSpace(s)
	return &v
}
