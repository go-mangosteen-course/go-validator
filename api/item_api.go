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
	HappenedAfter  time.Time `form:"happened_after" binding:"required"`
	HappenedBefore time.Time `form:"happened_before" binding:"required"`
	Kind           string    `form:"kind" binding:"required,oneof=expenses in_come"`
	GroupBy        string    `form:"group_by" binding:"required,oneof=happened_at tag_id"`
}

type GetSummaryByHappenedAtResponse struct {
	Groups []SummaryGroupByHappenedAt `json:"groups"`
	Total  int                        `json:"total"`
}

func NewGetSummaryByHappenedAtResponse() GetSummaryByHappenedAtResponse {
	return GetSummaryByHappenedAtResponse{
		Groups: []SummaryGroupByHappenedAt{},
		Total:  0,
	}
}

type GetSummaryByTagIDResponse struct {
	Groups []SummaryGroupByTagID `json:"groups"`
	Total  int                   `json:"total"`
}

func NewGetSummaryByTagIDResponse() GetSummaryByTagIDResponse {
	return GetSummaryByTagIDResponse{
		Groups: []SummaryGroupByTagID{},
		Total:  0,
	}
}

type SummaryGroupByHappenedAt struct {
	HappenedAt string `json:"happened_at"`
	Amount     int    `json:"amount"`
}
type SummaryGroupByTagID struct {
	TagID  int32       `json:"tag_id"`
	Tag    queries.Tag `json:"tag"`
	Amount int         `json:"amount"`
}
