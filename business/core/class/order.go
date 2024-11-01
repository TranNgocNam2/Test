package class

import "Backend/internal/order"

const (
	OrderByName      = "name"
	OrderByCode      = "code"
	OrderByStartDate = "start_date"
)

var orderByFields = map[string]string{
	OrderByName:      OrderByName,
	OrderByCode:      OrderByCode,
	OrderByStartDate: OrderByStartDate,
}

var DefaultOrderBy = order.NewBy(OrderByName, order.DESC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
