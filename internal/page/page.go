package page

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Page struct {
	Number int
	Size   int
}

func Parse(c *gin.Context) (Page, error) {
	page := c.DefaultQuery("page", "1")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return Page{}, err
	}

	size := c.DefaultQuery("size", "10")
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		return Page{}, err
	}

	return Page{
		Number: pageNumber,
		Size:   pageSize,
	}, nil
}
