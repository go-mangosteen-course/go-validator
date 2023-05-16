package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
//
//	@Summary	健康度检查
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Failure	500
//	@Router		/ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
