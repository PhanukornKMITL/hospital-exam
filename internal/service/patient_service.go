package service

import (
	"errors"
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
		// Log hospitalId from token for debugging
		println("[DEBUG] CreatePatient - HospitalID from token:", input.HospitalID.String())

	// Check duplicate national ID in the same hospital
	if input.NationalID != "" && strings.TrimSpace(input.NationalID) != "" {
		exists, err := s.repo.ExistsByNationalIDInHospital(input.HospitalID, strings.TrimSpace(input.NationalID))
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("nationalId already exists in this hospital")
		}
	}

	// Check duplicate passport ID in the same hospital
	if input.PassportID != "" && strings.TrimSpace(input.PassportID) != "" {
		exists, err := s.repo.ExistsByPassportIDInHospital(input.HospitalID, strings.TrimSpace(input.PassportID))
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("passportId already exists in this hospital")
		}
	}

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
