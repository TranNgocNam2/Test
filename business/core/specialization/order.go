package specialization

import "Backend/internal/order"

const (
	OrderByName = "name"
	OrderByCode = "code"
)

var orderByFields = map[string]string{
	OrderByName: "name",
	OrderByCode: "code",
}

var DefaultOrderBy = order.NewBy(OrderByCode, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
