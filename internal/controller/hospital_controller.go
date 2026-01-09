package controller

import (
	"net/http"

	"github.com/PhanukornKMITL/hospital-exam/internal/dto"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// UpdateHospital godoc
// @Summary Update a hospital
// @Description Update hospital information by ID
// @Tags hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID (UUID)"
// @Param request body dto.UpdateHospitalRequest true "Updated hospital details"
// @Success 200 {object} dto.HospitalResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospital/{id} [put]
func (h *HospitalController) UpdateHospital(c *gin.Context) {
	idParam := c.Param("id")
	hospitalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital id"})
		return
	}

	var req dto.UpdateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospital, err := h.service.UpdateHospital(hospitalID, service.HospitalUpdateInput{
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		if err.Error() == "hospital not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toHospitalResponse(*hospital))
}

// DeleteHospital godoc
// @Summary Delete a hospital
// @Description Delete a hospital by ID
// @Tags hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID (UUID)"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospital/{id} [delete]
func (h *HospitalController) DeleteHospital(c *gin.Context) {
	idParam := c.Param("id")
	hospitalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital id"})
		return
	}

	if err := h.service.DeleteHospital(hospitalID); err != nil {
		if err.Error() == "hospital not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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
