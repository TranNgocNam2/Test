package learner

import "Backend/internal/order"

const (
	OrderByFullName = "full_name"
)

var orderByFields = map[string]string{
	OrderByFullName: OrderByFullName,
}

var DefaultOrderBy = order.NewBy(OrderByFullName, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
