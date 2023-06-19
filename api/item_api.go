package api

import (
	"mangosteen/config/queries"
	"time"
)

type CreateItemRequest struct {
	Amount     int32     `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIds     []int32   `json:"tag_ids" binding:"required"`
}

type CreateItemResponse struct {
	Resource queries.Item
}

type GetPagedItemsRequest struct {
	Page           int32     `json:"page"`
	HappenedAfter  time.Time `json:"happened_after"`
	HappenedBefore time.Time `json:"happened_before"`
}

type GetPagesItemsResponse struct {
	Resources []queries.Item `json:"resources"`
	Pager     Pager          `json:"pager"`
}

type GetBalanceResponse struct {
	Income   int `json:"income"`
	Expenses int `json:"expenses"`
	Balance  int `json:"balance"`
}

type GetSummaryRequest struct {
	HappenedAfter  time.Time `form:"happened_after"`
	HappenedBefore time.Time `form:"happened_before"`
	Kind           string    `form:"kind"`
	GroupBy        string    `form:"group_by"`
}
