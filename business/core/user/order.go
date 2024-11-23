package user

import "Backend/internal/order"

const (
	OrderByName = "name"
)

var orderByFields = map[string]string{
	OrderByName: "name"}

var DefaultOrderBy = order.NewBy(OrderByName, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
