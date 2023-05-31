package controller

import (
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagController struct {
}

func (ctrl *TagController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/tags", ctrl.Create)
}

func (ctrl *TagController) Create(c *gin.Context) {
	var body api.CreateTagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误")
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.CreateTag(c, queries.CreateTagParams{
		UserID: user.ID,
		Name:   body.Name,
		Sign:   body.Sign,
		Kind:   body.Kind,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": tag,
	})
}

func (ctrl *TagController) Destroy(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}
