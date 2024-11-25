package order

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var (
	ErrInvalidDirection = errors.New("Không thể xác định được kiểu xác định: %s!")
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

	sortBy := c.DefaultQuery("sortBy", "ASC")
	strings.ToUpper(sortBy)

	var by By

	if _, exists := directions[sortBy]; !exists {
		return defaultOrder, fmt.Errorf(ErrInvalidDirection.Error(), by.Direction)
	}

	by = NewBy(strings.Trim(orderBy, " "), sortBy)

	return by, nil
}
