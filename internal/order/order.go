package order

import (
	"Backend/internal/validate"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

type By struct {
	Field     string `json:"sortBy"`
	Direction string `json:"orderBy"`
}

func NewBy(field string, direction string) By {
	if _, exists := directions[direction]; !exists {
		return By{
			Field:     field,
			Direction: ASC,
		}
	}

	return By{
		Field:     field,
		Direction: direction,
	}
}

func Parse(c *gin.Context, defaultOrder By) (By, error) {
	orderBy := c.Query("orderBy")

	parts := strings.Split(orderBy, ",")

	var by By

	switch len(parts) {
	case 1:
		by = NewBy(strings.TrimSpace(parts[0]), ASC)

	case 2:
		direction := strings.Trim(parts[1], " ")
		if _, exists := directions[direction]; !exists {
			return By{}, validate.NewFieldsError(orderBy, fmt.Errorf("unknown direction: %s", by.Direction))
		}

		by = NewBy(strings.Trim(parts[0], " "), direction)

	default:
		return By{}, validate.NewFieldsError(orderBy, errors.New("unknown order field"))
	}

	return by, nil
}
