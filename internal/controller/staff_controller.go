package controller

import (
	"net/http"

	"github.com/PhanukornKMITL/hospital-exam/internal/dto"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StaffController struct {
	service service.StaffService
}

func NewStaffController(service service.StaffService) *StaffController {
	return &StaffController{service: service}
}

func (s *StaffController) GetStaffs(c *gin.Context) {
	staffs, err := s.service.GetStaffs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toStaffResponses(staffs))
}

func (s *StaffController) CreateStaff(c *gin.Context) {
	var req dto.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirmPassword do not match"})
		return
	}

	staff, err := s.service.CreateStaff(service.StaffCreateInput{
		Username:   req.Username,
		Password:   req.Password,
		HospitalID: req.HospitalID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toStaffResponse(*staff))
}

func (s *StaffController) Login(c *gin.Context) {
	var req dto.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.service.Login(service.StaffLoginInput{
		Username:   req.Username,
		Password:   req.Password,
		HospitalID: req.HospitalID,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.StaffLoginResponse{Token: token})
}

func (s *StaffController) DeleteStaff(c *gin.Context) {
	idParam := c.Param("id")
	staffID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff id"})
		return
	}

	if err := s.service.DeleteStaff(staffID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func toStaffResponse(s entity.Staff) dto.StaffResponse {
	return dto.StaffResponse{
		ID:         s.ID,
		Username:   s.Username,
		HospitalID: s.HospitalID,
		CreatedAt:  s.CreatedAt,
	}
}

func toStaffResponses(staffs []entity.Staff) []dto.StaffResponse {
	resp := make([]dto.StaffResponse, 0, len(staffs))
	for _, st := range staffs {
		resp = append(resp, toStaffResponse(st))
	}
	return resp
}
