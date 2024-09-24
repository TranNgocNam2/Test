package schooldb

import (
	"Backend/business/core/school"
	"Backend/internal/order"
	"fmt"
)

var orderByFields = map[string]string{
	school.OrderByID:   "id",
	school.OrderByName: "name",
}

func orderByClause(orderBy order.By) error {
	_, exists := orderByFields[orderBy.Field]
	if !exists {
		_ = fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return nil
}
