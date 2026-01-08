package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhanukornKMITL/hospital-exam/internal/controller"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/PhanukornKMITL/hospital-exam/internal/config"
)

// SetupRoutes สร้าง repo -> service -> controller ภายในตัว route
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// สร้าง layers ของ Hospital
	hospitalRepo := repository.NewHospitalRepository(db)
	hospitalService := service.NewHospitalService(hospitalRepo)
	hospitalController := controller.NewHospitalController(hospitalService)

	// Routes ของ hospital
	r.GET("/hospitals", hospitalController.GetHospitals)
	r.POST("/hospitals", hospitalController.CreateHospital)

	// Health check route
	cfg := config.Load()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.AppEnv,
		})
	})
}
