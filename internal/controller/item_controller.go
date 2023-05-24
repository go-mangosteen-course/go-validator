package controller

import (
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemController struct{}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
}

func (ctrl *ItemController) Create(c *gin.Context) {
	var body struct {
		Amount     int32        `json:"amount" binding:"required"`
		Kind       queries.Kind `json:"kind" binding:"required"`
		HappenedAt time.Time    `json:"happened_at" binding:"required"`
		TagIds     []int32      `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误")
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	item, err := q.CreateItem(c, queries.CreateItemParams{
		UserID:     user.ID,
		Amount:     body.Amount,
		Kind:       body.Kind,
		HappenedAt: body.HappenedAt,
		TagIds:     body.TagIds,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": item,
	})
}

func (ctrl *ItemController) Destroy(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}
