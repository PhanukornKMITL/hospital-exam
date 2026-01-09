package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhanukornKMITL/hospital-exam/internal/controller"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
)

func SetupHospitalRoutes(r *gin.Engine, db *gorm.DB) {
	hospitalRepo := repository.NewHospitalRepository(db)
	hospitalService := service.NewHospitalService(hospitalRepo)
	hospitalController := controller.NewHospitalController(hospitalService)

	r.GET("/hospital", hospitalController.GetHospitals)
	r.POST("/hospital", hospitalController.CreateHospital)
	r.PUT("/hospital/:id", hospitalController.UpdateHospital)
	r.DELETE("/hospital/:id", hospitalController.DeleteHospital)
}
