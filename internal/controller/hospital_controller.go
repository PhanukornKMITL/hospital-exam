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
