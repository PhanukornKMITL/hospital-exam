package dto

import (
    "time"

    "github.com/google/uuid"
)

// CreateHospitalRequest defines input payload for creating a hospital.
type CreateHospitalRequest struct {
    Name    string `json:"name" binding:"required"`
    Address string `json:"address"`
}

// UpdateHospitalRequest defines input payload for updating a hospital.
type UpdateHospitalRequest struct {
    Name    string `json:"name" binding:"required"`
    Address string `json:"address"`
}

// HospitalResponse represents hospital data returned to clients.
type HospitalResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
