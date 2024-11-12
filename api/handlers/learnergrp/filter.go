package learnergrp

import (
	"Backend/business/core/learner"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	filterByName   = "name"
	filterBySchool = "school"
	filterByStatus = "status"
)

func parseFilter(ctx *gin.Context) (learner.QueryFilter, error) {

	var filter learner.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithFullName(name)
	}

	if school := ctx.Query(filterBySchool); school != "" {
		filter.WithSchoolName(school)
	}

	if statusStr := ctx.Query(filterByStatus); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			return filter, err
		}
		filter.WithStatus(status)
	}

	return filter, nil
}
