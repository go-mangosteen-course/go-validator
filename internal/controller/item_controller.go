package controller

import (
	"log"
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/nav-inc/datetime"
)

type ItemController struct {
	PerPage int32
}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
	v1.GET("/items", ctrl.GetPaged)
	v1.GET("/items/balance", ctrl.GetBalance)
	v1.GET("/items/summary", ctrl.GetSummary)
	ctrl.PerPage = 10
}

// CreateItem godoc
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//
//	@Param		amount		body		int						true	"金额（单位：分）"	example(100)
//	@Param		kind		body		string			true	"类型"
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
		return
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

func (ctrl *ItemController) GetBalance(c *gin.Context) {
	query := c.Request.URL.Query()
	happenedAfterString := query.Get("happened_after")
	happenedBeforeString := query.Get("happened_before")
	happenedAfter, err := datetime.Parse(happenedAfterString, time.Local)
	if err != nil {
		happenedAfter = time.Now().AddDate(-100, 0, 0)
	}
	happenedBefore, err := datetime.Parse(happenedBeforeString, time.Local)
	if err != nil {
		happenedBefore = time.Now().AddDate(1, 0, 0)
	}

	q := database.NewQuery()
	items, err := q.ListItemsHappenedBetween(c, queries.ListItemsHappenedBetweenParams{
		HappenedAfter:  happenedAfter,
		HappenedBefore: happenedBefore,
	})
	if err != nil {
		log.Printf("list items error: %v", err)
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	var r api.GetBalanceResponse
	for _, item := range items {
		if item.Kind == "in_come" {
			r.Income += int(item.Amount)
		} else {
			r.Expenses += int(item.Amount)
		}
	}
	r.Balance = r.Income - r.Expenses
	c.JSON(http.StatusOK, r)
}

// GetPagesItems godoc
//
//	@Summary	获取分页账目
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//
//	@Param		page			query		int							false	"页码"	example(100)
//	@Param		happened_after	query		string						false	"开始时间"	example(2000-01-01T01:01:01+0800)
//	@Param		happened_before	query		string						false	"结束时间"	example(2000-01-01T01:01:01+0800)
//
//	@Success	200				{object}	api.GetPagesItemsResponse	数据
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
		if t, err := datetime.Parse(happenedBefore, time.Local); err == nil {
			params.HappenedBefore = t
		}
	}

	happenedAfter, has := c.Params.Get("happened_after")
	if has {
		if t, err := datetime.Parse(happenedAfter, time.Local); err == nil {
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

func (ctrl *ItemController) GetSummary(c *gin.Context) {
	var query api.GetSummaryRequest
	if err := c.BindQuery(&query); err != nil {
		r := api.NewErrorResponse()
		switch x := err.(type) {
		case validator.ValidationErrors:
			for _, ve := range x {
				tag := ve.Tag()
				field := ve.Field()
				if r.Errors[field] == nil {
					r.Errors[field] = []string{}
				}
				r.Errors[field] = append(r.Errors[field], tag)
			}
			c.JSON(http.StatusUnprocessableEntity, r)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	me, _ := c.Get("me")
	user, _ := me.(queries.User)

	q := database.NewQuery()
	items, err := q.ListItemsByHappenedAtAndKind(c, queries.ListItemsByHappenedAtAndKindParams{
		HappenedAfter:  query.HappenedAfter,
		HappenedBefore: query.HappenedBefore,
		Kind:           query.Kind,
		UserID:         user.ID,
	})

	if err != nil {
		log.Printf("list items error: %v", err)
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	// ---------------

	if query.GroupBy == "happened_at" {
		res := api.NewGetSummaryByHappenedAtResponse()

		for _, item := range items {
			k := item.HappenedAt.Format("2006-01-02")
			res.Total += int(item.Amount)

			found := false
			for index, group := range res.Groups {
				if group.HappenedAt == k {
					found = true
					res.Groups[index].Amount += int(item.Amount)
				}
			}
			if !found {
				res.Groups = append(res.Groups, api.SummaryGroupByHappenedAt{
					HappenedAt: k,
					Amount:     int(item.Amount),
				})
			}
		}
		sort.Slice(res.Groups, func(i, j int) bool {
			return res.Groups[i].HappenedAt < res.Groups[j].HappenedAt
		})

		c.JSON(http.StatusOK, res)
	} else if query.GroupBy == "tag_id" {
		res := api.NewGetSummaryByTagIDResponse()

		for _, item := range items {
			if len(item.TagIds) == 0 {
				continue
			}
			k := item.TagIds[0]
			res.Total += int(item.Amount)
			found := false
			for index, group := range res.Groups {
				if group.TagID == k {
					found = true
					res.Groups[index].Amount += int(item.Amount)
				}
			}
			if !found {
				res.Groups = append(res.Groups, api.SummaryGroupByTagID{
					TagID:  k,
					Amount: int(item.Amount),
				})
			}
		}
		c.JSON(http.StatusOK, res)
	} else {
		c.String(http.StatusBadRequest, "参数错误")
	}
}
