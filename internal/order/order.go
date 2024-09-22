package order

import (
	"github.com/gin-gonic/gin"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

type By struct {
	Field     string `json:"sortBy"`
	Direction string `json:"orderBy"`
}

func NewBy(field string, direction string) By {
	return By{
		Field:     field,
		Direction: direction,
	}
}

func Parse(c *gin.Context, defaultOrder By) By {
	sortBy := c.Query("sortBy")
	orderBy := c.DefaultQuery("sortBy", ASC)

	if sortBy == "" {
		return defaultOrder
	}

	return NewBy(sortBy, orderBy)
}
