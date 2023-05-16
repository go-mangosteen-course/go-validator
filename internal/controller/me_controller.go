package controller

import (
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"net/http"

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
	if len(auth) < 8 {
		c.String(401, "JWT为空")
		return
	}
	jwtString := auth[7:]
	if len(jwtString) == 0 {
		c.String(401, "JWT为空")
		return
	}
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
	userID, ok := m["user_id"].(float64)
	if !ok {
		c.String(401, "无效的JWT")
		return
	}
	userIDInt := int32(userID)
	if err != nil {
		c.String(401, "无效的JWT")
		return
	}
	q := database.NewQuery()
	user, err := q.FindUser(c, userIDInt)
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
