package program

import "Backend/internal/order"

const (
	OrderByName      = "name"
	OrderByStartDate = "start_date"
)

var orderByFields = map[string]string{
	OrderByName:      OrderByName,
	OrderByStartDate: OrderByStartDate,
}

var DefaultOrderBy = order.NewBy(OrderByStartDate, order.DESC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
