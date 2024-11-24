package user

import "Backend/internal/order"

const (
	OrderByFullName  = "name"
	OrderByCreatedAt = "created_at"
)

var orderByFields = map[string]string{
	OrderByCreatedAt: "vl.created_at",
	OrderByFullName:  "u.full_name",
}

var DefaultOrderBy = order.NewBy(OrderByFullName, order.ASC)

func orderByClause(orderBy order.By) string {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		by = DefaultOrderBy.Field
	}
	return " ORDER BY " + by + " " + orderBy.Direction
}
