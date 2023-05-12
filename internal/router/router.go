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

	api := r.Group("/api")

	sc := controller.SessionController{}
	sc.RegisterRoutes(api)
	vcc := controller.ValidationCodeController{}
	vcc.RegisterRoutes(api)

	r.GET("/ping", controller.Ping)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
