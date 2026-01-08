package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
	patients, err := p.service.GetPatients()
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
		HospitalID:   req.HospitalID,
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
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toPatientResponse(*patient))
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
		NationalID:   p.NationalID,
		PassportID:   p.PassportID,
		PhoneNumber:  p.PhoneNumber,
		Email:        p.Email,
		Gender:       p.Gender,
		CreatedAt:    p.CreatedAt,
	}
}

func toPatientResponses(patients []entity.Patient) []dto.PatientResponse {
	resp := make([]dto.PatientResponse, 0, len(patients))
	for _, p := range patients {
		resp = append(resp, toPatientResponse(p))
	}
	return resp
}
