package router

import (
	"mangosteen/config"
	"mangosteen/internal/controller"
	"mangosteen/internal/database"

	"mangosteen/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swag
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func New() *gin.Engine {
	config.LoadAppConfig()
	r := gin.Default()
	docs.SwaggerInfo.Version = "1.0"

	database.Connect()

	r.GET("/api/v1/ping", controller.Ping)
	r.POST("/api/v1/validation_codes", controller.CreateValidationCode)
	r.POST("/api/v1/session", controller.CreateSession)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
