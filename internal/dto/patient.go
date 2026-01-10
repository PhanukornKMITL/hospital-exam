package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreatePatientRequest defines input payload for creating a patient.
type CreatePatientRequest struct {
	FirstNameTH  string `json:"firstNameTH" binding:"required"`
	MiddleNameTH string `json:"middleNameTH"`
	LastNameTH   string `json:"lastNameTH" binding:"required"`

	FirstNameEN  string `json:"firstNameEN" binding:"required"`
	MiddleNameEN string `json:"middleNameEN"`
	LastNameEN   string `json:"lastNameEN" binding:"required"`

	DateOfBirth *string `json:"dateOfBirth" binding:"required" example:"1990-01-15"` // expected format: YYYY-MM-DD

	NationalID  string `json:"nationalId"`
	PassportID  string `json:"passportId"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Email       string `json:"email"`
	Gender      string `json:"gender" binding:"required,oneof=M F"`
}

// UpdatePatientRequest allows partial updates; only provided fields will be applied.
type UpdatePatientRequest struct {
	FirstNameTH  *string `json:"firstNameTH"`
	MiddleNameTH *string `json:"middleNameTH"`
	LastNameTH   *string `json:"lastNameTH"`

	FirstNameEN  *string `json:"firstNameEN"`
	MiddleNameEN *string `json:"middleNameEN"`
	LastNameEN   *string `json:"lastNameEN"`

	DateOfBirth *string `json:"dateOfBirth" example:"1990-01-15"`

	NationalID  *string `json:"nationalId"`
	PassportID  *string `json:"passportId"`
	PhoneNumber *string `json:"phoneNumber"`
	Email       *string `json:"email"`
	Gender      *string `json:"gender"`
}

// PatientResponse represents patient data returned to clients.
type PatientResponse struct {
	ID         uuid.UUID `json:"id"`
	HospitalID uuid.UUID `json:"hospitalId"`
	PatientHN  string    `json:"patientHN"`

	FirstNameTH  string `json:"firstNameTH"`
	MiddleNameTH string `json:"middleNameTH"`
	LastNameTH   string `json:"lastNameTH"`

	FirstNameEN  string `json:"firstNameEN"`
	MiddleNameEN string `json:"middleNameEN"`
	LastNameEN   string `json:"lastNameEN"`

	DateOfBirth *time.Time `json:"dateOfBirth"`

	NationalID  string `json:"nationalId"`
	PassportID  string `json:"passportId"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`

	CreatedAt time.Time `json:"createdAt"`
}

// SearchPatientRequest defines filterable fields and pagination for patient search.
// All fields are optional; empty values are ignored.
type SearchPatientRequest struct {
	PatientHN    string `json:"patientHN"`
	FirstNameTH  string `json:"firstNameTH"`
	MiddleNameTH string `json:"middleNameTH"`
	LastNameTH   string `json:"lastNameTH"`

	FirstNameEN  string `json:"firstNameEN"`
	MiddleNameEN string `json:"middleNameEN"`
	LastNameEN   string `json:"lastNameEN"`

	DateOfBirth *string `json:"dateOfBirth" example:"1990-01-15"` // YYYY-MM-DD

	NationalID  string `json:"nationalId"`
	PassportID  string `json:"passportId"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`

	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedPatientsResponse struct {
	Data  []PatientResponse `json:"data"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Total int64             `json:"total"`
}
