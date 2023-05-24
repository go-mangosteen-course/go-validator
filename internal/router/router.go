package router

import (
	"mangosteen/config"
	"mangosteen/internal/controller"
	"mangosteen/internal/database"
	"mangosteen/internal/middleware"

	"mangosteen/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swag
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func loadControllers() []controller.Controller {
	return []controller.Controller{
		&controller.SessionController{},
		&controller.ValidationCodeController{},
	}
}

func New() *gin.Engine {
	config.LoadAppConfig()
	r := gin.Default()
	r.Use(middleware.Me())
	docs.SwaggerInfo.Version = "1.0"

	database.Connect()

	api := r.Group("/api")

	for _, ctrl := range loadControllers() {
		ctrl.RegisterRoutes(api)
	}

	r.GET("/ping", controller.Ping)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
