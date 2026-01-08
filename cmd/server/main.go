package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/PhanukornKMITL/hospital-exam/internal/config"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/route"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	db := config.NewPostgres()

	db.AutoMigrate(
		&entity.Hospital{},
		&entity.Staff{},
		&entity.Patient{},
		&entity.PatientSequence{},
	)

	r := gin.Default()

	route.SetupHospitalRoutes(r, db)
	route.SetupStaffRoutes(r, db)
	route.SetupPatientRoutes(r, db)
	route.SetupHealthRoute(r)

	r.Run(":" + cfg.AppPort)
}
