package controller

import (
	"log"
	"mangosteen/api"
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

// CreateSession godoc
//
//	@Summary	登录
//	@Accept		json
//	@Produce	json
//
//	@Param		email	body		string					true	"邮件地址"
//	@Param		code	body		string					true	"验证码"
//
//	@Success	200		{object}	api.CreateItemResponse	数据
//	@Failure	400		{string}	string					无效的验证码
//	@Failure	500		{string}	string					服务器错误
//	@Router		/api/v1/session [post]
func (ctrl *SessionController) Create(c *gin.Context) {
	var requestBody api.CreateSessionRequest
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
	res := api.CreateSessionResponse{
		Jwt:    jwt,
		UserID: user.ID,
	}
	c.JSON(http.StatusOK, res)
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
