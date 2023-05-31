package api

import "mangosteen/config/queries"

type CreateTagRequest struct {
	Name string       `json:"name" binding:"required"`
	Sign string       `json:"sign" binding:"required"`
	Kind queries.Kind `json:"kind" binding:"required"`
}
type CreateTagResponse struct {
	Resource queries.Tag `json:"resource"`
}
