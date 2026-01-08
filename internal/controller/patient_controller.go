package controller

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/PhanukornKMITL/hospital-exam/internal/dto"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
)

type PatientController struct {
	service service.PatientService
}

func NewPatientController(service service.PatientService) *PatientController {
	return &PatientController{service: service}
}

func (p *PatientController) GetPatients(c *gin.Context) {
	hospitalIDStr, exists := c.Get("hospitalId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "hospitalId not found in token"})
		return
	}

	hospitalID, err := uuid.Parse(hospitalIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalId format"})
		return
	}

	patients, err := p.service.GetPatientsByHospital(hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toPatientResponses(patients))
}

func (p *PatientController) CreatePatient(c *gin.Context) {
	var req dto.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staffHospitalIDStr, exists := c.Get("hospitalId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "hospitalId not found in token"})
		return
	}

	staffHospitalID, err := uuid.Parse(staffHospitalIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalId format in token"})
		return
	}

	if req.NationalID == "" && req.PassportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either nationalId or passportId must be provided"})
		return
	}

	var dob *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dateOfBirth format, expected YYYY-MM-DD"})
			return
		}
		dob = &parsed
	}

	patient, err := p.service.CreatePatient(service.PatientCreateInput{
		HospitalID:   staffHospitalID,
		FirstNameTH:  req.FirstNameTH,
		MiddleNameTH: req.MiddleNameTH,
		LastNameTH:   req.LastNameTH,
		FirstNameEN:  req.FirstNameEN,
		MiddleNameEN: req.MiddleNameEN,
		LastNameEN:   req.LastNameEN,
		DateOfBirth:  dob,
		NationalID:   req.NationalID,
		PassportID:   req.PassportID,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Gender:       req.Gender,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toPatientResponse(*patient))
}

func (p *PatientController) SearchPatient(c *gin.Context) {
	identifier := c.Param("id")

	hospitalIDStr, exists := c.Get("hospitalId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "hospitalId not found in token"})
		return
	}

	hospitalID, err := uuid.Parse(hospitalIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalId format"})
		return
	}

	patient, err := p.service.SearchPatientByIdentifier(hospitalID, identifier)
	if err != nil {
		if err.Error() == "identifier is required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// If not found, return 404
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, toPatientResponse(*patient))
}

// SearchPatients handles POST /patient/search with filter DTO and pagination.
func (p *PatientController) SearchPatients(c *gin.Context) {
	var req dto.SearchPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If body is empty (EOF), treat as no filters provided
		if !errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// proceed with zero-value req (no filters)
	}

	hospitalIDStr, exists := c.Get("hospitalId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "hospitalId not found in token"})
		return
	}

	hospitalID, err := uuid.Parse(hospitalIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospitalId format"})
		return
	}

	var dob *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dateOfBirth format, expected YYYY-MM-DD"})
			return
		}
		dob = &parsed
	}

	patients, total, err := p.service.SearchPatients(hospitalID, service.PatientSearchInput{
		PatientHN:    req.PatientHN,
		FirstNameTH:  req.FirstNameTH,
		MiddleNameTH: req.MiddleNameTH,
		LastNameTH:   req.LastNameTH,
		FirstNameEN:  req.FirstNameEN,
		MiddleNameEN: req.MiddleNameEN,
		LastNameEN:   req.LastNameEN,
		DateOfBirth:  dob,
		NationalID:   req.NationalID,
		PassportID:   req.PassportID,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Gender:       req.Gender,
	}, req.Page, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedPatientsResponse{
		Data:  toPatientResponses(patients),
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	})
}

func toPatientResponse(p entity.Patient) dto.PatientResponse {
	return dto.PatientResponse{
		ID:           p.ID,
		HospitalID:   p.HospitalID,
		PatientHN:    p.PatientHN,
		FirstNameTH:  p.FirstNameTH,
		MiddleNameTH: p.MiddleNameTH,
		LastNameTH:   p.LastNameTH,
		FirstNameEN:  p.FirstNameEN,
		MiddleNameEN: p.MiddleNameEN,
		LastNameEN:   p.LastNameEN,
		DateOfBirth:  p.DateOfBirth,
		NationalID:   derefString(p.NationalID),
		PassportID:   derefString(p.PassportID),
		PhoneNumber:  p.PhoneNumber,
		Email:        p.Email,
		Gender:       p.Gender,
		CreatedAt:    p.CreatedAt,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toPatientResponses(patients []entity.Patient) []dto.PatientResponse {
	resp := make([]dto.PatientResponse, 0, len(patients))
	for _, p := range patients {
		resp = append(resp, toPatientResponse(p))
	}
	return resp
}
