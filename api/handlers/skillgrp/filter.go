package skillgrp

import (
	"Backend/business/core/skill"
	"github.com/gin-gonic/gin"
)

const (
	filterByName = "name"
)

func parseFilter(ctx *gin.Context) (skill.QueryFilter, error) {

	var filter skill.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	return filter, nil
}
