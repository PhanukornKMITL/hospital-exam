package route

import (
    "github.com/gin-gonic/gin"

    "github.com/PhanukornKMITL/hospital-exam/internal/config"
)

func SetupHealthRoute(r *gin.Engine) {
    cfg := config.Load()
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "env":    cfg.AppEnv,
        })
    })
}
