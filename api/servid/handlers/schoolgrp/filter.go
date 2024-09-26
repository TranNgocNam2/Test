package schoolgrp

import (
	"Backend/business/core/school"

	"github.com/gin-gonic/gin"
)

func parseFilter(ctx *gin.Context) (school.QueryFilter, error) {
	const (
		filterByName = "name"
	)

	var filter school.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	return filter, nil
}
