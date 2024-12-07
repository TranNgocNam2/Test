package transcript

import "Backend/internal/order"

const (
	OrderByIndex = "t.index"
)

var orderByFields = map[string]string{
	OrderByIndex: OrderByIndex,
}

var DefaultOrderBy = order.NewBy(OrderByIndex, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
