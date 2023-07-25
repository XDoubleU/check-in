package dtos

type PaginatedResultDto[T any] struct {
	Data       []*T       `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
} //	@name	Pagination
