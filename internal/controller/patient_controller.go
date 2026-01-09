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

// GetPatients godoc
// @Summary Get patients by hospital
// @Description Retrieve patients only within the authenticated staff's hospital (no cross-hospital access)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PatientResponse
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /patients [get]
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

// CreatePatient godoc
// @Summary Create a new patient
// @Description Register a new patient only into the authenticated staff's hospital (no cross-hospital access)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePatientRequest true "Patient details"
// @Success 201 {object} dto.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /patients [post]
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

// SearchPatient godoc
// @Summary Search patient by identifier
// @Description Retrieve patient by identifier only within the authenticated staff's hospital (no cross-hospital access)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient identifier (nationalId or passportId)"
// @Success 200 {object} dto.PatientResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /patients/{id} [get]
func (p *PatientController) SearchPatientByID(c *gin.Context) {
	identifier := c.Param("id")

	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identifier is required"})
		return
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

	patient, err := p.service.SearchPatientByID(hospitalID, identifier)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, toPatientResponse(*patient))
}

// SearchPatients handles POST /patient/search with filter DTO and pagination.
// @Summary Search patients with filters
// @Description Search patients only within the authenticated staff's hospital (no cross-hospital access)
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchPatientRequest true "Search filters"
// @Success 200 {object} dto.PaginatedPatientsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /patients/search [post]
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
