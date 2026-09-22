package dto

type ListProductsRequestDTO struct {
	NextCursor string `json:"next_cursor"`
	Limit      int64  `json:"limit"`
}
