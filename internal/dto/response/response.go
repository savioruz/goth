package response

type Response[T any] struct {
	Data       T       `json:"data"`
	Pagination *Paging `json:"pagination,omitempty"`
}

type Paging struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPage  int `json:"total_page"`
	TotalCount int `json:"total_count"`
}

func NewResponse[T any](data T, pagination *Paging) *Response[T] {
	return &Response[T]{
		Data:       data,
		Pagination: pagination,
	}
}
