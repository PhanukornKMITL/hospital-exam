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
	)

	r := gin.Default()

	route.SetupRoutes(r, db)

	r.Run(":" + cfg.AppPort)
}
