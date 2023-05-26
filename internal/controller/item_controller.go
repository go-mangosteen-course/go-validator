package controller

import (
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	PerPage int32
}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
	v1.GET("/items", ctrl.GetPaged)
	ctrl.PerPage = 10
}

// CreateItem godoc
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//
//	@Param		amount		body		int						true	"金额（单位：分）"	example(100)
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

// GetPagesItems godoc
//
//	@Summary	获取分页账目
//	@Accept		json
//	@Produce	json
//
//	@Param		page		query		int						false	"页码"	example(100)
//	@Param		happened_after		query		string			false	"开始时间" example(2000-01-01T01:01:01+0800)
//	@Param		happened_before	query		string					false	"结束时间" example(2000-01-01T01:01:01+0800)
//
//	@Success	200			{object}	api.GetPagesItemsResponse	数据
//	@Failure	500
//	@Router		/api/v1/items [get]
func (ctrl *ItemController) GetPaged(c *gin.Context) {
	var params api.GetPagedItemsRequest
	pageStr, _ := c.Params.Get("page")
	if page, err := strconv.Atoi(pageStr); err == nil {
		params.Page = int32(page)
	}
	if params.Page == 0 {
		params.Page = 1
	}
	happenedBefore, has := c.Params.Get("happened_before")
	if has {
		if t, err := time.Parse(time.RFC3339, happenedBefore); err == nil {
			params.HappenedBefore = t
		}
	}

	happenedAfter, has := c.Params.Get("happened_after")
	if has {
		if t, err := time.Parse(time.RFC3339, happenedAfter); err == nil {
			params.HappenedAfter = t
		}
	}

	q := database.NewQuery()
	items, err := q.ListItems(c, queries.ListItemsParams{
		Offset: (params.Page - 1) * ctrl.PerPage,
		Limit:  ctrl.PerPage,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	count, err := q.CountItems(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	c.JSON(http.StatusOK, api.GetPagesItemsResponse{
		Resources: items,
		Pager: api.Pager{
			Page:    params.Page,
			PerPage: ctrl.PerPage,
			Count:   count,
		},
	})

}
