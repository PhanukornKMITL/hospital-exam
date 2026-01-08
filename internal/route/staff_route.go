package route

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/PhanukornKMITL/hospital-exam/internal/controller"
	"github.com/PhanukornKMITL/hospital-exam/internal/repository"
	"github.com/PhanukornKMITL/hospital-exam/internal/service"
)

func SetupStaffRoutes(r *gin.Engine, db *gorm.DB) {
	staffRepo := repository.NewStaffRepository(db)
	staffService := service.NewStaffService(staffRepo)
	staffController := controller.NewStaffController(staffService)

	// Routes ของ staff
	r.GET("/staff", staffController.GetStaffs)
	r.POST("/staff/create", staffController.CreateStaff)
	r.POST("/staff/login", staffController.Login)
	r.DELETE("/staff/:id", staffController.DeleteStaff)
}
