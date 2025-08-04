package pagination

type PageRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

type PageMetadata struct {
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
	CurrentPage  int   `json:"currentPage"`
	PageSize     int   `json:"pageSize"`
	NextPage     bool  `json:"nextPage"`
	PreviousPage bool  `json:"previousPage"`
}

type PageResponse[T any] struct {
	Data     []T          `json:"data"`
	Metadata PageMetadata `json:"metadata"`
}
