package entity

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	HospitalID uuid.UUID `gorm:"type:uuid;not null"`

	PatientHN string `gorm:"not null"`

	FirstNameTH  string `gorm:"not null"`
	MiddleNameTH string
	LastNameTH   string `gorm:"not null"`

	FirstNameEN  string `gorm:"not null"`
	MiddleNameEN string
	LastNameEN   string `gorm:"not null"`

	DateOfBirth *time.Time `gorm:"not null"`

	NationalID *string
	PassportID *string

	PhoneNumber string
	Email       string
	Gender      string `gorm:"type:varchar(1);not null;check:gender IN ('M', 'F')"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
