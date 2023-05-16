package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping
// @Summary      健康度检查
// @Description  如果 ping 返回 200 说明网站还在运行
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
