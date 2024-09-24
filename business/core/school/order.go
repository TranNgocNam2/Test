package school

import "Backend/internal/order"

var DefaultOrderBy = order.NewBy(OrderByName, order.ASC)

const (
	OrderByID   = "id"
	OrderByName = "name"
)
