package api

const (
	PageNumberKey     = "page_number"
	PageSizeKey       = "page_size"
	DefaultPageNumber = 1
	DefaultPageSize   = 100
)

type Pagination struct {
	PageNumber int
	PageSize   int
}

func (p Pagination) Limit() int {
	return p.PageSize
}

func (p Pagination) Offset() int {
	return p.Limit() * (p.PageNumber - 1)
}
