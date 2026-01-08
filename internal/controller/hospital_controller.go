package controller

import (
	"net/http"

	"github.com/PhanukornKMITL/hospital-exam/internal/dto"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/gin-gonic/gin"
)

type HospitalController struct {
	service service.HospitalService
}

func NewHospitalController(service service.HospitalService) *HospitalController {
	return &HospitalController{service: service}
}

// GetHospitals godoc
// @Summary Get all hospitals
// @Description Retrieve all hospitals in the system
// @Tags hospitals
// @Accept json
// @Produce json
// @Success 200 {array} dto.HospitalResponse
// @Failure 500 {object} map[string]string
// @Router /hospital [get]
func (h *HospitalController) GetHospitals(c *gin.Context) {
	hospitals, err := h.service.GetHospitals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, toHospitalResponses(hospitals))
}

// CreateHospital godoc
// @Summary Create a new hospital
// @Description Register a new hospital in the system
// @Tags hospitals
// @Accept json
// @Produce json
// @Param request body dto.CreateHospitalRequest true "Hospital details"
// @Success 201 {object} dto.HospitalResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospital [post]
func (h *HospitalController) CreateHospital(c *gin.Context) {
	var req dto.CreateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospital, err := h.service.CreateHospital(service.HospitalCreateInput{
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toHospitalResponse(*hospital))
}

func toHospitalResponse(h entity.Hospital) dto.HospitalResponse {
	return dto.HospitalResponse{
		ID:        h.ID,
		Name:      h.Name,
		Address:   h.Address,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func toHospitalResponses(hospitals []entity.Hospital) []dto.HospitalResponse {
	resp := make([]dto.HospitalResponse, 0, len(hospitals))
	for _, h := range hospitals {
		resp = append(resp, toHospitalResponse(h))
	}
	return resp
}
