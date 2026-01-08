package entity

import (
	"time"

	"github.com/google/uuid"
)

type Hospital struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"not null"`
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
