package entity

import "github.com/google/uuid"

type PatientSequence struct {
	HospitalID uuid.UUID `gorm:"type:uuid;primaryKey"`
	NextNumber int       `gorm:"not null"`
}
