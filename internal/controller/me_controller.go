package controller

import (
	"mangosteen/config/queries"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MeController struct {
}

func (ctrl *MeController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.GET("/me", ctrl.Get)
}

// GetMe godoc
//
//	@Summary	获取当前用户
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Success	200	{object}	api.GetMeResponse
//	@Failure	401	{string}	JWT为空	|	无效的JWT
//	@Router		/api/v1/me [get]
func (ctrl *MeController) Get(c *gin.Context) {
	me, _ := c.Get("me")
	if user, ok := me.(queries.User); !ok {
		c.Status(http.StatusUnauthorized)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"resource": user,
		})
	}
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
