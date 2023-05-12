package controller

import (
	"log"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionController struct {
}

func (ctrl *SessionController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("session", ctrl.Create)
}
func (ctrl *SessionController) Create(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.String(http.StatusBadRequest, "无效的参数")
		return
	}
	q := database.NewQuery()
	_, err := q.FindValidationCode(c, queries.FindValidationCodeParams{
		Email: requestBody.Email,
		Code:  requestBody.Code,
	})
	if err != nil {
		c.String(http.StatusBadRequest, "无效的验证码")
		return
	}
	user, err := q.FindUserByEmail(c, requestBody.Email)
	if err != nil {
		user, err = q.CreateUser(c, requestBody.Email)
		if err != nil {
			log.Println("CreateUser fail", err)
			c.String(http.StatusInternalServerError, "请稍后再试")
			return
		}
	}

	jwt, err := jwt_helper.GenerateJWT(int(user.ID))
	if err != nil {
		log.Println("GenerateJWT fail", err)
		c.String(http.StatusInternalServerError, "请稍后再试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"jwt":    jwt,
		"userId": user.ID,
	})
}

func (ctrl *SessionController) Destroy(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *SessionController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}
