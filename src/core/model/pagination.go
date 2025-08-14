package model

type Pagination struct {
	Total   int64 `json:"total,omitempty"`
	Limit   int64 `json:"limit,omitempty"`
	Page    int64 `json:"page,omitempty"`
	HasMore bool  `json:"has_more"`
}
