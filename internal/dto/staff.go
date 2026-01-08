package dto

import (
    "time"

    "github.com/google/uuid"
)

// CreateStaffRequest defines input payload for creating a staff.
type CreateStaffRequest struct {
    Username   string    `json:"username" binding:"required"`
    Password   string    `json:"password" binding:"required"`
    ConfirmPassword string `json:"confirmPassword" binding:"required"`
    HospitalID uuid.UUID `json:"hospitalId" binding:"required"`
}

// LoginStaffRequest defines input for staff login.
type LoginStaffRequest struct {
    Username   string    `json:"username" binding:"required"`
    Password   string    `json:"password" binding:"required"`
    HospitalID uuid.UUID `json:"hospitalId" binding:"required"`
}

// StaffResponse represents staff data returned to clients.
type StaffResponse struct {
    ID         uuid.UUID `json:"id"`
    Username   string    `json:"username"`
    HospitalID uuid.UUID `json:"hospitalId"`
    CreatedAt  time.Time `json:"createdAt"`
}

// StaffLoginResponse represents JWT token response.
type StaffLoginResponse struct {
    Token string `json:"token"`
}
