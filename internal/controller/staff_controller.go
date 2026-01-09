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

// GetStaffs godoc
// @Summary Get all staffs
// @Description Retrieve all staff members
// @Tags staffs
// @Accept json
// @Produce json
// @Success 200 {array} dto.StaffResponse
// @Failure 500 {object} map[string]string
// @Router /staff [get]
func (s *StaffController) GetStaffs(c *gin.Context) {
	staffs, err := s.service.GetStaffs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toStaffResponses(staffs))
}

// CreateStaff godoc
// @Summary Create a new staff
// @Description Create a new staff member with credentials
// @Tags staffs
// @Accept json
// @Produce json
// @Param request body dto.CreateStaffRequest true "Staff creation details"
// @Success 201 {object} dto.StaffResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /staff/create [post]
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

// Login godoc
// @Summary Staff login
// @Description Authenticate staff and receive JWT token
// @Tags staffs
// @Accept json
// @Produce json
// @Param credentials body dto.LoginStaffRequest true "Login credentials"
// @Success 200 {object} dto.StaffLoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /staff/login [post]
func (s *StaffController) Login(c *gin.Context) {
	var req dto.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.service.Login(service.StaffLoginInput{
		Username: req.Username,
		// DeleteStaff godoc
		// @Summary Delete a staff
		// @Description Delete staff member by ID
		// @Tags staffs
		// @Accept json
		// @Produce json
		// @Param id path string true "Staff ID (UUID)"
		// @Success 204
		// @Failure 400 {object} map[string]string
		// @Failure 500 {object} map[string]string
		// @Router /staffs/{id} [delete]
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
