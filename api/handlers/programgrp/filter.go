package programgrp

import (
	"Backend/business/core/program"
	"github.com/gin-gonic/gin"
)

const (
	filterByName = "name"
)

func parseFilter(ctx *gin.Context) (program.QueryFilter, error) {

	var filter program.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	return filter, nil
}
