package api

type Pager struct {
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	Count   int64 `json:"count"`
}
