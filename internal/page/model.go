package page

type Response[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func NewPageResponse[T any](items []T, total int, page int, size int) Response[T] {
	return Response[T]{
		Items: items,
		Page:  page,
		Size:  size,
		Total: total,
	}
}
