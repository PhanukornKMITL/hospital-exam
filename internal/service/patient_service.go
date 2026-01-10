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
	SearchPatientByID(hospitalID uuid.UUID, identifier string) (*entity.Patient, error)
	SearchPatients(hospitalID uuid.UUID, input PatientSearchInput, page, limit int) ([]entity.Patient, int64, error)
	UpdatePatient(hospitalID, patientID uuid.UUID, input PatientUpdateInput) (*entity.Patient, error)
	DeletePatient(hospitalID, patientID uuid.UUID) error
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

type PatientUpdateInput struct {
	FirstNameTH  *string
	MiddleNameTH *string
	LastNameTH   *string

	FirstNameEN  *string
	MiddleNameEN *string
	LastNameEN   *string

	DateOfBirth *time.Time

	NationalID  *string
	PassportID  *string
	PhoneNumber *string
	Email       *string
	Gender      *string
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

	// Validate gender - must be 'M' or 'F'
	if input.Gender != "M" && input.Gender != "F" {
		return nil, errors.New("gender must be either 'M' or 'F'")
	}

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

func (s *patientService) SearchPatientByID(hospitalID uuid.UUID, identifier string) (*entity.Patient, error) {
	id := strings.TrimSpace(identifier)
	return s.repo.FindByHospitalAndIdentifier(hospitalID, id)
}

type PatientSearchInput struct {
	PatientHN    string
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

func (s *patientService) SearchPatients(hospitalID uuid.UUID, input PatientSearchInput, page, limit int) ([]entity.Patient, int64, error) {
	toPtr := func(s string) *string {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		v := strings.TrimSpace(s)
		return &v
	}

	filters := repository.PatientSearchFilters{
		PatientHN:    toPtr(input.PatientHN),
		FirstNameTH:  toPtr(input.FirstNameTH),
		MiddleNameTH: toPtr(input.MiddleNameTH),
		LastNameTH:   toPtr(input.LastNameTH),
		FirstNameEN:  toPtr(input.FirstNameEN),
		MiddleNameEN: toPtr(input.MiddleNameEN),
		LastNameEN:   toPtr(input.LastNameEN),
		DateOfBirth:  input.DateOfBirth,
		NationalID:   toPtr(input.NationalID),
		PassportID:   toPtr(input.PassportID),
		PhoneNumber:  toPtr(input.PhoneNumber),
		Email:        toPtr(input.Email),
		Gender:       toPtr(input.Gender),
	}

	// Defaults for pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	patients, total, err := s.repo.FindByHospitalWithFilters(hospitalID, filters, page, limit)
	return patients, total, err
}

func (s *patientService) UpdatePatient(hospitalID, patientID uuid.UUID, input PatientUpdateInput) (*entity.Patient, error) {
	patient, err := s.repo.FindByHospitalAndID(hospitalID, patientID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("patient not found")
	}

	// Gender validation only when provided
	if input.Gender != nil {
		if *input.Gender != "M" && *input.Gender != "F" {
			return nil, errors.New("gender must be either 'M' or 'F'")
		}
	}

	// Apply field updates only if provided
	if input.FirstNameTH != nil {
		patient.FirstNameTH = *input.FirstNameTH
	}
	if input.MiddleNameTH != nil {
		patient.MiddleNameTH = *input.MiddleNameTH
	}
	if input.LastNameTH != nil {
		patient.LastNameTH = *input.LastNameTH
	}
	if input.FirstNameEN != nil {
		patient.FirstNameEN = *input.FirstNameEN
	}
	if input.MiddleNameEN != nil {
		patient.MiddleNameEN = *input.MiddleNameEN
	}
	if input.LastNameEN != nil {
		patient.LastNameEN = *input.LastNameEN
	}
	if input.DateOfBirth != nil {
		patient.DateOfBirth = input.DateOfBirth
	}
	if input.PhoneNumber != nil {
		patient.PhoneNumber = *input.PhoneNumber
	}
	if input.Email != nil {
		patient.Email = *input.Email
	}
	if input.Gender != nil {
		patient.Gender = *input.Gender
	}

	// Handle identifiers with duplicate checks only when provided
	if input.NationalID != nil {
		nationalPtr := normalizeOptionalPtr(*input.NationalID)
		if patient.NationalID == nil || (nationalPtr == nil || *patient.NationalID != *nationalPtr) {
			if nationalPtr != nil {
				exists, err := s.repo.ExistsByNationalIDInHospitalExcept(hospitalID, *nationalPtr, patientID)
				if err != nil {
					return nil, err
				}
				if exists {
					return nil, errors.New("nationalId already exists in this hospital")
				}
			}
		}
		patient.NationalID = nationalPtr
	}

	if input.PassportID != nil {
		passportPtr := normalizeOptionalPtr(*input.PassportID)
		if patient.PassportID == nil || (passportPtr == nil || *patient.PassportID != *passportPtr) {
			if passportPtr != nil {
				exists, err := s.repo.ExistsByPassportIDInHospitalExcept(hospitalID, *passportPtr, patientID)
				if err != nil {
					return nil, err
				}
				if exists {
					return nil, errors.New("passportId already exists in this hospital")
				}
			}
		}
		patient.PassportID = passportPtr
	}

	// Ensure at least one identifier remains
	if patient.NationalID == nil && patient.PassportID == nil {
		return nil, errors.New("either nationalId or passportId must be provided")
	}

	return s.repo.Update(patient)
}

func (s *patientService) DeletePatient(hospitalID, patientID uuid.UUID) error {
	patient, err := s.repo.FindByHospitalAndID(hospitalID, patientID)
	if err != nil {
		return err
	}
	if patient == nil {
		return errors.New("patient not found")
	}

	return s.repo.Delete(hospitalID, patientID)
}

// normalizeOptionalPtr trims and returns nil if empty/whitespace, else pointer.
func normalizeOptionalPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := strings.TrimSpace(s)
	return &v
}
