package learnergrp

import (
	"Backend/business/core/learner"
	"github.com/gin-gonic/gin"
)

const (
	filterByName   = "name"
	filterBySchool = "school"
)

func parseFilter(ctx *gin.Context) (learner.QueryFilter, error) {

	var filter learner.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithFullName(name)
	}

	if school := ctx.Query(filterBySchool); school != "" {
		filter.WithSchoolName(school)
	}

	return filter, nil
}
