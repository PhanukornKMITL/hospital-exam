package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhanukornKMITL/hospital-exam/internal/controller"
	"github.com/PhanukornKMITL/hospital-exam/internal/middleware"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
)

func SetupPatientRoutes(r *gin.Engine, db *gorm.DB) {
	patientRepo := repository.NewPatientRepository(db)
	patientService := service.NewPatientService(patientRepo)
	patientController := controller.NewPatientController(patientService)

	// Routes ของ patient - require authentication
	patientGroup := r.Group("/patient")
	patientGroup.Use(middleware.AuthMiddleware())
	{
		patientGroup.GET("", patientController.GetPatients)
		patientGroup.POST("/create", patientController.CreatePatient)
	}
}
