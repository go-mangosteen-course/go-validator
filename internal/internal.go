package internal

import (
	"mangosteen/config"
	"mangosteen/internal/database"
	"mangosteen/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	config.LoadAppConfig()
	database.Connect()
	r.Use(middleware.Me([]string{
		"/ping",
		"/api/v1/session",
		"/api/v1/validation_codes",
	}))
}
