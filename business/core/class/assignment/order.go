package assignment

import "Backend/internal/order"

const (
	OrderByDeadline = "deadline"
)

var orderByFields = map[string]string{
	OrderByDeadline: "deadline",
}

var DefaultOrderBy = order.NewBy(OrderByDeadline, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}

	return " ORDER BY " + by + " " + orderBy.Direction
}
