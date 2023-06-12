package controller

import (
	"mangosteen/api"
	"mangosteen/config/queries"
	"mangosteen/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagController struct {
}

func (ctrl *TagController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/tags", ctrl.Create)
	v1.PATCH("/tags/:id", ctrl.Update)
	v1.DELETE("/tags/:id", ctrl.Destroy)
}

// CreateTag godoc
//
//	@Summary	创建标签
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//
//	@Param		name		body		string						true	"标签名"	SchemaExample(通勤)
//	@Param		sign	  body		string					  true	"符号"    SchemaExample(👟)
//	@Param		kind		body		queries.Kind			true	"类型"
//
//	@Success	200			{object}	api.CreateTagResponse	数据
//	@Failure	422			{string}	string					参数错误
//	@Router		/api/v1/tags [post]
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
	c.JSON(http.StatusOK, api.CreateTagResponse{Resource: tag})
}

func (ctrl *TagController) Destroy(c *gin.Context) {
	idString, has := c.Params.Get("id")
	if has == false {
		c.String(422, "参数错误")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.String(422, "参数错误")
		return
	}
	q := database.NewQuery()
	err = q.DeleteTag(c, int32(id))
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (ctrl *TagController) Update(c *gin.Context) {
	var body api.UpdateTagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误")
		return
	}
	idString, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.String(422, "参数错误")
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.UpdateTag(c, queries.UpdateTagParams{
		ID:     int32(id),
		UserID: user.ID,
		Name:   body.Name,
		Sign:   body.Sign,
		Kind:   body.Kind,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(http.StatusOK, api.UpdateTagResponse{Resource: tag})
}

func (ctrl *TagController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *TagController) GetPaged(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}
