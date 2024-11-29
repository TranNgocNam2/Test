package page

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Page struct {
	Number int
	Size   int
}

func Parse(c *gin.Context) Page {
	page := c.DefaultQuery("page", "1")
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		return Page{
			Number: 1,
			Size:   10,
		}
	}

	size := c.DefaultQuery("size", "10")
	pageSize, err := strconv.Atoi(size)
	if err != nil {
		return Page{
			Number: pageNumber,
			Size:   10,
		}
	}

	return Page{
		Number: pageNumber,
		Size:   pageSize,
	}
}
