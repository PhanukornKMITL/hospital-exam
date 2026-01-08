package entity

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username   string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	HospitalID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
}
