package pagination

import "math"

func GetOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

func GetMetadata(totalRecords int64, page, limit int) PageMetadata {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	return PageMetadata{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     limit,
		NextPage:     page < totalPages,
		PreviousPage: page > 1,
	}
}
