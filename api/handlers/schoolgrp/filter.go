package schoolgrp

import (
	"Backend/business/core/school"

	"github.com/gin-gonic/gin"
)

const (
	filterByName = "name"
)

func parseFilter(ctx *gin.Context) (school.QueryFilter, error) {

	var filter school.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	return filter, nil
}
