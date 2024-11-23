package usergrp

import (
	"Backend/business/core/user"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	filterByName   = "name"
	filterBySchool = "school"
	filterByStatus = "status"
)

func parseFilter(ctx *gin.Context) (user.QueryFilter, error) {

	var filter user.QueryFilter

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
		filter.WithStatus(int16(status))
	}

	return filter, nil
}
