package entity

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	HospitalID uuid.UUID `gorm:"type:uuid;not null"`

	PatientHN string `gorm:"not null"`

	FirstNameTH  string
	MiddleNameTH string
	LastNameTH   string

	FirstNameEN  string
	MiddleNameEN string
	LastNameEN   string

	DateOfBirth *time.Time

	NationalID *string
	PassportID *string

	PhoneNumber string
	Email       string
	Gender      string

	CreatedAt time.Time
}
