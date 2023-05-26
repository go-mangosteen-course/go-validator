package api

import (
	"mangosteen/config/queries"
	"time"
)

type CreateItemRequest struct {
	Amount     int32        `json:"amount" binding:"required"`
	Kind       queries.Kind `json:"kind" binding:"required"`
	HappenedAt time.Time    `json:"happened_at" binding:"required"`
	TagIds     []int32      `json:"tag_ids" binding:"required"`
}

type CreateItemResponse struct {
	Resource queries.Item
}
