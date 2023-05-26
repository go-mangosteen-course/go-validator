package controller

import (
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemController struct{}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
}

// CreateItem godoc
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//
//	@Param		amount		body		int						true	"金额（单位：分）" example(100)
//	@Param		kind		body		queries.Kind			true	"类型"
//	@Param		happened_at	body		string					true	"发生时间"
//	@Param		tag_ids		body		[]string				true	"标签ID列表"
//
//	@Success	200			{object}	api.CreateItemResponse	数据
//	@Failure	422			{string}	string					参数错误
//	@Router		/api/v1/items [post]
func (ctrl *ItemController) Create(c *gin.Context) {
	var body api.CreateItemRequest
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
