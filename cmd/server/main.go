package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/PhanukornKMITL/hospital-exam/docs"
	"github.com/PhanukornKMITL/hospital-exam/internal/config"
	"github.com/PhanukornKMITL/hospital-exam/internal/entity"
	"github.com/PhanukornKMITL/hospital-exam/internal/route"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hospital Middleware API
// @version 1.0
// @description API for Hospital Middleware System - Patient Information Search
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@hospital.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

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

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	route.SetupHospitalRoutes(r, db)
	route.SetupStaffRoutes(r, db)
	route.SetupPatientRoutes(r, db)
	route.SetupHealthRoute(r)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.AppPort)
}
