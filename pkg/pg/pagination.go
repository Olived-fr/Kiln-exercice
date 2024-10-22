package pg

import "fmt"

const (
	LimitNone              = 0
	DefaultPaginationOrder = "DESC"
)

type Pagination struct {
	limit   int
	offset  int
	orderBy string
	order   string
}

func (p Pagination) Embed(query string) string {
	if p.OrderBy() != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", p.OrderBy(), p.Order())
	}

	if p.Limit() != LimitNone {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", p.Limit(), p.Offset())
	}

	return query
}

func NewPagination(limit, offset int, orderBy string) Pagination {
	return Pagination{
		limit:   limit,
		offset:  offset,
		orderBy: orderBy,
	}
}

func (p Pagination) Limit() int {
	return p.limit
}

func (p Pagination) Offset() int {
	return p.offset
}

func (p Pagination) OrderBy() string {
	return p.orderBy
}

func (p Pagination) Order() string {
	if p.order == "" {
		return DefaultPaginationOrder
	}

	return p.order
}
