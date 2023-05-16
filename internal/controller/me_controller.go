package controller

import (
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MeController struct {
}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.GET("/me", ctrl.Get)
}

func (ctrl *MeController) Get(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	jwtString := auth[7:]
	t, err := jwt_helper.Parse(jwtString)
	if err != nil {
		c.String(401, "无效的JWT")
		return
	}
	m, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		c.String(401, "无效的JWT")
		return
	}
	userID, ok := m["user_id"].(string)
	if !ok {
		c.String(401, "无效的JWT")
		return
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.String(401, "无效的JWT")
		return
	}
	q := database.NewQuery()
	userIDInt32 := int32(userIDInt)
	user, err := q.FindUser(c, userIDInt32)
	if err != nil {
		c.String(401, "无效的JWT")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": user,
	})
}

func (ctrl *MeController) Create(c *gin.Context) {
	panic("not implemented")
}

func (ctrl *MeController) Destroy(c *gin.Context) {
	panic("not implemented")
}

func (ctrl *MeController) Update(c *gin.Context) {
	panic("not implemented")
}

func (ctrl *MeController) GetPaged(c *gin.Context) {
	panic("not implemented")
}
